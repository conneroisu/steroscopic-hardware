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
    parameter SAD_SIZE = 16,     // $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1)
    parameter THREADS = 15
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
    parameter IMG_W = 64,           // Imwage resolution width // TODO: MIGHT HAVE TO CHANGE LATER TO CHUNK SIZE
    parameter MAX_DISP = 64,        // Max disparity for each pixel
    //\\ Calculated parameters //\\
    //\\ Otherwise localparam //\\
    parameter WIN_SIZE = 225,       // WIN * WIN,
    parameter SAD_BITS = 16,        // $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1),
    parameter DISP_BITS = 6,        // $clog2(MAX_DISP),
    parameter IMG_W_ARR = 6         // $clog2(IMG_W),
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

    // reg [CMP_IDX_T - 1 : 0] cmp_idx;
    reg [SAD_BITS - 1 : 0] best_sad;
    reg [DISP_BITS - 1 : 0] best_disp;
    reg [DISP_BITS - 1 : 0] disp_idx;

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

    // Flatten windows
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winL_flat;
    wire [DATA_SIZE * WIN_SIZE - 1 : 0] winR_flat; 
    generate
        for (r = 0; r < WIN; r = r + 1) begin : winL_row_unpack
            for (c = 0; c < WIN; c = c + 1) begin : winL_col_unpack
                // Create window reference
                assign winL_flat[DATA_SIZE * (r * WIN + c) +: DATA_SIZE] = 
                    img_block_L[r][col_index + c];

                // Shift right window by disp_idx
                assign winR_flat[ DATA_SIZE * (r * WIN + c) +: DATA_SIZE ] =
                    img_block_R[r][ col_index + c + disp_idx ];
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

    // Logic to iterate FSM
    always @(posedge clk or posedge rst) begin
        if (rst) begin
            state <= IDLE;
            disp_idx <= 0;
            best_sad <= {SAD_BITS{1'b1}};
            best_disp <= 0;
            done <= 0;
        end else begin
            state <= next_state;
            case (state)
                IDLE: begin
                    done <= 0;
                    if (input_ready) begin
                        disp_idx <= 0;
                        best_sad <= {SAD_BITS{1'b1}};
                        best_disp <= 0;
                    end
                end // END IDLE

                COMPUTE: begin
                    // Do nothing - give SAD a cycle
                end

                COMPARE: begin
                    // Perform comparison
                    if (sad_val < best_sad) begin
                        best_sad <= sad_val;
                        best_disp <= disp_idx;
                    end

                    // Batch comparisons finished (End of batch or min pixel found)
                    if (disp_idx == MAX_DISP - 1 || sad_val == 0) begin
                        // Go to DONE state
                    end else begin
                        disp_idx <= disp_idx + 1;
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
                if (disp_idx == MAX_DISP - 1 || sad_val == 0) begin
                    next_state = DONE; 
                end else begin
                    next_state = COMPUTE;
                end
            end // END PROC

            DONE: begin
                // No loop to IDLE, just wait for reset
            end // END COMPARE
        endcase

        end
endmodule