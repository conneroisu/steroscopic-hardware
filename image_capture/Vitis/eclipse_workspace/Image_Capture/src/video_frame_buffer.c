/*****************************************************************************
 * Joseph Zambreno
 * Phillip Jones
 *
 * Department of Electrical and Computer Engineering
 * Iowa State University
 *****************************************************************************/

/*****************************************************************************
 * video_frame_buffer.c - VDMA code. Uses higher-level AXI-VDMA functions
 * (as opposed to explicitly writing register values) to configure the
 * triple framebuffer VDMA mode.
 *
 *
 * NOTES:
 * 02/04/14 by JAZ::Design created.
 *****************************************************************************/

#include "camera_app.h"
#include "xscugic.h"


#define NUMBER_OF_READ_FRAMES    XPAR_AXIVDMA_0_NUM_FSTORES
#define NUMBER_OF_WRITE_FRAMES   XPAR_AXIVDMA_0_NUM_FSTORES

// Interrupts
static XScuGic_Config* gic_config;
XScuGic int_controller;

int vfb_common_init( u16 uDeviceId, XAxiVdma *pAxiVdma )
{
   int Status;
   XAxiVdma_Config *Config;

   Config = XAxiVdma_LookupConfig( uDeviceId );
   if (!Config)
   {
      xil_printf( "No video DMA found for ID %d\n\r", uDeviceId);
      return 1;
   }

   /* Initialize DMA engine */
   Status = XAxiVdma_CfgInitialize(pAxiVdma, Config, Config->BaseAddress);
   if (Status != XST_SUCCESS)
   {
      xil_printf( "Initialization failed %d\n\r", Status);
      return 1;
   }

   int status = 0;


  // Get config
  gic_config = XScuGic_LookupConfig(XPAR_PS7_SCUGIC_0_DEVICE_ID);
  if(gic_config == NULL)
  {
   xil_printf("ERROR: Could not get GIC config!\n\r");
  }
  // Initialize
  status = XScuGic_CfgInitialize(&int_controller, gic_config, gic_config->CpuBaseAddress);
  if(status != XST_SUCCESS)
  {
   xil_printf("ERROR: Could not initialize GIC!\n\r");
  }

  // Run self test
  status = XScuGic_SelfTest(&int_controller);
  if(status != XST_SUCCESS)
  {
   xil_printf("ERROR: GIC self test failed!\n\r");
  }

XScuGic_SetPriorityTriggerType(&int_controller, XPAR_FABRIC_AXI_VDMA_0_MM2S_INTROUT_INTR, 0xA0, 0x3);
XScuGic_SetPriorityTriggerType(&int_controller, XPAR_FABRIC_AXI_VDMA_0_S2MM_INTROUT_INTR, 0xA0, 0x3);

  Xil_ExceptionInit();
  Xil_ExceptionRegisterHandler(XIL_EXCEPTION_ID_INT, (Xil_ExceptionHandler) XScuGic_InterruptHandler, &int_controller);
  Xil_ExceptionEnable();

  // Connect ISR
  status = XScuGic_Connect(&int_controller, XPAR_FABRIC_AXI_VDMA_0_MM2S_INTROUT_INTR, (Xil_InterruptHandler)XAxiVdma_ReadIntrHandler, pAxiVdma);
  if(status != XST_SUCCESS)
  {
   xil_printf("ERROR: GIC could not connect the VDMA write interrupt!\n\r");
  }

  status = XScuGic_Connect(&int_controller, XPAR_FABRIC_AXI_VDMA_0_S2MM_INTROUT_INTR, (Xil_InterruptHandler)XAxiVdma_WriteIntrHandler, pAxiVdma);
  if(status != XST_SUCCESS)
  {
   xil_printf("ERROR: GIC could not connect the VDMA read interrupt!\n\r");
  }


	 // Enable interrupt on GIC
	 XScuGic_Enable(&int_controller, XPAR_FABRIC_AXI_VDMA_0_MM2S_INTROUT_INTR);
	 XScuGic_Enable(&int_controller, XPAR_FABRIC_AXI_VDMA_0_S2MM_INTROUT_INTR);

  // Set frame counter
  XAxiVdma_FrameCounter f;
  f.ReadDelayTimerCount = 0;
  f.WriteDelayTimerCount = 0;
  f.ReadFrameCount = 1;
  f.WriteFrameCount = 1;
  status = XAxiVdma_SetFrameCounter(pAxiVdma, &f);
  if(status != XST_SUCCESS)
  {
   xil_printf("ERROR: Could not set VDMA frame counter!\n\r");
  }


   return 0;
}

