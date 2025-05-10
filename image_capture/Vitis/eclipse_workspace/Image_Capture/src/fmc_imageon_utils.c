
#include "camera_app.h"
#include "xil_types.h"
#include "xil_cache.h"

// Main FMC-IMAGEON initialization function. Add your code here.
int fmc_imageon_enable( camera_config_t *config )
{
   int ret;

   config->bVerbose = 0;
   config->vita_aec = 0;       // off
   config->vita_again = 0;     // 1.0
   config->vita_dgain = 128;   // 1.0
   config->vita_exposure = 90; // 90% of frame period


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


   ret = fmc_iic_axi_init(&(config->fmc_imageon_iic), "FMC-IMAGEON I2C Controller", config->uBaseAddr_IIC_FmcImageon);
   if (!ret) {
      xil_printf( "ERROR: Failed to open FMC-IIC driver\n\r");
      exit(1);
   }

   fmc_imageon_init(&(config->fmc_imageon), "FMC-IMAGEON", &(config->fmc_imageon_iic));
   fmc_imageon_vclk_init( &(config->fmc_imageon) );
   fmc_imageon_vclk_config( &(config->fmc_imageon), FMC_IMAGEON_VCLK_FREQ_148_500_000);

   reset_dcms(config);

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

   config->hdmio_resolution = vres_detect(config->hdmio_width, config->hdmio_height);

   vgen_init( &(config->vtc_tpg), config->uDeviceId_VTC_tpg);
   vgen_config( &(config->vtc_tpg ), config->hdmio_resolution, 1);


   // FMC-IMAGEON HDMI Output Initialization
   ret = fmc_imageon_hdmio_init(&(config->fmc_imageon), 1, &(config->hdmio_timing), 0);
   if (!ret) {
      xil_printf("ERROR : Failed to init FMC-IMAGEON HDMI Output Interface\n\r");
      exit(0);
   }


   // FMC-IMAGEON VITA Camera Receiver Initialization
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
   sleep(1);


   // Initialize Input Side of AXI VDMA
   vfb_rx_init(
      &(config->vdma_hdmi),                   // pAxiVdma
      &(config->vdmacfg_hdmi_write),          // pWriteCfg
      config->hdmio_resolution,               // uVideoResolution
      config->hdmio_resolution,               // uStorageResolution
      config->uBaseAddr_MEM_HdmiFrameBuffer,  // uMemAddr
      config->uNumFrames_HdmiFrameBuffer      // uNumFrames
      );


   int vita_enabled_error = 0;
   do {
	   vita_enabled_error = fmc_imageon_enable_vita(config);
   } while(vita_enabled_error != 0);


   vfb_dump_registers( &(config->vdma_hdmi) );
   if ( vfb_check_errors( &(config->vdma_hdmi), 1/*clear errors, if any*/ ) ) {
      vfb_dump_registers( &(config->vdma_hdmi) );
   }

   return 0;
}



int fmc_imageon_enable_tpg( camera_config_t *config ) {
   // Define convenient volatile pointers for accessing TPG registers
   volatile uint32_t *TPG_CR       = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0);    // TPG Control
   volatile uint32_t *TPG_Act_H    = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x10); // Active Height
   volatile uint32_t *TPG_Act_W    = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x18); // Active Width
   volatile uint32_t *TPG_BGP      = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x20); // Background Pattern
   volatile uint32_t *TPG_FGP      = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x28); // Foreground Pattern
   volatile uint32_t *TPG_MS       = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x38); // Motion Speed
   volatile uint32_t *TPG_CF       = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x40); // TPG Color Format
   volatile uint32_t *TPG_BOX_SIZE = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x78);
   volatile uint32_t *TPG_BOX_COLOR_Y = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x80);
   volatile uint32_t *TPG_BOX_COLOR_U = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x88);
   volatile uint32_t *TPG_BOX_COLOR_V = (volatile uint32_t*) (config->uBaseAddr_TPG_PatternGenerator + 0x90);

   // Direct Memory Mapped access of TPG configuration registers
   // See TPG data sheet for configuring the TPG for other features
   TPG_Act_H[0]  = 0x438; // Active Height
   TPG_Act_W[0]  = 0x780; // Active Width
   TPG_BGP[0]    = 0x09;  // Background Pattern
   TPG_FGP[0]    = 0x01;  // Foreground Pattern
   TPG_MS[0]     = 0x04;  // Motion Speed
   TPG_BOX_SIZE[0] = 100;
   TPG_BOX_COLOR_Y[0] = 167;
   TPG_BOX_COLOR_U[0] = 120;
   TPG_BOX_COLOR_V[0] = 8;
   TPG_CF[0]     = 0x02;  // TPG Color Format
   TPG_CR[0]     = 0x81;  // TPG Control
   return 0;
}


// Disable Test Pattern Generator (TPG)
int fmc_imageon_disable_tpg( camera_config_t *config ) {

	// Define convenient volatile pointers for accessing TPG registers
   volatile unsigned int *TPG_CR  = config->uBaseAddr_TPG_PatternGenerator + 0;    // TPG Control

   TPG_CR[0]     = 0x00;  // TPG Control Register (Disable TPG)

   return 0;
}



int fmc_imageon_enable_vita( camera_config_t *config ) {
   int ret;

   // VITA-2000 Initialization
   ret = onsemi_vita_sensor_initialize( &(config->onsemi_vita), SENSOR_INIT_ENABLE, config->bVerbose );
   if (ret == 0) {
      xil_printf("VITA sensor failed to initialize ...\n\r");
      return -1;
   }

   onsemi_vita_sensor_initialize( &(config->onsemi_vita), SENSOR_INIT_STREAMON, config->bVerbose);
   sleep(1);

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

   if ((vita_width != 1920) || (vita_height != 1080) || (vita_rate == 0)) {
	   return 1;
   }

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
