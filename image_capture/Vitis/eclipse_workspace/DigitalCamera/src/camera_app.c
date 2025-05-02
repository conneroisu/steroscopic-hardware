
#include <stdlib.h>
#include "camera_app.h"
#include "xil_types.h"
#include "include/demosaicing.h"
#include "xscugic.h"
#include "xtime_l.h"
#include <stdio.h>

#define MAX_FPS_ENTRIES 10

typedef struct fps
{
	int next_to_write;
	int size;
	float entries[MAX_FPS_ENTRIES];
} fps_t;

// Timing globals
static fps_t fps;
static int sw_mode = 0;
static int snapshot_saved = 0;

static u8 back_buffer_frame = 2;
static u8 front_buffer_frame = 3;

// Frame that indicates that the write completed.
static u8 target_frame = 2;

static XTime tStart, tEnd = 0;

// In seconds
static float frame_time = 0;

// To store frame time message
static char* fps_msg = 0;

camera_config_t camera_config;

void fps_init(fps_t* f)
{
	f->next_to_write = 0;
	f->size = 0;

	for(int i = 0; i < MAX_FPS_ENTRIES; ++i)
	{
		f->entries[i] = 0;
	}
}

void fps_time_store(fps_t* f, float time)
{
	f->entries[f->next_to_write] = time;

	if(f->next_to_write == (MAX_FPS_ENTRIES - 1))
	{
		f->next_to_write = 0;
	}
	else
	{
		f->next_to_write++;
	}

	if(f->size != MAX_FPS_ENTRIES)
	{
		f->size++;
	}
}

float fps_calculate(fps_t* f)
{
	float sum = 0;

	for(int i = 0; i < f->size; ++i)
	{
		sum += f->entries[i];
	}

	return f->size / sum;
}

// Swaps the memory addresses associated with the frame pointers
void swap_frame_pointers(XAxiVdma* vdma, u16 dir, u8 a, u8 b)
{
	u32 offset = XAXIVDMA_TX_OFFSET ? dir == XAXIVDMA_WRITE : XAXIVDMA_RX_OFFSET;

#define START_ADDR(index) *((volatile u32*) (vdma->BaseAddr + offset + XAXIVDMA_START_ADDR_OFFSET + (index * 0x4)))
#define VSIZE *((volatile u32*) (vdma->BaseAddr + offset + XAXIVDMA_VSIZE_OFFSET))

	a &= 0x1F;
	b &= 0x1F;

	u32 start_a = START_ADDR(a);
	u32 start_b = START_ADDR(b);

	START_ADDR(a) = start_b;
	START_ADDR(b) = start_a;

	VSIZE = VSIZE;

#undef START_ADDR
#undef VSIZE
}

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

void error_isr(void* CallBackRef, u32 InterruptTypes)
{
	xil_printf("VDMA error %X occurred!!!\n\r", InterruptTypes);
}

void video_frame_output_isr(void* CallBackRef, u32 InterruptTypes)
{
	switch(InterruptTypes)
	{
		case XAXIVDMA_IXR_FRMCNT_MASK:
		{
			// Once we have read from the back buffer, we know that it is the new front buffer
			// so we must have just swapped.
			if(sw_mode && (get_current_frame_pointer((XAxiVdma*) CallBackRef, XAXIVDMA_READ) == target_frame) && snapshot_saved)
			{
				// We start and stop the timer on this ISR.
				// If this is the first frame, the timer won't be started, so start it up without ending any timer.
				if(!tStart)
				{
					XTime_GetTime(&tStart);
				}
				else
				{
					XTime_GetTime(&tEnd);
					snapshot_saved = 0;

					frame_time = (tEnd - tStart) / (float)COUNTS_PER_SECOND;

					fps_time_store(&fps, frame_time);

					sprintf(fps_msg, "Average FPS: %.5f", fps_calculate(&fps));

					xil_printf("%s\n\r", fps_msg);

					// Start the timer back up!
					XTime_GetTime(&tStart);
				}


			}
		}

		default:
		{
			break;
		}
	}

}

void camera_input_isr(void* CallBackRef, u32 InterruptTypes)
{
	switch(InterruptTypes)
	{
		case XAXIVDMA_IXR_FRMCNT_MASK:
		{
			if(sw_mode && get_current_frame_pointer((XAxiVdma*)CallBackRef, XAXIVDMA_WRITE) == 0 && !snapshot_saved)
			{
				snapshot_saved = 1;
				target_frame = back_buffer_frame;
			}
		}


		default:
		{
			break;
		}
	}

}


