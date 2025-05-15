#include "camera_app.h"
#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

camera_config_t camera_config;
#define ROW_SIZE 1080
#define COL_SIZE 1920
#define NUM_IMAGES 100

// Function prototypes
void camera_config_init(camera_config_t *config);
void camera_loop(camera_config_t *config);
void camera_interface(camera_config_t *config, Xuint16 *rawBayer, int *StoreimageIndex, int *imageIndex, int *playMode);

// Main function. Initializes the devices and configures VDMA
int main() {
    xil_printf("Initializing camera configuration...\r\n");
    camera_config_init(&camera_config);
    xil_printf("Camera configuration initialized.\r\n");

    xil_printf("Enabling FMC-IMAGEON...\r\n");
    if (fmc_imageon_enable(&camera_config) != 0) {
        xil_printf("ERROR: fmc_imageon_enable failed. Exiting.\r\n");
        return -1;
    }
    xil_printf("FMC-IMAGEON enabled successfully.\r\n");

    xil_printf("Starting camera loop...\r\n");
    camera_loop(&camera_config);

    xil_printf("Program completed successfully.\r\n");
    return 0;
}

// Initialize the camera configuration data structure
void camera_config_init(camera_config_t *config) {
    config->uBaseAddr_IIC_FmcIpmi =  XPAR_FMC_IPMI_ID_EEPROM_0_BASEADDR;   // Device for reading HDMI board IPMI EEPROM information
    config->uBaseAddr_IIC_FmcImageon = XPAR_FMC_IMAGEON_IIC_0_BASEADDR;    // Device for configuring the HDMI board

    // Uncomment when using VITA Camera for Video input
    config->uBaseAddr_VITA_SPI = XPAR_ONSEMI_VITA_SPI_0_S00_AXI_BASEADDR;  // Device for configuring the Camera sensor
    config->uBaseAddr_VITA_CAM = XPAR_ONSEMI_VITA_CAM_0_S00_AXI_BASEADDR;  // Device for receiving Camera sensor data

    // Uncomment when using the TPG for Video input
//    config->uBaseAddr_TPG_PatternGenerator = XPAR_V_TPG_0_S_AXI_CTRL_BASEADDR; // TPG Device

    config->uDeviceId_VTC_tpg   = XPAR_V_TC_0_DEVICE_ID;                        // Video Timer Controller (VTC) ID
    config->uDeviceId_VDMA_HdmiFrameBuffer = XPAR_AXI_VDMA_0_DEVICE_ID;         // VDMA ID
    config->uBaseAddr_MEM_HdmiFrameBuffer = XPAR_DDR_MEM_BASEADDR + 0x10000000; // VDMA base address for Frame buffers
    config->uNumFrames_HdmiFrameBuffer = XPAR_AXIVDMA_0_NUM_FSTORES;            // Number of VDMA Frame buffers
    return;
}

void camera_loop(camera_config_t *config) {
    Xuint32 parkptr;
    Xuint32 vdma_S2MM_DMACR, vdma_MM2S_DMACR;

    xil_printf("Entering main SW processing loop\r\n");

    int imageIndex = 0;
    int playMode = 0;
    int StoreimageIndex = 0;
    Xuint16 *rawBayer = (Xuint16 *)malloc(1920 * 1080 * sizeof(Xuint16) * NUM_IMAGES);

    UINTPTR SWIn = 0x41210000;
    UINTPTR ButtonIn = 0x41200000;

    xil_printf("Start processing frames!\r\n");

    while((Xil_In32(SWIn) & 0x00000002) != 0x00000002) {
        xil_printf("Loop!\r\n");
        camera_interface(config, rawBayer, &StoreimageIndex, &imageIndex, &playMode);
    }

    xil_printf("Main SW processing loop complete!\r\n");
    sleep(5);

    // Uncomment when using TPG for Video input
    // fmc_imageon_disable_tpg(config);

    sleep(1);
    return;
}

