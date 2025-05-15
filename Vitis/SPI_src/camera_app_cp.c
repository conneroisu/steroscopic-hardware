//#include "camera_app.h"
//#include <stdint.h>
//#include <stdbool.h>
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//
//camera_config_t camera_config;
//#define ROW_SIZE 1080
//#define COL_SIZE 1920
//
//// Main function. Initializes the devices and configures VDMA
//int main() {
//    xil_printf("Initializing camera configuration...\r\n");
//    camera_config_init(&camera_config);
//    xil_printf("Camera configuration initialized.\r\n");
//
//    xil_printf("Enabling FMC-IMAGEON...\r\n");
//    if (fmc_imageon_enable(&camera_config) != 0) {
//        xil_printf("ERROR: fmc_imageon_enable failed. Exiting.\r\n");
//        return -1;
//    }
//    xil_printf("FMC-IMAGEON enabled successfully.\r\n");
//
//    xil_printf("Starting camera loop...\r\n");
//    camera_loop(&camera_config);
//
//    xil_printf("Program completed successfully.\r\n");
//    return 0;
//}
//
//void demosaicing(uint8_t bw[ROW_SIZE][COL_SIZE], uint32_t color[ROW_SIZE][COL_SIZE], int rowSize, int colSize) {
//    for (int r = 1; r < rowSize - 1; r++) {
//        for (int c = 1; c < colSize - 1; c++) {
//            bool isBlueRow = (r % 2 != 0);
//            bool isGreenPixel = isBlueRow ? (c % 2 == 0) : (c % 2 != 0);
//            color[r][c] = 0;
//
//            if (isBlueRow && isGreenPixel) {
//                color[r][c] |= ((bw[r + 1][c] + bw[r - 1][c]) / 2) << 16;
//                color[r][c] |= ((bw[r][c] + bw[r - 1][c - 1] + bw[r + 1][c + 1] + bw[r + 1][c - 1] + bw[r - 1][c + 1]) / 5) << 8;
//                color[r][c] |= ((bw[r][c - 1] + bw[r][c + 1]) / 2);
//            } else if (isBlueRow && !isGreenPixel) {
//                color[r][c] |= ((bw[r - 1][c - 1] + bw[r + 1][c + 1] + bw[r + 1][c - 1] + bw[r - 1][c + 1]) / 4) << 16;
//                color[r][c] |= ((bw[r][c - 1] + bw[r][c + 1] + bw[r - 1][c] + bw[r + 1][c]) / 4) << 8;
//                color[r][c] |= bw[r][c];
//            } else if (!isBlueRow && isGreenPixel) {
//                color[r][c] |= ((bw[r][c + 1] + bw[r][c - 1]) / 2) << 16;
//                color[r][c] |= ((bw[r][c] + bw[r - 1][c - 1] + bw[r + 1][c + 1] + bw[r + 1][c - 1] + bw[r - 1][c + 1]) / 5) << 8;
//                color[r][c] |= ((bw[r - 1][c] + bw[r + 1][c]) / 2);
//            } else {
//                color[r][c] |= bw[r][c] << 16;
//                color[r][c] |= ((bw[r][c - 1] + bw[r][c + 1] + bw[r - 1][c] + bw[r + 1][c]) / 4) << 8;
//                color[r][c] |= ((bw[r - 1][c - 1] + bw[r + 1][c + 1] + bw[r + 1][c - 1] + bw[r - 1][c + 1]) / 4);
//            }
//        }
//    }
//}
//
//void rgbToYCbCr444(uint32_t data[ROW_SIZE][COL_SIZE], int rowSize, int colSize) {
//    for (int r = 0; r < rowSize; r++) {
//        for (int c = 0; c < colSize; c++) {
//            uint8_t red = (data[r][c] >> 16) & 0xFF;
//            uint8_t green = (data[r][c] >> 8) & 0xFF;
//            uint8_t blue = data[r][c] & 0xFF;
//            uint8_t Y = (uint8_t)((0.183 * red) + (0.614 * green) + (0.062 * blue) + 16);
//            uint8_t Cb = (uint8_t)((-0.101 * red) + (-0.338 * green) + (0.439 * blue) + 128);
//            uint8_t Cr = (uint8_t)((0.439 * red) + (-0.399 * green) + (-0.040 * blue) + 128);
//            data[r][c] = (Y << 16) | (Cb << 8) | Cr;
//        }
//    }
//}
//
//void YCbCr444to422(uint32_t data[ROW_SIZE][COL_SIZE], int rowSize, int colSize) {
//    for (int r = 0; r < rowSize; r++) {
//        for (int c = 0; c < colSize; c++) {
//            uint8_t Y = (data[r][c] >> 16) & 0xFF;
//            uint8_t Cb = (data[r][c] >> 8) & 0xFF;
//            uint8_t Cr = data[r][c] & 0xFF;
//            data[r][c] = (Cr << 24) | (Y << 16) | (Cb << 8) | Y;
//        }
//    }
//}
//
//// Initialize the camera configuration data structure
//void camera_config_init(camera_config_t *config) {
//    config->uBaseAddr_IIC_FmcIpmi =  XPAR_FMC_IPMI_ID_EEPROM_0_BASEADDR;   // Device for reading HDMI board IPMI EEPROM information
//    config->uBaseAddr_IIC_FmcImageon = XPAR_FMC_IMAGEON_IIC_0_BASEADDR;    // Device for configuring the HDMI board
//
//    // Uncomment when using VITA Camera for Video input
//    config->uBaseAddr_VITA_SPI = XPAR_ONSEMI_VITA_SPI_0_S00_AXI_BASEADDR;  // Device for configuring the Camera sensor
//    config->uBaseAddr_VITA_CAM = XPAR_ONSEMI_VITA_CAM_0_S00_AXI_BASEADDR;  // Device for receiving Camera sensor data
//
//
//    // Uncomment when using the TPG for Video input
////    config->uBaseAddr_TPG_PatternGenerator = XPAR_V_TPG_0_S_AXI_CTRL_BASEADDR; // TPG Device
//
//    config->uDeviceId_VTC_tpg   = XPAR_V_TC_0_DEVICE_ID;                        // Video Timer Controller (VTC) ID
//    config->uDeviceId_VDMA_HdmiFrameBuffer = XPAR_AXI_VDMA_0_DEVICE_ID;         // VDMA ID
//    config->uBaseAddr_MEM_HdmiFrameBuffer = XPAR_DDR_MEM_BASEADDR + 0x10000000; // VDMA base address for Frame buffers
//    config->uNumFrames_HdmiFrameBuffer = XPAR_AXIVDMA_0_NUM_FSTORES;            // NUmber of VDMA Frame buffers
//    return;
//}
//
//void camera_loop(camera_config_t *config) {
//	Xuint32 parkptr;
//	Xuint32 vdma_S2MM_DMACR, vdma_MM2S_DMACR;
////	int i, j;
//
//	xil_printf("Entering main SW processing loop\r\n");
//
//
//	// Grab the DMA parkptr, and update it to ensure that when parked, the S2MM side is on frame 0, and the MM2S side on frame 1
//	parkptr = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_PARKPTR_OFFSET);
//	parkptr &= ~XAXIVDMA_PARKPTR_READREF_MASK;
//	parkptr &= ~XAXIVDMA_PARKPTR_WRTREF_MASK;
//	parkptr |= 0x1;
//	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_PARKPTR_OFFSET, parkptr);
//
//
//	// Grab the DMA Control Registers, and clear circular park mode.
//	vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET);
//	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);
//	vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET);
//	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);
//
//
//	// Pointers to the S2MM memory frame and M2SS memory frame
//	volatile Xuint16 *pS2MM_Mem = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET);
//	volatile Xuint16 *pMM2S_Mem = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET+4);
//
////	uint8_t cfa[ROW_SIZE][COL_SIZE];
////    Xuint16 fullColor[ROW_SIZE][COL_SIZE];
//
//	uint32_t tempValue_R;
//	uint32_t tempValue_G;
//	uint32_t tempValue_B;
//
//    xil_printf("Start processing frames!\r\n");
//    for (int j = 0; j < 1000; j++) {
//    	for (int i = 0; i < (1080*1920); i+=2) {
//
//    		tempValue_R = 0x00000000;
//			tempValue_G = 0x00000000;
//			tempValue_B = 0x00000000;
//
//    		int r = i / 1920;
//    		int c = i % 1920;
//
//    		/// bggr
//    		// Red pixels
//			if ((r % 2 == 1) && (c % 2 == 1)) {
//				// Average blue neighbors (Diagonal)
//				tempValue_B += ((pS2MM_Mem[(r-1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r-1) * COL_SIZE + (c+1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c+1)])) / 4;
//				// Average green neighbors (Cardinal)
//				tempValue_G += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c]) + (pS2MM_Mem[r * COL_SIZE + (c - 1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 4;
//				// red stays the same
//				tempValue_R += (pS2MM_Mem[r * COL_SIZE + c]);
//
//			// Blue pixels
//			} else if ((r % 2 == 0) && (c % 2 == 0)) {
//				// blue stays the same
//				tempValue_B += (pS2MM_Mem[r * COL_SIZE + c]);
//				// Average green neighbors (Cardinal)
//				tempValue_G += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c]) + (pS2MM_Mem[r * COL_SIZE + (c - 1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 4;
//				// Average red neighbors (Diagonal)
//				tempValue_R += ((pS2MM_Mem[(r-1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r-1) * COL_SIZE + (c+1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c+1)])) / 4;
//			} else {
//					// Green stays same
//					tempValue_G += (pS2MM_Mem[r * COL_SIZE + c]);
//
//					// Blue col: odd, even
//					if ((r % 2 == 0) && (c % 2 == 1)) {
//						// Average red neighbors (left/right)
//						tempValue_R += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c])) / 2;
//						// Average blue neighbors (up/down)
//						tempValue_B += ((pS2MM_Mem[r * COL_SIZE + (c-1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 2;
//					} else {
//						// Average blue neighbors (left/right)
//						tempValue_B += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c])) / 2;
//						// Average red neighbors (up/down)
//						tempValue_R += ((pS2MM_Mem[r * COL_SIZE + (c-1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 2;
//					}
//			}
//
//			uint8_t Y = ((0.183 * tempValue_R) + (0.614 * tempValue_G) + (0.062 * tempValue_B) + 16);
//			uint8_t Cr = ((-0.101 * tempValue_R) + (-0.338 * tempValue_G) + (0.439 * tempValue_B) + 128);
//			uint8_t Cb = ((0.439 * tempValue_R) + (-0.399 * tempValue_G) + (-0.040 * tempValue_B) + 128);
//			// Set from YUV values.
//			pMM2S_Mem[i] = (Cb << 8) | (Y & 0xFF);
//			pMM2S_Mem[i + 1] = (Cr << 8) | ((Y & 0xFF00) >> 8);
//    	}
//    }
//
//
//    for (int j = 0; j < 1000; j++) {
//		xil_printf("Fuck");
//        // 1. Load CFA data from DMA memory (Bayer pattern)
////        for (int r = 0; r < ROW_SIZE-1; r++) {
////            for (int c = 0; c < COL_SIZE-1; c++) {
////        		xil_printf("Shit");
////                cfa[r][c] = 200;
////        		xil_printf("Damn");
////                		//pS2MM_Mem[r * COL_SIZE + c] & 0xFF;
//////                cfa[r][c] = pS2MM_Mem[r * COL_SIZE + c];
////            }
////        }
//		xil_printf("FuckYou");
//
////         2. Perform demosaicing (Bayer to RGB conversion)
////		        demosaicing(pS2MM_Mem, fullColor, ROW_SIZE, COL_SIZE);
//		        for (int r = 0; r < ROW_SIZE - 1; r++) {
//		            for (int c = 0; c < COL_SIZE - 1; c++) {
//						tempValue_R = 0x00000000;
//						tempValue_G = 0x00000000;
//						tempValue_B = 0x00000000;
////						rgb2YCbCr = 0x00000000;
//
//
////		            	if(r < 2 || r > (ROW_SIZE - 2) || c < 2 || c > (COL_SIZE-2)){
////		            		pMM2S_Mem[r *COL_SIZE + c] = 0x00000000;
////		            	}
////		            	else{
//
//
//
////			                bool isBlueRow = (r % 2 != 0);
////			                bool isGreenPixel = isBlueRow ? (c % 2 == 0) : (c % 2 != 0);
////
////
////
////
////			                if (isBlueRow && isGreenPixel) {
////			                	tempValue |= (((pS2MM_Mem[(r+1) * COL_SIZE + c] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c] & 0xFF)) / 2) << 16;
////			                	tempValue |= (((pS2MM_Mem[r * COL_SIZE + c] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c + 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c + 1] & 0xFF)) / 5) << 8;
////			                	tempValue |= (((pS2MM_Mem[r * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[r * COL_SIZE + c + 1] & 0xFF)) / 2);
////			                } else if (isBlueRow && !isGreenPixel) {
////			                	tempValue |= (((pS2MM_Mem[(r-1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c + 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c + 1] & 0xFF)) / 4) << 16;
////			                	tempValue |= (((pS2MM_Mem[r * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[r * COL_SIZE + c + 1] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c] & 0xFF)) / 4) << 8;
////			                	tempValue |= (pS2MM_Mem[r * COL_SIZE + c] & 0xFF);
////			                } else if (!isBlueRow && isGreenPixel) {
////			                	tempValue |= (((pS2MM_Mem[r * COL_SIZE + c + 1] & 0xFF) + (pS2MM_Mem[r * COL_SIZE + c - 1] & 0xFF)) / 2) << 16;
////			                	tempValue |= (((pS2MM_Mem[r * COL_SIZE + c] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c + 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c + 1] & 0xFF)) / 5) << 8;
////			                	tempValue |= (((pS2MM_Mem[(r-1) * COL_SIZE + c] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c] & 0xFF)) / 2);
////			                } else {
////			                	tempValue |= (pS2MM_Mem[r * COL_SIZE + c] & 0xFF) << 16;
////			                	tempValue |= (((pS2MM_Mem[r * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[r * COL_SIZE + c + 1] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c] & 0xFF)) / 4) << 8;
////			                	tempValue |= (((pS2MM_Mem[(r-1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c + 1] & 0xFF) + (pS2MM_Mem[(r+1) * COL_SIZE + c - 1] & 0xFF) + (pS2MM_Mem[(r-1) * COL_SIZE + c + 1] & 0xFF)) / 4);
////			                }
//
//
//							//// rggb
//							// Green Pixels eve, odd || odd, even
//	//		                if (((r % 2 == 0) && (c % 2 == 1)) || ((r % 2 == 0) && (c % 2 == 0)) ) {
//	//		                	// Average red neighbors in 3x3
//	//		                	tempValue |= ((pS2MM_Mem[(r-1) * COL_SIZE + c] << 16) + (pS2MM_Mem[(r+1) * COL_SIZE + c] << 16) + (pS2MM_Mem[r * COL_SIZE + (c - 1)] << 16) + (pS2MM_Mem[r * COL_SIZE + (c+1)] << 16)) / 4;
//	//							// Green stays same
//	//		                	tempValue |= (pS2MM_Mem[r * COL_SIZE + c] << 8);
//	//							// Average blue neighbors in 3x3
//	//							tempValue |= ((pS2MM_Mem[(r-1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r-1) * COL_SIZE + (c+1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c+1)])) / 4;
//	//		                } else if ((r % 2 == 1) && (c % 2 == 1)) {
//	//		                	// Red stays the same
//	//		                	tempValue |= (pS2MM_Mem[r * COL_SIZE + c] << 16);
//	//							// Average green neighbors
//	//		                	tempValue |= ((pS2MM_Mem[(r-1) * COL_SIZE + c] << 8) + (pS2MM_Mem[(r+1) * COL_SIZE + c] << 8) + (pS2MM_Mem[r * COL_SIZE + (c - 1)] << 8) + (pS2MM_Mem[r * COL_SIZE + (c+1)] << 8)) / 4;
//	//							// Average blue neighbors in 3x3
//	//							tempValue |= ((pS2MM_Mem[(r-1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r-1) * COL_SIZE + (c+1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c+1)])) / 4;
//	//		                } else {
//	//		                	// Average red neighbors in 3x3
//	//		                	tempValue |= ((pS2MM_Mem[(r-1) * COL_SIZE + (c-1)] << 16) + (pS2MM_Mem[(r-1) * COL_SIZE + (c+1)] << 16) + (pS2MM_Mem[(r+1) * COL_SIZE + (c-1)] << 16) + (pS2MM_Mem[(r+1) * COL_SIZE + (c+1)] << 16)) / 4;
//	//							// Average green neighbors
//	//		                	tempValue |= ((pS2MM_Mem[(r-1) * COL_SIZE + c] << 8) + (pS2MM_Mem[(r+1) * COL_SIZE + c] << 8) + (pS2MM_Mem[r * COL_SIZE + (c - 1)] << 8) + (pS2MM_Mem[r * COL_SIZE + (c+1)] << 8)) / 4;
//	//							// Average blue neighbors in 3x3
//	//							tempValue |= (pS2MM_Mem[r * COL_SIZE + c]);
//	//		                }
//
//							///// bggr
//	//		                // Red pixels
//							if ((r % 2 == 1) && (c % 2 == 1)) {
//								// Average blue neighbors (Diagonal)
//								tempValue_B += ((pS2MM_Mem[(r-1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r-1) * COL_SIZE + (c+1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c+1)])) / 4;
//								// Average green neighbors (Cardinal)
//								tempValue_G += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c]) + (pS2MM_Mem[r * COL_SIZE + (c - 1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 4;
//								// red stays the same
//								tempValue_R += (pS2MM_Mem[r * COL_SIZE + c]);
//
//							// Blue pixels
//							} else if ((r % 2 == 0) && (c % 2 == 0)) {
//								// blue stays the same
//								tempValue_B += (pS2MM_Mem[r * COL_SIZE + c]);
//								// Average green neighbors (Cardinal)
//								tempValue_G += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c]) + (pS2MM_Mem[r * COL_SIZE + (c - 1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 4;
//								// Average red neighbors (Diagonal)
//								tempValue_R += ((pS2MM_Mem[(r-1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r-1) * COL_SIZE + (c+1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c-1)]) + (pS2MM_Mem[(r+1) * COL_SIZE + (c+1)])) / 4;
//							} else {
//									// Green stays same
//									tempValue_G += (pS2MM_Mem[r * COL_SIZE + c]);
//
//									// Blue col: odd, even
//									if ((r % 2 == 1) && (c % 2 == 0)) {
//										// Average red neighbors (left/right)
//										tempValue_R += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c])) / 2;
//										// Average blue neighbors (up/down)
//										tempValue_B += ((pS2MM_Mem[r * COL_SIZE + (c-1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 2;
//									} else {
//										// Average blue neighbors (left/right)
//										tempValue_B += ((pS2MM_Mem[(r-1) * COL_SIZE + c]) + (pS2MM_Mem[(r+1) * COL_SIZE + c])) / 2;
//										// Average red neighbors (up/down)
//										tempValue_R += ((pS2MM_Mem[r * COL_SIZE + (c-1)]) + (pS2MM_Mem[r * COL_SIZE + (c+1)])) / 2;
//									}
//							}
////							uint8_t red = (tempValue >> 16);
////							uint8_t green = (tempValue >> 8);
////							uint8_t blue = tempValue;
//							uint8_t Y = ((0.183 * tempValue_R) + (0.614 * tempValue_G) + (0.062 * tempValue_B) + 16);
//							uint8_t Cr = ((-0.101 * tempValue_R) + (-0.338 * tempValue_G) + (0.439 * tempValue_B) + 128);
//							uint8_t Cb = ((0.439 * tempValue_R) + (-0.399 * tempValue_G) + (-0.040 * tempValue_B) + 128);
////							pMM2S_Mem[r *COL_SIZE + c] = (Y << 16) | (Cr << 8) | Cb;
//
////							Y = (tempValue >> 16);
////							Cb = (tempValue >> 8);
////							Cr = (tempValue);
////							pMM2S_Mem[r *COL_SIZE + c] = (Cr << 24) | (Y << 16) | (Cb << 8) | Y;
//							// Set from YUV values.
//							pMM2S_Mem[r * COL_SIZE + c] = (Cb << 8) | (Y & 0xFF);
//							pMM2S_Mem[(r * COL_SIZE + c) + 1] = (Cr << 8) | ((Y & 0xFF00) >> 8);
//
//		            	}
//
//		            }
//		        }
//
//
//
////				for (int r = 0; r < ROW_SIZE; r++) {
////					for (int c = 0; c < COL_SIZE; c++) {
////						uint8_t red = (pMM2S_Mem[r *COL_SIZE + c] >> 16) & 0xFF;
////						uint8_t green = (pMM2S_Mem[r *COL_SIZE + c] >> 8) & 0xFF;
////						uint8_t blue = pMM2S_Mem[r *COL_SIZE + c] & 0xFF;
////						uint8_t Y = (uint8_t)((0.183 * red) + (0.614 * green) + (0.062 * blue) + 16);
////						uint8_t Cb = (uint8_t)((-0.101 * red) + (-0.338 * green) + (0.439 * blue) + 128);
////						uint8_t Cr = (uint8_t)((0.439 * red) + (-0.399 * green) + (-0.040 * blue) + 128);
////						pMM2S_Mem[r *COL_SIZE + c] = (Y << 16) | (Cb << 8) | Cr;
////					}
////				}
////
////				for (int r = 0; r < ROW_SIZE; r++) {
////					for (int c = 0; c < COL_SIZE; c++) {
////						uint8_t Y = (pMM2S_Mem[r *COL_SIZE + c] >> 16) & 0xFF;
////						uint8_t Cb = (pMM2S_Mem[r *COL_SIZE + c] >> 8) & 0xFF;
////						uint8_t Cr = pMM2S_Mem[r *COL_SIZE + c] & 0xFF;
////						pMM2S_Mem[r *COL_SIZE + c] = (Cr << 24) | (Y << 16) | (Cb << 8) | Y;
////					}
////				}
//
//
////
////		for (int r = 0; r < ROW_SIZE; r++) {
////			for (int c = 0; c < COL_SIZE; c++) {
////				uint8_t Y = (pMM2S_Mem[r *COL_SIZE + c] >> 16) & 0xFF;
////				uint8_t Cb = (pMM2S_Mem[r *COL_SIZE + c] >> 8) & 0xFF;
////				uint8_t Cr = pMM2S_Mem[r *COL_SIZE + c] & 0xFF;
////				pMM2S_Mem[r *COL_SIZE + c] = (Cr << 24) | (Y << 16) | (Cb << 8) | Y;
////			}
////		}
//
//        // 3. Convert RGB to YCbCr 4:4:4
////        rgbToYCbCr444(fullColor, ROW_SIZE, COL_SIZE);
////
////        // 4. Convert YCbCr 4:4:4 to 4:2:2
////        YCbCr444to422(fullColor, ROW_SIZE, COL_SIZE);
//
//        // 5. Store the processed frame back to memory for HDMI output
////        for (int r = 0; r < ROW_SIZE; r++) {
////            for (int c = 0; c < COL_SIZE; c++) {
////                pMM2S_Mem[r * COL_SIZE + c] = pMM2S_Mem[r *COL_SIZE + c];
////            }
////        }
////    }
//
//
//	// Run for 1000 frames before going back to HW mode
////	for (j = 0; j < 1000; j++) {
////		for (i = 0; i < 1920*1080; i++) {
////	       pMM2S_Mem[i] = pS2MM_Mem[1920*1080-i-1];
////		}
////	}
//
//	// Grab the DMA Control Registers, and re-enable circular park mode.
//	vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET);
//	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);
//	vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET);
//	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);
//
//
//	xil_printf("Main SW processing loop complete!\r\n");
//
//	sleep(5);
//
//	// Uncomment when using TPG for Video input
//	// fmc_imageon_disable_tpg(config);
//
//	sleep(1);
//
//
//	return;
//
//}
