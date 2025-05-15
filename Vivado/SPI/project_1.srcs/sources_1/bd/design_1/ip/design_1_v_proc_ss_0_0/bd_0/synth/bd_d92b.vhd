--Copyright 1986-2020 Xilinx, Inc. All Rights Reserved.
----------------------------------------------------------------------------------
--Command: generate_target bd_d92b.bd
--Design : bd_d92b
--Purpose: IP block netlist
----------------------------------------------------------------------------------
library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
library UNISIM;
use UNISIM.VCOMPONENTS.ALL;
entity bd_d92b is
  port (
    aclk : in STD_LOGIC;
    aresetn : in STD_LOGIC;
    m_axis_tdata : out STD_LOGIC_VECTOR ( 23 downto 0 );
    m_axis_tdest : out STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_tid : out STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_tkeep : out STD_LOGIC_VECTOR ( 2 downto 0 );
    m_axis_tlast : out STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_tready : in STD_LOGIC;
    m_axis_tstrb : out STD_LOGIC_VECTOR ( 2 downto 0 );
    m_axis_tuser : out STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_tvalid : out STD_LOGIC;
    s_axi_ctrl_araddr : in STD_LOGIC_VECTOR ( 7 downto 0 );
    s_axi_ctrl_arready : out STD_LOGIC;
    s_axi_ctrl_arvalid : in STD_LOGIC;
    s_axi_ctrl_awaddr : in STD_LOGIC_VECTOR ( 7 downto 0 );
    s_axi_ctrl_awready : out STD_LOGIC;
    s_axi_ctrl_awvalid : in STD_LOGIC;
    s_axi_ctrl_bready : in STD_LOGIC;
    s_axi_ctrl_bresp : out STD_LOGIC_VECTOR ( 1 downto 0 );
    s_axi_ctrl_bvalid : out STD_LOGIC;
    s_axi_ctrl_rdata : out STD_LOGIC_VECTOR ( 31 downto 0 );
    s_axi_ctrl_rready : in STD_LOGIC;
    s_axi_ctrl_rresp : out STD_LOGIC_VECTOR ( 1 downto 0 );
    s_axi_ctrl_rvalid : out STD_LOGIC;
    s_axi_ctrl_wdata : in STD_LOGIC_VECTOR ( 31 downto 0 );
    s_axi_ctrl_wready : out STD_LOGIC;
    s_axi_ctrl_wstrb : in STD_LOGIC_VECTOR ( 3 downto 0 );
    s_axi_ctrl_wvalid : in STD_LOGIC;
    s_axis_tdata : in STD_LOGIC_VECTOR ( 23 downto 0 );
    s_axis_tdest : in STD_LOGIC_VECTOR ( 0 to 0 );
    s_axis_tid : in STD_LOGIC_VECTOR ( 0 to 0 );
    s_axis_tkeep : in STD_LOGIC_VECTOR ( 2 downto 0 );
    s_axis_tlast : in STD_LOGIC_VECTOR ( 0 to 0 );
    s_axis_tready : out STD_LOGIC;
    s_axis_tstrb : in STD_LOGIC_VECTOR ( 2 downto 0 );
    s_axis_tuser : in STD_LOGIC_VECTOR ( 0 to 0 );
    s_axis_tvalid : in STD_LOGIC
  );
  attribute CORE_GENERATION_INFO : string;
  attribute CORE_GENERATION_INFO of bd_d92b : entity is "bd_d92b,IP_Integrator,{x_ipVendor=xilinx.com,x_ipLibrary=BlockDiagram,x_ipName=bd_d92b,x_ipVersion=1.00.a,x_ipLanguage=VHDL,numBlks=1,numReposBlks=1,numNonXlnxBlks=0,numHierBlks=0,maxHierDepth=0,numSysgenBlks=0,numHlsBlks=0,numHdlrefBlks=0,numPkgbdBlks=0,bdsource=SBD,synth_mode=OOC_per_IP}";
  attribute HW_HANDOFF : string;
  attribute HW_HANDOFF of bd_d92b : entity is "design_1_v_proc_ss_0_0.hwdef";
end bd_d92b;

