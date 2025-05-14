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
    parameter WIN_SIZE = 225,    // WIN * WIN,
    parameter SAD_SIZE = 16     // $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1)
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
    parameter MAX_DISP = 64,        // Max disparity for each pixel
    //\\ Calculated parameters //\\
    //\\ Otherwise localparam //\\
    parameter WIN_SIZE = 225,       // WIN * WIN,
    parameter SAD_BITS = 16,        // $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1),
    parameter DISP_BITS = 6        // $clog2(MAX_DISP),
)(
    input wire [DATA_SIZE * MAX_DISP * WIN - 1 : 0] input_array_L,
    input wire [DATA_SIZE * MAX_DISP * WIN - 1 : 0] input_array_R,
    input wire clk,
    input wire rst, 
    input wire input_ready,
    input wire [DISP_BITS - 1 : 0] col_index,
    output reg [DISP_BITS - 1 : 0] output_disp,
    output reg done
);

    // FSM parameters
    reg [2:0] state, next_state;
    localparam IDLE = 3'd0;
    localparam PREPROC = 3'd1;
    localparam COMPUTE = 3'd2;
    localparam COMPARE = 3'd3;
    localparam DONE = 3'd4;

    reg [3:0] preload_cnt;

    // reg [CMP_IDX_T - 1 : 0] cmp_idx;
    reg [SAD_BITS - 1 : 0] best_sad;
    reg [DISP_BITS - 1 : 0] best_disp;
    reg [DISP_BITS - 1 : 0] disp_idx;

    //// PREPROCESSING (Unpacking and shit)
    // Unpack input_array into a 2D image block: img_block[row][col]
    wire [DATA_SIZE - 1 : 0] img_block_L [0 : WIN - 1][0 : MAX_DISP - 1];
    wire [DATA_SIZE - 1 : 0] img_block_R [0 : WIN - 1][0 : MAX_DISP - 1];

    genvar r, c;
    generate
        for (r = 0; r < WIN; r = r + 1) begin : row_unpack
            for (c = 0; c < MAX_DISP; c = c + 1) begin : col_unpack
                assign img_block_L[r][c] = input_array_L[
                    DATA_SIZE * (r * MAX_DISP + c) +: DATA_SIZE
                ];
                assign img_block_R[r][c] = input_array_R[
                    DATA_SIZE * (r * MAX_DISP + c) +: DATA_SIZE
                ];
            end
        end
    endgenerate

    // Sliding window FIFOs for each row
    reg [DATA_SIZE - 1 : 0] winL_buf [0 : WIN - 1][0 : WIN - 1];
    reg [DATA_SIZE - 1 : 0] winR_buf [0 : WIN - 1][0 : WIN - 1];

    // Flatten windows
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winL_flat;
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winR_flat; 
    generate
        for (r = 0; r < WIN; r = r + 1) begin : winL_row_unpack
            for (c = 0; c < WIN; c = c + 1) begin : winL_col_unpack
                assign winL_flat[DATA_SIZE * (r * WIN + c) +: DATA_SIZE] = winL_buf[r][c];
                assign winR_flat[DATA_SIZE * (r * WIN + c) +: DATA_SIZE] = winR_buf[r][c];
            end
        end
    endgenerate

    // Create SAD units for each thread
    wire [SAD_BITS - 1 : 0] sad_val;
    SAD #(
        .WIN(WIN),
        .DATA_SIZE(DATA_SIZE)
    ) sad_inst (
        .input_a(winL_flat),
        .input_b(winR_flat),
        .sad(sad_val)
    );
    //// END PREPROCESSING

    integer row, col;
    always @(posedge clk or posedge rst) begin
        if (rst) begin 
            state <= IDLE;
            disp_idx <= 0;
            best_sad <= {SAD_BITS{1'b1}};
            best_disp <= 0;
            done <= 0;
            preload_cnt <= 0;
            // Clean shift reg buffers
            for (row = 0; row < WIN; row = row + 1) begin
                for (col = 0; col < WIN; col = col + 1) begin
                    winL_buf[row][col] <= {DATA_SIZE{1'b0}};
                    winR_buf[row][col] <= {DATA_SIZE{1'b0}};
                end
            end 
        end else begin 
            state <= next_state;  
            case(state)
                IDLE: begin
                    done <= 0;
                    if (input_ready) begin
                        disp_idx <= 0;
                        best_sad <= {SAD_BITS{1'b1}};
                        best_disp <= 0;
                        preload_cnt <= 0;
                    end
                end // END IDLE

                PREPROC: begin
                    if (preload_cnt < WIN) begin
                    
                        for (row = 0; row < WIN; row = row + 1) begin
                            for (col = 0; col < WIN - 1; col = col + 1) begin
                                winL_buf[row][col] <= winL_buf[row][col + 1];
                                winR_buf[row][col] <= winR_buf[row][col + 1];
                            end
                            winL_buf[row][WIN - 1] <= img_block_L[row][col_index + preload_cnt];
                            winR_buf[row][WIN - 1] <= img_block_R[row][col_index + preload_cnt]; // disp_idx == 0
                        end
                        preload_cnt <= preload_cnt + 1;

                        if (preload_cnt == 0) begin
                            // Reset bests
                            disp_idx  <= 0;
                            best_sad  <= {SAD_BITS{1'b1}};
                            best_disp <= 0; 
                        end

                    end
                end // END PREPROC

                COMPUTE: begin
                    // Do nothing - give SAD a cycle
                end // END COMPUTE

                COMPARE: begin 
                    if (sad_val < best_sad) begin
                        best_sad <= sad_val;
                        best_disp <= disp_idx;
                    end

                    if (disp_idx < MAX_DISP - 1 && sad_val != 0) begin
                        disp_idx <= disp_idx + 1;
                        // Shift left
                        for (row = 0; row < WIN; row = row + 1) begin
                            for (col = 0; col < WIN - 1; col = col + 1) begin
                                winR_buf[row][col] <= winR_buf[row][col + 1];
                            end
                            winR_buf[row][WIN - 1] <= img_block_R[row][col_index + disp_idx + 1];
                        end
                    end
                end // END COMPARE

                DONE: begin 
                    output_disp <= best_disp - (WIN - 1);
                    done <= 1;
                end

                default: ; // TODO LATER
            endcase
        end
    end

    always @* begin
        next_state = state;

        case(state)
            IDLE: begin
                if (input_ready) begin
                    next_state = PREPROC;
                end
            end

            PREPROC: begin
                if (preload_cnt == WIN) begin
                    next_state = COMPUTE;
                end else begin
                    next_state = PREPROC;
                end
            end

            COMPUTE: begin
                next_state = COMPARE;
            end

            COMPARE: begin
                // Batch comparisons finished (End of batch or min pixel found)
                if (disp_idx == MAX_DISP - 1 || best_sad == 0) begin
                    next_state = DONE; 
                end else begin
                    next_state = COMPUTE;
                end
            end

            DONE: begin
                // No loop to IDLE, just wait for reset
            end // END DONE
        endcase
    end

endmodule

//// Read inputs from framebuffer and compile disparity map
//module compute_max_disp #(
//    parameter WIN = 15,             // Window size/block size; Win * Win
//    parameter DATA_SIZE = 8,        // Data size in bits
//    parameter IMG_W = 640,
//    parameter IMG_H = 480
//    parameter MAX_DISP = 64,        // Max disparity for each pixel
//    //\\ Calculated parameters  //\\
//    //\\ Otherwise localparam   //\\
//    parameter WIN_SIZE = 225,       // WIN * WIN,
//    parameter SAD_BITS = 16,        // $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1),
//    parameter DISP_BITS = 6,        // $clog2(MAX_DISP),
//    // parameter IMG_W_ARR = 6         // $clog2(IMG_W),
//)(
//    input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_L,
//    input wire [DATA_SIZE * IMG_W * WIN - 1 : 0] input_array_R,
//    input wire clk,
//    input wire rst, 
//    input wire input_ready,
//    input wire [IMG_W_ARR - 1 : 0] col_index,
//    output reg [DISP_BITS - 1 : 0] output_disp,
//    output reg done
//);



//endmodule 