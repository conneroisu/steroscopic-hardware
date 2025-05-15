/*****************************************************************************
 * Joseph Zambreno
 * Phillip Jones
 *
 * Department of Electrical and Computer Engineering
 * Iowa State University
 *****************************************************************************/

/*****************************************************************************
 * fmc_imageon_utils.c - main initialization and configuration routines
 * for the Avnet FMC-IMAGEON board. Most of the bugs are their fault.
 *
 *
 * NOTES:
 * 02/04/14 by JAZ::Design created.
 *****************************************************************************/

#include "camera_app.h"
#include "xil_cache.h"


// Main FMC-IMAGEON initialization function. Add your code here.
int fmc_imageon_enable( camera_config_t *config )
{
   int ret;

   xil_printf("\n\r");
   xil_printf("------------------------------------------------\n\r");
   xil_printf("--    FMC-IMAGEON Camera Application (MP-2)   --\n\r");
   xil_printf("------------------------------------------------\n\r");
   xil_printf("\n\r");

   config->bVerbose = 1;
   config->vita_aec = 0;       // off
   config->vita_again = 0;     // 1.0
   config->vita_dgain = 128;   // 1.0
   config->vita_exposure = 90; // 90% of frame period


   xil_printf("FMC-IPMI Initialization ...\n\r");
   ret = fmc_iic_axi_init(&(config->fmc_ipmi_iic), "FMC-IPMI I2C Controller", config->uBaseAddr_IIC_FmcIpmi);
   if (!ret ) {
      xil_printf("ERROR: Failed to open FMC-IIC driver,\n\r");
      exit(1);
   }


   // FMC Module Validation
   if (fmc_ipmi_detect(&(config->fmc_ipmi_iic), "FMC-IMAGEON", FMC_ID_ALL)) {
      fmc_ipmi_enable( &(config->fmc_ipmi_iic), FMC_ID_SLOT1 );
   }
   else {
	   xil_printf("ERROR: Failed to validate FMC-IPMI I2C Controller.\n\r");
       exit(1);
   }


   xil_printf("FMC-IMAGEON I2C Initialization ...\n\r");
   ret = fmc_iic_axi_init(&(config->fmc_imageon_iic), "FMC-IMAGEON I2C Controller", config->uBaseAddr_IIC_FmcImageon);
   if (!ret) {
      xil_printf( "ERROR: Failed to open FMC-IIC driver\n\r");
      exit(1);
   }

   xil_printf("FMC-IMAGEON Video Clock Initialization ...\n\r");
   fmc_imageon_init(&(config->fmc_imageon), "FMC-IMAGEON", &(config->fmc_imageon_iic));
   fmc_imageon_vclk_init( &(config->fmc_imageon) );
   fmc_imageon_vclk_config( &(config->fmc_imageon), FMC_IMAGEON_VCLK_FREQ_148_500_000);

   xil_printf("Resetting clock generator ...\n\r");
   reset_dcms(config);

   // Initialize Video Output Timing
   xil_printf("Initializing Video Output for 1080P60 ...\n\r");

   config->hdmio_width  = 1920;
   config->hdmio_height = 1080;
   config->hdmio_timing.IsHDMI        = 0; // DVI Mode
   config->hdmio_timing.IsEncrypted   = 0;
   config->hdmio_timing.IsInterlaced  = 0;
   config->hdmio_timing.ColorDepth    = 8;
   config->hdmio_timing.HActiveVideo  = 1920;
   config->hdmio_timing.HFrontPorch   =   88;
   config->hdmio_timing.HSyncWidth    =   44;
   config->hdmio_timing.HSyncPolarity =    1;
   config->hdmio_timing.HBackPorch    =  148;
   config->hdmio_timing.VActiveVideo  = 1080;
   config->hdmio_timing.VFrontPorch   =    4;
   config->hdmio_timing.VSyncWidth    =    5;
   config->hdmio_timing.VSyncPolarity =    1;
   config->hdmio_timing.VBackPorch    =   36;

   if (config->bVerbose) {
  	   xil_printf("ADV7511 Video Output Information\n\r");
    	   xil_printf( "\tHSYNC Timing     = hav=%04d, hfp=%02d, hsw=%02d(hsp=%d), hbp=%03d\n\r",
      			   config->hdmio_timing.HActiveVideo,
      			   config->hdmio_timing.HFrontPorch,
      			   config->hdmio_timing.HSyncWidth, config->hdmio_timing.HSyncPolarity,
      			   config->hdmio_timing.HBackPorch);
    	   xil_printf( "\tVSYNC Timing     = vav=%04d, vfp=%02d, vsw=%02d(vsp=%d), vbp=%03d\n\r",
    			   config->hdmio_timing.VActiveVideo,
    			   config->hdmio_timing.VFrontPorch,
    			   config->hdmio_timing.VSyncWidth, config->hdmio_timing.VSyncPolarity,
    			   config->hdmio_timing.VBackPorch);
    	   xil_printf( "\tVideo Dimensions = %d x %d\n\r", config->hdmio_width, config->hdmio_height );

   }

   config->hdmio_resolution = vres_detect(config->hdmio_width, config->hdmio_height);
   xil_printf("\tVideo Resolution = %s\n\r", vres_get_name(config->hdmio_resolution));

   xil_printf( "Video Generator Configuration ...\n\r");
   vgen_init( &(config->vtc_tpg), config->uDeviceId_VTC_tpg);
   vgen_config( &(config->vtc_tpg ), config->hdmio_resolution, 1);


   // FMC-IMAGEON HDMI Output Initialization
   xil_printf( "FMC-IMAGEON HDMI Output Initialization ...\n\r" );
   ret = fmc_imageon_hdmio_init(&(config->fmc_imageon), 1, &(config->hdmio_timing), 0);
   if (!ret) {
      xil_printf("ERROR : Failed to init FMC-IMAGEON HDMI Output Interface\n\r");
      exit(0);
   }


   // FMC-IMAGEON VITA Camera Receiver Initialization
   xil_printf( "FMC-IMAGEON VITA Camera Initialization ...\n\r" );
   onsemi_vita_init( &(config->onsemi_vita), "VITA-2000", config->uBaseAddr_VITA_SPI, config->uBaseAddr_VITA_CAM);
   config->onsemi_vita.uManualTap = 25;
   // Assuming a 75 MHz AXI-Lite SPI bus
   onsemi_vita_spi_config( &(config->onsemi_vita), (75000000/10000000) ); // AXI-Lite SPI Speed (HZ) / 10,000,000 Hz



   // Enable spread-spectrum clocking (SSC)
   enable_ssc(config);

   // Clear frame stores
   Xuint32 i;
   Xuint32 storage_size = config->uNumFrames_HdmiFrameBuffer * ((1920*1080)<<1);
   volatile Xuint32 *pStorageMem = (Xuint32 *)config->uBaseAddr_MEM_HdmiFrameBuffer;

   xil_printf("pStorageMem = %X\n\r", pStorageMem);

   // Frame #1 - Red pixels
   for (i = 0; i < storage_size / config->uNumFrames_HdmiFrameBuffer; i += 4) {
	  *pStorageMem++ = 0xF0525A52;  // Red
   }
   // Frame #2 - Green pixels
   for (i = 0; i < storage_size / config->uNumFrames_HdmiFrameBuffer; i += 4) {
	  *pStorageMem++ = 0x36912291; // Green
   }
   // Frame #3 - Blue pixels
   for (i = 0; i < storage_size / config->uNumFrames_HdmiFrameBuffer; i += 4) {
	  *pStorageMem++ = 0x6E29F029;  // Blue
   }

   Xil_DCacheFlush(); // Flush Cache

   // Initialize Output Side of AXI VDMA
   xil_printf( "Video DMA (Output Side) Initialization ...\n\r" );
   vfb_common_init(
      config->uDeviceId_VDMA_HdmiFrameBuffer, // uDeviceId
      &(config->vdma_hdmi)                    // pAxiVdma
      );
   vfb_tx_init(
      &(config->vdma_hdmi),                   // pAxiVdma
      &(config->vdmacfg_hdmi_read),           // pReadCfg
      config->hdmio_resolution,               // uVideoResolution
      config->hdmio_resolution,               // uStorageResolution
      config->uBaseAddr_MEM_HdmiFrameBuffer,  // uMemAddr
      config->uNumFrames_HdmiFrameBuffer      // uNumFrames
      );


   // Output static Frame buffer for 5 seconds
   xil_printf( "Output static Frame buffer for 5 seconds\n\r" );
   sleep(5);


   // Initialize Input Side of AXI VDMA
   xil_printf( "Video DMA (Input Side) Initialization ...\n\r" );
   vfb_rx_init(
      &(config->vdma_hdmi),                   // pAxiVdma
      &(config->vdmacfg_hdmi_write),          // pWriteCfg
      config->hdmio_resolution,               // uVideoResolution
      config->hdmio_resolution,               // uStorageResolution
      config->uBaseAddr_MEM_HdmiFrameBuffer,  // uMemAddr
      config->uNumFrames_HdmiFrameBuffer      // uNumFrames
      );


   // Choose Video Source  ( 1. TPG or 2. onsemi VITA Camera)

   // 1. Uncomment for TPG
//   fmc_imageon_enable_vita(config);

   // 2. Uncomment for onsemi VITA Camera
   int vita_enabled_error = 0;
   int vita_enable_attempt=1;
   do {
	   xil_printf("\r\n\n\nFMC_IMAGEON_ENABLE_VITA, attempt %d\r\n\n\n", vita_enable_attempt++);
	   vita_enabled_error = fmc_imageon_enable_vita(config);
   } while(vita_enabled_error != 0);


     // Uncomment to enable HW Video processing pipeling (last part of lab)
     // You need to complete implmentation of this function before enabling
   //// TODO:
   fmc_imageon_enable_ipipe(config);


   // Output Video input source in Hardware mode for 10 seconds
   xil_printf( "Output Video input source in Hardware mode for 10 seconds\n\r" );
   sleep(10);


   // Status of AXI VDMA
   vfb_dump_registers( &(config->vdma_hdmi) );
   if ( vfb_check_errors( &(config->vdma_hdmi), 1/*clear errors, if any*/ ) ) {
      vfb_dump_registers( &(config->vdma_hdmi) );
   }

   xil_printf("\n\r");
   xil_printf( "Done\n\r" );
   xil_printf("\n\r");

   return 0;
}