architecture STRUCTURE of bd_d92b is
  component bd_d92b_csc_0 is
  port (
    s_axi_CTRL_AWADDR : in STD_LOGIC_VECTOR ( 7 downto 0 );
    s_axi_CTRL_AWVALID : in STD_LOGIC;
    s_axi_CTRL_AWREADY : out STD_LOGIC;
    s_axi_CTRL_WDATA : in STD_LOGIC_VECTOR ( 31 downto 0 );
    s_axi_CTRL_WSTRB : in STD_LOGIC_VECTOR ( 3 downto 0 );
    s_axi_CTRL_WVALID : in STD_LOGIC;
    s_axi_CTRL_WREADY : out STD_LOGIC;
    s_axi_CTRL_BRESP : out STD_LOGIC_VECTOR ( 1 downto 0 );
    s_axi_CTRL_BVALID : out STD_LOGIC;
    s_axi_CTRL_BREADY : in STD_LOGIC;
    s_axi_CTRL_ARADDR : in STD_LOGIC_VECTOR ( 7 downto 0 );
    s_axi_CTRL_ARVALID : in STD_LOGIC;
    s_axi_CTRL_ARREADY : out STD_LOGIC;
    s_axi_CTRL_RDATA : out STD_LOGIC_VECTOR ( 31 downto 0 );
    s_axi_CTRL_RRESP : out STD_LOGIC_VECTOR ( 1 downto 0 );
    s_axi_CTRL_RVALID : out STD_LOGIC;
    s_axi_CTRL_RREADY : in STD_LOGIC;
    ap_clk : in STD_LOGIC;
    ap_rst_n : in STD_LOGIC;
    interrupt : out STD_LOGIC;
    s_axis_video_TVALID : in STD_LOGIC;
    s_axis_video_TREADY : out STD_LOGIC;
    s_axis_video_TDATA : in STD_LOGIC_VECTOR ( 23 downto 0 );
    s_axis_video_TKEEP : in STD_LOGIC_VECTOR ( 2 downto 0 );
    s_axis_video_TSTRB : in STD_LOGIC_VECTOR ( 2 downto 0 );
    s_axis_video_TUSER : in STD_LOGIC_VECTOR ( 0 to 0 );
    s_axis_video_TLAST : in STD_LOGIC_VECTOR ( 0 to 0 );
    s_axis_video_TID : in STD_LOGIC_VECTOR ( 0 to 0 );
    s_axis_video_TDEST : in STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_video_TVALID : out STD_LOGIC;
    m_axis_video_TREADY : in STD_LOGIC;
    m_axis_video_TDATA : out STD_LOGIC_VECTOR ( 23 downto 0 );
    m_axis_video_TKEEP : out STD_LOGIC_VECTOR ( 2 downto 0 );
    m_axis_video_TSTRB : out STD_LOGIC_VECTOR ( 2 downto 0 );
    m_axis_video_TUSER : out STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_video_TLAST : out STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_video_TID : out STD_LOGIC_VECTOR ( 0 to 0 );
    m_axis_video_TDEST : out STD_LOGIC_VECTOR ( 0 to 0 )
  );
  end component bd_d92b_csc_0;
  signal intf_net_bdry_in_s_axi_ctrl_ARADDR : STD_LOGIC_VECTOR ( 7 downto 0 );
  signal intf_net_bdry_in_s_axi_ctrl_ARREADY : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_ARVALID : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_AWADDR : STD_LOGIC_VECTOR ( 7 downto 0 );
  signal intf_net_bdry_in_s_axi_ctrl_AWREADY : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_AWVALID : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_BREADY : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_BRESP : STD_LOGIC_VECTOR ( 1 downto 0 );
  signal intf_net_bdry_in_s_axi_ctrl_BVALID : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_RDATA : STD_LOGIC_VECTOR ( 31 downto 0 );
  signal intf_net_bdry_in_s_axi_ctrl_RREADY : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_RRESP : STD_LOGIC_VECTOR ( 1 downto 0 );
  signal intf_net_bdry_in_s_axi_ctrl_RVALID : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_WDATA : STD_LOGIC_VECTOR ( 31 downto 0 );
  signal intf_net_bdry_in_s_axi_ctrl_WREADY : STD_LOGIC;
  signal intf_net_bdry_in_s_axi_ctrl_WSTRB : STD_LOGIC_VECTOR ( 3 downto 0 );
  signal intf_net_bdry_in_s_axi_ctrl_WVALID : STD_LOGIC;
  signal intf_net_bdry_in_s_axis_TDATA : STD_LOGIC_VECTOR ( 23 downto 0 );
  signal intf_net_bdry_in_s_axis_TDEST : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_bdry_in_s_axis_TID : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_bdry_in_s_axis_TKEEP : STD_LOGIC_VECTOR ( 2 downto 0 );
  signal intf_net_bdry_in_s_axis_TLAST : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_bdry_in_s_axis_TREADY : STD_LOGIC;
  signal intf_net_bdry_in_s_axis_TSTRB : STD_LOGIC_VECTOR ( 2 downto 0 );
  signal intf_net_bdry_in_s_axis_TUSER : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_bdry_in_s_axis_TVALID : STD_LOGIC;
  signal intf_net_csc_m_axis_video_TDATA : STD_LOGIC_VECTOR ( 23 downto 0 );
  signal intf_net_csc_m_axis_video_TDEST : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_csc_m_axis_video_TID : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_csc_m_axis_video_TKEEP : STD_LOGIC_VECTOR ( 2 downto 0 );
  signal intf_net_csc_m_axis_video_TLAST : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_csc_m_axis_video_TREADY : STD_LOGIC;
  signal intf_net_csc_m_axis_video_TSTRB : STD_LOGIC_VECTOR ( 2 downto 0 );
  signal intf_net_csc_m_axis_video_TUSER : STD_LOGIC_VECTOR ( 0 to 0 );
  signal intf_net_csc_m_axis_video_TVALID : STD_LOGIC;
  signal net_bdry_in_aclk : STD_LOGIC;
  signal net_bdry_in_aresetn : STD_LOGIC;
  signal NLW_csc_interrupt_UNCONNECTED : STD_LOGIC;
  attribute X_INTERFACE_INFO : string;
  attribute X_INTERFACE_INFO of aclk : signal is "xilinx.com:signal:clock:1.0 CLK.ACLK CLK";
  attribute X_INTERFACE_PARAMETER : string;
  attribute X_INTERFACE_PARAMETER of aclk : signal is "XIL_INTERFACENAME CLK.ACLK, ASSOCIATED_BUSIF m_axis:s_axi_ctrl:s_axis, ASSOCIATED_RESET aresetn, CLK_DOMAIN /clk_wiz_0_clk_out1, FREQ_HZ 148500000, FREQ_TOLERANCE_HZ 0, INSERT_VIP 0, PHASE 0.0";
  attribute X_INTERFACE_INFO of aresetn : signal is "xilinx.com:signal:reset:1.0 RST.ARESETN RST";
  attribute X_INTERFACE_PARAMETER of aresetn : signal is "XIL_INTERFACENAME RST.ARESETN, INSERT_VIP 0, POLARITY ACTIVE_LOW";
  attribute X_INTERFACE_INFO of m_axis_tready : signal is "xilinx.com:interface:axis:1.0 m_axis TREADY";
  attribute X_INTERFACE_INFO of m_axis_tvalid : signal is "xilinx.com:interface:axis:1.0 m_axis TVALID";
  attribute X_INTERFACE_INFO of s_axi_ctrl_arready : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl ARREADY";
  attribute X_INTERFACE_INFO of s_axi_ctrl_arvalid : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl ARVALID";
  attribute X_INTERFACE_INFO of s_axi_ctrl_awready : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl AWREADY";
  attribute X_INTERFACE_INFO of s_axi_ctrl_awvalid : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl AWVALID";
  attribute X_INTERFACE_INFO of s_axi_ctrl_bready : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl BREADY";
  attribute X_INTERFACE_INFO of s_axi_ctrl_bvalid : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl BVALID";
  attribute X_INTERFACE_INFO of s_axi_ctrl_rready : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl RREADY";
  attribute X_INTERFACE_INFO of s_axi_ctrl_rvalid : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl RVALID";
  attribute X_INTERFACE_INFO of s_axi_ctrl_wready : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl WREADY";
  attribute X_INTERFACE_INFO of s_axi_ctrl_wvalid : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl WVALID";
  attribute X_INTERFACE_INFO of s_axis_tready : signal is "xilinx.com:interface:axis:1.0 s_axis TREADY";
  attribute X_INTERFACE_INFO of s_axis_tvalid : signal is "xilinx.com:interface:axis:1.0 s_axis TVALID";
  attribute X_INTERFACE_INFO of m_axis_tdata : signal is "xilinx.com:interface:axis:1.0 m_axis TDATA";
  attribute X_INTERFACE_PARAMETER of m_axis_tdata : signal is "XIL_INTERFACENAME m_axis, CLK_DOMAIN /clk_wiz_0_clk_out1, FREQ_HZ 148500000, HAS_TKEEP 1, HAS_TLAST 1, HAS_TREADY 1, HAS_TSTRB 1, INSERT_VIP 0, LAYERED_METADATA undef, PHASE 0.0, TDATA_NUM_BYTES 3, TDEST_WIDTH 1, TID_WIDTH 1, TUSER_WIDTH 1";
  attribute X_INTERFACE_INFO of m_axis_tdest : signal is "xilinx.com:interface:axis:1.0 m_axis TDEST";
  attribute X_INTERFACE_INFO of m_axis_tid : signal is "xilinx.com:interface:axis:1.0 m_axis TID";
  attribute X_INTERFACE_INFO of m_axis_tkeep : signal is "xilinx.com:interface:axis:1.0 m_axis TKEEP";
  attribute X_INTERFACE_INFO of m_axis_tlast : signal is "xilinx.com:interface:axis:1.0 m_axis TLAST";
  attribute X_INTERFACE_INFO of m_axis_tstrb : signal is "xilinx.com:interface:axis:1.0 m_axis TSTRB";
  attribute X_INTERFACE_INFO of m_axis_tuser : signal is "xilinx.com:interface:axis:1.0 m_axis TUSER";
  attribute X_INTERFACE_INFO of s_axi_ctrl_araddr : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl ARADDR";
  attribute X_INTERFACE_PARAMETER of s_axi_ctrl_araddr : signal is "XIL_INTERFACENAME s_axi_ctrl, ADDR_WIDTH 16, ARUSER_WIDTH 0, AWUSER_WIDTH 0, BUSER_WIDTH 0, CLK_DOMAIN /clk_wiz_0_clk_out1, DATA_WIDTH 32, FREQ_HZ 148500000, HAS_BRESP 1, HAS_BURST 0, HAS_CACHE 0, HAS_LOCK 0, HAS_PROT 0, HAS_QOS 0, HAS_REGION 0, HAS_RRESP 1, HAS_WSTRB 1, ID_WIDTH 0, INSERT_VIP 0, MAX_BURST_LENGTH 1, NUM_READ_OUTSTANDING 1, NUM_READ_THREADS 1, NUM_WRITE_OUTSTANDING 1, NUM_WRITE_THREADS 1, PHASE 0.0, PROTOCOL AXI4LITE, READ_WRITE_MODE READ_WRITE, RUSER_BITS_PER_BYTE 0, RUSER_WIDTH 0, SUPPORTS_NARROW_BURST 0, WUSER_BITS_PER_BYTE 0, WUSER_WIDTH 0";
  attribute X_INTERFACE_INFO of s_axi_ctrl_awaddr : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl AWADDR";
  attribute X_INTERFACE_INFO of s_axi_ctrl_bresp : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl BRESP";
  attribute X_INTERFACE_INFO of s_axi_ctrl_rdata : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl RDATA";
  attribute X_INTERFACE_INFO of s_axi_ctrl_rresp : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl RRESP";
  attribute X_INTERFACE_INFO of s_axi_ctrl_wdata : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl WDATA";
  attribute X_INTERFACE_INFO of s_axi_ctrl_wstrb : signal is "xilinx.com:interface:aximm:1.0 s_axi_ctrl WSTRB";
  attribute X_INTERFACE_INFO of s_axis_tdata : signal is "xilinx.com:interface:axis:1.0 s_axis TDATA";
  attribute X_INTERFACE_PARAMETER of s_axis_tdata : signal is "XIL_INTERFACENAME s_axis, CLK_DOMAIN /clk_wiz_0_clk_out1, FREQ_HZ 148500000, HAS_TKEEP 1, HAS_TLAST 1, HAS_TREADY 1, HAS_TSTRB 1, INSERT_VIP 0, LAYERED_METADATA undef, PHASE 0.0, TDATA_NUM_BYTES 3, TDEST_WIDTH 1, TID_WIDTH 1, TUSER_WIDTH 1";
  attribute X_INTERFACE_INFO of s_axis_tdest : signal is "xilinx.com:interface:axis:1.0 s_axis TDEST";
  attribute X_INTERFACE_INFO of s_axis_tid : signal is "xilinx.com:interface:axis:1.0 s_axis TID";
  attribute X_INTERFACE_INFO of s_axis_tkeep : signal is "xilinx.com:interface:axis:1.0 s_axis TKEEP";
  attribute X_INTERFACE_INFO of s_axis_tlast : signal is "xilinx.com:interface:axis:1.0 s_axis TLAST";
  attribute X_INTERFACE_INFO of s_axis_tstrb : signal is "xilinx.com:interface:axis:1.0 s_axis TSTRB";
  attribute X_INTERFACE_INFO of s_axis_tuser : signal is "xilinx.com:interface:axis:1.0 s_axis TUSER";
