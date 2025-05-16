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
