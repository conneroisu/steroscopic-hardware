// ==============================================================
// Vivado(TM) HLS - High-Level Synthesis from C, C++ and SystemC v2020.1 (64-bit)
// Copyright 1986-2020 Xilinx, Inc. All Rights Reserved.
// ==============================================================
`timescale 1ns/1ps
module bd_d92b_csc_0_v_csc_CTRL_s_axi
#(parameter
    C_S_AXI_ADDR_WIDTH = 8,
    C_S_AXI_DATA_WIDTH = 32
)(
    input  wire                          ACLK,
    input  wire                          ARESET,
    input  wire                          ACLK_EN,
    input  wire [C_S_AXI_ADDR_WIDTH-1:0] AWADDR,
    input  wire                          AWVALID,
    output wire                          AWREADY,
    input  wire [C_S_AXI_DATA_WIDTH-1:0] WDATA,
    input  wire [C_S_AXI_DATA_WIDTH/8-1:0] WSTRB,
    input  wire                          WVALID,
    output wire                          WREADY,
    output wire [1:0]                    BRESP,
    output wire                          BVALID,
    input  wire                          BREADY,
    input  wire [C_S_AXI_ADDR_WIDTH-1:0] ARADDR,
    input  wire                          ARVALID,
    output wire                          ARREADY,
    output wire [C_S_AXI_DATA_WIDTH-1:0] RDATA,
    output wire [1:0]                    RRESP,
    output wire                          RVALID,
    input  wire                          RREADY,
    output wire                          interrupt,
    output wire [7:0]                    HwReg_InVideoFormat,
    output wire [7:0]                    HwReg_OutVideoFormat,
    output wire [15:0]                   HwReg_width,
    output wire [15:0]                   HwReg_height,
    output wire [15:0]                   HwReg_K11,
    output wire [15:0]                   HwReg_K12,
    output wire [15:0]                   HwReg_K13,
    output wire [15:0]                   HwReg_K21,
    output wire [15:0]                   HwReg_K22,
    output wire [15:0]                   HwReg_K23,
    output wire [15:0]                   HwReg_K31,
    output wire [15:0]                   HwReg_K32,
    output wire [15:0]                   HwReg_K33,
    output wire [9:0]                    HwReg_ROffset_V,
    output wire [9:0]                    HwReg_GOffset_V,
    output wire [9:0]                    HwReg_BOffset_V,
    output wire [7:0]                    HwReg_ClampMin_V,
    output wire [7:0]                    HwReg_ClipMax_V,
    output wire                          ap_start,
    input  wire                          ap_done,
    input  wire                          ap_ready,
    input  wire                          ap_idle
);
//------------------------Address Info-------------------
// 0x00 : Control signals
//        bit 0  - ap_start (Read/Write/COH)
//        bit 1  - ap_done (Read/COR)
//        bit 2  - ap_idle (Read)
//        bit 3  - ap_ready (Read)
//        bit 7  - auto_restart (Read/Write)
//        others - reserved
// 0x04 : Global Interrupt Enable Register
//        bit 0  - Global Interrupt Enable (Read/Write)
//        others - reserved
// 0x08 : IP Interrupt Enable Register (Read/Write)
//        bit 0  - Channel 0 (ap_done)
//        bit 1  - Channel 1 (ap_ready)
//        others - reserved
// 0x0c : IP Interrupt Status Register (Read/TOW)
//        bit 0  - Channel 0 (ap_done)
//        bit 1  - Channel 1 (ap_ready)
//        others - reserved
// 0x10 : Data signal of HwReg_InVideoFormat
//        bit 7~0 - HwReg_InVideoFormat[7:0] (Read/Write)
//        others  - reserved
// 0x14 : reserved
// 0x18 : Data signal of HwReg_OutVideoFormat
//        bit 7~0 - HwReg_OutVideoFormat[7:0] (Read/Write)
//        others  - reserved
// 0x1c : reserved
// 0x20 : Data signal of HwReg_width
//        bit 15~0 - HwReg_width[15:0] (Read/Write)
//        others   - reserved
// 0x24 : reserved
// 0x28 : Data signal of HwReg_height
//        bit 15~0 - HwReg_height[15:0] (Read/Write)
//        others   - reserved
// 0x2c : reserved
// 0x50 : Data signal of HwReg_K11
//        bit 15~0 - HwReg_K11[15:0] (Read/Write)
//        others   - reserved
// 0x54 : reserved
// 0x58 : Data signal of HwReg_K12
//        bit 15~0 - HwReg_K12[15:0] (Read/Write)
//        others   - reserved
// 0x5c : reserved
// 0x60 : Data signal of HwReg_K13
//        bit 15~0 - HwReg_K13[15:0] (Read/Write)
//        others   - reserved
// 0x64 : reserved
// 0x68 : Data signal of HwReg_K21
//        bit 15~0 - HwReg_K21[15:0] (Read/Write)
//        others   - reserved
// 0x6c : reserved
// 0x70 : Data signal of HwReg_K22
//        bit 15~0 - HwReg_K22[15:0] (Read/Write)
//        others   - reserved
// 0x74 : reserved
// 0x78 : Data signal of HwReg_K23
//        bit 15~0 - HwReg_K23[15:0] (Read/Write)
//        others   - reserved
// 0x7c : reserved
// 0x80 : Data signal of HwReg_K31
//        bit 15~0 - HwReg_K31[15:0] (Read/Write)
//        others   - reserved
// 0x84 : reserved
// 0x88 : Data signal of HwReg_K32
//        bit 15~0 - HwReg_K32[15:0] (Read/Write)
//        others   - reserved
// 0x8c : reserved
// 0x90 : Data signal of HwReg_K33
//        bit 15~0 - HwReg_K33[15:0] (Read/Write)
//        others   - reserved
// 0x94 : reserved
// 0x98 : Data signal of HwReg_ROffset_V
//        bit 9~0 - HwReg_ROffset_V[9:0] (Read/Write)
//        others  - reserved
// 0x9c : reserved
// 0xa0 : Data signal of HwReg_GOffset_V
//        bit 9~0 - HwReg_GOffset_V[9:0] (Read/Write)
//        others  - reserved
// 0xa4 : reserved
// 0xa8 : Data signal of HwReg_BOffset_V
//        bit 9~0 - HwReg_BOffset_V[9:0] (Read/Write)
//        others  - reserved
// 0xac : reserved
// 0xb0 : Data signal of HwReg_ClampMin_V
//        bit 7~0 - HwReg_ClampMin_V[7:0] (Read/Write)
//        others  - reserved
// 0xb4 : reserved
// 0xb8 : Data signal of HwReg_ClipMax_V
//        bit 7~0 - HwReg_ClipMax_V[7:0] (Read/Write)
//        others  - reserved
// 0xbc : reserved
// (SC = Self Clear, COR = Clear on Read, TOW = Toggle on Write, COH = Clear on Handshake)

//------------------------Parameter----------------------
localparam
    ADDR_AP_CTRL                     = 8'h00,
    ADDR_GIE                         = 8'h04,
    ADDR_IER                         = 8'h08,
    ADDR_ISR                         = 8'h0c,
    ADDR_HWREG_INVIDEOFORMAT_DATA_0  = 8'h10,
    ADDR_HWREG_INVIDEOFORMAT_CTRL    = 8'h14,
    ADDR_HWREG_OUTVIDEOFORMAT_DATA_0 = 8'h18,
    ADDR_HWREG_OUTVIDEOFORMAT_CTRL   = 8'h1c,
    ADDR_HWREG_WIDTH_DATA_0          = 8'h20,
    ADDR_HWREG_WIDTH_CTRL            = 8'h24,
    ADDR_HWREG_HEIGHT_DATA_0         = 8'h28,
    ADDR_HWREG_HEIGHT_CTRL           = 8'h2c,
    ADDR_HWREG_K11_DATA_0            = 8'h50,
    ADDR_HWREG_K11_CTRL              = 8'h54,
    ADDR_HWREG_K12_DATA_0            = 8'h58,
    ADDR_HWREG_K12_CTRL              = 8'h5c,
    ADDR_HWREG_K13_DATA_0            = 8'h60,
    ADDR_HWREG_K13_CTRL              = 8'h64,
    ADDR_HWREG_K21_DATA_0            = 8'h68,
    ADDR_HWREG_K21_CTRL              = 8'h6c,
    ADDR_HWREG_K22_DATA_0            = 8'h70,
    ADDR_HWREG_K22_CTRL              = 8'h74,
    ADDR_HWREG_K23_DATA_0            = 8'h78,
    ADDR_HWREG_K23_CTRL              = 8'h7c,
    ADDR_HWREG_K31_DATA_0            = 8'h80,
    ADDR_HWREG_K31_CTRL              = 8'h84,
    ADDR_HWREG_K32_DATA_0            = 8'h88,
    ADDR_HWREG_K32_CTRL              = 8'h8c,
    ADDR_HWREG_K33_DATA_0            = 8'h90,
    ADDR_HWREG_K33_CTRL              = 8'h94,
    ADDR_HWREG_ROFFSET_V_DATA_0      = 8'h98,
    ADDR_HWREG_ROFFSET_V_CTRL        = 8'h9c,
    ADDR_HWREG_GOFFSET_V_DATA_0      = 8'ha0,
    ADDR_HWREG_GOFFSET_V_CTRL        = 8'ha4,
    ADDR_HWREG_BOFFSET_V_DATA_0      = 8'ha8,
    ADDR_HWREG_BOFFSET_V_CTRL        = 8'hac,
    ADDR_HWREG_CLAMPMIN_V_DATA_0     = 8'hb0,
    ADDR_HWREG_CLAMPMIN_V_CTRL       = 8'hb4,
    ADDR_HWREG_CLIPMAX_V_DATA_0      = 8'hb8,
    ADDR_HWREG_CLIPMAX_V_CTRL        = 8'hbc,
    WRIDLE                           = 2'd0,
    WRDATA                           = 2'd1,
    WRRESP                           = 2'd2,
    WRRESET                          = 2'd3,
    RDIDLE                           = 2'd0,
    RDDATA                           = 2'd1,
    RDRESET                          = 2'd2,
    ADDR_BITS         = 8;

//------------------------Local signal-------------------
    reg  [1:0]                    wstate = WRRESET;
    reg  [1:0]                    wnext;
    reg  [ADDR_BITS-1:0]          waddr;
    wire [31:0]                   wmask;
    wire                          aw_hs;
    wire                          w_hs;
    reg  [1:0]                    rstate = RDRESET;
    reg  [1:0]                    rnext;
    reg  [31:0]                   rdata;
    wire                          ar_hs;
    wire [ADDR_BITS-1:0]          raddr;
    // internal registers
    reg                           int_ap_idle;
    reg                           int_ap_ready;
    reg                           int_ap_done = 1'b0;
    reg                           int_ap_start = 1'b0;
    reg                           int_auto_restart = 1'b0;
    reg                           int_gie = 1'b0;
    reg  [1:0]                    int_ier = 2'b0;
    reg  [1:0]                    int_isr = 2'b0;
    reg  [7:0]                    int_HwReg_InVideoFormat = 'b0;
    reg  [7:0]                    int_HwReg_OutVideoFormat = 'b0;
    reg  [15:0]                   int_HwReg_width = 'b0;
    reg  [15:0]                   int_HwReg_height = 'b0;
    reg  [15:0]                   int_HwReg_K11 = 'b0;
    reg  [15:0]                   int_HwReg_K12 = 'b0;
    reg  [15:0]                   int_HwReg_K13 = 'b0;
    reg  [15:0]                   int_HwReg_K21 = 'b0;
    reg  [15:0]                   int_HwReg_K22 = 'b0;
    reg  [15:0]                   int_HwReg_K23 = 'b0;
    reg  [15:0]                   int_HwReg_K31 = 'b0;
    reg  [15:0]                   int_HwReg_K32 = 'b0;
    reg  [15:0]                   int_HwReg_K33 = 'b0;
    reg  [9:0]                    int_HwReg_ROffset_V = 'b0;
    reg  [9:0]                    int_HwReg_GOffset_V = 'b0;
    reg  [9:0]                    int_HwReg_BOffset_V = 'b0;
    reg  [7:0]                    int_HwReg_ClampMin_V = 'b0;
    reg  [7:0]                    int_HwReg_ClipMax_V = 'b0;

//------------------------Instantiation------------------

//------------------------AXI write fsm------------------
assign AWREADY = (wstate == WRIDLE);
assign WREADY  = (wstate == WRDATA);
assign BRESP   = 2'b00;  // OKAY
assign BVALID  = (wstate == WRRESP);
assign wmask   = { {8{WSTRB[3]}}, {8{WSTRB[2]}}, {8{WSTRB[1]}}, {8{WSTRB[0]}} };
assign aw_hs   = AWVALID & AWREADY;
assign w_hs    = WVALID & WREADY;

// wstate
always @(posedge ACLK) begin
    if (ARESET)
        wstate <= WRRESET;
    else if (ACLK_EN)
        wstate <= wnext;
end

// wnext
always @(*) begin
    case (wstate)
        WRIDLE:
            if (AWVALID)
                wnext = WRDATA;
            else
                wnext = WRIDLE;
        WRDATA:
            if (WVALID)
                wnext = WRRESP;
            else
                wnext = WRDATA;
        WRRESP:
            if (BREADY)
                wnext = WRIDLE;
            else
                wnext = WRRESP;
        default:
            wnext = WRIDLE;
    endcase
end

// waddr
always @(posedge ACLK) begin
    if (ACLK_EN) begin
        if (aw_hs)
            waddr <= AWADDR[ADDR_BITS-1:0];
    end
end

//------------------------AXI read fsm-------------------
assign ARREADY = (rstate == RDIDLE);
assign RDATA   = rdata;
assign RRESP   = 2'b00;  // OKAY
assign RVALID  = (rstate == RDDATA);
assign ar_hs   = ARVALID & ARREADY;
assign raddr   = ARADDR[ADDR_BITS-1:0];

// rstate
always @(posedge ACLK) begin
    if (ARESET)
        rstate <= RDRESET;
    else if (ACLK_EN)
        rstate <= rnext;
end

// rnext
always @(*) begin
    case (rstate)
        RDIDLE:
            if (ARVALID)
                rnext = RDDATA;
            else
                rnext = RDIDLE;
        RDDATA:
            if (RREADY & RVALID)
                rnext = RDIDLE;
            else
                rnext = RDDATA;
        default:
            rnext = RDIDLE;
    endcase
end

// rdata
always @(posedge ACLK) begin
    if (ACLK_EN) begin
        if (ar_hs) begin
            rdata <= 1'b0;
            case (raddr)
                ADDR_AP_CTRL: begin
                    rdata[0] <= int_ap_start;
                    rdata[1] <= int_ap_done;
                    rdata[2] <= int_ap_idle;
                    rdata[3] <= int_ap_ready;
                    rdata[7] <= int_auto_restart;
                end
                ADDR_GIE: begin
                    rdata <= int_gie;
                end
                ADDR_IER: begin
                    rdata <= int_ier;
                end
                ADDR_ISR: begin
                    rdata <= int_isr;
                end
                ADDR_HWREG_INVIDEOFORMAT_DATA_0: begin
                    rdata <= int_HwReg_InVideoFormat[7:0];
                end
                ADDR_HWREG_OUTVIDEOFORMAT_DATA_0: begin
                    rdata <= int_HwReg_OutVideoFormat[7:0];
                end
                ADDR_HWREG_WIDTH_DATA_0: begin
                    rdata <= int_HwReg_width[15:0];
                end
                ADDR_HWREG_HEIGHT_DATA_0: begin
                    rdata <= int_HwReg_height[15:0];
                end
                ADDR_HWREG_K11_DATA_0: begin
                    rdata <= int_HwReg_K11[15:0];
                end
                ADDR_HWREG_K12_DATA_0: begin
                    rdata <= int_HwReg_K12[15:0];
                end
                ADDR_HWREG_K13_DATA_0: begin
                    rdata <= int_HwReg_K13[15:0];
                end
                ADDR_HWREG_K21_DATA_0: begin
                    rdata <= int_HwReg_K21[15:0];
                end
                ADDR_HWREG_K22_DATA_0: begin
                    rdata <= int_HwReg_K22[15:0];
                end
                ADDR_HWREG_K23_DATA_0: begin
                    rdata <= int_HwReg_K23[15:0];
                end
                ADDR_HWREG_K31_DATA_0: begin
                    rdata <= int_HwReg_K31[15:0];
                end
                ADDR_HWREG_K32_DATA_0: begin
                    rdata <= int_HwReg_K32[15:0];
                end
                ADDR_HWREG_K33_DATA_0: begin
                    rdata <= int_HwReg_K33[15:0];
                end
                ADDR_HWREG_ROFFSET_V_DATA_0: begin
                    rdata <= int_HwReg_ROffset_V[9:0];
                end
                ADDR_HWREG_GOFFSET_V_DATA_0: begin
                    rdata <= int_HwReg_GOffset_V[9:0];
                end
                ADDR_HWREG_BOFFSET_V_DATA_0: begin
                    rdata <= int_HwReg_BOffset_V[9:0];
                end
                ADDR_HWREG_CLAMPMIN_V_DATA_0: begin
                    rdata <= int_HwReg_ClampMin_V[7:0];
                end
                ADDR_HWREG_CLIPMAX_V_DATA_0: begin
                    rdata <= int_HwReg_ClipMax_V[7:0];
                end
            endcase
        end
    end
end


//------------------------Register logic-----------------
assign interrupt            = int_gie & (|int_isr);
assign ap_start             = int_ap_start;
assign HwReg_InVideoFormat  = int_HwReg_InVideoFormat;
assign HwReg_OutVideoFormat = int_HwReg_OutVideoFormat;
assign HwReg_width          = int_HwReg_width;
assign HwReg_height         = int_HwReg_height;
assign HwReg_K11            = int_HwReg_K11;
assign HwReg_K12            = int_HwReg_K12;
assign HwReg_K13            = int_HwReg_K13;
assign HwReg_K21            = int_HwReg_K21;
assign HwReg_K22            = int_HwReg_K22;
assign HwReg_K23            = int_HwReg_K23;
assign HwReg_K31            = int_HwReg_K31;
assign HwReg_K32            = int_HwReg_K32;
assign HwReg_K33            = int_HwReg_K33;
assign HwReg_ROffset_V      = int_HwReg_ROffset_V;
assign HwReg_GOffset_V      = int_HwReg_GOffset_V;
assign HwReg_BOffset_V      = int_HwReg_BOffset_V;
assign HwReg_ClampMin_V     = int_HwReg_ClampMin_V;
assign HwReg_ClipMax_V      = int_HwReg_ClipMax_V;
// int_ap_start
always @(posedge ACLK) begin
    if (ARESET)
        int_ap_start <= 1'b0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_AP_CTRL && WSTRB[0] && WDATA[0])
            int_ap_start <= 1'b1;
        else if (ap_ready)
            int_ap_start <= int_auto_restart; // clear on handshake/auto restart
    end
end

// int_ap_done
always @(posedge ACLK) begin
    if (ARESET)
        int_ap_done <= 1'b0;
    else if (ACLK_EN) begin
        if (ap_done)
            int_ap_done <= 1'b1;
        else if (ar_hs && raddr == ADDR_AP_CTRL)
            int_ap_done <= 1'b0; // clear on read
    end
end

// int_ap_idle
always @(posedge ACLK) begin
    if (ARESET)
        int_ap_idle <= 1'b0;
    else if (ACLK_EN) begin
            int_ap_idle <= ap_idle;
    end
end

// int_ap_ready
always @(posedge ACLK) begin
    if (ARESET)
        int_ap_ready <= 1'b0;
    else if (ACLK_EN) begin
            int_ap_ready <= ap_ready;
    end
end

// int_auto_restart
always @(posedge ACLK) begin
    if (ARESET)
        int_auto_restart <= 1'b0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_AP_CTRL && WSTRB[0])
            int_auto_restart <=  WDATA[7];
    end
end

// int_gie
always @(posedge ACLK) begin
    if (ARESET)
        int_gie <= 1'b0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_GIE && WSTRB[0])
            int_gie <= WDATA[0];
    end
end

// int_ier
always @(posedge ACLK) begin
    if (ARESET)
        int_ier <= 1'b0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_IER && WSTRB[0])
            int_ier <= WDATA[1:0];
    end
end

// int_isr[0]
always @(posedge ACLK) begin
    if (ARESET)
        int_isr[0] <= 1'b0;
    else if (ACLK_EN) begin
        if (int_ier[0] & ap_done)
            int_isr[0] <= 1'b1;
        else if (w_hs && waddr == ADDR_ISR && WSTRB[0])
            int_isr[0] <= int_isr[0] ^ WDATA[0]; // toggle on write
    end
end

// int_isr[1]
always @(posedge ACLK) begin
    if (ARESET)
        int_isr[1] <= 1'b0;
    else if (ACLK_EN) begin
        if (int_ier[1] & ap_ready)
            int_isr[1] <= 1'b1;
        else if (w_hs && waddr == ADDR_ISR && WSTRB[0])
            int_isr[1] <= int_isr[1] ^ WDATA[1]; // toggle on write
    end
end

// int_HwReg_InVideoFormat[7:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_InVideoFormat[7:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_INVIDEOFORMAT_DATA_0)
            int_HwReg_InVideoFormat[7:0] <= (WDATA[31:0] & wmask) | (int_HwReg_InVideoFormat[7:0] & ~wmask);
    end
end

// int_HwReg_OutVideoFormat[7:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_OutVideoFormat[7:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_OUTVIDEOFORMAT_DATA_0)
            int_HwReg_OutVideoFormat[7:0] <= (WDATA[31:0] & wmask) | (int_HwReg_OutVideoFormat[7:0] & ~wmask);
    end
end

// int_HwReg_width[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_width[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_WIDTH_DATA_0)
            int_HwReg_width[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_width[15:0] & ~wmask);
    end
end

// int_HwReg_height[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_height[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_HEIGHT_DATA_0)
            int_HwReg_height[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_height[15:0] & ~wmask);
    end
end

// int_HwReg_K11[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K11[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K11_DATA_0)
            int_HwReg_K11[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K11[15:0] & ~wmask);
    end
end

// int_HwReg_K12[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K12[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K12_DATA_0)
            int_HwReg_K12[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K12[15:0] & ~wmask);
    end
end

// int_HwReg_K13[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K13[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K13_DATA_0)
            int_HwReg_K13[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K13[15:0] & ~wmask);
    end
end

// int_HwReg_K21[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K21[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K21_DATA_0)
            int_HwReg_K21[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K21[15:0] & ~wmask);
    end
end

// int_HwReg_K22[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K22[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K22_DATA_0)
            int_HwReg_K22[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K22[15:0] & ~wmask);
    end
end

// int_HwReg_K23[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K23[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K23_DATA_0)
            int_HwReg_K23[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K23[15:0] & ~wmask);
    end
end

// int_HwReg_K31[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K31[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K31_DATA_0)
            int_HwReg_K31[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K31[15:0] & ~wmask);
    end
end

// int_HwReg_K32[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K32[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K32_DATA_0)
            int_HwReg_K32[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K32[15:0] & ~wmask);
    end
end

// int_HwReg_K33[15:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_K33[15:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_K33_DATA_0)
            int_HwReg_K33[15:0] <= (WDATA[31:0] & wmask) | (int_HwReg_K33[15:0] & ~wmask);
    end
end

// int_HwReg_ROffset_V[9:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_ROffset_V[9:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_ROFFSET_V_DATA_0)
            int_HwReg_ROffset_V[9:0] <= (WDATA[31:0] & wmask) | (int_HwReg_ROffset_V[9:0] & ~wmask);
    end
end

// int_HwReg_GOffset_V[9:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_GOffset_V[9:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_GOFFSET_V_DATA_0)
            int_HwReg_GOffset_V[9:0] <= (WDATA[31:0] & wmask) | (int_HwReg_GOffset_V[9:0] & ~wmask);
    end
end

// int_HwReg_BOffset_V[9:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_BOffset_V[9:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_BOFFSET_V_DATA_0)
            int_HwReg_BOffset_V[9:0] <= (WDATA[31:0] & wmask) | (int_HwReg_BOffset_V[9:0] & ~wmask);
    end
end

// int_HwReg_ClampMin_V[7:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_ClampMin_V[7:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_CLAMPMIN_V_DATA_0)
            int_HwReg_ClampMin_V[7:0] <= (WDATA[31:0] & wmask) | (int_HwReg_ClampMin_V[7:0] & ~wmask);
    end
end

// int_HwReg_ClipMax_V[7:0]
always @(posedge ACLK) begin
    if (ARESET)
        int_HwReg_ClipMax_V[7:0] <= 0;
    else if (ACLK_EN) begin
        if (w_hs && waddr == ADDR_HWREG_CLIPMAX_V_DATA_0)
            int_HwReg_ClipMax_V[7:0] <= (WDATA[31:0] & wmask) | (int_HwReg_ClipMax_V[7:0] & ~wmask);
    end
end


//------------------------Memory logic-------------------

endmodule
