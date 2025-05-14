// Recieves two arrays that represent the respective window around a centre pixel
// Expects the input as such (ex 3x3):
// [0,1,2,
//  3,4,5,
//  6,7,8]
// Outputs a calculated SAD value of that window
module SAD #(
    parameter WIN = 15,
    parameter DATA_SIZE = 8,
    //\\ Otherwise parameters //\\
    parameter WIN_SIZE = WIN * WIN,
    parameter SAD_SIZE = 16; // $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1)
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
// MAX_DISP windows must be searched serially - chunked into G = MAX_DISP / DIPS_THREADS
module compute_max_disp #(
    parameter WIN = 15,             // Window size/block size; Win * Win
    parameter DATA_SIZE = 8,        // Data size in bits
    parameter IMG_W = 64,           // Imwage resolution width
    parameter MAX_DISP = 64,        // Max disparity for each pixel
    parameter DISP_THREADS = 8,    // Number of disparities in parallel (MAX_DISP has to be divisible by this)
    //\\ Calculated parameters //\\
    //\\ Otherwise parameters //\\
    parameter G = MAX_DISP / DISP_THREADS, // How many iterations we need to do
    parameter WIN_SIZE = WIN * WIN,
    parameter SAD_BITS = 16 // $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1),
    parameter DISP_BITS = 6; // $clog2(MAX_DISP),
    parameter CYCLE_BITS = 3; // $clog2(G),
    parameter IMG_W_ARR = 6; // $clog2(IMG_W),
    parameter CMP_IDX_T = 3; // $clog2(DISP_THREADS)
)(
    input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_L,
    input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_R,
    input wire clk,
    input wire rst, 
    input wire input_ready,
    input wire [IMG_W_ARR - 1 : 0] col_index,
    output reg [DISP_BITS - 1 : 0] output_disp,
    output reg done
);

    // FSM parameters
    reg [1:0] state, next_state;
    localparam IDLE = 2'd0;
    localparam COMPUTE = 2'd1;
    localparam COMPARE = 2'd2;
    localparam DONE = 2'd3;

    reg [$clog2(DISP_THREADS) - 1 : 0] cmp_idx;
    reg [SAD_BITS - 1 : 0] best_sad;
    reg [DISP_BITS - 1 : 0] best_disp;

    //// PREPROCESSING (Unpacking and shit)
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
    //// END PREPROCESSING

    // Logic to iterate FSM
    always @(posedge clk or posedge rst) begin
        if (rst) begin
            state <= IDLE;
            // next_state <= IDLE;
            cmp_idx <= 0;
            best_sad <= {SAD_BITS{1'b1}};
            best_disp <= 0;
            best_disp <= 0;
            group_idx <= 0;
            done <= 0;
        end else begin
            state <= next_state;
            case (state)
                IDLE: begin
                    // Do nothing in sequential
                    // Waiting for input_ready signal
                end // END IDLE

                COMPUTE: begin
                    // Do nothing - give SAD a cycle
                end

                COMPARE: begin
                    // Perform comparison
                    if (sad_val[cmp_idx] < best_sad) begin
                        best_sad <= sad_val[cmp_idx];
                        best_disp <= base_disp + cmp_idx;
                    end
                    
                    // Batch comparisons finished (End of batch or min pixel found)
                    if (cmp_idx == DISP_THREADS - 1 || best_sad == 0 ) begin
                        group_idx <= group_idx + 1;
                        cmp_idx <= 0;
                    end else begin
                        cmp_idx <= cmp_idx + 1;
                    end
                end // END COMPUTE

                DONE: begin
                    output_disp <= best_disp;
                    done <= 1;
                end // END COMPARE
            endcase
        end
    end 

    // Combinational FSM
    always @* begin
        // Default values
        next_state  = state;

        case(state)
            IDLE: begin
                if (input_ready) begin
                    next_state = COMPUTE;
                end
            end // END IDLE

            COMPUTE: begin
                // Do nothing - give SAD a cycle
                next_state = COMPARE;
            end

            COMPARE: begin

                // Batch comparisons finished (End of batch or min pixel found)
                if (cmp_idx == DISP_THREADS - 1 || best_sad == 0 ) begin
                    // Exit conditions for pixel (End of group or min pixel found)
                    if (group_idx == G - 1 || best_sad == 0) begin
                        next_state = DONE;
                    end else begin
                        next_state = COMPUTE;
                    end
                end 
            end // END PROC

            DONE: begin
                // No loop to IDLE, just wait for reset
            end // END COMPARE
        endcase

        end
endmodule

// module top_disparity #(
//     parameter WIN = 15,             // Window size/block size; Win * Win
//     parameter DATA_SIZE = 8,        // Data size in bits
//     parameter IMG_W = 64,           // Imwage resolution width
//     parameter MAX_DISP = 64,        // Max disparity for each pixel
//     parameter DISP_THREADS = 16,    // Number of disparities in parallel (MAX_DISP has to be divisible by this)
//     // Calculated parameters \\
//     parameter G = MAX_DISP / DISP_THREADS, // How many iterations we need to do
//     parameter WIN_SIZE = WIN * WIN,
//     parameter SAD_BITS = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1),
//     parameter DISP_BITS = $clog2(MAX_DISP),
//     parameter CYCLE_BITS = $clog2(G),
//     parameter IMG_W_ARR = $clog2(IMG_W),
//     parameter CMP_IDX_T = $clog2(DISP_THREADS)
// )(
//     input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_L,
//     input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_R,
//     input wire clk,
//     input wire rst, 
//     // input wire input_ready,
//     // input wire [IMG_W_ARR - 1 : 0] col_index,
//     output reg [DATA_SIZE * IMG_W : 0] output_row,
//     output reg done
// );





// endmodule