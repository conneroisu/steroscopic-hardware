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
    output wire [SAD_SIZE: 0] sad 
    ); 

    wire [DATA_SIZE - 1 : 0] array_a [0 : WIN_SIZE - 1];
    wire [DATA_SIZE - 1 : 0] array_b [0 : WIN_SIZE - 1];

    reg [SAD_SIZE : 0] sad_accum;
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

// // Assumes the number of rows loaded into 
// // design is >= half window + 1.
// // Performs SAD per row in parallel.
// module compute_SAD #(
//     parameter WIN = 15,
//     parameter WIN_SIZE = WIN * WIN,
//     // parameter MIN_NUM_ROWS = (WIN / 2) + 1,
//     parameter DATA_SIZE = 8,
//     parameter IMG_W = 640,
//     parameter IMG_H = 480,
//     parameter MAX_DISP = 64,
//     // // Max dec from window = window_size * (2^data_size - 1)
//     // parameter MAX_SIZE_ROW = 640 * ((1 << DATA_SIZE) - 1),
//     // // Dec to bin for output width
//     // parameter SAD_SIZE_ROW = $clog2(MAX_SIZE + 1) 
//     parameter SAD_SIZE = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1)
// )(
//     input wire [DATA_SIZE * IMG_W * WIN : 0] input_array,
//     output wire [DATA_SIZE * IMG_W - 1 : 0] output_row
// );

// wire [SAD_SIZE : 0] sad_value [0 : IMG_W - 1] [0 : MAX_DISP - 1];
// reg [DATA_SIZE - 1 : 0] max_value;

// // Makes IMG_W * MAX_DISP SAD engines in parallel
// // So, 640 * 64 = 40960 engines in this case
// // Probs need at least 510 LUTs for each engine
// // 53,200 LUT < 20,889,600 LUTs lmao
// genvar c, d;
// generate
//     for (c = 0; c < IMG_W; c = c + 1) begin : column_loop
//         max_value = 0;
//         for (d = 0; d < MAX_DISP; d = d + 1) begin : disparity_loop    
//             wire [DATA_SIZE*WIN_SIZE-1:0] winL_flat;
//             wire [DATA_SIZE*WIN_SIZE-1:0] winR_flat;

//             // fill winL_flat and winR_flat from img_block
//             // based on c, d, and current window

//             SAD #(.WIN(WIN), .DATA_SIZE(DATA_SIZE)) sad_inst (
//                 .input_a(winL_flat),
//                 .input_b(winR_flat),
//                 .sad(sad_value[c][d])
//             );

//             if (max_value < sad_value[c][d]) begin
//                 max_value = sad_value[c][d];
//             end

//         end
//         output_row[]
//     end
// endgenerate



// endmodule


module compute_SAD #(
    parameter WIN       = 15,
    parameter WIN_SIZE  = WIN * WIN,
    parameter DATA_SIZE = 8,
    parameter IMG_W     = 640,
    parameter MAX_DISP  = 64,
    parameter SAD_SIZE  = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1)
)(
    // Flattened block of WIN rows, each with IMG_W pixels
    input  wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array,
    output wire [DATA_SIZE * IMG_W - 1 : 0]       output_row
);

    // Flattened SAD outputs
    wire [SAD_SIZE:0] sad_value [0:IMG_W-1][0:MAX_DISP-1];

    // Best disparity result per column (to be assigned to output_row)
    reg [SAD_SIZE-1:0]   min_sad     [0:IMG_W-1];
    reg [$clog2(MAX_DISP)-1:0] best_disp  [0:IMG_W-1];

    // Unpack input_array into a 2D image block: img_block[row][col]
    wire [DATA_SIZE-1:0] img_block [0:WIN-1][0:IMG_W-1];

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

    // Instantiate SAD units for each column and disparity
    genvar x, d, w;
    generate
        for (x = 0; x < IMG_W; x = x + 1) begin : column_loop
            for (d = 0; d < MAX_DISP; d = d + 1) begin : disparity_loop
                wire [DATA_SIZE*WIN_SIZE-1:0] winL_flat;
                wire [DATA_SIZE*WIN_SIZE-1:0] winR_flat;

                // Build 15x15 windows around (row, col) for L and R
                for (w = 0; w < WIN_SIZE; w = w + 1) begin : pack_windows
                    localparam integer WY = w / WIN;
                    localparam integer WX = w % WIN;

                    wire [DATA_SIZE-1:0] pix_L = img_block[WY][x + WX - (WIN >> 1)];
                    wire [DATA_SIZE-1:0] pix_R = img_block[WY][x + WX - (WIN >> 1) - d];

                    assign winL_flat[DATA_SIZE * w +: DATA_SIZE] = (x + WX >= (WIN >> 1) && x + WX < IMG_W - (WIN >> 1)) ? pix_L : 0;
                    assign winR_flat[DATA_SIZE * w +: DATA_SIZE] = (x + WX >= (WIN >> 1) + d && x + WX < IMG_W - (WIN >> 1)) ? pix_R : 0;
                end

                SAD #(.WIN(WIN), .DATA_SIZE(DATA_SIZE)) sad_inst (
                    .input_a(winL_flat),
                    .input_b(winR_flat),
                    .sad(sad_value[x][d])
                );
            end
        end
    endgenerate

    // Min SAD and best disparity per column
    integer i, j;
    always @* begin
        for (i = 0; i < IMG_W; i = i + 1) begin
            min_sad[i]   = {SAD_SIZE{1'b1}}; // Max value
            best_disp[i] = 0;

            for (j = 0; j < MAX_DISP; j = j + 1) begin
                if (sad_value[i][j] < min_sad[i]) begin
                    min_sad[i]   = sad_value[i][j];
                    best_disp[i] = j;
                end
            end
        end
    end

    // Assign best disparity per column to output_row
    genvar out_col;
    generate
        for (out_col = 0; out_col < IMG_W; out_col = out_col + 1) begin : output_gen
            assign output_row[DATA_SIZE*out_col +: DATA_SIZE] = best_disp[out_col];
        end
    endgenerate

endmodule