// Main function. Initializes the devices and configures VDMA
int main() {

	// Init the FPS struct.
	fps_init(&fps);

	fps_msg = (char*) calloc(64, sizeof(char*));

	camera_config_init(&camera_config);
	fmc_imageon_enable(&camera_config);
	camera_loop(&camera_config);

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
    //config->uBaseAddr_TPG_PatternGenerator = XPAR_V_TPG_0_S_AXI_CTRL_BASEADDR; // TPG Device

    config->uDeviceId_VTC_tpg   = XPAR_V_TC_0_DEVICE_ID;                        // Video Timer Controller (VTC) ID
    config->uDeviceId_VDMA_HdmiFrameBuffer = XPAR_AXI_VDMA_0_DEVICE_ID;         // VDMA ID
    config->uBaseAddr_MEM_HdmiFrameBuffer = XPAR_DDR_MEM_BASEADDR + 0x10000000; // VDMA base address for Frame buffers
    config->uNumFrames_HdmiFrameBuffer = XPAR_AXIVDMA_0_NUM_FSTORES;            // NUmber of VDMA Frame buffers

    return;
}


// Main (SW) processing loop. Recommended to have an explicit exit condition
void camera_loop(camera_config_t *config) {

	Xuint32 vdma_S2MM_DMACR, vdma_MM2S_DMACR;


	xil_printf("Entering main SW processing loop\r\n");

	// Grab the DMA parkptr, and update it to ensure that when parked, the S2MM side is on frame 1, and the MM2S side on frame 2
	set_park_frame(&(config->vdma_hdmi), 1, XAXIVDMA_WRITE);
	set_park_frame(&(config->vdma_hdmi), 2, XAXIVDMA_READ);

	// Grab the DMA Control Registers, and clear circular park mode.
	vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET);
	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);
	vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET);
	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR & ~XAXIVDMA_CR_TAIL_EN_MASK);

	sw_mode = 1;

	// Pointers to the S2MM memory frame and M2SS memory frame
	volatile Xuint16 *pS2MM_Mem = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_S2MM_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET);
	volatile Xuint16 *pMM2S_Mem = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET + (front_buffer_frame * 0x4));


	xil_printf("Start processing 1000 frames!\r\n");
	xil_printf("pS2MM_Mem = %X\n\r", pS2MM_Mem);
	xil_printf("pMM2S_Mem = %X\n\r", pMM2S_Mem);

	// Run for 100 frames before going back to HW mode
	for (int j = 0; j < 1000; j++)
	{
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

		// Apply CFA
		run_demosaicing((uint16_t*)pS2MM_Mem, (uint16_t*)pMM2S_Mem);

		// Swap back and front buffers
		u8 temp = back_buffer_frame;

		back_buffer_frame = front_buffer_frame;
		front_buffer_frame = temp;

		// Have the read side park on the new front buffer
		set_park_frame(&(config->vdma_hdmi), front_buffer_frame, XAXIVDMA_READ);

		// Wait for park frame to update.
		while(get_current_frame_pointer(&(config->vdma_hdmi), XAXIVDMA_READ) != front_buffer_frame)
		{

		}

		// Update pMM2S_Mem to point to back buffer.
		pMM2S_Mem = (Xuint16 *)XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_MM2S_ADDR_OFFSET+XAXIVDMA_START_ADDR_OFFSET + (back_buffer_frame * 0x4));

	}

	sw_mode = 0;

	// Grab the DMA Control Registers, and re-enable circular park mode.
	vdma_MM2S_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET);
	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_TX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_MM2S_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);
	vdma_S2MM_DMACR = XAxiVdma_ReadReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET);
	XAxiVdma_WriteReg(config->vdma_hdmi.BaseAddr, XAXIVDMA_RX_OFFSET+XAXIVDMA_CR_OFFSET, vdma_S2MM_DMACR | XAXIVDMA_CR_TAIL_EN_MASK);


	xil_printf("Main SW processing loop complete!\r\n");

	sleep(5);

	//fmc_imageon_disable_tpg(config);

	sleep(1);


	return;
}
