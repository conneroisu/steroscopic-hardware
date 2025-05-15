--Copyright 1986-2020 Xilinx, Inc. All Rights Reserved.
----------------------------------------------------------------------------------
--Tool Version: Vivado v.2020.1 (win64) Build 2902540 Wed May 27 19:54:49 MDT 2020
--Date        : Tue May 13 18:36:56 2025
--Host        : CO2041-06 running 64-bit major release  (build 9200)
--Command     : generate_target design_1_wrapper.bd
--Design      : design_1_wrapper
--Purpose     : IP block netlist
----------------------------------------------------------------------------------
library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
library UNISIM;
use UNISIM.VCOMPONENTS.ALL;
entity design_1_wrapper is
  port (
    DDR_addr : inout STD_LOGIC_VECTOR ( 14 downto 0 );
    DDR_ba : inout STD_LOGIC_VECTOR ( 2 downto 0 );
    DDR_cas_n : inout STD_LOGIC;
    DDR_ck_n : inout STD_LOGIC;
    DDR_ck_p : inout STD_LOGIC;
    DDR_cke : inout STD_LOGIC;
    DDR_cs_n : inout STD_LOGIC;
    DDR_dm : inout STD_LOGIC_VECTOR ( 3 downto 0 );
    DDR_dq : inout STD_LOGIC_VECTOR ( 31 downto 0 );
    DDR_dqs_n : inout STD_LOGIC_VECTOR ( 3 downto 0 );
    DDR_dqs_p : inout STD_LOGIC_VECTOR ( 3 downto 0 );
    DDR_odt : inout STD_LOGIC;
    DDR_ras_n : inout STD_LOGIC;
    DDR_reset_n : inout STD_LOGIC;
    DDR_we_n : inout STD_LOGIC;
    FIXED_IO_ddr_vrn : inout STD_LOGIC;
    FIXED_IO_ddr_vrp : inout STD_LOGIC;
    FIXED_IO_mio : inout STD_LOGIC_VECTOR ( 53 downto 0 );
    FIXED_IO_ps_clk : inout STD_LOGIC;
    FIXED_IO_ps_porb : inout STD_LOGIC;
    FIXED_IO_ps_srstb : inout STD_LOGIC;
    IO_HDMII_spdif : in STD_LOGIC;
    IO_HDMIO_clk : out STD_LOGIC;
    IO_HDMIO_data : out STD_LOGIC_VECTOR ( 15 downto 0 );
    IO_HDMIO_spdif : out STD_LOGIC;
    IO_VITA_CAM_clk_out_n : in STD_LOGIC;
    IO_VITA_CAM_clk_out_p : in STD_LOGIC;
    IO_VITA_CAM_clk_pll : out STD_LOGIC;
    IO_VITA_CAM_data_n : in STD_LOGIC_VECTOR ( 3 downto 0 );
    IO_VITA_CAM_data_p : in STD_LOGIC_VECTOR ( 3 downto 0 );
    IO_VITA_CAM_monitor : in STD_LOGIC_VECTOR ( 1 downto 0 );
    IO_VITA_CAM_reset_n : out STD_LOGIC;
    IO_VITA_CAM_sync_n : in STD_LOGIC;
    IO_VITA_CAM_sync_p : in STD_LOGIC;
    IO_VITA_CAM_trigger : out STD_LOGIC_VECTOR ( 2 downto 0 );
    IO_VITA_SPI_spi_miso : in STD_LOGIC;
    IO_VITA_SPI_spi_mosi : out STD_LOGIC;
    IO_VITA_SPI_spi_sclk : out STD_LOGIC;
    IO_VITA_SPI_spi_ssel_n : out STD_LOGIC;
    PMOD_io0_io : inout STD_LOGIC;
    PMOD_io1_io : inout STD_LOGIC;
    PMOD_sck_io : inout STD_LOGIC;
    PMOD_ss_io : inout STD_LOGIC_VECTOR ( 0 to 0 );
    btns_5bits_tri_i : in STD_LOGIC_VECTOR ( 4 downto 0 );
    fmc_imageon_iic_rst_n : out STD_LOGIC_VECTOR ( 0 to 0 );
    fmc_imageon_iic_scl_io : inout STD_LOGIC;
    fmc_imageon_iic_sda_io : inout STD_LOGIC;
    fmc_imageon_vclk : in STD_LOGIC;
    fmc_ipmi_id_eeprom_scl_io : inout STD_LOGIC;
    fmc_ipmi_id_eeprom_sda_io : inout STD_LOGIC;
    sws_8bits_tri_i : in STD_LOGIC_VECTOR ( 7 downto 0 )
  );
