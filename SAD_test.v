module SAD_test;

    parameter WIN = 15;
    parameter WIN_SIZE = WIN * WIN;
    parameter DATA_SIZE = 8;
    parameter SAD_SIZE = $clog2(WIN_SIZE * ((1 << DATA_SIZE) - 1) + 1);
    parameter FLAT_WIDTH = DATA_SIZE * WIN_SIZE;

    reg  [DATA_SIZE-1:0] winL [0:WIN_SIZE-1];
    reg  [DATA_SIZE-1:0] winR [0:WIN_SIZE-1];
    wire [SAD_SIZE:0] sad;

    // Flattened wires to connect to DUT
    wire [FLAT_WIDTH-1:0] input_a_flat;
    wire [FLAT_WIDTH-1:0] input_b_flat;

    // Flatten the arrays
    genvar i;
    generate
        for (i = 0; i < WIN_SIZE; i = i + 1) begin : flatten_inputs
            assign input_a_flat[DATA_SIZE*i +: DATA_SIZE] = winL[i];
            assign input_b_flat[DATA_SIZE*i +: DATA_SIZE] = winR[i];
        end
    endgenerate

    // Instantiate the DUT
    SAD #(.WIN(15)) dut (
        .input_a(input_a_flat),
        .input_b(input_b_flat),
        .sad(sad)
    );

    integer j;
    integer expected_sad;

    initial begin
        // Load winL and winR values from external files
        $readmemh("C:/Users/jaxie/Desktop/cpre488/steroscopic-hardware/winL.mem", winL);
        $readmemh("C:/Users/jaxie/Desktop/cpre488/steroscopic-hardware/winR.mem", winR);

        // Wait for the assignments to propagate
        #10;

        // Compute expected SAD
        expected_sad = 0;
        for (j = 0; j < WIN_SIZE; j = j + 1) begin
            if (winL[j] > winR[j])
                expected_sad = expected_sad + (winL[j] - winR[j]);
            else
                expected_sad = expected_sad + (winR[j] - winL[j]);
        end

        $display("DUT SAD      = %0d", sad);
        $display("Expected SAD = %0d", expected_sad);

        if (sad !== expected_sad)
            $fatal("SAD mismatch");
        else
            $display("Test passed");

        $finish;
    end
endmodule