int vfb_rx_init( XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pWriteCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames )
{
   int Status;
   XAxiVdma_FrameCounter FrameCfg;
   u32 uBaseAddr;
   u32 uDMACR;

   /* Setup the write channel */
   Status = vfb_rx_setup(pAxiVdma,pWriteCfg,uVideoResolution,uStorageResolution,uMemAddr,uNumFrames);
   if (Status != XST_SUCCESS) {
           xdbg_printf(XDBG_DEBUG_ERROR,
               "Write channel setup failed %d\r\n", Status);

           return 1;
   }

   // Configure RX int handler.
   XAxiVdma_SetCallBack(pAxiVdma, XAXIVDMA_HANDLER_GENERAL, video_frame_output_isr, (void*) pAxiVdma, XAXIVDMA_READ);
   XAxiVdma_SetCallBack(pAxiVdma, XAXIVDMA_HANDLER_ERROR, error_isr, (void*) pAxiVdma, XAXIVDMA_READ);
   XAxiVdma_IntrEnable(pAxiVdma, XAXIVDMA_IXR_FRMCNT_MASK, XAXIVDMA_READ);

   /* Start the DMA engine to transfer
	*/
   Status = vfb_rx_start(pAxiVdma);
   if (Status != XST_SUCCESS) {
		   return 1;
   }

	XAxiVdma_FsyncSrcSelect(pAxiVdma, XAXIVDMA_S2MM_TUSER_FSYNC, XAXIVDMA_WRITE);

}

int vfb_tx_init( XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pReadCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames )
{
   int Status;
   XAxiVdma_FrameCounter FrameCfg;
   u32 uBaseAddr;
   u32 uDMACR;

   /* Setup the read channel */
   Status = vfb_tx_setup(pAxiVdma,pReadCfg,uVideoResolution,uStorageResolution,uMemAddr,uNumFrames);
   if (Status != XST_SUCCESS) {
		   xdbg_printf(XDBG_DEBUG_ERROR,
			   "Read channel setup failed %d\n\r", Status);

		   return 1;
   }

   // Configure RX int handler.
   XAxiVdma_SetCallBack(pAxiVdma, XAXIVDMA_HANDLER_GENERAL, camera_input_isr, (void*) pAxiVdma, XAXIVDMA_WRITE);
   XAxiVdma_SetCallBack(pAxiVdma, XAXIVDMA_HANDLER_ERROR, error_isr, (void*) pAxiVdma, XAXIVDMA_WRITE);
   XAxiVdma_IntrEnable(pAxiVdma, XAXIVDMA_IXR_FRMCNT_MASK, XAXIVDMA_WRITE);

   /* Start the DMA engine to transfer
	*/
   Status = vfb_tx_start(pAxiVdma);
   if (Status != XST_SUCCESS) {
		   return 1;
   }

#if 0
	// This function returns prematurely due to (!Channel->GenLock) evaluating to false
	XAxiVdma_GenLockSourceSelect(pAxiVdma, XAXIVDMA_INTERNAL_GENLOCK, XAXIVDMA_READ);
#else
	uBaseAddr = pAxiVdma->BaseAddr;
	//xil_printf( "\t MM2S_DMACR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET )) );
	uDMACR = *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET ));
	uDMACR |= 0x00000080;
	*((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET )) = uDMACR;
	//xil_printf( "\t MM2S_DMACR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET )) );
#endif

}

