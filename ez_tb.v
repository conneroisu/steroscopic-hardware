`timescale 1ns/1ps

module ez_tb;

    parameter WIN       = 15;
    parameter WIN_SIZE  = WIN * WIN;
    parameter DATA_SIZE = 8;
    parameter IMG_W     = 1;
    parameter MAX_DISP  = 3;

    parameter IN_WIDTH  = DATA_SIZE * IMG_W * WIN;
    parameter OUT_WIDTH = DATA_SIZE * IMG_W;

    reg  [IN_WIDTH-1:0] input_array;
    wire [OUT_WIDTH-1:0] output_row;

    // DUT instantiation
    compute_SAD #(
        .WIN(WIN),
        .DATA_SIZE(DATA_SIZE),
        .IMG_W(IMG_W),
        .MAX_DISP(MAX_DISP)
    ) dut (
        .input_array(input_array),
        .output_row(output_row)
    );

    integer i, d;
    reg [7:0] imgL [0:WIN-1];
    reg [7:0] imgR [0:WIN-1][0:MAX_DISP-1];
    integer sad;
    integer min_sad;
    integer best_disp;

    initial begin
        // -------------------------------
        // Test Case Setup
        // imgL = 0,1,2,...14
        // imgR[d] = shifted variants
        // -------------------------------
        for (i = 0; i < WIN; i = i + 1) begin
            imgL[i] = i;
            imgR[i][0] = i;               // SAD = 0
            imgR[i][1] = WIN - 1 - i;     // SAD = larger
            imgR[i][2] = 0;               // SAD = large
        end

        // Flatten into input_array:
        // Left image is row WIN/2 (middle row), fixed
        // Right images are simulated by shift in indexing inside compute_SAD

        for (i = 0; i < WIN; i = i + 1) begin
            input_array[i * DATA_SIZE +: DATA_SIZE] = imgL[i];
        end

        // Wait for propagation
        #10;

        // Manual golden model
        min_sad = 9999;
        best_disp = 0;
        for (d = 0; d < MAX_DISP; d = d + 1) begin
            sad = 0;
            for (i = 0; i < WIN; i = i + 1) begin
                sad = sad + (imgL[i] > imgR[i][d] ? (imgL[i] - imgR[i][d]) : (imgR[i][d] - imgL[i]));
            end

            if (sad < min_sad) begin
                min_sad = sad;
                best_disp = d;
            end
        end

        // Compare
        $display("Best disparity (DUT) = %0d", output_row);
        $display("Expected disparity    = %0d", best_disp);

        if (output_row !== best_disp)
            $fatal("Disparity mismatch");
        else
            $display("Disparity match");

        $finish;
    end

endmodule