end design_1_wrapper;

architecture STRUCTURE of design_1_wrapper is
  component design_1 is
  port (
    fmc_imageon_iic_rst_n : out STD_LOGIC_VECTOR ( 0 to 0 );
    IO_HDMII_spdif : in STD_LOGIC;
    fmc_imageon_vclk : in STD_LOGIC;
    btns_5bits_tri_i : in STD_LOGIC_VECTOR ( 4 downto 0 );
    sws_8bits_tri_i : in STD_LOGIC_VECTOR ( 7 downto 0 );
    IO_VITA_SPI_spi_sclk : out STD_LOGIC;
    IO_VITA_SPI_spi_ssel_n : out STD_LOGIC;
    IO_VITA_SPI_spi_mosi : out STD_LOGIC;
    IO_VITA_SPI_spi_miso : in STD_LOGIC;
    DDR_cas_n : inout STD_LOGIC;
    DDR_cke : inout STD_LOGIC;
    DDR_ck_n : inout STD_LOGIC;
    DDR_ck_p : inout STD_LOGIC;
    DDR_cs_n : inout STD_LOGIC;
    DDR_reset_n : inout STD_LOGIC;
    DDR_odt : inout STD_LOGIC;
    DDR_ras_n : inout STD_LOGIC;
    DDR_we_n : inout STD_LOGIC;
    DDR_ba : inout STD_LOGIC_VECTOR ( 2 downto 0 );
    DDR_addr : inout STD_LOGIC_VECTOR ( 14 downto 0 );
    DDR_dm : inout STD_LOGIC_VECTOR ( 3 downto 0 );
    DDR_dq : inout STD_LOGIC_VECTOR ( 31 downto 0 );
    DDR_dqs_n : inout STD_LOGIC_VECTOR ( 3 downto 0 );
    DDR_dqs_p : inout STD_LOGIC_VECTOR ( 3 downto 0 );
    fmc_ipmi_id_eeprom_scl_i : in STD_LOGIC;
    fmc_ipmi_id_eeprom_scl_o : out STD_LOGIC;
    fmc_ipmi_id_eeprom_scl_t : out STD_LOGIC;
    fmc_ipmi_id_eeprom_sda_i : in STD_LOGIC;
    fmc_ipmi_id_eeprom_sda_o : out STD_LOGIC;
    fmc_ipmi_id_eeprom_sda_t : out STD_LOGIC;
    IO_HDMIO_clk : out STD_LOGIC;
    IO_HDMIO_data : out STD_LOGIC_VECTOR ( 15 downto 0 );
    IO_HDMIO_spdif : out STD_LOGIC;
    FIXED_IO_mio : inout STD_LOGIC_VECTOR ( 53 downto 0 );
    FIXED_IO_ddr_vrn : inout STD_LOGIC;
    FIXED_IO_ddr_vrp : inout STD_LOGIC;
    FIXED_IO_ps_srstb : inout STD_LOGIC;
    FIXED_IO_ps_clk : inout STD_LOGIC;
    FIXED_IO_ps_porb : inout STD_LOGIC;
    IO_VITA_CAM_data_p : in STD_LOGIC_VECTOR ( 3 downto 0 );
    IO_VITA_CAM_sync_p : in STD_LOGIC;
    IO_VITA_CAM_sync_n : in STD_LOGIC;
    IO_VITA_CAM_reset_n : out STD_LOGIC;
    IO_VITA_CAM_trigger : out STD_LOGIC_VECTOR ( 2 downto 0 );
    IO_VITA_CAM_monitor : in STD_LOGIC_VECTOR ( 1 downto 0 );
    IO_VITA_CAM_clk_pll : out STD_LOGIC;
    IO_VITA_CAM_data_n : in STD_LOGIC_VECTOR ( 3 downto 0 );
    IO_VITA_CAM_clk_out_p : in STD_LOGIC;
    IO_VITA_CAM_clk_out_n : in STD_LOGIC;
    fmc_imageon_iic_scl_i : in STD_LOGIC;
    fmc_imageon_iic_scl_o : out STD_LOGIC;
    fmc_imageon_iic_scl_t : out STD_LOGIC;
    fmc_imageon_iic_sda_i : in STD_LOGIC;
    fmc_imageon_iic_sda_o : out STD_LOGIC;
    fmc_imageon_iic_sda_t : out STD_LOGIC;
    PMOD_io0_i : in STD_LOGIC;
    PMOD_io0_o : out STD_LOGIC;
    PMOD_io0_t : out STD_LOGIC;
    PMOD_io1_i : in STD_LOGIC;
    PMOD_io1_o : out STD_LOGIC;
    PMOD_io1_t : out STD_LOGIC;
    PMOD_ss_i : in STD_LOGIC_VECTOR ( 0 to 0 );
    PMOD_ss_o : out STD_LOGIC_VECTOR ( 0 to 0 );
    PMOD_ss_t : out STD_LOGIC;
    PMOD_sck_i : in STD_LOGIC;
    PMOD_sck_o : out STD_LOGIC;
    PMOD_sck_t : out STD_LOGIC
  );
  end component design_1;
  component IOBUF is
  port (
    I : in STD_LOGIC;
    O : out STD_LOGIC;
    T : in STD_LOGIC;
    IO : inout STD_LOGIC
  );
  end component IOBUF;
  signal PMOD_io0_i : STD_LOGIC;
  signal PMOD_io0_o : STD_LOGIC;
  signal PMOD_io0_t : STD_LOGIC;
  signal PMOD_io1_i : STD_LOGIC;
  signal PMOD_io1_o : STD_LOGIC;
  signal PMOD_io1_t : STD_LOGIC;
  signal PMOD_sck_i : STD_LOGIC;
  signal PMOD_sck_o : STD_LOGIC;
  signal PMOD_sck_t : STD_LOGIC;
  signal PMOD_ss_i_0 : STD_LOGIC_VECTOR ( 0 to 0 );
  signal PMOD_ss_io_0 : STD_LOGIC_VECTOR ( 0 to 0 );
  signal PMOD_ss_o_0 : STD_LOGIC_VECTOR ( 0 to 0 );
  signal PMOD_ss_t : STD_LOGIC;
  signal fmc_imageon_iic_scl_i : STD_LOGIC;
  signal fmc_imageon_iic_scl_o : STD_LOGIC;
  signal fmc_imageon_iic_scl_t : STD_LOGIC;
  signal fmc_imageon_iic_sda_i : STD_LOGIC;
  signal fmc_imageon_iic_sda_o : STD_LOGIC;
  signal fmc_imageon_iic_sda_t : STD_LOGIC;
  signal fmc_ipmi_id_eeprom_scl_i : STD_LOGIC;
  signal fmc_ipmi_id_eeprom_scl_o : STD_LOGIC;
  signal fmc_ipmi_id_eeprom_scl_t : STD_LOGIC;
  signal fmc_ipmi_id_eeprom_sda_i : STD_LOGIC;
  signal fmc_ipmi_id_eeprom_sda_o : STD_LOGIC;
  signal fmc_ipmi_id_eeprom_sda_t : STD_LOGIC;