int vfb_rx_setup(XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pWriteCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames )
{
	int i;
	u32 Addr;
	int Status;

	Xuint32 video_width, video_height;
	Xuint32 storage_width, storage_height, storage_stride, storage_size, storage_offset;

	// Get Video dimensions
	video_height = vres_get_height( uVideoResolution );      // in lines
	video_width  = vres_get_width ( uVideoResolution ) << 1; // in bytes

	// Get Storage dimensions
	storage_height = vres_get_height( uStorageResolution );      // in lines
	storage_width  = vres_get_width ( uStorageResolution ) << 1; // in bytes
	storage_stride = storage_width;
	storage_size   = storage_width * storage_height;
	storage_offset = ((storage_height-video_height)>>1)*storage_width + ((storage_width-video_width)>>1);

	pWriteCfg->VertSizeInput = video_height;
	pWriteCfg->HoriSizeInput = video_width;
	pWriteCfg->Stride        = storage_stride;

	pWriteCfg->FrameDelay = 0;  /* This example does not test frame delay */

	pWriteCfg->EnableCircularBuf = 1;
	pWriteCfg->EnableSync = 1;

	pWriteCfg->PointNum = 1;
	pWriteCfg->EnableFrameCounter = 0; /* Endless transfers */

	pWriteCfg->FixedFrameStoreAddr = 0; /* We are not doing parking */

	Status = XAxiVdma_DmaConfig(pAxiVdma, XAXIVDMA_WRITE, pWriteCfg);
	if (Status != XST_SUCCESS) {
			xdbg_printf(XDBG_DEBUG_ERROR,
				"Write channel config failed %d\r\n", Status);

			return XST_FAILURE;
	}

	/* Initialize buffer addresses
	 *
	 * Use physical addresses
	 */
	Addr = uMemAddr + storage_offset;
	for(i = 0; i < uNumFrames; i++)
	{
		pWriteCfg->FrameStoreStartAddr[i] = Addr;

		Addr += storage_size;
	}

	/* Set the buffer addresses for transfer in the DMA engine
	 */
	Status = XAxiVdma_DmaSetBufferAddr(pAxiVdma, XAXIVDMA_WRITE,
			pWriteCfg->FrameStoreStartAddr);
	if (Status != XST_SUCCESS) {
			xdbg_printf(XDBG_DEBUG_ERROR,
				"Write channel set buffer address failed %d\r\n", Status);

			return XST_FAILURE;
	}


	return XST_SUCCESS;
}

int vfb_tx_setup(XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pReadCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames )
{
	int i;
	u32 Addr;
	int Status;

	Xuint32 video_width, video_height;
	Xuint32 storage_width, storage_height, storage_stride, storage_size, storage_offset;

	// Get Video dimensions
	video_height = vres_get_height( uVideoResolution );      // in lines
	video_width  = vres_get_width ( uVideoResolution ) << 1; // in bytes

	// Get Storage dimensions
	storage_height = vres_get_height( uStorageResolution );      // in lines
	storage_width  = vres_get_width ( uStorageResolution ) << 1; // in bytes
	storage_stride = storage_width;
	storage_size   = storage_width * storage_height;
	storage_offset = ((storage_height-video_height)>>1)*storage_width + ((storage_width-video_width)>>1);

	pReadCfg->VertSizeInput = video_height;
	pReadCfg->HoriSizeInput = video_width;
	pReadCfg->Stride        = storage_stride;

	pReadCfg->FrameDelay = 0;  /* This example does not test frame delay */

	pReadCfg->EnableCircularBuf = 1;
	pReadCfg->EnableSync = 1;

	pReadCfg->PointNum = 1;
	pReadCfg->EnableFrameCounter = 0; /* Endless transfers */

	pReadCfg->FixedFrameStoreAddr = 0; /* We are not doing parking */

	Status = XAxiVdma_DmaConfig(pAxiVdma, XAXIVDMA_READ, pReadCfg);
	if (Status != XST_SUCCESS) {
			xdbg_printf(XDBG_DEBUG_ERROR,
				"Read channel config failed %d\n\r", Status);

			return XST_FAILURE;
	}

	/* Initialize buffer addresses
	 *
	 * These addresses are physical addresses
	 */
	Addr = uMemAddr + storage_offset;
	for(i = 0; i < uNumFrames; i++)
	{
		pReadCfg->FrameStoreStartAddr[i] = Addr;

		Addr += storage_size;
	}

	/* Set the buffer addresses for transfer in the DMA engine
	 * The buffer addresses are physical addresses
	 */
	Status = XAxiVdma_DmaSetBufferAddr(pAxiVdma, XAXIVDMA_READ,
			pReadCfg->FrameStoreStartAddr);
	if (Status != XST_SUCCESS) {
			xdbg_printf(XDBG_DEBUG_ERROR,
				"Read channel set buffer address failed %d\n\r", Status);

			return XST_FAILURE;
	}

	return XST_SUCCESS;
}

