/*****************************************************************************
 * Joseph Zambreno
 * Phillip Jones
 *
 * Department of Electrical and Computer Engineering
 * Iowa State University
 *****************************************************************************/

/*****************************************************************************
 * xtpg_app.h - Customized macros and routines for the TPG test pattern
 * generator module.
 *
 *
 * NOTES:
 * 02/04/14 by JAZ::Design created.
 *****************************************************************************/

#ifndef XTPG_APP_H
#define XTPG_APP_H

/***************************** Include Files *******************************/

// !!! UNDER CONSTRUCTION: Only use as a high-level structural reference !!!

#include "xbasic_types.h"
#include "xstatus.h"
#include "xil_io.h"

void xTPG_main (int width, int height, int TPG_BASEADDR);
void ZPlate_config(int width, int height, int TPG_BASEADDR);


//#ifdef XPAR_AXI_TPG_0_BASEADDR

/************************** Constant Definitions ***************************/

#define XTPG_USER_SLV_SPACE_OFFSET (0x00000000)
#define XTPG_SLV_REG0_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000000)
#define XTPG_SLV_REG1_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000004)
#define XTPG_SLV_REG2_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000008)
#define XTPG_SLV_REG3_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000000C)
#define XTPG_SLV_REG4_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000010)
#define XTPG_SLV_REG5_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000014)
#define XTPG_SLV_REG6_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000018)
#define XTPG_SLV_REG7_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000001C)
#define XTPG_SLV_REG8_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000020)
#define XTPG_SLV_REG9_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000024)
#define XTPG_SLV_REG10_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000028)
#define XTPG_SLV_REG11_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000002C)
#define XTPG_SLV_REG12_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000030)
#define XTPG_SLV_REG13_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000034)
#define XTPG_SLV_REG14_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000038)
#define XTPG_SLV_REG15_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000003C)
#define XTPG_SLV_REG16_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000040)
#define XTPG_SLV_REG17_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000044)
#define XTPG_SLV_REG18_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000048)
#define XTPG_SLV_REG19_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000004C)
#define XTPG_SLV_REG20_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000050)
#define XTPG_SLV_REG21_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000054)
#define XTPG_SLV_REG22_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000058)
#define XTPG_SLV_REG23_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000005C)
#define XTPG_SLV_REG24_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000060)
#define XTPG_SLV_REG25_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000064)
#define XTPG_SLV_REG26_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000068)
#define XTPG_SLV_REG27_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000006C)
#define XTPG_SLV_REG28_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000070)
#define XTPG_SLV_REG29_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000074)
#define XTPG_SLV_REG30_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000078)
#define XTPG_SLV_REG31_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x0000007C)
#define XTPG_SLV_REG32_OFFSET (XTPG_USER_SLV_SPACE_OFFSET + 0x00000080)


/**
 * Software Reset Space Register Offsets
 * -- RST : software reset register
 */
#define XTPG_SOFT_RST_SPACE_OFFSET (0x00000100)
#define XTPG_RST_REG_OFFSET (XTPG_SOFT_RST_SPACE_OFFSET + 0x00000000)

/**
 * Software Reset Masks
 * -- SOFT_RESET : software reset
 */
#define SOFT_RESET (0x0000000A)

/**
 * Interrupt Controller Space Offsets
 * -- INTR_DISR  : device (peripheral) interrupt status register
 * -- INTR_DIPR  : device (peripheral) interrupt pending register
 * -- INTR_DIER  : device (peripheral) interrupt enable register
 * -- INTR_DIIR  : device (peripheral) interrupt id (priority encoder) register
 * -- INTR_DGIER : device (peripheral) global interrupt enable register
 * -- INTR_ISR   : ip (user logic) interrupt status register
 * -- INTR_IER   : ip (user logic) interrupt enable register
 */
#define XTPG_INTR_CNTRL_SPACE_OFFSET (0x00000200)
#define XTPG_INTR_DISR_OFFSET (XTPG_INTR_CNTRL_SPACE_OFFSET + 0x00000000)
#define XTPG_INTR_DIPR_OFFSET (XTPG_INTR_CNTRL_SPACE_OFFSET + 0x00000004)
#define XTPG_INTR_DIER_OFFSET (XTPG_INTR_CNTRL_SPACE_OFFSET + 0x00000008)
#define XTPG_INTR_DIIR_OFFSET (XTPG_INTR_CNTRL_SPACE_OFFSET + 0x00000018)
#define XTPG_INTR_DGIER_OFFSET (XTPG_INTR_CNTRL_SPACE_OFFSET + 0x0000001C)
#define XTPG_INTR_IPISR_OFFSET (XTPG_INTR_CNTRL_SPACE_OFFSET + 0x00000020)
#define XTPG_INTR_IPIER_OFFSET (XTPG_INTR_CNTRL_SPACE_OFFSET + 0x00000028)

/**
 * Interrupt Controller Masks
 * -- INTR_TERR_MASK : transaction error
 * -- INTR_DPTO_MASK : data phase time-out
 * -- INTR_IPIR_MASK : ip interrupt requeset
 * -- INTR_RFDL_MASK : read packet fifo deadlock interrupt request
 * -- INTR_WFDL_MASK : write packet fifo deadlock interrupt request
 * -- INTR_IID_MASK  : interrupt id
 * -- INTR_GIE_MASK  : global interrupt enable
 * -- INTR_NOPEND    : the DIPR has no pending interrupts
 */
