module SAD_test;
    reg [8:0] winL;
    reg [8:0] winR;
    wire [11:0] sad;

    SAD dut (
        .input_array_a(winL),
        .input_array_b(winR),
        .sad(sad)
    );

    initial begin
        winL[0] = 0;
        winL[1] = 10;
        winL[2] = 20;
        winL[3] = 30;
        winL[4] = 0;
        winL[5] = 10;
        winL[6] = 20;
        winL[7] = 30;
        winL[8] = 0;

        winR[0] = 110;
        winR[1] = 100;
        winR[2] = 90;
        winR[3] = 80;
        winR[4] = 70;
        winR[5] = 60;
        winR[6] = 50;
        winR[7] = 40;
        winR[8] = 30;

        #10
        $display("SAD output = %d", sad);
        $display("Expected 510");
        $finish;
    end 
endmodule