int fmc_imageon_enable_tpg( camera_config_t *config ) {
   Xuint32 Reg32Value;
   Xuint32 CurrentSpeed;

   // Define convenient volatile pointers for accessing TPG registers
   volatile unsigned int *TPG_CR       = config->uBaseAddr_TPG_PatternGenerator + 0;    // TPG Control
   volatile unsigned int *TPG_Act_H    = config->uBaseAddr_TPG_PatternGenerator + 0x10; // Active Height
   volatile unsigned int *TPG_Act_W    = config->uBaseAddr_TPG_PatternGenerator + 0x18; // Active Width
   volatile unsigned int *TPG_BGP      = config->uBaseAddr_TPG_PatternGenerator + 0x20; // Background Pattern
   volatile unsigned int *TPG_FGP      = config->uBaseAddr_TPG_PatternGenerator + 0x28; // Foreground Pattern
   volatile unsigned int *TPG_MS       = config->uBaseAddr_TPG_PatternGenerator + 0x38; // Motion Speed
   volatile unsigned int *TPG_CF       = config->uBaseAddr_TPG_PatternGenerator + 0x40; // TPG Color Format


   xil_printf("Test Pattern Generator Initialization ...\n\r");

   // Direct Memory Mapped access of TPG configuration registers
   // See TPG data sheet for configuring the TPG for other features
   TPG_Act_H[0]  = 0x438; // Active Height
   TPG_Act_W[0]  = 0x780; // Active Width
   TPG_BGP[0]    = 0x0A;  // Background Pattern
   TPG_FGP[0]    = 0x00;  // Foreground Pattern
   TPG_MS[0]     = 0x04;  // Motion Speed
   TPG_CF[0]     = 0x02;  // TPG Color Format
   TPG_CR[0]     = 0x81;  // TPG Control


   // (UNDER CONSTRUCTION: DO NOT USE!!) Useful portable wrapper functions for accessing TPG configuration registers
   /*
   // See TPG data sheet for configuring the TPG for other features

   xTPG_config(config->hdmio_width, config->hdmio_height, config->uBaseAddr_TPG_PatternGenerator, 1);

   // Default to Zone Plate
   Reg32Value = XTPG_mReadSlaveReg0(config->uBaseAddr_TPG_PatternGenerator, 0);
   XTPG_mWriteSlaveReg0(config->uBaseAddr_TPG_PatternGenerator, 0, (Reg32Value & 0xFFFFFFF0) | 10);
   ZPlate_config(config->hdmio_width, config->hdmio_height, config->uBaseAddr_TPG_PatternGenerator);

   // Enable motion
   Reg32Value = XTPG_mReadSlaveReg1(config->uBaseAddr_TPG_PatternGenerator, 0);
   if ((Reg32Value & 0x00000001)==0)
	   XTPG_mWriteSlaveReg1(config->uBaseAddr_TPG_PatternGenerator, 0, (Reg32Value | 0x00000001));
   else
	   XTPG_mWriteSlaveReg1(config->uBaseAddr_TPG_PatternGenerator, 0, (Reg32Value & 0xFFFFFFFE));

   // Speed = 4
   Reg32Value = XTPG_mReadSlaveReg1(config->uBaseAddr_TPG_PatternGenerator, 0);
   CurrentSpeed = (Reg32Value&0x000001fe)>>1;
   CurrentSpeed = 4;
   XTPG_mWriteSlaveReg1(config->uBaseAddr_TPG_PatternGenerator, 0, (CurrentSpeed<<1)|(Reg32Value & 0x00000001));

   */

   return 0;
}


