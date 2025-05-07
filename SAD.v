// Recieves two arrays that represent the respective window around a centre pixel
// Expects the input as such (ex 3x3):
// [0,1,2,
//  3,4,5,
//  6,7,8]
// Outputs a calculated SAD value of that window

module SAD #(
    parameter WIN = 3,
    parameter WIN_SIZE = WIN * WIN,
    parameter DATA_SIZE = 8
)(
    // input wire clk
    // Flattened input since 2001 doesn't allow for port arrays
    input wire [DATA_SIZE * WIN_SIZE - 1 : 0] input_a,
    input wire [DATA_SIZE * WIN_SIZE - 1 : 0] input_b,
    output wire [DATA_SIZE - 1 : 0] sad 
); 

wire [DATA_SIZE - 1 : 0] array_a [0 : WIN_SIZE - 1];
wire [DATA_SIZE - 1 : 0] array_b [0 : WIN_SIZE - 1];

wire [DATA_SIZE - 1 : 0] diff [0 : WIN_SIZE - 1];

// Unpacking flattened vector to arrays
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

// Can't find a way to generate the adder tree
assign sad = diff[0] + diff[1] + diff[2] + diff[3] + diff[4] + diff[5] + diff[6] + diff[7] + diff[8];

endmodule 