begin
  intf_net_bdry_in_s_axi_ctrl_ARADDR(7 downto 0) <= s_axi_ctrl_araddr(7 downto 0);
  intf_net_bdry_in_s_axi_ctrl_ARVALID <= s_axi_ctrl_arvalid;
  intf_net_bdry_in_s_axi_ctrl_AWADDR(7 downto 0) <= s_axi_ctrl_awaddr(7 downto 0);
  intf_net_bdry_in_s_axi_ctrl_AWVALID <= s_axi_ctrl_awvalid;
  intf_net_bdry_in_s_axi_ctrl_BREADY <= s_axi_ctrl_bready;
  intf_net_bdry_in_s_axi_ctrl_RREADY <= s_axi_ctrl_rready;
  intf_net_bdry_in_s_axi_ctrl_WDATA(31 downto 0) <= s_axi_ctrl_wdata(31 downto 0);
  intf_net_bdry_in_s_axi_ctrl_WSTRB(3 downto 0) <= s_axi_ctrl_wstrb(3 downto 0);
  intf_net_bdry_in_s_axi_ctrl_WVALID <= s_axi_ctrl_wvalid;
  intf_net_bdry_in_s_axis_TDATA(23 downto 0) <= s_axis_tdata(23 downto 0);
  intf_net_bdry_in_s_axis_TDEST(0) <= s_axis_tdest(0);
  intf_net_bdry_in_s_axis_TID(0) <= s_axis_tid(0);
  intf_net_bdry_in_s_axis_TKEEP(2 downto 0) <= s_axis_tkeep(2 downto 0);
  intf_net_bdry_in_s_axis_TLAST(0) <= s_axis_tlast(0);
  intf_net_bdry_in_s_axis_TSTRB(2 downto 0) <= s_axis_tstrb(2 downto 0);
  intf_net_bdry_in_s_axis_TUSER(0) <= s_axis_tuser(0);
  intf_net_bdry_in_s_axis_TVALID <= s_axis_tvalid;
  intf_net_csc_m_axis_video_TREADY <= m_axis_tready;
  m_axis_tdata(23 downto 0) <= intf_net_csc_m_axis_video_TDATA(23 downto 0);
  m_axis_tdest(0) <= intf_net_csc_m_axis_video_TDEST(0);
  m_axis_tid(0) <= intf_net_csc_m_axis_video_TID(0);
  m_axis_tkeep(2 downto 0) <= intf_net_csc_m_axis_video_TKEEP(2 downto 0);
  m_axis_tlast(0) <= intf_net_csc_m_axis_video_TLAST(0);
  m_axis_tstrb(2 downto 0) <= intf_net_csc_m_axis_video_TSTRB(2 downto 0);
  m_axis_tuser(0) <= intf_net_csc_m_axis_video_TUSER(0);
  m_axis_tvalid <= intf_net_csc_m_axis_video_TVALID;
  net_bdry_in_aclk <= aclk;
  net_bdry_in_aresetn <= aresetn;
  s_axi_ctrl_arready <= intf_net_bdry_in_s_axi_ctrl_ARREADY;
  s_axi_ctrl_awready <= intf_net_bdry_in_s_axi_ctrl_AWREADY;
  s_axi_ctrl_bresp(1 downto 0) <= intf_net_bdry_in_s_axi_ctrl_BRESP(1 downto 0);
  s_axi_ctrl_bvalid <= intf_net_bdry_in_s_axi_ctrl_BVALID;
  s_axi_ctrl_rdata(31 downto 0) <= intf_net_bdry_in_s_axi_ctrl_RDATA(31 downto 0);
  s_axi_ctrl_rresp(1 downto 0) <= intf_net_bdry_in_s_axi_ctrl_RRESP(1 downto 0);
  s_axi_ctrl_rvalid <= intf_net_bdry_in_s_axi_ctrl_RVALID;
  s_axi_ctrl_wready <= intf_net_bdry_in_s_axi_ctrl_WREADY;
  s_axis_tready <= intf_net_bdry_in_s_axis_TREADY;
