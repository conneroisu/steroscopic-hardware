-- Copyright 1986-2020 Xilinx, Inc. All Rights Reserved.
-- --------------------------------------------------------------------------------
-- Tool Version: Vivado v.2020.1 (win64) Build 2902540 Wed May 27 19:54:49 MDT 2020
-- Date        : Fri Feb 28 01:55:54 2025
-- Host        : DESKTOP-V4CC4PR running 64-bit major release  (build 9200)
-- Command     : write_vhdl -force -mode synth_stub -rename_top decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix -prefix
--               decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix_ design_1_onsemi_vita_cam_0_0_stub.vhdl
-- Design      : design_1_onsemi_vita_cam_0_0
-- Purpose     : Stub declaration of top-level module interface
-- Device      : xc7z020clg484-1
-- --------------------------------------------------------------------------------
library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

entity decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix is
  Port ( 
    clk200 : in STD_LOGIC;
    clk : in STD_LOGIC;
    reset : in STD_LOGIC;
    oe : in STD_LOGIC;
    io_vita_clk_pll : out STD_LOGIC;
    io_vita_reset_n : out STD_LOGIC;
    io_vita_trigger : out STD_LOGIC_VECTOR ( 2 downto 0 );
    io_vita_monitor : in STD_LOGIC_VECTOR ( 1 downto 0 );
    io_vita_clk_out_p : in STD_LOGIC;
    io_vita_clk_out_n : in STD_LOGIC;
    io_vita_sync_p : in STD_LOGIC;
    io_vita_sync_n : in STD_LOGIC;
    io_vita_data_p : in STD_LOGIC_VECTOR ( 3 downto 0 );
    io_vita_data_n : in STD_LOGIC_VECTOR ( 3 downto 0 );
    trigger1 : in STD_LOGIC;
    fsync : out STD_LOGIC;
    video_vsync : out STD_LOGIC;
    video_hsync : out STD_LOGIC;
    video_vblank : out STD_LOGIC;
    video_hblank : out STD_LOGIC;
    video_active_video : out STD_LOGIC;
    video_data : out STD_LOGIC_VECTOR ( 7 downto 0 );
    s00_axi_aclk : in STD_LOGIC;
    s00_axi_aresetn : in STD_LOGIC;
    s00_axi_awaddr : in STD_LOGIC_VECTOR ( 7 downto 0 );
    s00_axi_awprot : in STD_LOGIC_VECTOR ( 2 downto 0 );
    s00_axi_awvalid : in STD_LOGIC;
    s00_axi_awready : out STD_LOGIC;
    s00_axi_wdata : in STD_LOGIC_VECTOR ( 31 downto 0 );
    s00_axi_wstrb : in STD_LOGIC_VECTOR ( 3 downto 0 );
    s00_axi_wvalid : in STD_LOGIC;
    s00_axi_wready : out STD_LOGIC;
    s00_axi_bresp : out STD_LOGIC_VECTOR ( 1 downto 0 );
    s00_axi_bvalid : out STD_LOGIC;
    s00_axi_bready : in STD_LOGIC;
    s00_axi_araddr : in STD_LOGIC_VECTOR ( 7 downto 0 );
    s00_axi_arprot : in STD_LOGIC_VECTOR ( 2 downto 0 );
    s00_axi_arvalid : in STD_LOGIC;
    s00_axi_arready : out STD_LOGIC;
    s00_axi_rdata : out STD_LOGIC_VECTOR ( 31 downto 0 );
    s00_axi_rresp : out STD_LOGIC_VECTOR ( 1 downto 0 );
    s00_axi_rvalid : out STD_LOGIC;
    s00_axi_rready : in STD_LOGIC
  );

end decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix;

architecture stub of decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix is
attribute syn_black_box : boolean;
attribute black_box_pad_pin : string;
attribute syn_black_box of stub : architecture is true;
attribute black_box_pad_pin of stub : architecture is "clk200,clk,reset,oe,io_vita_clk_pll,io_vita_reset_n,io_vita_trigger[2:0],io_vita_monitor[1:0],io_vita_clk_out_p,io_vita_clk_out_n,io_vita_sync_p,io_vita_sync_n,io_vita_data_p[3:0],io_vita_data_n[3:0],trigger1,fsync,video_vsync,video_hsync,video_vblank,video_hblank,video_active_video,video_data[7:0],s00_axi_aclk,s00_axi_aresetn,s00_axi_awaddr[7:0],s00_axi_awprot[2:0],s00_axi_awvalid,s00_axi_awready,s00_axi_wdata[31:0],s00_axi_wstrb[3:0],s00_axi_wvalid,s00_axi_wready,s00_axi_bresp[1:0],s00_axi_bvalid,s00_axi_bready,s00_axi_araddr[7:0],s00_axi_arprot[2:0],s00_axi_arvalid,s00_axi_arready,s00_axi_rdata[31:0],s00_axi_rresp[1:0],s00_axi_rvalid,s00_axi_rready";
attribute x_core_info : string;
attribute x_core_info of stub : architecture is "onsemi_vita_cam_v3_1,Vivado 2020.1";
begin
end;