begin
PMOD_io0_iobuf: component IOBUF
     port map (
      I => PMOD_io0_o,
      IO => PMOD_io0_io,
      O => PMOD_io0_i,
      T => PMOD_io0_t
    );
PMOD_io1_iobuf: component IOBUF
     port map (
      I => PMOD_io1_o,
      IO => PMOD_io1_io,
      O => PMOD_io1_i,
      T => PMOD_io1_t
    );
PMOD_sck_iobuf: component IOBUF
     port map (
      I => PMOD_sck_o,
      IO => PMOD_sck_io,
      O => PMOD_sck_i,
      T => PMOD_sck_t
    );
PMOD_ss_iobuf_0: component IOBUF
     port map (
      I => PMOD_ss_o_0(0),
      IO => PMOD_ss_io(0),
      O => PMOD_ss_i_0(0),
      T => PMOD_ss_t
    );
design_1_i: component design_1
     port map (
      DDR_addr(14 downto 0) => DDR_addr(14 downto 0),
      DDR_ba(2 downto 0) => DDR_ba(2 downto 0),
      DDR_cas_n => DDR_cas_n,
      DDR_ck_n => DDR_ck_n,
      DDR_ck_p => DDR_ck_p,
      DDR_cke => DDR_cke,
      DDR_cs_n => DDR_cs_n,
      DDR_dm(3 downto 0) => DDR_dm(3 downto 0),
      DDR_dq(31 downto 0) => DDR_dq(31 downto 0),
      DDR_dqs_n(3 downto 0) => DDR_dqs_n(3 downto 0),
      DDR_dqs_p(3 downto 0) => DDR_dqs_p(3 downto 0),
      DDR_odt => DDR_odt,
      DDR_ras_n => DDR_ras_n,
      DDR_reset_n => DDR_reset_n,
      DDR_we_n => DDR_we_n,
      FIXED_IO_ddr_vrn => FIXED_IO_ddr_vrn,
      FIXED_IO_ddr_vrp => FIXED_IO_ddr_vrp,
      FIXED_IO_mio(53 downto 0) => FIXED_IO_mio(53 downto 0),
      FIXED_IO_ps_clk => FIXED_IO_ps_clk,
      FIXED_IO_ps_porb => FIXED_IO_ps_porb,
      FIXED_IO_ps_srstb => FIXED_IO_ps_srstb,
      IO_HDMII_spdif => IO_HDMII_spdif,
      IO_HDMIO_clk => IO_HDMIO_clk,
      IO_HDMIO_data(15 downto 0) => IO_HDMIO_data(15 downto 0),
      IO_HDMIO_spdif => IO_HDMIO_spdif,
      IO_VITA_CAM_clk_out_n => IO_VITA_CAM_clk_out_n,
      IO_VITA_CAM_clk_out_p => IO_VITA_CAM_clk_out_p,
      IO_VITA_CAM_clk_pll => IO_VITA_CAM_clk_pll,
      IO_VITA_CAM_data_n(3 downto 0) => IO_VITA_CAM_data_n(3 downto 0),
      IO_VITA_CAM_data_p(3 downto 0) => IO_VITA_CAM_data_p(3 downto 0),
      IO_VITA_CAM_monitor(1 downto 0) => IO_VITA_CAM_monitor(1 downto 0),
      IO_VITA_CAM_reset_n => IO_VITA_CAM_reset_n,
      IO_VITA_CAM_sync_n => IO_VITA_CAM_sync_n,
      IO_VITA_CAM_sync_p => IO_VITA_CAM_sync_p,
      IO_VITA_CAM_trigger(2 downto 0) => IO_VITA_CAM_trigger(2 downto 0),
      IO_VITA_SPI_spi_miso => IO_VITA_SPI_spi_miso,
      IO_VITA_SPI_spi_mosi => IO_VITA_SPI_spi_mosi,
      IO_VITA_SPI_spi_sclk => IO_VITA_SPI_spi_sclk,
      IO_VITA_SPI_spi_ssel_n => IO_VITA_SPI_spi_ssel_n,
      PMOD_io0_i => PMOD_io0_i,
      PMOD_io0_o => PMOD_io0_o,
      PMOD_io0_t => PMOD_io0_t,
      PMOD_io1_i => PMOD_io1_i,
      PMOD_io1_o => PMOD_io1_o,
      PMOD_io1_t => PMOD_io1_t,
      PMOD_sck_i => PMOD_sck_i,
      PMOD_sck_o => PMOD_sck_o,
      PMOD_sck_t => PMOD_sck_t,
      PMOD_ss_i(0) => PMOD_ss_i_0(0),
      PMOD_ss_o(0) => PMOD_ss_o_0(0),
      PMOD_ss_t => PMOD_ss_t,
      btns_5bits_tri_i(4 downto 0) => btns_5bits_tri_i(4 downto 0),
      fmc_imageon_iic_rst_n(0) => fmc_imageon_iic_rst_n(0),
      fmc_imageon_iic_scl_i => fmc_imageon_iic_scl_i,
      fmc_imageon_iic_scl_o => fmc_imageon_iic_scl_o,
      fmc_imageon_iic_scl_t => fmc_imageon_iic_scl_t,
      fmc_imageon_iic_sda_i => fmc_imageon_iic_sda_i,
      fmc_imageon_iic_sda_o => fmc_imageon_iic_sda_o,
      fmc_imageon_iic_sda_t => fmc_imageon_iic_sda_t,
      fmc_imageon_vclk => fmc_imageon_vclk,
      fmc_ipmi_id_eeprom_scl_i => fmc_ipmi_id_eeprom_scl_i,
      fmc_ipmi_id_eeprom_scl_o => fmc_ipmi_id_eeprom_scl_o,
      fmc_ipmi_id_eeprom_scl_t => fmc_ipmi_id_eeprom_scl_t,
      fmc_ipmi_id_eeprom_sda_i => fmc_ipmi_id_eeprom_sda_i,
      fmc_ipmi_id_eeprom_sda_o => fmc_ipmi_id_eeprom_sda_o,
      fmc_ipmi_id_eeprom_sda_t => fmc_ipmi_id_eeprom_sda_t,
      sws_8bits_tri_i(7 downto 0) => sws_8bits_tri_i(7 downto 0)
    );
