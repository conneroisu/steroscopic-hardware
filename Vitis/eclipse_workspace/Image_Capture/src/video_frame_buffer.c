#include "camera_app.h"

#define NUMBER_OF_READ_FRAMES XPAR_AXIVDMA_0_NUM_FSTORES
#define NUMBER_OF_WRITE_FRAMES XPAR_AXIVDMA_0_NUM_FSTORES

int vfb_common_init(u16 uDeviceId, XAxiVdma *pAxiVdma)
{
	int Status;
	XAxiVdma_Config *Config;

	Config = XAxiVdma_LookupConfig(uDeviceId);
	if (!Config)
	{
		xil_printf("No video DMA found for ID %d\n\r", uDeviceId);
		return 1;
	}

	/* Initialize DMA engine */
	Status = XAxiVdma_CfgInitialize(pAxiVdma, Config, Config->BaseAddress);
	if (Status != XST_SUCCESS)
	{
		xil_printf("Initialization failed %d\n\r", Status);
		return 1;
	}

	int status = 0;

	// Set frame counter
	XAxiVdma_FrameCounter f;
	f.ReadDelayTimerCount = 0;
	f.WriteDelayTimerCount = 0;
	f.ReadFrameCount = 1;
	f.WriteFrameCount = 1;
	status = XAxiVdma_SetFrameCounter(pAxiVdma, &f);
	if (status != XST_SUCCESS)
	{
		xil_printf("ERROR: Could not set VDMA frame counter!\n\r");
	}

	return 0;
}

int vfb_rx_init(XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pWriteCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames)
{
	int Status;

	/* Setup the write channel */
	Status = vfb_rx_setup(pAxiVdma, pWriteCfg, uVideoResolution, uStorageResolution, uMemAddr, uNumFrames);
	if (Status != XST_SUCCESS)
	{
		xdbg_printf(XDBG_DEBUG_ERROR,
					"Write channel setup failed %d\r\n", Status);

		return 1;
	}

	/* Start the DMA engine to transfer
	 */
	Status = vfb_rx_start(pAxiVdma);
	if (Status != XST_SUCCESS)
	{
		return 1;
	}

	XAxiVdma_FsyncSrcSelect(pAxiVdma, XAXIVDMA_S2MM_TUSER_FSYNC, XAXIVDMA_WRITE);

	return 0;
}

int vfb_tx_init(XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pReadCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames)
{
	int Status;
	u32 uBaseAddr;
	u32 uDMACR;

	/* Setup the read channel */
	Status = vfb_tx_setup(pAxiVdma, pReadCfg, uVideoResolution, uStorageResolution, uMemAddr, uNumFrames);
	if (Status != XST_SUCCESS)
	{
		xdbg_printf(XDBG_DEBUG_ERROR,
					"Read channel setup failed %d\n\r", Status);

		return 1;
	}

	/* Start the DMA engine to transfer
	 */
	Status = vfb_tx_start(pAxiVdma);
	if (Status != XST_SUCCESS)
	{
		return 1;
	}

#if 0
	// This function returns prematurely due to (!Channel->GenLock) evaluating to false
	XAxiVdma_GenLockSourceSelect(pAxiVdma, XAXIVDMA_INTERNAL_GENLOCK, XAXIVDMA_READ);
#else
	uBaseAddr = pAxiVdma->BaseAddr;
	// xil_printf( "\t MM2S_DMACR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET )) );
	uDMACR = *((volatile int *)(uBaseAddr + XAXIVDMA_TX_OFFSET + XAXIVDMA_CR_OFFSET));
	uDMACR |= 0x00000080;
	*((volatile int *)(uBaseAddr + XAXIVDMA_TX_OFFSET + XAXIVDMA_CR_OFFSET)) = uDMACR;
	// xil_printf( "\t MM2S_DMACR       = 0x%08X\n\r", *((volatile int *)(uBaseAddr+XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET )) );
#endif

	return 0;
}

