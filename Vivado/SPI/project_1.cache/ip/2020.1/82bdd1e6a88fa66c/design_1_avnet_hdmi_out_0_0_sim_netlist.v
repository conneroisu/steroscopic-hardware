// Copyright 1986-2020 Xilinx, Inc. All Rights Reserved.
// --------------------------------------------------------------------------------
// Tool Version: Vivado v.2020.1 (win64) Build 2902540 Wed May 27 19:54:49 MDT 2020
// Date        : Fri Feb 28 01:51:31 2025
// Host        : DESKTOP-V4CC4PR running 64-bit major release  (build 9200)
// Command     : write_verilog -force -mode funcsim -rename_top decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix -prefix
//               decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix_ design_1_avnet_hdmi_out_0_0_sim_netlist.v
// Design      : design_1_avnet_hdmi_out_0_0
// Purpose     : This verilog netlist is a functional simulation representation of the design and should not be modified
//               or synthesized. This netlist cannot be used for SDF annotated simulation.
// Device      : xc7z020clg484-1
// --------------------------------------------------------------------------------
`timescale 1 ps / 1 ps

module decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix_adv7511_embed_syncs
   (\video_data_df_reg[14]_0 ,
    \video_data_d6_reg[7]_0 ,
    \video_data_d6_reg[15]_0 ,
    clk,
    D,
    \vblank_d_reg[1]_0 ,
    Q);
  output [13:0]\video_data_df_reg[14]_0 ;
  output \video_data_d6_reg[7]_0 ;
  output \video_data_d6_reg[15]_0 ;
  input clk;
  input [0:0]D;
  input [0:0]\vblank_d_reg[1]_0 ;
  input [15:0]Q;

  wire [0:0]D;
  wire [15:0]Q;
  wire clk;
  wire \hblank_d_reg_n_0_[1] ;
  wire \hblank_d_reg_n_0_[2] ;
  wire \hblank_d_reg_n_0_[3] ;
  wire \hblank_d_reg_n_0_[4] ;
  wire \hblank_d_reg_n_0_[5] ;
  wire \hblank_d_reg_n_0_[6] ;
  wire [1:0]vblank_d;
  wire [0:0]\vblank_d_reg[1]_0 ;
  wire [15:0]video_data_d1;
  wire [15:0]video_data_d2;
  wire \video_data_d2[0]_i_1_n_0 ;
  wire \video_data_d2[10]_i_1_n_0 ;
  wire \video_data_d2[11]_i_1_n_0 ;
  wire \video_data_d2[11]_i_2_n_0 ;
  wire \video_data_d2[12]_i_1_n_0 ;
  wire \video_data_d2[13]_i_1_n_0 ;
  wire \video_data_d2[13]_i_2_n_0 ;
  wire \video_data_d2[15]_i_1_n_0 ;
  wire \video_data_d2[15]_i_2_n_0 ;
  wire \video_data_d2[1]_i_1_n_0 ;
  wire \video_data_d2[2]_i_1_n_0 ;
  wire \video_data_d2[3]_i_1_n_0 ;
  wire \video_data_d2[4]_i_1_n_0 ;
  wire \video_data_d2[5]_i_1_n_0 ;
  wire \video_data_d2[7]_i_1_n_0 ;
  wire \video_data_d2[8]_i_1_n_0 ;
  wire \video_data_d2[9]_i_1_n_0 ;
  wire \video_data_d5_reg[0]_srl3_n_0 ;
  wire \video_data_d5_reg[10]_srl3_n_0 ;
  wire \video_data_d5_reg[11]_srl3_n_0 ;
  wire \video_data_d5_reg[12]_srl3_n_0 ;
  wire \video_data_d5_reg[13]_srl3_n_0 ;
  wire \video_data_d5_reg[14]_srl4_i_1_n_0 ;
  wire \video_data_d5_reg[14]_srl4_n_0 ;
  wire \video_data_d5_reg[15]_srl3_n_0 ;
  wire \video_data_d5_reg[1]_srl3_n_0 ;
  wire \video_data_d5_reg[2]_srl3_n_0 ;
  wire \video_data_d5_reg[3]_srl3_n_0 ;
  wire \video_data_d5_reg[4]_srl3_n_0 ;
  wire \video_data_d5_reg[5]_srl3_n_0 ;
  wire \video_data_d5_reg[6]_srl4_i_1_n_0 ;
  wire \video_data_d5_reg[6]_srl4_n_0 ;
  wire \video_data_d5_reg[7]_srl3_n_0 ;
  wire \video_data_d5_reg[8]_srl3_n_0 ;
  wire \video_data_d5_reg[9]_srl3_n_0 ;
  wire [15:0]video_data_d6;
  wire \video_data_d6_reg[15]_0 ;
  wire \video_data_d6_reg[7]_0 ;
  wire \video_data_df[0]_i_1_n_0 ;
  wire \video_data_df[10]_i_1_n_0 ;
  wire \video_data_df[11]_i_1_n_0 ;
  wire \video_data_df[12]_i_1_n_0 ;
  wire \video_data_df[13]_i_1_n_0 ;
  wire \video_data_df[13]_i_2_n_0 ;
  wire \video_data_df[14]_i_1_n_0 ;
  wire \video_data_df[14]_i_2_n_0 ;
  wire \video_data_df[1]_i_1_n_0 ;
  wire \video_data_df[2]_i_1_n_0 ;
  wire \video_data_df[3]_i_1_n_0 ;
  wire \video_data_df[4]_i_1_n_0 ;
  wire \video_data_df[5]_i_1_n_0 ;
  wire \video_data_df[6]_i_1_n_0 ;
  wire \video_data_df[8]_i_1_n_0 ;
  wire \video_data_df[9]_i_1_n_0 ;
  wire [13:0]\video_data_df_reg[14]_0 ;

  LUT6 #(
    .INIT(64'hAAAAE0ABAAAAAAAA)) 
    \VIDEO_PORTS_16BIT_GEN.video_r_reg[15]_srl2_i_1 
       (.I0(video_data_d6[15]),
        .I1(\hblank_d_reg_n_0_[4] ),
        .I2(\hblank_d_reg_n_0_[3] ),
        .I3(\hblank_d_reg_n_0_[5] ),
        .I4(\hblank_d_reg_n_0_[2] ),
        .I5(\hblank_d_reg_n_0_[6] ),
        .O(\video_data_d6_reg[15]_0 ));
  LUT6 #(
    .INIT(64'hAAAAE0ABAAAAAAAA)) 
    \VIDEO_PORTS_16BIT_GEN.video_r_reg[7]_srl2_i_1 
       (.I0(video_data_d6[7]),
        .I1(\hblank_d_reg_n_0_[4] ),
        .I2(\hblank_d_reg_n_0_[3] ),
        .I3(\hblank_d_reg_n_0_[5] ),
        .I4(\hblank_d_reg_n_0_[2] ),
        .I5(\hblank_d_reg_n_0_[6] ),
        .O(\video_data_d6_reg[7]_0 ));
  FDRE \hblank_d_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(D),
        .Q(\hblank_d_reg_n_0_[1] ),
        .R(1'b0));
  FDRE \hblank_d_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(\hblank_d_reg_n_0_[1] ),
        .Q(\hblank_d_reg_n_0_[2] ),
        .R(1'b0));
  FDRE \hblank_d_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(\hblank_d_reg_n_0_[2] ),
        .Q(\hblank_d_reg_n_0_[3] ),
        .R(1'b0));
  FDRE \hblank_d_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(\hblank_d_reg_n_0_[3] ),
        .Q(\hblank_d_reg_n_0_[4] ),
        .R(1'b0));
  FDRE \hblank_d_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(\hblank_d_reg_n_0_[4] ),
        .Q(\hblank_d_reg_n_0_[5] ),
        .R(1'b0));
  FDRE \hblank_d_reg[6] 
       (.C(clk),
        .CE(1'b1),
        .D(\hblank_d_reg_n_0_[5] ),
        .Q(\hblank_d_reg_n_0_[6] ),
        .R(1'b0));
  FDRE \vblank_d_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(\vblank_d_reg[1]_0 ),
        .Q(vblank_d[0]),
        .R(1'b0));
  FDRE \vblank_d_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(vblank_d[0]),
        .Q(vblank_d[1]),
        .R(1'b0));
  FDRE \video_data_d1_reg[0] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[0]),
        .Q(video_data_d1[0]),
        .R(1'b0));
  FDRE \video_data_d1_reg[10] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[10]),
        .Q(video_data_d1[10]),
        .R(1'b0));
  FDRE \video_data_d1_reg[11] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[11]),
        .Q(video_data_d1[11]),
        .R(1'b0));
  FDRE \video_data_d1_reg[12] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[12]),
        .Q(video_data_d1[12]),
        .R(1'b0));
  FDRE \video_data_d1_reg[13] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[13]),
        .Q(video_data_d1[13]),
        .R(1'b0));
  FDRE \video_data_d1_reg[14] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[14]),
        .Q(video_data_d1[14]),
        .R(1'b0));
  FDRE \video_data_d1_reg[15] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[15]),
        .Q(video_data_d1[15]),
        .R(1'b0));
  FDRE \video_data_d1_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[1]),
        .Q(video_data_d1[1]),
        .R(1'b0));
  FDRE \video_data_d1_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[2]),
        .Q(video_data_d1[2]),
        .R(1'b0));
  FDRE \video_data_d1_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[3]),
        .Q(video_data_d1[3]),
        .R(1'b0));
  FDRE \video_data_d1_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[4]),
        .Q(video_data_d1[4]),
        .R(1'b0));
  FDRE \video_data_d1_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[5]),
        .Q(video_data_d1[5]),
        .R(1'b0));
  FDRE \video_data_d1_reg[6] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[6]),
        .Q(video_data_d1[6]),
        .R(1'b0));
  FDRE \video_data_d1_reg[7] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[7]),
        .Q(video_data_d1[7]),
        .R(1'b0));
  FDRE \video_data_d1_reg[8] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[8]),
        .Q(video_data_d1[8]),
        .R(1'b0));
  FDRE \video_data_d1_reg[9] 
       (.C(clk),
        .CE(1'b1),
        .D(Q[9]),
        .Q(video_data_d1[9]),
        .R(1'b0));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[0]_i_1 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[0]),
        .O(\video_data_d2[0]_i_1_n_0 ));
  LUT5 #(
    .INIT(32'hFBFF0000)) 
    \video_data_d2[10]_i_1 
       (.I0(\hblank_d_reg_n_0_[4] ),
        .I1(\hblank_d_reg_n_0_[1] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[2] ),
        .I4(video_data_d1[10]),
        .O(\video_data_d2[10]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'h1000000000000030)) 
    \video_data_d2[11]_i_1 
       (.I0(vblank_d[1]),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[1] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(\hblank_d_reg_n_0_[2] ),
        .O(\video_data_d2[11]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[11]_i_2 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[11]),
        .O(\video_data_d2[11]_i_2_n_0 ));
  LUT5 #(
    .INIT(32'hFBFF0000)) 
    \video_data_d2[12]_i_1 
       (.I0(\hblank_d_reg_n_0_[4] ),
        .I1(\hblank_d_reg_n_0_[1] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[2] ),
        .I4(video_data_d1[12]),
        .O(\video_data_d2[12]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'h2000000000000030)) 
    \video_data_d2[13]_i_1 
       (.I0(vblank_d[1]),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[1] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(\hblank_d_reg_n_0_[2] ),
        .O(\video_data_d2[13]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[13]_i_2 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[13]),
        .O(\video_data_d2[13]_i_2_n_0 ));
  LUT5 #(
    .INIT(32'h20000002)) 
    \video_data_d2[15]_i_1 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[3] ),
        .I4(\hblank_d_reg_n_0_[4] ),
        .O(\video_data_d2[15]_i_1_n_0 ));
  LUT5 #(
    .INIT(32'hFBFF0000)) 
    \video_data_d2[15]_i_2 
       (.I0(\hblank_d_reg_n_0_[4] ),
        .I1(\hblank_d_reg_n_0_[1] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[2] ),
        .I4(video_data_d1[15]),
        .O(\video_data_d2[15]_i_2_n_0 ));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[1]_i_1 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[1]),
        .O(\video_data_d2[1]_i_1_n_0 ));
  LUT5 #(
    .INIT(32'hFBFF0000)) 
    \video_data_d2[2]_i_1 
       (.I0(\hblank_d_reg_n_0_[4] ),
        .I1(\hblank_d_reg_n_0_[1] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[2] ),
        .I4(video_data_d1[2]),
        .O(\video_data_d2[2]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[3]_i_1 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[3]),
        .O(\video_data_d2[3]_i_1_n_0 ));
  LUT5 #(
    .INIT(32'hFBFF0000)) 
    \video_data_d2[4]_i_1 
       (.I0(\hblank_d_reg_n_0_[4] ),
        .I1(\hblank_d_reg_n_0_[1] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[2] ),
        .I4(video_data_d1[4]),
        .O(\video_data_d2[4]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[5]_i_1 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[5]),
        .O(\video_data_d2[5]_i_1_n_0 ));
  LUT5 #(
    .INIT(32'hFBFF0000)) 
    \video_data_d2[7]_i_1 
       (.I0(\hblank_d_reg_n_0_[4] ),
        .I1(\hblank_d_reg_n_0_[1] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[2] ),
        .I4(video_data_d1[7]),
        .O(\video_data_d2[7]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[8]_i_1 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[8]),
        .O(\video_data_d2[8]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hDFDFFFDF00000000)) 
    \video_data_d2[9]_i_1 
       (.I0(\hblank_d_reg_n_0_[1] ),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[2] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d1[9]),
        .O(\video_data_d2[9]_i_1_n_0 ));
  FDSE \video_data_d2_reg[0] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[0]_i_1_n_0 ),
        .Q(video_data_d2[0]),
        .S(\video_data_d2[11]_i_1_n_0 ));
  FDSE \video_data_d2_reg[10] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[10]_i_1_n_0 ),
        .Q(video_data_d2[10]),
        .S(\video_data_d2[15]_i_1_n_0 ));
  FDSE \video_data_d2_reg[11] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[11]_i_2_n_0 ),
        .Q(video_data_d2[11]),
        .S(\video_data_d2[11]_i_1_n_0 ));
  FDSE \video_data_d2_reg[12] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[12]_i_1_n_0 ),
        .Q(video_data_d2[12]),
        .S(\video_data_d2[15]_i_1_n_0 ));
  FDSE \video_data_d2_reg[13] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[13]_i_2_n_0 ),
        .Q(video_data_d2[13]),
        .S(\video_data_d2[13]_i_1_n_0 ));
  FDSE \video_data_d2_reg[15] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[15]_i_2_n_0 ),
        .Q(video_data_d2[15]),
        .S(\video_data_d2[15]_i_1_n_0 ));
  FDSE \video_data_d2_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[1]_i_1_n_0 ),
        .Q(video_data_d2[1]),
        .S(\video_data_d2[13]_i_1_n_0 ));
  FDSE \video_data_d2_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[2]_i_1_n_0 ),
        .Q(video_data_d2[2]),
        .S(\video_data_d2[15]_i_1_n_0 ));
  FDSE \video_data_d2_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[3]_i_1_n_0 ),
        .Q(video_data_d2[3]),
        .S(\video_data_d2[11]_i_1_n_0 ));
  FDSE \video_data_d2_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[4]_i_1_n_0 ),
        .Q(video_data_d2[4]),
        .S(\video_data_d2[15]_i_1_n_0 ));
  FDSE \video_data_d2_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[5]_i_1_n_0 ),
        .Q(video_data_d2[5]),
        .S(\video_data_d2[13]_i_1_n_0 ));
  FDSE \video_data_d2_reg[7] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[7]_i_1_n_0 ),
        .Q(video_data_d2[7]),
        .S(\video_data_d2[15]_i_1_n_0 ));
  FDSE \video_data_d2_reg[8] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[8]_i_1_n_0 ),
        .Q(video_data_d2[8]),
        .S(\video_data_d2[11]_i_1_n_0 ));
  FDSE \video_data_d2_reg[9] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d2[9]_i_1_n_0 ),
        .Q(video_data_d2[9]),
        .S(\video_data_d2[13]_i_1_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[0]_srl3 " *) 
  SRL16E \video_data_d5_reg[0]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[0]),
        .Q(\video_data_d5_reg[0]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[10]_srl3 " *) 
  SRL16E \video_data_d5_reg[10]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[10]),
        .Q(\video_data_d5_reg[10]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[11]_srl3 " *) 
  SRL16E \video_data_d5_reg[11]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[11]),
        .Q(\video_data_d5_reg[11]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[12]_srl3 " *) 
  SRL16E \video_data_d5_reg[12]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[12]),
        .Q(\video_data_d5_reg[12]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[13]_srl3 " *) 
  SRL16E \video_data_d5_reg[13]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[13]),
        .Q(\video_data_d5_reg[13]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[14]_srl4 " *) 
  SRL16E \video_data_d5_reg[14]_srl4 
       (.A0(1'b1),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(\video_data_d5_reg[14]_srl4_i_1_n_0 ),
        .Q(\video_data_d5_reg[14]_srl4_n_0 ));
  LUT6 #(
    .INIT(64'h8A8AAA8AAAAAAABA)) 
    \video_data_d5_reg[14]_srl4_i_1 
       (.I0(video_data_d1[14]),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[1] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(\hblank_d_reg_n_0_[2] ),
        .O(\video_data_d5_reg[14]_srl4_i_1_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[15]_srl3 " *) 
  SRL16E \video_data_d5_reg[15]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[15]),
        .Q(\video_data_d5_reg[15]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[1]_srl3 " *) 
  SRL16E \video_data_d5_reg[1]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[1]),
        .Q(\video_data_d5_reg[1]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[2]_srl3 " *) 
  SRL16E \video_data_d5_reg[2]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[2]),
        .Q(\video_data_d5_reg[2]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[3]_srl3 " *) 
  SRL16E \video_data_d5_reg[3]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[3]),
        .Q(\video_data_d5_reg[3]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[4]_srl3 " *) 
  SRL16E \video_data_d5_reg[4]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[4]),
        .Q(\video_data_d5_reg[4]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[5]_srl3 " *) 
  SRL16E \video_data_d5_reg[5]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[5]),
        .Q(\video_data_d5_reg[5]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[6]_srl4 " *) 
  SRL16E \video_data_d5_reg[6]_srl4 
       (.A0(1'b1),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(\video_data_d5_reg[6]_srl4_i_1_n_0 ),
        .Q(\video_data_d5_reg[6]_srl4_n_0 ));
  LUT6 #(
    .INIT(64'h8A8AAA8AAAAAAABA)) 
    \video_data_d5_reg[6]_srl4_i_1 
       (.I0(video_data_d1[6]),
        .I1(\hblank_d_reg_n_0_[5] ),
        .I2(\hblank_d_reg_n_0_[1] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(\hblank_d_reg_n_0_[2] ),
        .O(\video_data_d5_reg[6]_srl4_i_1_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[7]_srl3 " *) 
  SRL16E \video_data_d5_reg[7]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[7]),
        .Q(\video_data_d5_reg[7]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[8]_srl3 " *) 
  SRL16E \video_data_d5_reg[8]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[8]),
        .Q(\video_data_d5_reg[8]_srl3_n_0 ));
  (* srl_bus_name = "\U0/embed_syncs_l/video_data_d5_reg " *) 
  (* srl_name = "\U0/embed_syncs_l/video_data_d5_reg[9]_srl3 " *) 
  SRL16E \video_data_d5_reg[9]_srl3 
       (.A0(1'b0),
        .A1(1'b1),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(video_data_d2[9]),
        .Q(\video_data_d5_reg[9]_srl3_n_0 ));
  FDRE \video_data_d6_reg[0] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[0]_srl3_n_0 ),
        .Q(video_data_d6[0]),
        .R(1'b0));
  FDRE \video_data_d6_reg[10] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[10]_srl3_n_0 ),
        .Q(video_data_d6[10]),
        .R(1'b0));
  FDRE \video_data_d6_reg[11] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[11]_srl3_n_0 ),
        .Q(video_data_d6[11]),
        .R(1'b0));
  FDRE \video_data_d6_reg[12] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[12]_srl3_n_0 ),
        .Q(video_data_d6[12]),
        .R(1'b0));
  FDRE \video_data_d6_reg[13] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[13]_srl3_n_0 ),
        .Q(video_data_d6[13]),
        .R(1'b0));
  FDRE \video_data_d6_reg[14] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[14]_srl4_n_0 ),
        .Q(video_data_d6[14]),
        .R(1'b0));
  FDRE \video_data_d6_reg[15] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[15]_srl3_n_0 ),
        .Q(video_data_d6[15]),
        .R(1'b0));
  FDRE \video_data_d6_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[1]_srl3_n_0 ),
        .Q(video_data_d6[1]),
        .R(1'b0));
  FDRE \video_data_d6_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[2]_srl3_n_0 ),
        .Q(video_data_d6[2]),
        .R(1'b0));
  FDRE \video_data_d6_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[3]_srl3_n_0 ),
        .Q(video_data_d6[3]),
        .R(1'b0));
  FDRE \video_data_d6_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[4]_srl3_n_0 ),
        .Q(video_data_d6[4]),
        .R(1'b0));
  FDRE \video_data_d6_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[5]_srl3_n_0 ),
        .Q(video_data_d6[5]),
        .R(1'b0));
  FDRE \video_data_d6_reg[6] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[6]_srl4_n_0 ),
        .Q(video_data_d6[6]),
        .R(1'b0));
  FDRE \video_data_d6_reg[7] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[7]_srl3_n_0 ),
        .Q(video_data_d6[7]),
        .R(1'b0));
  FDRE \video_data_d6_reg[8] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[8]_srl3_n_0 ),
        .Q(video_data_d6[8]),
        .R(1'b0));
  FDRE \video_data_d6_reg[9] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d5_reg[9]_srl3_n_0 ),
        .Q(video_data_d6[9]),
        .R(1'b0));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[0]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[0]),
        .O(\video_data_df[0]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[10]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[10]),
        .O(\video_data_df[10]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[11]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[11]),
        .O(\video_data_df[11]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[12]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[12]),
        .O(\video_data_df[12]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'h0000000081800000)) 
    \video_data_df[13]_i_1 
       (.I0(\hblank_d_reg_n_0_[5] ),
        .I1(\hblank_d_reg_n_0_[4] ),
        .I2(\hblank_d_reg_n_0_[3] ),
        .I3(vblank_d[1]),
        .I4(\hblank_d_reg_n_0_[6] ),
        .I5(\hblank_d_reg_n_0_[2] ),
        .O(\video_data_df[13]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[13]_i_2 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[13]),
        .O(\video_data_df[13]_i_2_n_0 ));
  LUT5 #(
    .INIT(32'h00800000)) 
    \video_data_df[14]_i_1 
       (.I0(\hblank_d_reg_n_0_[4] ),
        .I1(\hblank_d_reg_n_0_[3] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[2] ),
        .I4(\hblank_d_reg_n_0_[6] ),
        .O(\video_data_df[14]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[14]_i_2 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[14]),
        .O(\video_data_df[14]_i_2_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[1]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[1]),
        .O(\video_data_df[1]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[2]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[2]),
        .O(\video_data_df[2]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[3]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[3]),
        .O(\video_data_df[3]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[4]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[4]),
        .O(\video_data_df[4]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[5]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[5]),
        .O(\video_data_df[5]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[6]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[6]),
        .O(\video_data_df[6]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[8]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[8]),
        .O(\video_data_df[8]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'hFFFFBFBB00000000)) 
    \video_data_df[9]_i_1 
       (.I0(\hblank_d_reg_n_0_[2] ),
        .I1(\hblank_d_reg_n_0_[6] ),
        .I2(\hblank_d_reg_n_0_[5] ),
        .I3(\hblank_d_reg_n_0_[4] ),
        .I4(\hblank_d_reg_n_0_[3] ),
        .I5(video_data_d6[9]),
        .O(\video_data_df[9]_i_1_n_0 ));
  FDSE \video_data_df_reg[0] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[0]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [0]),
        .S(\video_data_df[13]_i_1_n_0 ));
  FDSE \video_data_df_reg[10] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[10]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [9]),
        .S(\video_data_df[14]_i_1_n_0 ));
  FDSE \video_data_df_reg[11] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[11]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [10]),
        .S(\video_data_df[13]_i_1_n_0 ));
  FDSE \video_data_df_reg[12] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[12]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [11]),
        .S(\video_data_df[14]_i_1_n_0 ));
  FDSE \video_data_df_reg[13] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[13]_i_2_n_0 ),
        .Q(\video_data_df_reg[14]_0 [12]),
        .S(\video_data_df[13]_i_1_n_0 ));
  FDSE \video_data_df_reg[14] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[14]_i_2_n_0 ),
        .Q(\video_data_df_reg[14]_0 [13]),
        .S(\video_data_df[14]_i_1_n_0 ));
  FDSE \video_data_df_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[1]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [1]),
        .S(\video_data_df[13]_i_1_n_0 ));
  FDSE \video_data_df_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[2]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [2]),
        .S(\video_data_df[14]_i_1_n_0 ));
  FDSE \video_data_df_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[3]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [3]),
        .S(\video_data_df[13]_i_1_n_0 ));
  FDSE \video_data_df_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[4]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [4]),
        .S(\video_data_df[14]_i_1_n_0 ));
  FDSE \video_data_df_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[5]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [5]),
        .S(\video_data_df[13]_i_1_n_0 ));
  FDSE \video_data_df_reg[6] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[6]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [6]),
        .S(\video_data_df[14]_i_1_n_0 ));
  FDSE \video_data_df_reg[8] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[8]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [7]),
        .S(\video_data_df[13]_i_1_n_0 ));
  FDSE \video_data_df_reg[9] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_df[9]_i_1_n_0 ),
        .Q(\video_data_df_reg[14]_0 [8]),
        .S(\video_data_df[13]_i_1_n_0 ));
endmodule

module decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix_avnet_hdmi_out
   (io_hdmio_spdif,
    io_hdmio_video,
    io_hdmio_clk,
    clk,
    audio_spdif,
    video_hblank,
    video_vblank,
    video_data,
    oe);
  output io_hdmio_spdif;
  output [15:0]io_hdmio_video;
  output io_hdmio_clk;
  input clk;
  input audio_spdif;
  input video_hblank;
  input video_vblank;
  input [15:0]video_data;
  input oe;

  wire D1;
  wire I;
  wire T;
  wire \VIDEO_PORTS_16BIT_GEN.video_r_reg[15]_srl2_n_0 ;
  wire \VIDEO_PORTS_16BIT_GEN.video_r_reg[7]_srl2_n_0 ;
  wire audio_spdif;
  wire clk;
  wire embed_syncs_l_n_14;
  wire embed_syncs_l_n_15;
  wire hblank_d;
  wire hdmio_clk_o;
  wire hdmio_clk_t;
  wire [15:0]hdmio_video_o;
  wire io_hdmio_clk;
  wire io_hdmio_spdif;
  wire [15:0]io_hdmio_video;
  wire oe;
  wire vblank_d_reg_n_0;
  wire [15:0]video_data;
  wire [15:0]video_data_d;
  wire \video_data_d[8]_i_1_n_0 ;
  wire \video_data_d[8]_i_2_n_0 ;
  wire [14:0]video_data_df;
  wire video_hblank;
  wire [14:0]video_r;
  wire video_vblank;
  wire \NLW_V6_GEN.ODDR_hdmio_clk_o_R_UNCONNECTED ;
  wire \NLW_V6_GEN.ODDR_hdmio_clk_o_S_UNCONNECTED ;
  wire \NLW_V6_GEN.ODDR_hdmio_clk_t_R_UNCONNECTED ;
  wire \NLW_V6_GEN.ODDR_hdmio_clk_t_S_UNCONNECTED ;

  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[0].OBUFT_hdmio_video 
       (.I(hdmio_video_o[0]),
        .O(io_hdmio_video[0]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[10].OBUFT_hdmio_video 
       (.I(hdmio_video_o[10]),
        .O(io_hdmio_video[10]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[11].OBUFT_hdmio_video 
       (.I(hdmio_video_o[11]),
        .O(io_hdmio_video[11]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[12].OBUFT_hdmio_video 
       (.I(hdmio_video_o[12]),
        .O(io_hdmio_video[12]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[13].OBUFT_hdmio_video 
       (.I(hdmio_video_o[13]),
        .O(io_hdmio_video[13]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[14].OBUFT_hdmio_video 
       (.I(hdmio_video_o[14]),
        .O(io_hdmio_video[14]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[15].OBUFT_hdmio_video 
       (.I(hdmio_video_o[15]),
        .O(io_hdmio_video[15]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[1].OBUFT_hdmio_video 
       (.I(hdmio_video_o[1]),
        .O(io_hdmio_video[1]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[2].OBUFT_hdmio_video 
       (.I(hdmio_video_o[2]),
        .O(io_hdmio_video[2]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[3].OBUFT_hdmio_video 
       (.I(hdmio_video_o[3]),
        .O(io_hdmio_video[3]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[4].OBUFT_hdmio_video 
       (.I(hdmio_video_o[4]),
        .O(io_hdmio_video[4]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[5].OBUFT_hdmio_video 
       (.I(hdmio_video_o[5]),
        .O(io_hdmio_video[5]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[6].OBUFT_hdmio_video 
       (.I(hdmio_video_o[6]),
        .O(io_hdmio_video[6]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[7].OBUFT_hdmio_video 
       (.I(hdmio_video_o[7]),
        .O(io_hdmio_video[7]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[8].OBUFT_hdmio_video 
       (.I(hdmio_video_o[8]),
        .O(io_hdmio_video[8]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    \IO1[9].OBUFT_hdmio_video 
       (.I(hdmio_video_o[9]),
        .O(io_hdmio_video[9]),
        .T(T));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    OBUFT_hdmio_clk
       (.I(hdmio_clk_o),
        .O(io_hdmio_clk),
        .T(hdmio_clk_t));
  (* CAPACITANCE = "DONT_CARE" *) 
  (* box_type = "PRIMITIVE" *) 
  OBUFT #(
    .IOSTANDARD("DEFAULT")) 
    OBUFT_hdmio_spdif
       (.I(I),
        .O(io_hdmio_spdif),
        .T(T));
  (* OPT_MODIFIED = "MLO" *) 
  (* __SRVAL = "TRUE" *) 
  (* box_type = "PRIMITIVE" *) 
  ODDR #(
    .DDR_CLK_EDGE("SAME_EDGE"),
    .INIT(1'b1),
    .IS_C_INVERTED(1'b0),
    .IS_D1_INVERTED(1'b0),
    .IS_D2_INVERTED(1'b0),
    .SRTYPE("ASYNC")) 
    \V6_GEN.ODDR_hdmio_clk_o 
       (.C(clk),
        .CE(1'b1),
        .D1(1'b0),
        .D2(1'b1),
        .Q(hdmio_clk_o),
        .R(\NLW_V6_GEN.ODDR_hdmio_clk_o_R_UNCONNECTED ),
        .S(\NLW_V6_GEN.ODDR_hdmio_clk_o_S_UNCONNECTED ));
  (* OPT_MODIFIED = "MLO" *) 
  (* __SRVAL = "TRUE" *) 
  (* box_type = "PRIMITIVE" *) 
  ODDR #(
    .DDR_CLK_EDGE("SAME_EDGE"),
    .INIT(1'b1),
    .IS_C_INVERTED(1'b0),
    .IS_D1_INVERTED(1'b1),
    .IS_D2_INVERTED(1'b1),
    .SRTYPE("ASYNC")) 
    \V6_GEN.ODDR_hdmio_clk_t 
       (.C(clk),
        .CE(1'b1),
        .D1(oe),
        .D2(oe),
        .Q(hdmio_clk_t),
        .R(\NLW_V6_GEN.ODDR_hdmio_clk_t_R_UNCONNECTED ),
        .S(\NLW_V6_GEN.ODDR_hdmio_clk_t_S_UNCONNECTED ));
  LUT1 #(
    .INIT(2'h1)) 
    \V6_GEN.ODDR_hdmio_clk_t_i_1 
       (.I0(oe),
        .O(D1));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[0] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[0]),
        .Q(video_r[0]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[10] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[10]),
        .Q(video_r[10]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[11] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[11]),
        .Q(video_r[11]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[12] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[12]),
        .Q(video_r[12]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[13] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[13]),
        .Q(video_r[13]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[14] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[14]),
        .Q(video_r[14]),
        .R(1'b0));
  (* srl_bus_name = "\U0/VIDEO_PORTS_16BIT_GEN.video_r_reg " *) 
  (* srl_name = "\U0/VIDEO_PORTS_16BIT_GEN.video_r_reg[15]_srl2 " *) 
  SRL16E \VIDEO_PORTS_16BIT_GEN.video_r_reg[15]_srl2 
       (.A0(1'b1),
        .A1(1'b0),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(embed_syncs_l_n_15),
        .Q(\VIDEO_PORTS_16BIT_GEN.video_r_reg[15]_srl2_n_0 ));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[1]),
        .Q(video_r[1]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[2]),
        .Q(video_r[2]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[3]),
        .Q(video_r[3]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[4]),
        .Q(video_r[4]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[5]),
        .Q(video_r[5]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[6] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[6]),
        .Q(video_r[6]),
        .R(1'b0));
  (* srl_bus_name = "\U0/VIDEO_PORTS_16BIT_GEN.video_r_reg " *) 
  (* srl_name = "\U0/VIDEO_PORTS_16BIT_GEN.video_r_reg[7]_srl2 " *) 
  SRL16E \VIDEO_PORTS_16BIT_GEN.video_r_reg[7]_srl2 
       (.A0(1'b1),
        .A1(1'b0),
        .A2(1'b0),
        .A3(1'b0),
        .CE(1'b1),
        .CLK(clk),
        .D(embed_syncs_l_n_14),
        .Q(\VIDEO_PORTS_16BIT_GEN.video_r_reg[7]_srl2_n_0 ));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[8] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[8]),
        .Q(video_r[8]),
        .R(1'b0));
  FDRE \VIDEO_PORTS_16BIT_GEN.video_r_reg[9] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data_df[9]),
        .Q(video_r[9]),
        .R(1'b0));
  decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix_adv7511_embed_syncs embed_syncs_l
       (.D(hblank_d),
        .Q(video_data_d),
        .clk(clk),
        .\vblank_d_reg[1]_0 (vblank_d_reg_n_0),
        .\video_data_d6_reg[15]_0 (embed_syncs_l_n_15),
        .\video_data_d6_reg[7]_0 (embed_syncs_l_n_14),
        .\video_data_df_reg[14]_0 ({video_data_df[14:8],video_data_df[6:0]}));
  FDRE hblank_d_reg
       (.C(clk),
        .CE(1'b1),
        .D(video_hblank),
        .Q(hblank_d),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE hdmio_spdif_o_reg
       (.C(clk),
        .CE(1'b1),
        .D(audio_spdif),
        .Q(I),
        .R(1'b0));
  (* OPT_MODIFIED = "MLO" *) 
  FDRE hdmio_spdif_t_reg
       (.C(clk),
        .CE(1'b1),
        .D(D1),
        .Q(T),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[0] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[0]),
        .Q(hdmio_video_o[0]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[10] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[10]),
        .Q(hdmio_video_o[10]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[11] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[11]),
        .Q(hdmio_video_o[11]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[12] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[12]),
        .Q(hdmio_video_o[12]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[13] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[13]),
        .Q(hdmio_video_o[13]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[14] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[14]),
        .Q(hdmio_video_o[14]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[15] 
       (.C(clk),
        .CE(1'b1),
        .D(\VIDEO_PORTS_16BIT_GEN.video_r_reg[15]_srl2_n_0 ),
        .Q(hdmio_video_o[15]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[1]),
        .Q(hdmio_video_o[1]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[2]),
        .Q(hdmio_video_o[2]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[3]),
        .Q(hdmio_video_o[3]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[4]),
        .Q(hdmio_video_o[4]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[5]),
        .Q(hdmio_video_o[5]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[6] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[6]),
        .Q(hdmio_video_o[6]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[7] 
       (.C(clk),
        .CE(1'b1),
        .D(\VIDEO_PORTS_16BIT_GEN.video_r_reg[7]_srl2_n_0 ),
        .Q(hdmio_video_o[7]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[8] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[8]),
        .Q(hdmio_video_o[8]),
        .R(1'b0));
  (* IOB = "TRUE" *) 
  FDRE \hdmio_video_o_reg[9] 
       (.C(clk),
        .CE(1'b1),
        .D(video_r[9]),
        .Q(hdmio_video_o[9]),
        .R(1'b0));
  FDRE vblank_d_reg
       (.C(clk),
        .CE(1'b1),
        .D(video_vblank),
        .Q(vblank_d_reg_n_0),
        .R(1'b0));
  LUT3 #(
    .INIT(8'hF4)) 
    \video_data_d[8]_i_1 
       (.I0(video_data[15]),
        .I1(\video_data_d[8]_i_2_n_0 ),
        .I2(video_data[8]),
        .O(\video_data_d[8]_i_1_n_0 ));
  LUT6 #(
    .INIT(64'h0000000000000001)) 
    \video_data_d[8]_i_2 
       (.I0(video_data[13]),
        .I1(video_data[10]),
        .I2(video_data[9]),
        .I3(video_data[11]),
        .I4(video_data[12]),
        .I5(video_data[14]),
        .O(\video_data_d[8]_i_2_n_0 ));
  FDRE \video_data_d_reg[0] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[0]),
        .Q(video_data_d[0]),
        .R(1'b0));
  FDRE \video_data_d_reg[10] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[10]),
        .Q(video_data_d[10]),
        .R(1'b0));
  FDRE \video_data_d_reg[11] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[11]),
        .Q(video_data_d[11]),
        .R(1'b0));
  FDRE \video_data_d_reg[12] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[12]),
        .Q(video_data_d[12]),
        .R(1'b0));
  FDRE \video_data_d_reg[13] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[13]),
        .Q(video_data_d[13]),
        .R(1'b0));
  FDRE \video_data_d_reg[14] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[14]),
        .Q(video_data_d[14]),
        .R(1'b0));
  FDRE \video_data_d_reg[15] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[15]),
        .Q(video_data_d[15]),
        .R(1'b0));
  FDRE \video_data_d_reg[1] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[1]),
        .Q(video_data_d[1]),
        .R(1'b0));
  FDRE \video_data_d_reg[2] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[2]),
        .Q(video_data_d[2]),
        .R(1'b0));
  FDRE \video_data_d_reg[3] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[3]),
        .Q(video_data_d[3]),
        .R(1'b0));
  FDRE \video_data_d_reg[4] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[4]),
        .Q(video_data_d[4]),
        .R(1'b0));
  FDRE \video_data_d_reg[5] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[5]),
        .Q(video_data_d[5]),
        .R(1'b0));
  FDRE \video_data_d_reg[6] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[6]),
        .Q(video_data_d[6]),
        .R(1'b0));
  FDRE \video_data_d_reg[7] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[7]),
        .Q(video_data_d[7]),
        .R(1'b0));
  FDRE \video_data_d_reg[8] 
       (.C(clk),
        .CE(1'b1),
        .D(\video_data_d[8]_i_1_n_0 ),
        .Q(video_data_d[8]),
        .R(1'b0));
  FDRE \video_data_d_reg[9] 
       (.C(clk),
        .CE(1'b1),
        .D(video_data[9]),
        .Q(video_data_d[9]),
        .R(1'b0));
endmodule

(* CHECK_LICENSE_TYPE = "design_1_avnet_hdmi_out_0_0,avnet_hdmi_out,{}" *) (* downgradeipidentifiedwarnings = "yes" *) (* x_core_info = "avnet_hdmi_out,Vivado 2020.1" *) 
(* NotValidForBitStream *)
module decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix
   (clk,
    reset,
    oe,
    embed_syncs,
    audio_spdif,
    video_vblank,
    video_hblank,
    video_de,
    video_data,
    io_hdmio_spdif,
    io_hdmio_video,
    io_hdmio_clk);
  input clk;
  input reset;
  input oe;
  input embed_syncs;
  input audio_spdif;
  (* x_interface_info = "xilinx.com:interface:vid_io:1.0 VID_IO_IN VBLANK" *) input video_vblank;
  (* x_interface_info = "xilinx.com:interface:vid_io:1.0 VID_IO_IN HBLANK" *) input video_hblank;
  (* x_interface_info = "xilinx.com:interface:vid_io:1.0 VID_IO_IN ACTIVE_VIDEO" *) input video_de;
  (* x_interface_info = "xilinx.com:interface:vid_io:1.0 VID_IO_IN DATA" *) input [15:0]video_data;
  (* x_interface_info = "avnet.com:interface:avnet_hdmi:2.0 IO_HDMIO SPDIF" *) output io_hdmio_spdif;
  (* x_interface_info = "avnet.com:interface:avnet_hdmi:2.0 IO_HDMIO DATA" *) output [15:0]io_hdmio_video;
  (* x_interface_info = "avnet.com:interface:avnet_hdmi:2.0 IO_HDMIO CLK" *) output io_hdmio_clk;

  wire audio_spdif;
  wire clk;
  (* DRIVE = "12" *) (* SLEW = "SLOW" *) wire io_hdmio_clk;
  (* DRIVE = "12" *) (* SLEW = "SLOW" *) wire io_hdmio_spdif;
  (* DRIVE = "12" *) (* SLEW = "SLOW" *) wire [15:0]io_hdmio_video;
  wire oe;
  wire [15:0]video_data;
  wire video_hblank;
  wire video_vblank;

  decalper_eb_ot_sdeen_pot_pi_dehcac_xnilix_avnet_hdmi_out U0
       (.audio_spdif(audio_spdif),
        .clk(clk),
        .io_hdmio_clk(io_hdmio_clk),
        .io_hdmio_spdif(io_hdmio_spdif),
        .io_hdmio_video(io_hdmio_video),
        .oe(oe),
        .video_data(video_data),
        .video_hblank(video_hblank),
        .video_vblank(video_vblank));
endmodule
`ifndef GLBL
`define GLBL
`timescale  1 ps / 1 ps