// Disable Test Pattern Generator (TPG)
int fmc_imageon_disable_tpg( camera_config_t *config ) {

	// Define convenient volatile pointers for accessing TPG registers
   volatile unsigned int *TPG_CR  = config->uBaseAddr_TPG_PatternGenerator + 0;    // TPG Control

   xil_printf("Test Pattern Generator Disable ...\n\r");
   TPG_CR[0]     = 0x00;  // TPG Control Register (Disable TPG)

   return 0;
}



int fmc_imageon_enable_vita( camera_config_t *config ) {
   int ret;

   // VITA-2000 Initialization
   xil_printf( "FMC-IMAGEON VITA Initialization ...\n\r" );
   ret = onsemi_vita_sensor_initialize( &(config->onsemi_vita), SENSOR_INIT_ENABLE, config->bVerbose );
   if (ret == 0) {
      xil_printf("VITA sensor failed to initialize ...\n\r");
      return -1;
   }

   onsemi_vita_sensor_initialize( &(config->onsemi_vita), SENSOR_INIT_STREAMON, config->bVerbose);
   sleep(1);

   xil_printf("FMC-IMAGEON VITA Configuration for 1080P60 timing ...\n\r");
   ret = onsemi_vita_sensor_1080P60( &(config->onsemi_vita), config->bVerbose);
   if (ret == 0 ) {
      xil_printf( "VITA sensor failed to configure for 1080P60 timing ...\n\r" );
      return -1;
   }
   sleep(1);

   onsemi_vita_get_status( &(config->onsemi_vita), &(config->vita_status_t1), 0/*config->bVerbose*/);
   sleep(1);
   onsemi_vita_get_status( &(config->onsemi_vita), &(config->vita_status_t2), 0/*config->bVerbose*/);

   int vita_width, vita_height, vita_rate, vita_crc;
   vita_width = config->vita_status_t1.cntImagePixels * 4;
   vita_height = config->vita_status_t1.cntImageLines;
   vita_rate = config->vita_status_t2.cntFrames - config->vita_status_t1.cntFrames;
   vita_crc = config->vita_status_t2.crcStatus;
   xil_printf("VITA Status = \n\r");
   xil_printf("\tImage Width  = %d\n\r", vita_width);
   xil_printf("\tImage Height = %d\n\r", vita_height);
   xil_printf("\tFrame Rate   = %d frames/sec\n\r", vita_rate);
   xil_printf("\tCRC = %d\n\r", vita_crc);

   if ( config->bVerbose )
   {
     onsemi_vita_get_status( &(config->onsemi_vita), &(config->vita_status_t2), 1 );
   }


   if ((vita_width != 1920) || (vita_height != 1080) || (vita_rate == 0)) {
	   return 1;
   }

   return 0;
}


