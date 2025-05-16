`timescale 1ns/1ps

module disparity_top_tb;
  //--------------------------------------------------------------------------
  //  PARAMETERS — must match your disparity_top
  //--------------------------------------------------------------------------
  localparam WIN        = 15;
  localparam DATA_SIZE  = 8;
  localparam IMG_W      = 128;
  localparam IMG_H      = 128;
  localparam MAX_DISP   = 64;
  localparam DISP_BITS  = 6; //$clog2(MAX_DISP);

  localparam TEST_CASES = 4;

  //--------------------------------------------------------------------------
  //  SIGNALS
  //--------------------------------------------------------------------------
  reg                     clk, rst;
  // Left image AXI-Stream slave
  reg        [DATA_SIZE-1:0] sL_tdata;
  reg                         sL_tvalid;
  wire                        sL_tready;
  reg                         sL_tlast;
  // Right image AXI-Stream slave
  reg        [DATA_SIZE-1:0] sR_tdata;
  reg                         sR_tvalid;
  wire                        sR_tready;
  reg                         sR_tlast;

  // Output AXI-Stream master
  wire       [DATA_SIZE-1:0] m_tdata;
  wire                        m_tvalid;
  reg                         m_tready = 1;
  wire                        m_tlast;

  // Memories loaded from .mem files
  reg [DATA_SIZE-1:0] memL [0:TEST_CASES*IMG_W*IMG_H-1];
  reg [DATA_SIZE-1:0] memR [0:TEST_CASES*IMG_W*IMG_H-1];
  reg [DATA_SIZE-1:0] expD [0:TEST_CASES*IMG_W*IMG_H-1];

  // Buffer to capture outputs
  reg [DATA_SIZE-1:0] gotD [0:IMG_W*IMG_H-1];

  integer tc, x, y, idx_in, idx_out, count;

  //--------------------------------------------------------------------------
  //  CLOCK
  //--------------------------------------------------------------------------
  initial clk = 0;
  always #5 clk = ~clk;  // 100 MHz

  //--------------------------------------------------------------------------
  //  DUT INSTANTIATION
  //--------------------------------------------------------------------------
  // Make sure your disparity_top has been modified to:
  //   input  wire [7:0] sL_tdata, sL_tvalid, sL_tready, sL_tlast,
  //   input  wire [7:0] sR_tdata, sR_tvalid, sR_tready, sR_tlast,
  //   output wire [7:0] m_tdata, m_tvalid, m_tlast, input m_tready.
  disparity_top #(
    .WIN      (WIN),
    .DATA_SIZE(DATA_SIZE),
    .MAX_DISP (MAX_DISP),
    .IMG_W    (IMG_W),
    .IMG_H    (IMG_H),
    .HALF_WIN (WIN>>1),
    .DISP_BITS(DISP_BITS)
  ) dut (
    .clk        (clk),
    .rst        (rst),

    .sL_tdata   (sL_tdata),
    .sL_tvalid  (sL_tvalid),
    .sL_tready  (sL_tready),
    .sL_tlast   (sL_tlast),

    .sR_tdata   (sR_tdata),
    .sR_tvalid  (sR_tvalid),
    .sR_tready  (sR_tready),
    .sR_tlast   (sR_tlast),

    .m_axis_tdata    (m_tdata),
    .m_axis_tvalid   (m_tvalid),
    .m_axis_tready   (m_tready),
    .m_axis_tlast    (m_tlast)
  );

  //--------------------------------------------------------------------------
  //  LOAD MEM FILES
  //--------------------------------------------------------------------------
  initial begin
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//Hardware_shit//mems//left_image_p.mem",  memL);
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//Hardware_shit//mems//right_image_p.mem", memR);
    $readmemh("C://Users//jaxie//Desktop//cpre488//steroscopic-hardware//Hardware_shit//mems//exp_disp_p.mem",    expD);
  end

  //--------------------------------------------------------------------------
  //  TEST SEQUENCE
  //--------------------------------------------------------------------------
  initial begin
    // 1) Reset
    rst = 1; #20;
    rst = 0; #20;

    for (tc = 0; tc < TEST_CASES; tc = tc + 1) begin
      // 2) Stream-in one patch (128×128) of left & right pixels
      idx_in = tc * IMG_W * IMG_H;
      for (y = 0; y < IMG_H; y = y + 1) begin
        for (x = 0; x < IMG_W; x = x + 1) begin
          // drive tdata & tvalid
          @(posedge clk);
          sL_tdata  <= memL[idx_in];
          sR_tdata  <= memR[idx_in];
          sL_tvalid <= 1;
          sR_tvalid <= 1;
          sL_tlast  <= (x == IMG_W-1);
          sR_tlast  <= (x == IMG_W-1);
          idx_in    <= idx_in + 1;
          // wait for handshake
          wait (sL_tready && sR_tready);
        end
      end
      // deassert valid after frame
      @(posedge clk);
      sL_tvalid <= 0;
      sR_tvalid <= 0;

      // 3) Capture exactly IMG_W×IMG_H output pixels
      idx_out = 0;
      count   = 0;
      while (count < IMG_W*IMG_H) begin
        @(posedge clk);
        if (m_tvalid) begin
          gotD[count] = m_tdata;
          count = count + 1;
        end
      end

      // 4) Compare against golden
      for (idx_out = 0; idx_out < IMG_W*IMG_H; idx_out = idx_out + 1) begin
        if (gotD[idx_out] !== expD[tc*IMG_W*IMG_H + idx_out]) begin
          $error("TC %0d: pixel %0d mismatch: got=%0d exp=%0d",
                 tc, idx_out, gotD[idx_out], expD[tc*IMG_W*IMG_H + idx_out]);
        end
      end
      $display("TC %0d passed", tc);

      // small pause between cases
      #100;
    end

    $display("All %0d test cases PASSED!", TEST_CASES);
    $finish;
  end

endmodule