module glbl ();

    parameter ROC_WIDTH = 100000;
    parameter TOC_WIDTH = 0;
    parameter GRES_WIDTH = 10000;
    parameter GRES_START = 10000;

//--------   STARTUP Globals --------------
    wire GSR;
    wire GTS;
    wire GWE;
    wire PRLD;
    wire GRESTORE;
    tri1 p_up_tmp;
    tri (weak1, strong0) PLL_LOCKG = p_up_tmp;

    wire PROGB_GLBL;
    wire CCLKO_GLBL;
    wire FCSBO_GLBL;
    wire [3:0] DO_GLBL;
    wire [3:0] DI_GLBL;
   
    reg GSR_int;
    reg GTS_int;
    reg PRLD_int;
    reg GRESTORE_int;

//--------   JTAG Globals --------------
    wire JTAG_TDO_GLBL;
    wire JTAG_TCK_GLBL;
    wire JTAG_TDI_GLBL;
    wire JTAG_TMS_GLBL;
    wire JTAG_TRST_GLBL;

    reg JTAG_CAPTURE_GLBL;
    reg JTAG_RESET_GLBL;
    reg JTAG_SHIFT_GLBL;
    reg JTAG_UPDATE_GLBL;
    reg JTAG_RUNTEST_GLBL;

    reg JTAG_SEL1_GLBL = 0;
    reg JTAG_SEL2_GLBL = 0 ;
    reg JTAG_SEL3_GLBL = 0;
    reg JTAG_SEL4_GLBL = 0;

    reg JTAG_USER_TDO1_GLBL = 1'bz;
    reg JTAG_USER_TDO2_GLBL = 1'bz;
    reg JTAG_USER_TDO3_GLBL = 1'bz;
    reg JTAG_USER_TDO4_GLBL = 1'bz;

    assign (strong1, weak0) GSR = GSR_int;
    assign (strong1, weak0) GTS = GTS_int;
    assign (weak1, weak0) PRLD = PRLD_int;
    assign (strong1, weak0) GRESTORE = GRESTORE_int;

    initial begin
	GSR_int = 1'b1;
	PRLD_int = 1'b1;
	#(ROC_WIDTH)
	GSR_int = 1'b0;
	PRLD_int = 1'b0;
    end

    initial begin
	GTS_int = 1'b1;
	#(TOC_WIDTH)
	GTS_int = 1'b0;
    end

    initial begin 
	GRESTORE_int = 1'b0;
	#(GRES_START);
	GRESTORE_int = 1'b1;
	#(GRES_WIDTH);
	GRESTORE_int = 1'b0;
    end

endmodule
`endif
