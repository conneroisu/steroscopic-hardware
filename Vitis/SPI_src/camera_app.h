/*****************************************************************************
 * Joseph Zambreno
 * Phillip Jones
 *
 * Department of Electrical and Computer Engineering
 * Iowa State University
 *****************************************************************************/

/*****************************************************************************
 * camera_app.h - header file for the main camera application code.
 *
 * NOTES:
 * 02/04/14 by JAZ::Design created.
 *****************************************************************************/



#ifndef __CAMERA_APP_H__
#define __CAMERA_APP_H__


#include <xparameters.h>
#include <xbasic_types.h>
#include <stdlib.h>
#include <xil_printf.h>
#include <sleep.h>
#include "fmc_iic.h"
#include "fmc_ipmi.h"
#include "fmc_imageon.h"
#include "onsemi_vita_sw.h"
#include "xv_demosaic.h"   // Uncomment when using Demosaic IP core
#include "xvprocss.h"      // Uncomment when using Video Processing Subsystem IP cores
#include "xvtc.h"
#include "xaxivdma.h"
#include "xtpg_app.h"


// Constants for library code
#define ZED_FMC_IMAGEON_GETTING_STARTED_HW
#define ADV7511_ADDR   0x72


// This structure contains the configuration context for the
// camera peripherals
struct struct_camera_config_t {

	// IP base addresses
	Xuint32 uBaseAddr_MUX_VideoSource;
	Xuint32 uBaseAddr_IIC_FmcIpmi;
	Xuint32 uBaseAddr_IIC_FmcImageon;
	Xuint32 uBaseAddr_IIC_HdmiOut;
	Xuint32 uBaseAddr_VITA_SPI;
	Xuint32 uBaseAddr_VITA_CAM;
	Xuint32 uBaseAddr_TPG_PatternGenerator;

	Xuint32 uDeviceId_VTC_ipipe;
	Xuint32 uDeviceId_VTC_tpg;

	Xuint32 uDeviceId_CFA;
	Xuint32 uDeviceId_CRES;
	Xuint32 uDeviceId_RGBYCC;

	// Frame Buffer memory addresses
	Xuint32 uDeviceId_VDMA_HdmiFrameBuffer;
	Xuint32 uBaseAddr_VDMA_HdmiFrameBuffer;
	Xuint32 uBaseAddr_MEM_HdmiFrameBuffer;
	Xuint32 uNumFrames_HdmiFrameBuffer;

	// VDMA data structures
	XAxiVdma vdma_hdmi;
	XAxiVdma_DmaSetup vdmacfg_hdmi_read;
	XAxiVdma_DmaSetup vdmacfg_hdmi_write;


	fmc_iic_t fmc_ipmi_iic;
	fmc_iic_t fmc_imageon_iic;
	fmc_imageon_t fmc_imageon;

	onsemi_vita_t onsemi_vita;
	onsemi_vita_status_t vita_status_t1;
	onsemi_vita_status_t vita_status_t2;

	XVtc vtc_ipipe;
	XVtc vtc_tpg;

	Xuint32 vita_aec;
	Xuint32 vita_again;
	Xuint32 vita_dgain;
	Xuint32 vita_exposure;

	Xuint32 bVerbose;

	// HDMI Output settings
	Xuint32 hdmio_width;
	Xuint32 hdmio_height;
	Xuint32 hdmio_resolution;
	fmc_imageon_video_timing_t hdmio_timing;
}; typedef struct struct_camera_config_t camera_config_t;


// Video resolution macros and structure
#define VIDEO_RESOLUTION_VGA       0
#define VIDEO_RESOLUTION_NTSC      1
#define VIDEO_RESOLUTION_SVGA      2
#define VIDEO_RESOLUTION_XGA       3
#define VIDEO_RESOLUTION_720P      4
#define VIDEO_RESOLUTION_SXGA      5
#define VIDEO_RESOLUTION_1080P     6
#define VIDEO_RESOLUTION_UXGA      7
#define NUM_VIDEO_RESOLUTIONS      8

struct struct_vres_timing_t {
	char *pName;
	Xuint32 VActiveVideo;
	Xuint32 VFrontPorch;
	Xuint32 VSyncWidth;
	Xuint32 VBackPorch;
	Xuint32 VSyncPolarity;
	Xuint32 HActiveVideo;
	Xuint32 HFrontPorch;
	Xuint32 HSyncWidth;
	Xuint32 HBackPorch;
	Xuint32 HSyncPolarity;
}; typedef struct struct_vres_timing_t vres_timing_t;


// Function prototypes (camera_app.c)
void camera_config_init(camera_config_t *config);
void camera_loop(camera_config_t *config);

// Function prototypes (fmc_imageon_utils.c)
int fmc_imageon_enable(camera_config_t *config);
int fmc_imageon_enable_tpg(camera_config_t *config);
int fmc_imageon_disable_tpg(camera_config_t *config);
int fmc_imageon_enable_vita(camera_config_t *config);
int fmc_imageon_enable_ipipe(camera_config_t *config);
void reset_dcms(camera_config_t *config);
void enable_ssc(camera_config_t *config);



// Function prototypes (video_resolution.c)
char * vres_get_name(Xuint32 resolutionId);
Xuint32 vres_get_width(Xuint32 resolutionId);
Xuint32 vres_get_height(Xuint32 resolutionId);
Xuint32 vres_get_timing(Xuint32 resolutionId, vres_timing_t *pTiming);
Xint32 vres_detect(Xuint32 width, Xuint32 height);

// Function prototypes (video_generator.c)
int vgen_init(XVtc *pVtc, u16 VtcDeviceID);
int vgen_config(XVtc *pVtc, int ResolutionId, int bVerbose);

// Function prototypes (video_detector.c)
int vdet_init(XVtc *pVtc, u16 VtcDeviceID);
int vdet_detect(XVtc *pVtc, int bVerbose);
int vdet_config(XVtc *pVtc, int ResolutionId, int bVerbose);

// Function prototypes (video_frame_buffer.c)
int vfb_common_init( u16 uDeviceId, XAxiVdma * InstancePtr );
int vfb_rx_init( XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pWriteCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames );
int vfb_rx_setup( XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pWriteCfg, Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames );
int vfb_rx_start( XAxiVdma *pAxiVdma );
int vfb_rx_stop ( XAxiVdma *pAxiVdma );
int vfb_tx_init( XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pReadCfg , Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames );
int vfb_tx_setup( XAxiVdma *pAxiVdma, XAxiVdma_DmaSetup *pReadCfg , Xuint32 uVideoResolution, Xuint32 uStorageResolution, Xuint32 uMemAddr, Xuint32 uNumFrames );
int vfb_tx_start( XAxiVdma *pAxiVdma );
int vfb_tx_stop ( XAxiVdma *pAxiVdma );
int vfb_dump_registers( XAxiVdma *pAxiVdma);
int vfb_check_errors( XAxiVdma *pAxiVdma, u8 bClearErrors );

//New Stuff
void rgb_to_ycbcr(Xuint8 r, Xuint8 g, Xuint8 b, Xuint8 *y, Xuint8 *cb, Xuint8 *cr);



#endif // __CAMERA_APP_H__