#define INTR_TERR_MASK (0x00000001UL)
#define INTR_DPTO_MASK (0x00000002UL)
#define INTR_IPIR_MASK (0x00000004UL)
#define INTR_RFDL_MASK (0x00000020UL)
#define INTR_WFDL_MASK (0x00000040UL)
#define INTR_IID_MASK (0x000000FFUL)
#define INTR_GIE_MASK (0x80000000UL)
#define INTR_NOPEND (0x80)

/**************************** Type Definitions *****************************/


/***************** Macros (Inline Functions) Definitions *******************/

/**
 *
 * Write a value to a XTPG register. A 32 bit write is performed.
 * If the component is implemented in a smaller width, only the least
 * significant data is written.
 *
 * @param   BaseAddress is the base address of the XTPG device.
 * @param   RegOffset is the register offset from the base to write to.
 * @param   Data is the data written to the register.
 *
 * @return  None.
 *
 * @note
 * C-style signature:
 *  void XTPG_mWriteReg(Xuint32 BaseAddress, unsigned RegOffset, Xuint32 Data)
 *
 */
#define XTPG_mWriteReg(BaseAddress, RegOffset, Data) \
    Xil_Out32((BaseAddress) + (RegOffset), (Xuint32)(Data))

/**
 *
 * Read a value from a XTPG register. A 32 bit read is performed.
 * If the component is implemented in a smaller width, only the least
 * significant data is read from the register. The most significant data
 * will be read as 0.
 *
 * @param   BaseAddress is the base address of the XTPG device.
 * @param   RegOffset is the register offset from the base to write to.
 *
 * @return  Data is the data from the register.
 *
 * @note
 * C-style signature:
 *  Xuint32 XTPG_mReadReg(Xuint32 BaseAddress, unsigned RegOffset)
 *
 */
#define XTPG_mReadReg(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (RegOffset))


/**
 *
 * Write/Read 32 bit value to/from XTPG user logic slave registers.
 *
 * @param   BaseAddress is the base address of the XTPG device.
 * @param   RegOffset is the offset from the slave register to write to or read from.
 * @param   Value is the data written to the register.
 *
 * @return  Data is the data from the user logic slave register.
 *
 * @note
 * C-style signature:
 *  void XTPG_mWriteSlaveRegn(Xuint32 BaseAddress, unsigned RegOffset, Xuint32 Value)
 *  Xuint32 XTPG_mReadSlaveRegn(Xuint32 BaseAddress, unsigned RegOffset)
 *
 */
#define XTPG_mWriteSlaveReg0(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG0_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg1(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG1_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg2(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG2_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg3(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG3_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg4(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG4_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg5(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG5_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg6(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG6_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg7(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG7_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg8(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG8_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg9(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG9_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg10(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG10_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg11(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG11_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg12(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG12_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg13(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG13_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg14(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG14_OFFSET) + (RegOffset), (Xuint32)(Value))
#define XTPG_mWriteSlaveReg15(BaseAddress, RegOffset, Value) \
    Xil_Out32((BaseAddress) + (XTPG_SLV_REG15_OFFSET) + (RegOffset), (Xuint32)(Value))

#define XTPG_mReadSlaveReg0(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG0_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg1(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG1_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg2(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG2_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg3(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG3_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg4(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG4_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg5(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG5_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg6(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG6_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg7(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG7_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg8(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG8_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg9(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG9_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg10(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG10_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg11(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG11_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg12(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG12_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg13(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG13_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg14(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG14_OFFSET) + (RegOffset))
#define XTPG_mReadSlaveReg15(BaseAddress, RegOffset) \
    Xil_In32((BaseAddress) + (XTPG_SLV_REG15_OFFSET) + (RegOffset))


/**
 *
 * Reset XTPG via software.
 *
 * @param   BaseAddress is the base address of the XTPG device.
 *
 * @return  None.
 *
 * @note
 * C-style signature:
 *  void XTPG_mReset(Xuint32 BaseAddress)
 *
 */
#define XTPG_mReset(BaseAddress) \
    Xil_Out32((BaseAddress)+(XTPG_RST_REG_OFFSET), SOFT_RESET)

/************************** Function Prototypes ****************************/


/**
 *
 * Enable all possible interrupts from XTPG device.
 *
 * @param   baseaddr_p is the base address of the XTPG device.
 *
 * @return  None.
 *
 * @note    None.
 *
 */
void XTPG_EnableInterrupt(void * baseaddr_p);

/**
 *
 * Example interrupt controller handler.
 *
 * @param   baseaddr_p is the base address of the XTPG device.
 *
 * @return  None.
 *
 * @note    None.
 *
 */

void xTPG_config(int width, int height, int TPG_BASEADDR, int reset);

void XTPG_Intr_DefaultHandler(void * baseaddr_p);

/**
 *
 * Run a self-test on the driver/device. Note this may be a destructive test if
 * resets of the device are performed.
 *
 * If the hardware system is not built correctly, this function may never
 * return to the caller.
 *
 * @param   baseaddr_p is the base address of the XTPG instance to be worked on.
 *
 * @return
 *
 *    - XST_SUCCESS   if all self-test code passed
 *    - XST_FAILURE   if any self-test code failed
 *
 * @note    Caching must be turned off for this function to work.
 * @note    Self test may fail if data memory and device are not on the same bus.
 *
 */
XStatus XTPG_SelfTest(void * baseaddr_p);
//#endif
#endif // XTPG_H
