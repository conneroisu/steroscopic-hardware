`timescale 1ns/1ps

module compute_max_disp_tb;
  //--------------------------------------------------------------------------
  //  PARAMETERS
  //--------------------------------------------------------------------------
  localparam WIN          = 15;
  localparam DATA_SIZE    = 8;
  localparam IMG_W        = 64;
  localparam MAX_DISP     = 64;
  localparam DISP_THREADS = 16;
  localparam G            = MAX_DISP/DISP_THREADS;
  localparam WIN_SIZE     = WIN*WIN;
  localparam SAD_BITS     = $clog2(WIN_SIZE*((1<<DATA_SIZE)-1)+1);
  localparam DISP_BITS    = $clog2(MAX_DISP);
  localparam IMG_W_ARR    = $clog2(IMG_W);

  localparam TEST_CASES   = 4;

  //--------------------------------------------------------------------------
  //  SIGNALS
  //--------------------------------------------------------------------------
  reg                                 clk;
  reg                                 rst;
  reg                                 input_ready;
  reg  [DATA_SIZE*IMG_W*WIN-1 : 0]    input_array_L;
  reg  [DATA_SIZE*IMG_W*WIN-1 : 0]    input_array_R;
  wire [DISP_BITS-1 : 0]              output_disp;
  wire                                done;

  // Memory arrays
  reg [DATA_SIZE-1:0] memL [0:TEST_CASES*IMG_W*WIN-1];
  reg [DATA_SIZE-1:0] memR [0:TEST_CASES*IMG_W*WIN-1];
  reg [DISP_BITS-1:0] exp_disp [0:TEST_CASES-1];
  reg [IMG_W_ARR - 1 : 0] col_index; 
  integer i, k;

  //--------------------------------------------------------------------------
  //  CLOCK
  //--------------------------------------------------------------------------
  initial clk = 0;
  always #5 clk = ~clk;  // 100 MHz

  //--------------------------------------------------------------------------
  //  DUT
  //--------------------------------------------------------------------------
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
    .input_ready(input_ready),
    .col_index(col_index),
    .output_disp(output_disp),
    .done(done)
  );

  //--------------------------------------------------------------------------
  //  LOAD MEM FILES
  //--------------------------------------------------------------------------
  initial begin
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//left_image.mem",  memL);
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//right_image.mem", memR);
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//exp_disp.mem",    exp_disp);
  end

  //--------------------------------------------------------------------------
  //  TEST SEQUENCE
  //--------------------------------------------------------------------------
  initial begin
    rst = 1; input_ready = 0;
    #20;
    rst = 0;
    #20;

    col_index = 0;

    for (i = 0; i < TEST_CASES; i = i + 1) begin
      // Pack L/R arrays
      input_array_L = {DATA_SIZE*IMG_W*WIN{1'b0}};
      input_array_R = {DATA_SIZE*IMG_W*WIN{1'b0}};
      for (k = 0; k < IMG_W*WIN; k = k + 1) begin
        input_array_L[k*DATA_SIZE +: DATA_SIZE] = memL[i*(IMG_W*WIN) + k];
        input_array_R[k*DATA_SIZE +: DATA_SIZE] = memR[i*(IMG_W*WIN) + k];
      end

      // Start one‐shot
      input_ready = 1;
      #10;
      input_ready = 0;

      // Wait for done
      wait (done);
      #1;
      if (output_disp !== exp_disp[i]) begin
        $error("FAIL case %0d: got %0d expected %0d", 
               i, output_disp, exp_disp[i]);
      end else begin
        $display("PASS case %0d: got %0d and expected %0d", i, output_disp, exp_disp[i]);
      end

      // Clear done (if your module holds it)
      #10;
      rst = 1; #10; rst = 0;
      #10;
    end

    $display("All %0d test cases completed.", TEST_CASES);
    $finish;
  end
endmodule