int vfb_rx_start(XAxiVdma *pAxiVdma)
{
   int Status;

   // S2MM Startup
   Status = XAxiVdma_DmaStart(pAxiVdma, XAXIVDMA_WRITE);
   if (Status != XST_SUCCESS)
   {
      xil_printf( "Start Write transfer failed %d\r\n", Status);
      return XST_FAILURE;
   }

   return XST_SUCCESS;
}

int vfb_tx_start(XAxiVdma *pAxiVdma)
{
   int Status;

   // MM2S Startup
   Status = XAxiVdma_DmaStart(pAxiVdma, XAXIVDMA_READ);
   if (Status != XST_SUCCESS)
   {
      xil_printf("Start read transfer failed %d\n\r", Status);
      return XST_FAILURE;
   }

   return XST_SUCCESS;
}

int vfb_rx_stop(XAxiVdma *pAxiVdma)
{
   // S2MM Stop
   XAxiVdma_DmaStop(pAxiVdma, XAXIVDMA_WRITE);

   return XST_SUCCESS;
}

int vfb_tx_stop(XAxiVdma *pAxiVdma)
{
   // MM2S Stop
   XAxiVdma_DmaStop(pAxiVdma, XAXIVDMA_READ);

   return XST_SUCCESS;
}


int vfb_dump_registers(XAxiVdma *pAxiVdma)
{
   u32 uBaseAddr = pAxiVdma->BaseAddr;

   // Partial Register Dump
   xil_printf( "AXI_VDMA - Partial Register Dump (uBaseAddr = 0x%08X):\n\r", uBaseAddr );
   xil_printf( "\t PARKPTR          = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_PARKPTR_OFFSET )) );
   xil_printf( "\t ----------------\n\r" );
   xil_printf( "\t S2MM_DMACR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET )) );
   xil_printf( "\t S2MM_DMASR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_RX_OFFSET+XAXIVDMA_SR_OFFSET )) );
   xil_printf( "\t S2MM_STRD_FRMDLY = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_STRD_FRMDLY_OFFSET)) );
   xil_printf( "\t S2MM_START_ADDR0 = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET+0)) );
   xil_printf( "\t S2MM_START_ADDR1 = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET+4)) );
   xil_printf( "\t S2MM_START_ADDR2 = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET+8)) );
   xil_printf( "\t S2MM_HSIZE       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_HSIZE_OFFSET)) );
   xil_printf( "\t S2MM_VSIZE       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_VSIZE_OFFSET)) );
   xil_printf( "\t ----------------\n\r" );
   xil_printf( "\t MM2S_DMACR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET )) );
   xil_printf( "\t MM2S_DMASR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_SR_OFFSET )) );
   xil_printf( "\t MM2S_STRD_FRMDLY = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_STRD_FRMDLY_OFFSET)) );
   xil_printf( "\t MM2S_START_ADDR0 = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET+0)) );
   xil_printf( "\t MM2S_START_ADDR1 = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET+4)) );
   xil_printf( "\t MM2S_START_ADDR2 = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET+8)) );
   xil_printf( "\t MM2S_HSIZE       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_HSIZE_OFFSET)) );
   xil_printf( "\t MM2S_VSIZE       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_VSIZE_OFFSET)) );
   xil_printf( "\t ----------------\n\r" );
   xil_printf( "\t S2MM_HSIZE_STATUS= 0x%08X\n\r", *((volatile int *)(uBaseAddr+0xF0 )) );
   xil_printf( "\t S2MM_VSIZE_STATUS= 0x%08X\n\r", *((volatile int *)(uBaseAddr+0xF4 )) );
   xil_printf( "\t ----------------\n\r" );

   return 0;
}


int vfb_check_errors(XAxiVdma *pAxiVdma, u8 bClearErrors )
{
   u32 uBaseAddr = pAxiVdma->BaseAddr;
   Xuint32 inErrors;
   Xuint32 outErrors;
   Xuint32 Errors;

   // Get Status of Error Flags
   //inErrors = XAxiVdma_GetStatus(pAxiVdma, XAXIVDMA_WRITE) & XAXIVDMA_SR_ERR_ALL_MASK;
   //outErrors = XAxiVdma_GetStatus(pAxiVdma, XAXIVDMA_READ) & XAXIVDMA_SR_ERR_ALL_MASK;
   inErrors  = *((volatile int *)(uBaseAddr+XAXIVDMA_RX_OFFSET+XAXIVDMA_SR_OFFSET )) & 0x0000CFF0;
   outErrors = *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_SR_OFFSET )) & 0x000046F0;
   xil_printf( "AXI_VDMA - Checking Error Flags\n\r" );

   Errors = (inErrors << 16) | (outErrors);

   if ( Errors )
   {
	   if ( inErrors & 0x00004000 )
	   {
          xil_printf( "\tS2MM_DMASR - ErrIrq\n\r" );
	   }
	   if ( inErrors & 0x00008000 )
	   {
          xil_printf( "\tS2MM_DMASR - EOLLateErr\n\r" );
	   }
	   if ( inErrors & 0x00000800 )
	   {
          xil_printf( "\tS2MM_DMASR - SOFLateErr\n\r" );
	   }
	   if ( inErrors & 0x00000400 )
	   {
          xil_printf( "\tS2MM_DMASR - SGDecErr\n\r" );
	   }
	   if ( inErrors & 0x00000200 )
	   {
          xil_printf( "\tS2MM_DMASR - SGSlvErr\n\r" );
	   }
	   if ( inErrors & 0x00000100 )
	   {
          xil_printf( "\tS2MM_DMASR - EOLEarlyErr\n\r" );
	   }
	   if ( inErrors & 0x00000080 )
	   {
          xil_printf( "\tS2MM_DMASR - SOFEarlyErr\n\r" );
	   }
	   if ( inErrors & 0x00000040 )
	   {
          xil_printf( "\tS2MM_DMASR - DMADecErr\n\r" );
	   }
	   if ( inErrors & 0x00000020 )
	   {
          xil_printf( "\tS2MM_DMASR - DMASlvErr\n\r" );
	   }
	   if ( inErrors & 0x00000010 )
	   {
          xil_printf( "\tS2MM_DMASR - DMAIntErr\n\r" );
	   }

	   if ( outErrors & 0x00004000 )
	   {
          xil_printf( "\tMM2S_DMASR - ErrIrq\n\r" );
	   }
	   if ( outErrors & 0x00000400 )
	   {
          xil_printf( "\tMM2S_DMASR - SGDecErr\n\r" );
	   }
	   if ( outErrors & 0x00000200 )
	   {
          xil_printf( "\tMM2S_DMASR - SGSlvErr\n\r" );
	   }
	   if ( outErrors & 0x00000080 )
	   {
          xil_printf( "\tMM2S_DMASR - SOFEarlyErr\n\r" );
	   }
	   if ( outErrors & 0x00000040 )
	   {
          xil_printf( "\tMM2S_DMASR - DMADecErr\n\r" );
	   }
	   if ( outErrors & 0x00000020 )
	   {
          xil_printf( "\tMM2S_DMASR - DMASlvErr\n\r" );
	   }
	   if ( outErrors & 0x00000010 )
	   {
          xil_printf( "\tMM2S_DMASR - DMAIntErr\n\r" );
	   }

	   // Clear error flags
	   xil_printf( "AXI_VDMA - Clearing Error Flags\n\r" );
	   *((volatile int *)(uBaseAddr+XAXIVDMA_RX_OFFSET+XAXIVDMA_SR_OFFSET )) = 0x0000CFF0; //XAXIVDMA_SR_ERR_ALL_MASK;
	   *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_SR_OFFSET )) = 0x000046F0; //XAXIVDMA_SR_ERR_ALL_MASK;
   }

   return Errors;
}

