`timescale 1ns/1ps

module ez_tb;

    parameter WIN       = 15;
    parameter WIN_SIZE  = WIN * WIN;
    parameter DATA_SIZE = 8;
    parameter IMG_W     = 64;  // reduced for simulation feasibility
    parameter MAX_DISP  = 16;
    parameter MAX_DISP_ARR = $clog2(MAX_DISP) - 1;

    parameter IN_WIDTH  = DATA_SIZE * IMG_W * WIN;

    reg  [IN_WIDTH-1:0] input_array;
    wire [MAX_DISP_ARR:0] output_disp;
    reg clk, rst;
    wire done;

    reg [7:0] img_block [0:IMG_W * WIN - 1];

    // DUT instantiation
    compute_max_disp #(
        .WIN(WIN),
        .DATA_SIZE(DATA_SIZE),
        .IMG_W(IMG_W),
        .MAX_DISP(MAX_DISP)
    ) dut (
        .input_array(input_array),
        .clk(clk),
        .rst(rst),
        .output_disp(output_disp),
        .done(done)
    );

    // Clock generation
    initial clk = 0;
    always #5 clk = ~clk;

    integer i;
    initial begin
        $readmemh("C:/Users/jaxie/Desktop/cpre488/steroscopic-hardware/compute_SAD.mem", img_block);

        // Reset
        rst = 1;
        #20;
        rst = 0;

        // Flatten 2D memory block into input_array
        for (i = 0; i < IMG_W * WIN; i = i + 1) begin
            input_array[i * DATA_SIZE +: DATA_SIZE] = img_block[i];
        end

        wait (done);

        $display("Done reached, last output_disp = %0d", output_disp);
        $finish;
    end

endmodule