int vfb_rx_setup(XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pWriteCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames)
{
	int i;
	u32 Addr;
	int Status;

	Xuint32 video_width, video_height;
	Xuint32 storage_width, storage_height, storage_stride, storage_size, storage_offset;

	// Get Video dimensions
	video_height = vres_get_height(uVideoResolution);	 // in lines
	video_width = vres_get_width(uVideoResolution) << 1; // in bytes

	// Get Storage dimensions
	storage_height = vres_get_height(uStorageResolution);	 // in lines
	storage_width = vres_get_width(uStorageResolution) << 1; // in bytes
	storage_stride = storage_width;
	storage_size = storage_width * storage_height;
	storage_offset = ((storage_height - video_height) >> 1) * storage_width + ((storage_width - video_width) >> 1);

	pWriteCfg->VertSizeInput = video_height;
	pWriteCfg->HoriSizeInput = video_width;
	pWriteCfg->Stride = storage_stride;

	pWriteCfg->FrameDelay = 0; /* This example does not test frame delay */

	pWriteCfg->EnableCircularBuf = 1;
	pWriteCfg->EnableSync = 1;

	pWriteCfg->PointNum = 1;
	pWriteCfg->EnableFrameCounter = 0; /* Endless transfers */

	pWriteCfg->FixedFrameStoreAddr = 0; /* We are not doing parking */

	Status = XAxiVdma_DmaConfig(pAxiVdma, XAXIVDMA_WRITE, pWriteCfg);
	if (Status != XST_SUCCESS)
	{
		xdbg_printf(XDBG_DEBUG_ERROR,
					"Write channel config failed %d\r\n", Status);

		return XST_FAILURE;
	}

	/* Initialize buffer addresses
	 *
	 * Use physical addresses
	 */
	Addr = uMemAddr + storage_offset;
	for (i = 0; i < uNumFrames; i++)
	{
		pWriteCfg->FrameStoreStartAddr[i] = Addr;

		Addr += storage_size;
	}

	/* Set the buffer addresses for transfer in the DMA engine
	 */
	Status = XAxiVdma_DmaSetBufferAddr(pAxiVdma, XAXIVDMA_WRITE,
									   pWriteCfg->FrameStoreStartAddr);
	if (Status != XST_SUCCESS)
	{
		xdbg_printf(XDBG_DEBUG_ERROR,
					"Write channel set buffer address failed %d\r\n", Status);

		return XST_FAILURE;
	}

	return XST_SUCCESS;
}

int vfb_tx_setup(XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pReadCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames)
{
	int i;
	u32 Addr;
	int Status;

	Xuint32 video_width, video_height;
	Xuint32 storage_width, storage_height, storage_stride, storage_size, storage_offset;

	// Get Video dimensions
	video_height = vres_get_height(uVideoResolution);	 // in lines
	video_width = vres_get_width(uVideoResolution) << 1; // in bytes

	// Get Storage dimensions
	storage_height = vres_get_height(uStorageResolution);	 // in lines
	storage_width = vres_get_width(uStorageResolution) << 1; // in bytes
	storage_stride = storage_width;
	storage_size = storage_width * storage_height;
	storage_offset = ((storage_height - video_height) >> 1) * storage_width + ((storage_width - video_width) >> 1);

	pReadCfg->VertSizeInput = video_height;
	pReadCfg->HoriSizeInput = video_width;
	pReadCfg->Stride = storage_stride;

	pReadCfg->FrameDelay = 0; /* This example does not test frame delay */

	pReadCfg->EnableCircularBuf = 1;
	pReadCfg->EnableSync = 1;

	pReadCfg->PointNum = 1;
	pReadCfg->EnableFrameCounter = 0; /* Endless transfers */

	pReadCfg->FixedFrameStoreAddr = 0; /* We are not doing parking */

	Status = XAxiVdma_DmaConfig(pAxiVdma, XAXIVDMA_READ, pReadCfg);
	if (Status != XST_SUCCESS)
	{
		xdbg_printf(XDBG_DEBUG_ERROR,
					"Read channel config failed %d\n\r", Status);

		return XST_FAILURE;
	}

	/* Initialize buffer addresses
	 *
	 * These addresses are physical addresses
	 */
	Addr = uMemAddr + storage_offset;
	for (i = 0; i < uNumFrames; i++)
	{
		pReadCfg->FrameStoreStartAddr[i] = Addr;

		Addr += storage_size;
	}

	/* Set the buffer addresses for transfer in the DMA engine
	 * The buffer addresses are physical addresses
	 */
	Status = XAxiVdma_DmaSetBufferAddr(pAxiVdma, XAXIVDMA_READ,
									   pReadCfg->FrameStoreStartAddr);
	if (Status != XST_SUCCESS)
	{
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
		xil_printf("Start Write transfer failed %d\r\n", Status);
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

	return 0;
}

int vfb_check_errors(XAxiVdma *pAxiVdma, u8 bClearErrors)
{
	u32 uBaseAddr = pAxiVdma->BaseAddr;
	Xuint32 inErrors;
	Xuint32 outErrors;
	Xuint32 Errors;

	// Get Status of Error Flags
	inErrors = *((volatile int *)(uBaseAddr + XAXIVDMA_RX_OFFSET + XAXIVDMA_SR_OFFSET)) & 0x0000CFF0;
	outErrors = *((volatile int *)(uBaseAddr + XAXIVDMA_TX_OFFSET + XAXIVDMA_SR_OFFSET)) & 0x000046F0;

	Errors = (inErrors << 16) | (outErrors);

	if (Errors)
	{
		*((volatile int *)(uBaseAddr + XAXIVDMA_RX_OFFSET + XAXIVDMA_SR_OFFSET)) = 0x0000CFF0; // XAXIVDMA_SR_ERR_ALL_MASK;
		*((volatile int *)(uBaseAddr + XAXIVDMA_TX_OFFSET + XAXIVDMA_SR_OFFSET)) = 0x000046F0; // XAXIVDMA_SR_ERR_ALL_MASK;
	}

	return Errors;
}