// Uncomment for Hardware Image color Pipeline
 XVprocSs proc_ss_RGB_YCrCb_444;  // To hold info for the Video Processing Subsystem: Color Conversion Only IP core: See xv_procss.h, and xv_csc_l2.h
 XVprocSs_Config *Config_ptr;

int fmc_imageon_enable_ipipe( camera_config_t *config ) {


   xil_printf("Hardware Image Processing Pipeline (iPIPE) Initialization ...\n\r" );

   //// TODO:
   // Video Processing Subsystem (Only re-sampling) 4:4:4 to 4:2:2  (See IP documentation for register details)
   // TODO Add additional register assignments here to fully configure this core. See Video Processing Subsystem IP documentation for register details.
   // Hint 1: You will need to configure 4 additional registers, 
   // Hint 2: You will need to dig through some header files for some of the values,
   // Hint 3: Set up the Demosaic IP block first. It is a good warm up.
   
   // Add assignments here

   Xil_Out32(XPAR_V_PROC_SS_1_BASEADDR + XV_HCRESAMPLER_CTRL_ADDR_HWREG_WIDTH_DATA, 1920);  // Width = 1920 pixels
   Xil_Out32(XPAR_V_PROC_SS_1_BASEADDR + XV_HCRESAMPLER_CTRL_ADDR_HWREG_HEIGHT_DATA, 1080); // Height = 1080 lines

   Xil_Out32(XPAR_V_PROC_SS_1_BASEADDR + XV_HCRESAMPLER_CTRL_ADDR_HWREG_INPUT_VIDEO_FORMAT_DATA, 0x01);  // Input format = 4:4:4
   Xil_Out32(XPAR_V_PROC_SS_1_BASEADDR + XV_HCRESAMPLER_CTRL_ADDR_HWREG_OUTPUT_VIDEO_FORMAT_DATA, 0x02); // Output format = 4:2:2


   //Xil_Out32((XPAR_V_PROC_SS_1_BASEADDR) + 0x30, (u32)(0x0));  // Set chroma phase
   xil_printf("4:4:4 to 4:2:2 Re-sampling IP Configuration and Enable done\r\n");


   // Uncomment as part of Re-sampler IP setup
   Xil_Out32((XPAR_V_PROC_SS_1_BASEADDR) + (XV_HCRESAMPLER_CTRL_ADDR_AP_CTRL), (u32)(0x81));  // Control
   xil_printf("4:4:4 to 4:2:2 Re-sampling IP Configuration and Enable done\r\n");

   // Video Processing Subsystem Hardware IP (configured for Only Color Conversion) from 24-bit RGB to YCrCb 4:4:4
   // For this IP core there was a High level API for setting it up.  Trace through these calls to see how much work they do.
   // This could have been set up with direct register writes, but there are about 20 registers that need to be set for this IP block

   // Uncomment to setup Color Conversion IP
   Config_ptr = XVprocSs_LookupConfig(XPAR_XVPROCSS_0_DEVICE_ID);
   XVprocSs_CfgInitialize(&proc_ss_RGB_YCrCb_444, Config_ptr, XPAR_XVPROCSS_0_BASEADDR);
   XVprocSs_SetSubsystemConfig(&proc_ss_RGB_YCrCb_444);
   XV_CscSetColorspace(proc_ss_RGB_YCrCb_444.CscPtr, XVIDC_CSF_RGB, XVIDC_CSF_YCRCB_444, XVIDC_BT_709, XVIDC_BT_709, XVIDC_CR_0_255);
   XVprocSs_Start(&proc_ss_RGB_YCrCb_444);
   xil_printf("RGB to 4:4:4 IP Configuration and Enable done\r\n");


   // Demosaic to convert sensor Bayer pattern into 24-bit RGB
   // TODO Add additional register assignments here to fully configure this core. See Demosaic IP documentation for register details.
   // Hint 1: You will need to configure 3 addition registers

   // Add assignment here
   Xil_Out32((XPAR_XV_DEMOSAIC_0_S_AXI_CTRL_BASEADDR) + 0x0028, (u32)(0x0));   // Set Bayer pattern
   Xil_Out32((XPAR_XV_DEMOSAIC_0_S_AXI_CTRL_BASEADDR) + 0x0010, (u32)(1920));    // Set image width
   Xil_Out32((XPAR_XV_DEMOSAIC_0_S_AXI_CTRL_BASEADDR) + 0x0018, (u32)(1080));   // Set image height

   // Uncomment as part of Demosaic IP setup
   Xil_Out32((XPAR_XV_DEMOSAIC_0_S_AXI_CTRL_BASEADDR) + (XV_DEMOSAIC_CTRL_ADDR_AP_CTRL), (u32)(0x81));// 0b10000001 means start and freerun mode (page 16 in PG286)
   xil_printf("Demosaic IP Configuring and Enable done\r\n"); // RGRG sensor pattern

   return 0;
}