void camera_interface(camera_config_t *config, Xuint16 *rawBayer, int *StoreimageIndex, int *imageIndex, int *playMode) {
    UINTPTR SWIn = 0x41210000;
    UINTPTR ButtonIn = 0x41200000;
	Xuint32 parkptr;
	Xuint32 vdma_S2MM_DMACR, vdma_MM2S_DMACR;
	int i, j;



    volatile Xuint16 *pS2MM_Mem = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_S2MM_ADDR_OFFSET + XAXIVDMA_START_ADDR_OFFSET);
    volatile Xuint16 *pMM2S_Mem = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_MM2S_ADDR_OFFSET + XAXIVDMA_START_ADDR_OFFSET + 4);

    // Capture frame on middle button press
    if ((Xil_In32(ButtonIn) & 0x00000010) == 0x00000010) {
        if (*StoreimageIndex < NUM_IMAGES - 1) {
            xil_printf("Capture\r\n");

            // Wait for frame completion
//            while (!(XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_SR_OFFSET) & 0x00000001)) {
//                // Wait for frame completion
//            }

            // Copy the current frame to the rawBayer array
            for (int i = 0; i < (1080 * 1920); i++) {
                rawBayer[i + 1920 * 1080 * (*StoreimageIndex)] = pS2MM_Mem[i];
            }


            (*StoreimageIndex)++;
        } else {
            xil_printf("No more space for images\r\n");
        }
    }

    // Playback mode
    if ((Xil_In32(SWIn) & 0x00000001)) {
        if (!(*playMode)) {
            xil_printf("Play mode\r\n");
            *playMode = 1;
        }

        // Navigate through images with left and right buttons
        if ((Xil_In32(ButtonIn) & 0x00000001) == 0x00000001) {
            if (*imageIndex < *StoreimageIndex - 1) {
                (*imageIndex)++;
            }
        }
        if ((Xil_In32(ButtonIn) & 0x00000002) == 0x00000002) {
            if (*imageIndex > 0) {
                (*imageIndex)--;
            }
        }

        // Display the current image in playback mode


    	parkptr = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_PARKPTR_OFFSET);
    	parkptr &= ~XAXIVDMA_PARKPTR_READREF_MASK;
    	parkptr &= ~XAXIVDMA_PARKPTR_WRTREF_MASK;
    	parkptr |= 0x1;
    	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_PARKPTR_OFFSET, parkptr);


    	// Grab the DMA Control Registers, and clear circular park mode.
    	vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET);
    	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);
    	vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET);
    	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);



        int count = 0;
        while(count < 1){
            for (int i = 0; i < (1080 * 1920); i++) {
                pMM2S_Mem[i] = rawBayer[i + 1920 * 1080 * (*imageIndex)];
            }
            count++;
        }
        sleep(.002);


    	vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET);
    	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);
    	vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET);
    	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);

//	    Xuint32 parkptr;
//	    Xuint32 vdma_S2MM_DMACR, vdma_MM2S_DMACR;
//	    // Grab the DMA parkptr, and update it to ensure that when parked, the S2MM side is on frame 0, and the MM2S side on frame 1
//	    parkptr = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_PARKPTR_OFFSET);
//	    parkptr &= ~XAXIVDMA_PARKPTR_READREF_MASK;
//	    parkptr &= ~XAXIVDMA_PARKPTR_WRTREF_MASK;
//	    parkptr |= 0x1;
//	    XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_PARKPTR_OFFSET, parkptr);
//	    // Grab the DMA Control Registers, and clear circular park mode.
//	    vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET + XAXIVDMA_CR_OFFSET);
//	    XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET + XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);
//	    vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET + XAXIVDMA_CR_OFFSET);
//	    XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET + XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);
//
//
//
//
//	    sleep(2);
//
//
//
//
//	    // Grab the DMA Control Registers, and re-enable circular park mode.
//	    vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET + XAXIVDMA_CR_OFFSET);
//	    XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET + XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);
//	    vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET + XAXIVDMA_CR_OFFSET);
//	    XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET + XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);










    } else {
        *playMode = 0;
    }
}