fmc_imageon_iic_scl_iobuf: component IOBUF
     port map (
      I => fmc_imageon_iic_scl_o,
      IO => fmc_imageon_iic_scl_io,
      O => fmc_imageon_iic_scl_i,
      T => fmc_imageon_iic_scl_t
    );
fmc_imageon_iic_sda_iobuf: component IOBUF
     port map (
      I => fmc_imageon_iic_sda_o,
      IO => fmc_imageon_iic_sda_io,
      O => fmc_imageon_iic_sda_i,
      T => fmc_imageon_iic_sda_t
    );
fmc_ipmi_id_eeprom_scl_iobuf: component IOBUF
     port map (
      I => fmc_ipmi_id_eeprom_scl_o,
      IO => fmc_ipmi_id_eeprom_scl_io,
      O => fmc_ipmi_id_eeprom_scl_i,
      T => fmc_ipmi_id_eeprom_scl_t
    );
fmc_ipmi_id_eeprom_sda_iobuf: component IOBUF
     port map (
      I => fmc_ipmi_id_eeprom_sda_o,
      IO => fmc_ipmi_id_eeprom_sda_io,
      O => fmc_ipmi_id_eeprom_sda_i,
      T => fmc_ipmi_id_eeprom_sda_t
    );
end STRUCTURE;
