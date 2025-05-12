`timescale 1ns/1ps

module compute_max_disp_tb;
  //------------------------------------------------------------------------------
  // PARAMETERS: match your DUT (or shrink for quick tests)
  //------------------------------------------------------------------------------
  localparam WIN          = 15;
  localparam DATA_SIZE    = 8;
  localparam IMG_W        = 64;
  localparam MAX_DISP     = 64;
  localparam DISP_THREADS = 16;
  localparam G            = MAX_DISP / DISP_THREADS;
  localparam WIN_SIZE     = WIN * WIN;
  localparam SAD_BITS     = $clog2(WIN_SIZE * ((1<<DATA_SIZE)-1) + 1);
  localparam DISP_BITS    = $clog2(MAX_DISP);
  localparam col_index    = 0;

  localparam TEST_CASES   = 1;  // ← number of test vectors

  //------------------------------------------------------------------------------
  // SIGNAL DECLARATIONS
  //------------------------------------------------------------------------------
  reg clk = 0;
  reg rst;
  reg [DATA_SIZE*IMG_W*WIN-1:0] input_array_L;
  reg [DATA_SIZE*IMG_W*WIN-1:0] input_array_R;
  wire [DISP_BITS-1:0] output_disp;
  wire done;

  // Memory arrays for vectors
  reg [DATA_SIZE-1:0] memL [0:TEST_CASES*IMG_W*WIN-1];
  reg [DATA_SIZE-1:0] memR [0:TEST_CASES*IMG_W*WIN-1];
  reg [DISP_BITS-1:0] exp_disp [0:TEST_CASES-1];

  integer i, k;

  //------------------------------------------------------------------------------
  // CLOCK
  //------------------------------------------------------------------------------
  always #5 clk = ~clk;  // 100 MHz

  //------------------------------------------------------------------------------
  // DUT INSTANTIATION
  //------------------------------------------------------------------------------
  compute_max_disp #(
    .WIN(WIN),
    .DATA_SIZE(DATA_SIZE),
    .IMG_W(IMG_W),
    .MAX_DISP(MAX_DISP),
    .DISP_THREADS(DISP_THREADS)
  ) dut (
    .input_array_L(input_array_L),
    .input_array_R(input_array_R),
    .clk(clk),
    .rst(rst),
    .col_index(col_index),
    .output_disp(output_disp),
    .done(done)
  );

  //------------------------------------------------------------------------------
  // INITIAL: load all .mem files
  //------------------------------------------------------------------------------
  initial begin
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//left_image.mem",  memL);
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//right_image.mem", memR);
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//exp_disp.mem",    exp_disp);
  end

  //------------------------------------------------------------------------------
  // TEST SEQUENCE
  //------------------------------------------------------------------------------
  initial begin
    // give time for mem load
    #1;
    for (i = 0; i < TEST_CASES; i = i + 1) begin
      // pack flattened input_array from memL/R
      input_array_L = {DATA_SIZE*IMG_W*WIN{1'b0}};
      input_array_R = {DATA_SIZE*IMG_W*WIN{1'b0}};
      for (k = 0; k < IMG_W*WIN; k = k + 1) begin
        input_array_L[k*DATA_SIZE +: DATA_SIZE] = memL[i*(IMG_W*WIN) + k];
        input_array_R[k*DATA_SIZE +: DATA_SIZE] = memR[i*(IMG_W*WIN) + k];
      end

      // pulse reset
      rst = 1; #10;
      rst = 0;

      // wait for done
      wait (done);
      #1;
      if (output_disp !== exp_disp[i]) begin
        $error("Test %0d FAILED: got %0d, expected %0d", i, output_disp, exp_disp[i]);
      end else begin
        $display("Test %0d PASS: disp=%0d", i, output_disp);
      end
      // small pause before next vector
      #5;
    end

    $display("All %0d tests completed.", TEST_CASES);
    $finish;
  end

endmodule