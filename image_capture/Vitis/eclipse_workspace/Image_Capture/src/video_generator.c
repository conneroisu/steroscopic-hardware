
#include "camera_app.h"

/*****************************************************************************/
/**
*
* This function sets up the Video Timing Controller Signal configuration.
*
* @param	None.
*
* @return	None.
*
* @note		None.
*
****************************************************************************/
static void SignalSetup( XVtc *pVtc, Xuint32 ResolutionId, XVtc_Signal *SignalCfgPtr )
{
	vres_timing_t VideoTiming;

	int HFrontPorch;
	int HSyncWidth;
	int HBackPorch;
	int VFrontPorch;
	int VSyncWidth;
	int VBackPorch;
	int LineWidth;
	int FrameHeight;

	vres_get_timing(ResolutionId, &VideoTiming);

	HFrontPorch = VideoTiming.HFrontPorch;
	HSyncWidth  = VideoTiming.HSyncWidth;
	HBackPorch  = VideoTiming.HBackPorch;
	VFrontPorch = VideoTiming.VFrontPorch;
	VSyncWidth  = VideoTiming.VSyncWidth;
	VBackPorch  = VideoTiming.VBackPorch;
	LineWidth   = VideoTiming.HActiveVideo;
	FrameHeight = VideoTiming.VActiveVideo;

	/* Clear the VTC Signal config structure */

	memset((void *)SignalCfgPtr, 0, sizeof(XVtc_Signal));

	/* Populate the VTC Signal config structure. Ignore the Field 1 */

	SignalCfgPtr->HFrontPorchStart = LineWidth;
	SignalCfgPtr->HTotal = HFrontPorch + HSyncWidth + HBackPorch
				+ LineWidth;
	SignalCfgPtr->HBackPorchStart = LineWidth + HFrontPorch + HSyncWidth;
	SignalCfgPtr->HSyncStart = LineWidth + HFrontPorch;
	SignalCfgPtr->HActiveStart = 0;

	SignalCfgPtr->V0FrontPorchStart = FrameHeight;
	SignalCfgPtr->V0Total = VFrontPorch + VSyncWidth + VBackPorch
				+ FrameHeight;
	SignalCfgPtr->V0BackPorchStart = FrameHeight + VFrontPorch + VSyncWidth;
	SignalCfgPtr->V0SyncStart = FrameHeight + VFrontPorch;
	SignalCfgPtr->V0ChromaStart = 0;
	SignalCfgPtr->V0ActiveStart = 0;

	 return;
}

/*****************************************************************************/
/**
*
* vgen_init
* - initializes the VTC detector
*
* @param	VtcDeviceID is the device ID of the Video Timing Controller core.
*           pVtc is a pointer to a VTC instance

*
* @return	0 if all tests pass, 1 otherwise.
*
* @note		None.
*
******************************************************************************/
int vgen_init(XVtc *pVtc, u16 VtcDeviceID)
{
	int Status;
	XVtc_Config *VtcCfgPtr;

	Xuint32 Width;
	Xuint32 Height;
	int ResolutionId;

	/* Look for the device configuration info for the Video Timing
	 * Controller.
	 */
	VtcCfgPtr = XVtc_LookupConfig( VtcDeviceID );
	if (VtcCfgPtr == NULL) {
		return 1;
	}

	/* Initialize the Video Timing Controller instance */

	Status = XVtc_CfgInitialize(pVtc, VtcCfgPtr,
		VtcCfgPtr->BaseAddress);
	if (Status != XST_SUCCESS) {
		return 1;
	}

	XVtc_DisableSync(pVtc);

	sleep(1);

	/* Enable the generator module */

	// phjones update to 1 arg.  XVtc_Enable(pVtc, XVTC_EN_GENERATOR);
	XVtc_EnableGenerator(pVtc);


	//	XVtc_DisableSync(pVtc);

	return 0;
}


/*****************************************************************************/
/**
*
* vgen_config
* - configures the generator to generate missing syncs
*
* @param	pVtc is a pointer to an initialized VTC instance
*           ResolutionId identified a video resolution
*           vVerbose = 0 no verbose, 1 minimal verbose, 2 most verbose
*
* @return	0 if all tests pass, 1 otherwise.
*
* @note		None.
*
******************************************************************************/
int vgen_config(XVtc *pVtc, int ResolutionId, int bVerbose)
{
	int Status;

	XVtc_Signal Signal;		/* VTC Signal configuration */
	XVtc_Polarity Polarity;		/* Polarity configuration */
	XVtc_HoriOffsets HoriOffsets;  /* Horizontal offsets configuration */
	XVtc_SourceSelect SourceSelect;	/* Source Selection configuration */

	sleep(5);

    if ( bVerbose )
    {
		xil_printf( "\tVideo Resolution = %s\n\r", vres_get_name(ResolutionId) );
	}

    /* Set up Polarity of all outputs */

	memset((void *)&Polarity, 0, sizeof(Polarity));
	Polarity.ActiveChromaPol = 1;
	Polarity.ActiveVideoPol = 1;
	Polarity.FieldIdPol = 0;
	Polarity.VBlankPol = 1;
	Polarity.VSyncPol = 1;
	Polarity.HBlankPol = 1;
	Polarity.HSyncPol = 1;

	XVtc_SetPolarity(pVtc, &Polarity);

	/* Set up Generator */

	memset((void *)&HoriOffsets, 0, sizeof(HoriOffsets));
	HoriOffsets.V0BlankHoriEnd = 1920;
	HoriOffsets.V0BlankHoriStart = 1920;
	HoriOffsets.V0SyncHoriEnd = 1920;
	HoriOffsets.V0SyncHoriStart = 1920;

	XVtc_SetGeneratorHoriOffset(pVtc, &HoriOffsets);

	SignalSetup(pVtc,ResolutionId, &Signal);


	XVtc_SetGenerator(pVtc, &Signal);

	/* Set up source select */

	memset((void *)&SourceSelect, 0, sizeof(SourceSelect));
	SourceSelect.VChromaSrc = 0;
	SourceSelect.VActiveSrc = 1;
	SourceSelect.VBackPorchSrc = 1;
	SourceSelect.VSyncSrc = 1;
	SourceSelect.VFrontPorchSrc = 1;
	SourceSelect.VTotalSrc = 1;
	SourceSelect.HActiveSrc = 1;
	SourceSelect.HBackPorchSrc = 1;
	SourceSelect.HSyncSrc = 1;
	SourceSelect.HFrontPorchSrc = 1;
	SourceSelect.HTotalSrc = 1;

	XVtc_SetSource(pVtc, &SourceSelect);


	/* Return success */

	return 0;
}