csc: component bd_d92b_csc_0
     port map (
      ap_clk => net_bdry_in_aclk,
      ap_rst_n => net_bdry_in_aresetn,
      interrupt => NLW_csc_interrupt_UNCONNECTED,
      m_axis_video_TDATA(23 downto 0) => intf_net_csc_m_axis_video_TDATA(23 downto 0),
      m_axis_video_TDEST(0) => intf_net_csc_m_axis_video_TDEST(0),
      m_axis_video_TID(0) => intf_net_csc_m_axis_video_TID(0),
      m_axis_video_TKEEP(2 downto 0) => intf_net_csc_m_axis_video_TKEEP(2 downto 0),
      m_axis_video_TLAST(0) => intf_net_csc_m_axis_video_TLAST(0),
      m_axis_video_TREADY => intf_net_csc_m_axis_video_TREADY,
      m_axis_video_TSTRB(2 downto 0) => intf_net_csc_m_axis_video_TSTRB(2 downto 0),
      m_axis_video_TUSER(0) => intf_net_csc_m_axis_video_TUSER(0),
      m_axis_video_TVALID => intf_net_csc_m_axis_video_TVALID,
      s_axi_CTRL_ARADDR(7 downto 0) => intf_net_bdry_in_s_axi_ctrl_ARADDR(7 downto 0),
      s_axi_CTRL_ARREADY => intf_net_bdry_in_s_axi_ctrl_ARREADY,
      s_axi_CTRL_ARVALID => intf_net_bdry_in_s_axi_ctrl_ARVALID,
      s_axi_CTRL_AWADDR(7 downto 0) => intf_net_bdry_in_s_axi_ctrl_AWADDR(7 downto 0),
      s_axi_CTRL_AWREADY => intf_net_bdry_in_s_axi_ctrl_AWREADY,
      s_axi_CTRL_AWVALID => intf_net_bdry_in_s_axi_ctrl_AWVALID,
      s_axi_CTRL_BREADY => intf_net_bdry_in_s_axi_ctrl_BREADY,
      s_axi_CTRL_BRESP(1 downto 0) => intf_net_bdry_in_s_axi_ctrl_BRESP(1 downto 0),
      s_axi_CTRL_BVALID => intf_net_bdry_in_s_axi_ctrl_BVALID,
      s_axi_CTRL_RDATA(31 downto 0) => intf_net_bdry_in_s_axi_ctrl_RDATA(31 downto 0),
      s_axi_CTRL_RREADY => intf_net_bdry_in_s_axi_ctrl_RREADY,
      s_axi_CTRL_RRESP(1 downto 0) => intf_net_bdry_in_s_axi_ctrl_RRESP(1 downto 0),
      s_axi_CTRL_RVALID => intf_net_bdry_in_s_axi_ctrl_RVALID,
      s_axi_CTRL_WDATA(31 downto 0) => intf_net_bdry_in_s_axi_ctrl_WDATA(31 downto 0),
      s_axi_CTRL_WREADY => intf_net_bdry_in_s_axi_ctrl_WREADY,
      s_axi_CTRL_WSTRB(3 downto 0) => intf_net_bdry_in_s_axi_ctrl_WSTRB(3 downto 0),
      s_axi_CTRL_WVALID => intf_net_bdry_in_s_axi_ctrl_WVALID,
      s_axis_video_TDATA(23 downto 0) => intf_net_bdry_in_s_axis_TDATA(23 downto 0),
      s_axis_video_TDEST(0) => intf_net_bdry_in_s_axis_TDEST(0),
      s_axis_video_TID(0) => intf_net_bdry_in_s_axis_TID(0),
      s_axis_video_TKEEP(2 downto 0) => intf_net_bdry_in_s_axis_TKEEP(2 downto 0),
      s_axis_video_TLAST(0) => intf_net_bdry_in_s_axis_TLAST(0),
      s_axis_video_TREADY => intf_net_bdry_in_s_axis_TREADY,
      s_axis_video_TSTRB(2 downto 0) => intf_net_bdry_in_s_axis_TSTRB(2 downto 0),
      s_axis_video_TUSER(0) => intf_net_bdry_in_s_axis_TUSER(0),
      s_axis_video_TVALID => intf_net_bdry_in_s_axis_TVALID
    );
end STRUCTURE;
