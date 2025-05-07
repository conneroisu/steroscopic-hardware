module SAD_test;

    reg  [7:0] winL [0:8];
    reg  [7:0] winR [0:8];
    wire [7:0] sad;

    // Flattened wires to connect to DUT
    wire [71:0] input_a_flat;
    wire [71:0] input_b_flat;

    // Flatten the arrays
    genvar i;
    generate
        for (i = 0; i < 9; i = i + 1) begin : flatten_inputs
            assign input_a_flat[8*i +: 8] = winL[i];
            assign input_b_flat[8*i +: 8] = winR[i];
        end
    endgenerate

    // Instantiate the DUT
    SAD #(.WIN(3)) dut (
        .input_a(input_a_flat),
        .input_b(input_b_flat),
        .sad(sad)
    );

    initial begin
        winL[0] = 0;   winL[1] = 10;  winL[2] = 20;
        winL[3] = 30;  winL[4] = 0;   winL[5] = 10;
        winL[6] = 20;  winL[7] = 30;  winL[8] = 0;

        winR[0] = 110; winR[1] = 100; winR[2] = 90;
        winR[3] = 80;  winR[4] = 70;  winR[5] = 60;
        winR[6] = 50;  winR[7] = 40;  winR[8] = 30;

        #10;
        $display("SAD output = %d", sad);  // should be 510
        $display("Expected     = 510");
        $finish;
    end
endmodule