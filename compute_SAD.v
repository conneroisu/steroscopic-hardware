// Recieves two arrays that represent the respective window around a centre pixel
// Expects the input as such (ex 3x3):
// [0,1,2,
//  3,4,5,
//  6,7,8]
// Outputs a calculated SAD value of that window
module SAD #(
    parameter WIN = 15,
    parameter WIN_SIZE = WIN * WIN,
    parameter DATA_SIZE = 8,
    // Max dec from window = window_size * (2^data_size - 1)
    // parameter MAX_SIZE = WIN_SIZE * ((1 << DATA_SIZE) - 1),
    // // Dec to bin for output width
    // parameter SAD_SIZE = $clog2(MAX_SIZE + 1) 
    parameter SAD_SIZE = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1)
    )(
    // input wire clk
    // Flattened input since 2001 doesn't allow for port arrays
    input wire [DATA_SIZE * WIN_SIZE - 1 : 0] input_a,
    input wire [DATA_SIZE * WIN_SIZE - 1 : 0] input_b,
    output wire [SAD_SIZE - 1 : 0] sad 
    ); 

    wire [DATA_SIZE - 1 : 0] array_a [0 : WIN_SIZE - 1];
    wire [DATA_SIZE - 1 : 0] array_b [0 : WIN_SIZE - 1];

    reg [SAD_SIZE - 1 : 0] sad_accum;
    wire [DATA_SIZE - 1 : 0] diff [0 : WIN_SIZE - 1];

    // Unpacking flattened vector to arrays
    // Assuming data dependencies forces unpacking to occur before SAD comp
    genvar r;
    generate
        for (r = 0; r < WIN_SIZE; r = r + 1) begin : unpack
            assign array_a[r] = input_a[DATA_SIZE * r +: DATA_SIZE];
            assign array_b[r] = input_b[DATA_SIZE * r +: DATA_SIZE];
        end
    endgenerate

    // Compute SAD outputs in parallel
    genvar i;
    generate 
        for (i = 0; i < WIN_SIZE; i = i + 1) begin : abs_diff
            assign diff[i] = (array_a[i] >= array_b[i]) ? (array_a[i] - array_b[i]) : (array_b[i] - array_a[i]);
        end
    endgenerate

    
    // Wildcard should only be looking at diff[w]
    // Generates sum
    integer w;
    always @* begin
        sad_accum = 0;
        for (w = 0; w < WIN_SIZE; w = w + 1) 
            sad_accum = sad_accum + diff[w];
    end

    // Pass output to wire
    assign sad = sad_accum; 

endmodule 

