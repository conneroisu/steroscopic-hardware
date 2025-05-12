// Recieves two arrays that represent the respective window around a centre pixel
// Expects the input as such (ex 3x3):
// [0,1,2,
//  3,4,5,
//  6,7,8]
// Outputs a calculated SAD value of that window
module SAD #(
    parameter WIN = 15,
    localparam WIN_SIZE = WIN * WIN,
    parameter DATA_SIZE = 8,
    // Max dec from window = window_size * (2^data_size - 1)
    // parameter MAX_SIZE = WIN_SIZE * ((1 << DATA_SIZE) - 1),
    // // Dec to bin for output width
    // parameter SAD_SIZE = $clog2(MAX_SIZE + 1) 
    localparam SAD_SIZE = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1)
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
/// I think the likelyhood of an earlier disparity map to match is pretty high - so probs just do a simple comparator for now
///     and exit if 0 is found
module compute_max_disp #(
    parameter WIN = 15,             // Window size/block size; Win * Win
    parameter DATA_SIZE = 8,        // Data size in bits
    parameter IMG_W = 64,          // Imwage resolution width
    parameter MAX_DISP = 64,        // Max disparity for each pixel
    parameter DISP_THREADS = 16,    // Number of disparities in parallel (MAX_DISP has to be divisible by this)
    // Calculated parameters \\
    localparam G = MAX_DISP / DISP_THREADS,
    localparam WIN_SIZE = WIN * WIN,
    localparam SAD_BITS = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1),
    localparam DISP_BITS = $clog2(MAX_DISP),
    localparam CYCLE_BITS = $clog2(G),
    localparam IMG_W_ARR = $clog2(IMG_W)
)(
    input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_L,
    input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_R,
    input wire clk,
    input wire rst, 
    input reg col_index, 
    output reg [DISP_BITS - 1 : 0] output_disp,
    output reg done
);

    // Unpack input_array into a 2D image block: img_block[row][col]
    wire [DATA_SIZE - 1 : 0] img_block_L [0 : WIN - 1][0 : IMG_W - 1];
    wire [DATA_SIZE - 1 : 0] img_block_R [0 : WIN - 1][0 : IMG_W - 1];

    genvar r, c;
    generate
        for (r = 0; r < WIN; r = r + 1) begin : row_unpack
            for (c = 0; c < IMG_W; c = c + 1) begin : col_unpack
                assign img_block_L[r][c] = input_array_L[
                    DATA_SIZE * (r * IMG_W + c) +: DATA_SIZE
                ];
                assign img_block_R[r][c] = input_array_R[
                    DATA_SIZE * (r * IMG_W + c) +: DATA_SIZE
                ];
            end
        end
    endgenerate
    
    // Flatten refence (left in our case) window buffer
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winL_flat; 
    
    generate
        for (r = 0; r < WIN; r = r + 1) begin : winL_row_unpack
            for (c = 0; c < WIN; c = c + 1) begin : winL_col_unpack
                assign winL_flat[DATA_SIZE * (r * WIN + c) +: DATA_SIZE] = 
                    img_block_L[r][col_index + c];
            end
        end    
    endgenerate

    // Create windows per DISP_THREADS, itereate by WIN_SIZE pixels shifted by base disp + d
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winR_flat [0 : DISP_THREADS - 1];
    reg [CYCLE_BITS - 1 : 0] group_idx;
    wire [DISP_BITS - 1 : 0] base_disp = group_idx * DISP_THREADS;
    
    genvar d;
    generate
        for (d = 0; d < DISP_THREADS; d = d + 1) begin : winR_flatten
            for (r = 0; r < WIN; r = r + 1) begin : winR_row_unpack
                for (c = 0; c < WIN; c = c + 1) begin : winR_col_unpack
                    assign winR_flat[d][DATA_SIZE * (r * WIN + c) +: DATA_SIZE] = 
                        img_block_R[r][col_index + c + base_disp + d];
                end
            end   
        end
    endgenerate

    // Create SAD units for each thread
    wire [SAD_BITS - 1 : 0] sad_val [0 : DISP_THREADS - 1];
    generate
        for (d = 0; d < DISP_THREADS; d = d + 1) begin : create_saddness
            SAD #(
                .WIN(WIN),
                .DATA_SIZE(DATA_SIZE)
            ) sad_inst (
                .input_a(winL_flat),
                .input_b(winR_flat[d]),
                .sad(sad_val[d])
            );
        end
    endgenerate

    // Find minimum SAD value between disparities
    reg [SAD_BITS - 1 : 0] best_sad;
    reg [DISP_BITS - 1 : 0] best_disp;
    integer i;
    always @(posedge clk) begin
        if (rst) begin
            group_idx <= 0;
            best_sad <= {SAD_BITS{1'b1}};
            best_disp <= 0;
            done <= 0;
        end else begin
            // Compare batch of Threads
            for (i = 0; i < DISP_THREADS; i = i + 1) begin
                if (sad_val[i] < best_sad) begin
                    best_sad <= sad_val[i];
                    best_disp <= base_disp + i;
                end
            end

            // Once last group was processed, exit
            // Or if best_disp finds a min value - 0 in this case
            if (group_idx == G-1 || best_sad == 0) begin
                output_disp <= best_disp;
                done <= 1;
                group_idx <= 0;
            end else begin
                group_idx <= group_idx + 1;
            end

        end
    end


endmodule