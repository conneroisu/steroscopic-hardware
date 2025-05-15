// Copyright 1986-2020 Xilinx, Inc. All Rights Reserved.
// --------------------------------------------------------------------------------
// Tool Version: Vivado v.2020.1 (win64) Build 2902540 Wed May 27 19:54:49 MDT 2020
// Date        : Fri Feb 28 01:55:54 2025
// Host        : DESKTOP-V4CC4PR running 64-bit major release  (build 9200)
// Command     : write_verilog -force -mode synth_stub -rename_top design_1_onsemi_vita_cam_0_0 -prefix
//               design_1_onsemi_vita_cam_0_0_ design_1_onsemi_vita_cam_0_0_stub.v
// Design      : design_1_onsemi_vita_cam_0_0
// Purpose     : Stub declaration of top-level module interface
// Device      : xc7z020clg484-1
// --------------------------------------------------------------------------------

// This empty module with port declaration file causes synthesis tools to infer a black box for IP.
// The synthesis directives are for Synopsys Synplify support to prevent IO buffer insertion.
// Please paste the declaration into a Verilog source file or add the file as an additional source.
(* x_core_info = "onsemi_vita_cam_v3_1,Vivado 2020.1" *)
module design_1_onsemi_vita_cam_0_0(clk200, clk, reset, oe, io_vita_clk_pll, 
  io_vita_reset_n, io_vita_trigger, io_vita_monitor, io_vita_clk_out_p, io_vita_clk_out_n, 
  io_vita_sync_p, io_vita_sync_n, io_vita_data_p, io_vita_data_n, trigger1, fsync, video_vsync, 
  video_hsync, video_vblank, video_hblank, video_active_video, video_data, s00_axi_aclk, 
  s00_axi_aresetn, s00_axi_awaddr, s00_axi_awprot, s00_axi_awvalid, s00_axi_awready, 
  s00_axi_wdata, s00_axi_wstrb, s00_axi_wvalid, s00_axi_wready, s00_axi_bresp, 
  s00_axi_bvalid, s00_axi_bready, s00_axi_araddr, s00_axi_arprot, s00_axi_arvalid, 
  s00_axi_arready, s00_axi_rdata, s00_axi_rresp, s00_axi_rvalid, s00_axi_rready)
/* synthesis syn_black_box black_box_pad_pin="clk200,clk,reset,oe,io_vita_clk_pll,io_vita_reset_n,io_vita_trigger[2:0],io_vita_monitor[1:0],io_vita_clk_out_p,io_vita_clk_out_n,io_vita_sync_p,io_vita_sync_n,io_vita_data_p[3:0],io_vita_data_n[3:0],trigger1,fsync,video_vsync,video_hsync,video_vblank,video_hblank,video_active_video,video_data[7:0],s00_axi_aclk,s00_axi_aresetn,s00_axi_awaddr[7:0],s00_axi_awprot[2:0],s00_axi_awvalid,s00_axi_awready,s00_axi_wdata[31:0],s00_axi_wstrb[3:0],s00_axi_wvalid,s00_axi_wready,s00_axi_bresp[1:0],s00_axi_bvalid,s00_axi_bready,s00_axi_araddr[7:0],s00_axi_arprot[2:0],s00_axi_arvalid,s00_axi_arready,s00_axi_rdata[31:0],s00_axi_rresp[1:0],s00_axi_rvalid,s00_axi_rready" */;
  input clk200;
  input clk;
  input reset;
  input oe;
  output io_vita_clk_pll;
  output io_vita_reset_n;
  output [2:0]io_vita_trigger;
  input [1:0]io_vita_monitor;
  input io_vita_clk_out_p;
  input io_vita_clk_out_n;
  input io_vita_sync_p;
  input io_vita_sync_n;
  input [3:0]io_vita_data_p;
  input [3:0]io_vita_data_n;
  input trigger1;
  output fsync;
  output video_vsync;
  output video_hsync;
  output video_vblank;
  output video_hblank;
  output video_active_video;
  output [7:0]video_data;
  input s00_axi_aclk;
  input s00_axi_aresetn;
  input [7:0]s00_axi_awaddr;
  input [2:0]s00_axi_awprot;
  input s00_axi_awvalid;
  output s00_axi_awready;
  input [31:0]s00_axi_wdata;
  input [3:0]s00_axi_wstrb;
  input s00_axi_wvalid;
  output s00_axi_wready;
  output [1:0]s00_axi_bresp;
  output s00_axi_bvalid;
  input s00_axi_bready;
  input [7:0]s00_axi_araddr;
  input [2:0]s00_axi_arprot;
  input s00_axi_arvalid;
  output s00_axi_arready;
  output [31:0]s00_axi_rdata;
  output [1:0]s00_axi_rresp;
  output s00_axi_rvalid;
  input s00_axi_rready;
endmodule