// Compute disparity of a pixel
// Perform SAD computation in parallel
// MAX_DISP windows must be searched serially - could be chunked later
module compute_max_disp #(
    parameter WIN = 15,
    parameter DATA_SIZE = 8,
    parameter IMG_W = 640,
    parameter MAX_DISP = 64,
    parameter DISP_PARTS = 8,
    localparam G = MAX_DISP / DISP_PARTS,
    localparam WIN_SIZE = WIN * WIN,
    localparam SAD_BITS = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1),
    // localparam MAX_DISP_ARR = $clog2(MAX_DISP) - 1,
    localparam DISP_BITS = $clog2(MAX_DISP),
    localparam CYCLE_BITS = $clog2(G),
    localparam IMG_W_ARR = $clog2(IMG_W)
)(
    // Flattened block of WIN rows, each with IMG_W pixels
    input  wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array,
    input clk,
    input rst, 
    output wire [DISP_BITS - 1 : 0] output_disp,
    output reg done
);

    // Unpack input_array into a 2D image block: img_block[row][col]
    wire [DATA_SIZE - 1 : 0] img_block [0 : WIN - 1][0 : IMG_W - 1];
    reg [IMG_W_ARR : 0] col_index;

    genvar r, c;
    generate
        for (r = 0; r < WIN; r = r + 1) begin : row_unpack
            for (c = 0; c < IMG_W; c = c + 1) begin : col_unpack
                assign img_block[r][c] = input_array[
                    DATA_SIZE * (r * IMG_W + c) +: DATA_SIZE
                ];
            end
        end
    endgenerate

    // Create window buffers for # of disparities
    // Partition the windows 
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winL_flat [0 : DISP_PARTS - 1][0 : G - 1];
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winR_flat [0 : DISP_PARTS - 1][0 : G - 1];

    // Flattened SAD outputs
    wire [SAD_SIZE - 1 : 0] sad_value [0 : DISP_PARTS - 1];

    // Indicies and local minima
    integer i;
    reg [IMG_W_ARR - 1 : 0] col_index;
    reg [CYCLE_BITS - 1 : 0] cycle;
    reg [SAD_BITS - 1 : 0] local_min [0 : DISP_PARTS - 1];
    reg [DISP_BITS - 1 : 0] local_disp [0 : DISP_PARTS - 1];
    reg [SAD_BITS - 1 : 0] best_value;
    reg [DISP_BITS - 1 : 0] best_disp;
    
    // Instantiate SAD units for each column and disparity
    genvar x, w;
    generate

        // Iterate through disparities for a pixel
        for (x = 0; x < MAX_DISP; x = x + 1) begin : disparity_loop

            // Build 15x15 windows
            for (w = 0; w < WIN_SIZE; w = w + 1) begin : pack_windows
                // Not sure how these synthesize yet
                // Determine X, Y coordinates for Window
                localparam integer WY = w / WIN;
                localparam integer WX = w % WIN;

                // Find pixel values (and check boundaries)
                // col_index + WX >= (WIN >> 1) prevents underflow [pixel beyond left boundary]
                // col_index + WX < IMG_W - (WIN >> 1) prevents overflow [pixel beyond right boundary]
                wire [DATA_SIZE - 1 : 0] pix_L = ((col_index + WX >= (WIN >> 1)) && (col_index + WX < IMG_W - (WIN >> 1))) ? img_block[WY][col_index + WX - (WIN >> 1)] : 0;
                // Right image offset by disparity
                wire [DATA_SIZE - 1 : 0] pix_R = ((col_index + WX >= (WIN >> 1) + x) && (col_index + WX < IMG_W - (WIN >> 1))) ? img_block[WY][col_index + WX - (WIN >> 1) - x] : 0;

                // Assign to flatten arry to be passed to SAD 
                assign winL_flat[x][DATA_SIZE * w +: DATA_SIZE] = pix_L;
                assign winR_flat[x][DATA_SIZE * w +: DATA_SIZE] = pix_R;
            end
        end
    endgenerate

    // Max disp now in parallel - iterate pixel rightwards every clk cycle
    integer q;
    always @(posedge clk or posedge rst) begin
        min_sad <= {SAD_SIZE{1'b1}};
        best_disp <= 0;

        if (rst) begin  
            col_index <= 0;
            done <= 0;
        end else if (!done) begin
            for (q = 0; q < MAX_DISP; q = q + 1) begin
                if (sad_value[q] < min_sad[q]) begin
                    min_sad <= sad_value[q];
                    best_disp <= q;
                end
            end
            
            col_index <= col_index + 1;
            if (col_index == IMG_W) begin
                done <= 1;
            end
        end
    end

    assign output_disp = best_disp;

    // // Min SAD and best disparity per column
    // integer i, j;
    // always @* begin
    //     for (i = 0; i < IMG_W; i = i + 1) begin
    //         min_sad[i] = {SAD_SIZE{1'b1}}; // Max value
    //         best_disp[i] = 0;

    //         for (j = 0; j < MAX_DISP; j = j + 1) begin
    //             if (sad_value[i][j] < min_sad[i]) begin
    //                 min_sad[i] = sad_value[i][j];
    //                 best_disp[i] = j;
    //             end
    //         end
    //     end
    // end

    // Assign best disparity per column to output_row
    // genvar out_col;
    // generate
    //     for (out_col = 0; out_col < IMG_W; out_col = out_col + 1) begin : output_gen
    //         assign output_row[DATA_SIZE * out_col +: DATA_SIZE] = best_disp[out_col];
    //         // assign output_row[DATA_SIZE*out_col +: DATA_SIZE] = disparity;
    //     end
    // endgenerate

endmodule