// Enables Spread-Spectrum Clocking (SSC)
void enable_ssc(camera_config_t *config) {

	int i;

	Xuint8 iic_cdce913_ssc_on[3][2] = {
			{0x10, 0x6D}, // SSC = 011 (0.75%)
			{0x11, 0xB6}, //
			{0x12, 0xDB}  //
	};

	xil_printf( "Enabling spread-spectrum clocking (SSC)\n\r" );
	xil_printf( "\ttype=down-spread, amount=-0.75%%\n\r" );
	fmc_imageon_iic_mux( &(config->fmc_imageon), FMC_IMAGEON_I2C_SELECT_VID_CLK );

	for ( i = 0; i < 3; i++ ) {
		config->fmc_imageon.pIIC->fpIicWrite( config->fmc_imageon.pIIC, FMC_IMAGEON_VID_CLK_ADDR,
				(0x80 | iic_cdce913_ssc_on[i][0]), &(iic_cdce913_ssc_on[i][1]), 1);
	}

	return;
}



// Toggles the reset on the DCM core (clock generator)
void reset_dcms(camera_config_t *config) {

	Xuint32 value;

	// Force reset high
	config->fmc_ipmi_iic.fpGpoRead( &(config->fmc_ipmi_iic), &value );
	value = value | 0x00000004; // Force bit 2 to 1
	config->fmc_ipmi_iic.fpGpoWrite( &(config->fmc_ipmi_iic), value );
    usleep(200000);

    // Force reset low
    config->fmc_ipmi_iic.fpGpoRead( &(config->fmc_ipmi_iic), &value );
    value = value & ~0x00000004; // Force bit 2 to 0
    config->fmc_ipmi_iic.fpGpoWrite( &(config->fmc_ipmi_iic), value );
    usleep(500000);
}
