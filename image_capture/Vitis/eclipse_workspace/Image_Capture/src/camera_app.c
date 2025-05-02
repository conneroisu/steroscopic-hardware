
#include <stdlib.h>
#include "camera_app.h"
#include "xil_types.h"
#include "xtime_l.h"
#include <stdio.h>
#include "xil_printf.h"

// Opcodes
#define STOP_OP_CODE 0
#define SEND_IMAGE_OP_CODE 1

#define IMAGE_WIDTH 1920
#define IMAGE_HEIGHT 1080

camera_config_t camera_config;

u8 get_current_frame_pointer(XAxiVdma* vdma, u16 dir)
{

	u8 result = 0;

	u32 mask = 0;
	u32 shift_amt = 0;

	if(dir == XAXIVDMA_READ)
	{
		mask = 0x1F0000;
		shift_amt = 16;
	}
	else if(dir == XAXIVDMA_WRITE)
	{
		mask = 0x1F00000;
		shift_amt = 24;
	}

	result = (*((volatile u32*) (vdma->BaseAddr + XAXIVDMA_PARKPTR_OFFSET)) & mask) >> shift_amt;

	return result;
}

void set_park_frame(XAxiVdma* vdma, u8 frame, u16 dir)
{
#define	PARK *((volatile u32*) (vdma->BaseAddr + XAXIVDMA_PARKPTR_OFFSET))

	u32 mask = 0;
	u32 shift_amt = 0;

	if(dir == XAXIVDMA_READ)
	{
		mask = ~0x1F;
	}
	else if(dir == XAXIVDMA_WRITE)
	{
		mask = ~0x1F0;
		shift_amt = 8;
	}

	PARK = (PARK & mask) | ((u32)(frame & 0x1F) << shift_amt);


#undef PARK
}


// Main function. Initializes the devices and configures VDMA
int main()
{
	camera_config_init(&camera_config);
	fmc_imageon_enable(&camera_config);
	camera_loop(&camera_config);

	return 0;
}


// Initialize the camera configuration data structure
void camera_config_init(camera_config_t *config)
{

    config->uBaseAddr_IIC_FmcIpmi =  XPAR_FMC_IPMI_ID_EEPROM_0_BASEADDR;   // Device for reading HDMI board IPMI EEPROM information
    config->uBaseAddr_IIC_FmcImageon = XPAR_FMC_IMAGEON_IIC_0_BASEADDR;    // Device for configuring the HDMI board

     config->uBaseAddr_VITA_SPI = XPAR_ONSEMI_VITA_SPI_0_S00_AXI_BASEADDR;  // Device for configuring the Camera sensor
     config->uBaseAddr_VITA_CAM = XPAR_ONSEMI_VITA_CAM_0_S00_AXI_BASEADDR;  // Device for receiving Camera sensor data

    config->uDeviceId_VTC_tpg   = XPAR_V_TC_0_DEVICE_ID;                        // Video Timer Controller (VTC) ID
    config->uDeviceId_VDMA_HdmiFrameBuffer = XPAR_AXI_VDMA_0_DEVICE_ID;         // VDMA ID
    config->uBaseAddr_MEM_HdmiFrameBuffer = XPAR_DDR_MEM_BASEADDR + 0x10000000; // VDMA base address for Frame buffers
    config->uNumFrames_HdmiFrameBuffer = XPAR_AXIVDMA_0_NUM_FSTORES;            // NUmber of VDMA Frame buffers

    return;
}


void camera_loop(camera_config_t *config)
{

	Xuint32 vdma_S2MM_DMACR, vdma_MM2S_DMACR;

	// Read should park on 1, this will stream the greyscale image.
	// Start write on 1 for now.
	set_park_frame(&(config->vdma_hdmi), 1, XAXIVDMA_WRITE);
	set_park_frame(&(config->vdma_hdmi), 1, XAXIVDMA_READ);

	// Grab the DMA Control Registers, and clear circular park mode.
	vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET);
	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);
	vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET);
	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);


	// Pointer to snapshot.
	// Each sensor reading is 16 bits.
	volatile Xuint16 *snapshot = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET);

	// Main loop
	while(1)
	{
		char opcode = inbyte();

		// Send image
		if(opcode == SEND_IMAGE_OP_CODE)
		{
			// Reply with the opcode byte to let the host PC know that we got the opcode
			outbyte(opcode);

			// Save a snapshot to frame 0.
			set_park_frame(&(config->vdma_hdmi), 0, XAXIVDMA_WRITE);

			// Wait until frame zero is being written to.
			while(get_current_frame_pointer(&(config->vdma_hdmi), XAXIVDMA_WRITE))
			{

			}

			set_park_frame(&(config->vdma_hdmi), 1, XAXIVDMA_WRITE);

			// Wait until frame one is being written to.
			while(!get_current_frame_pointer(&(config->vdma_hdmi), XAXIVDMA_WRITE))
			{

			}

			// Send the entire greyscale image over UART.
			// Note: Two bytes need to be sent for a single sensor reading
			for(int i = 0; i < (IMAGE_WIDTH * IMAGE_HEIGHT); i++)
			{
				outbyte((char) (snapshot[i] & 0xFF));
			}
		}
		else if(!opcode)
		{
			// Terminate
			break;
		}
	}

}
