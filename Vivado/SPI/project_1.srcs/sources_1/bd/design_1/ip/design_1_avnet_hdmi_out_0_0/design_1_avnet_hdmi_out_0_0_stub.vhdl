-- Copyright 1986-2020 Xilinx, Inc. All Rights Reserved.
-- --------------------------------------------------------------------------------
-- Tool Version: Vivado v.2020.1 (win64) Build 2902540 Wed May 27 19:54:49 MDT 2020
-- Date        : Fri Feb 28 01:51:31 2025
-- Host        : DESKTOP-V4CC4PR running 64-bit major release  (build 9200)
-- Command     : write_vhdl -force -mode synth_stub -rename_top design_1_avnet_hdmi_out_0_0 -prefix
--               design_1_avnet_hdmi_out_0_0_ design_1_avnet_hdmi_out_0_0_stub.vhdl
-- Design      : design_1_avnet_hdmi_out_0_0
-- Purpose     : Stub declaration of top-level module interface
-- Device      : xc7z020clg484-1
-- --------------------------------------------------------------------------------
library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

entity design_1_avnet_hdmi_out_0_0 is
  Port ( 
    clk : in STD_LOGIC;
    reset : in STD_LOGIC;
    oe : in STD_LOGIC;
    embed_syncs : in STD_LOGIC;
    audio_spdif : in STD_LOGIC;
    video_vblank : in STD_LOGIC;
    video_hblank : in STD_LOGIC;
    video_de : in STD_LOGIC;
    video_data : in STD_LOGIC_VECTOR ( 15 downto 0 );
    io_hdmio_spdif : out STD_LOGIC;
    io_hdmio_video : out STD_LOGIC_VECTOR ( 15 downto 0 );
    io_hdmio_clk : out STD_LOGIC
  );

end design_1_avnet_hdmi_out_0_0;

architecture stub of design_1_avnet_hdmi_out_0_0 is
attribute syn_black_box : boolean;
attribute black_box_pad_pin : string;
attribute syn_black_box of stub : architecture is true;
attribute black_box_pad_pin of stub : architecture is "clk,reset,oe,embed_syncs,audio_spdif,video_vblank,video_hblank,video_de,video_data[15:0],io_hdmio_spdif,io_hdmio_video[15:0],io_hdmio_clk";
attribute x_core_info : string;
attribute x_core_info of stub : architecture is "avnet_hdmi_out,Vivado 2020.1";
begin
end;
