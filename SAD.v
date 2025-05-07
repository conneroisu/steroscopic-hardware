// Recieves two arrays that represent the respective window around a centre pixel
// Expects the input as such (ex 3x3):
// [0,1,2,
//  3,4,5,
//  6,7,8]
// Outputs a calculated SAD value of that window

module SAD (
    input_array_a,
    input_array_b,
    sad
);

// input wire clk
parameter WIN = 3,
          WIN_SIZE = WIN * WIN;

input wire [WIN_SIZE - 1 : 0] input_array_a;
input wire [WIN_SIZE - 1 : 0] input_array_b;
output wire sad; 

reg [WIN_SIZE : 0] diff;

genvar i;
generate 
    for (i = 0; i < WIN_SIZE; i = i + 1) begin : abs_diff
        assign diff[i] = (input_array_a[i] >= input_array_b[i]) ? (input_array_a[i] - input_array_b[i]) : (input_array_b[i] - input_array_a[i]);
    end
endgenerate

// Can't find a way to generate the adder tree
assign sad = diff[0] + diff[1] + diff[2] + diff[3] + diff[4] + diff[5] + diff[6] + diff[7] + diff[8];

endmodule 