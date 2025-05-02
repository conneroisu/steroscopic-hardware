/*****************************************************************************
 * Joseph Zambreno
 * Phillip Jones
 *
 * Department of Electrical and Computer Engineering
 * Iowa State University
 *****************************************************************************/

/*****************************************************************************
 * xtpg_app.c - Customized macros and routines for the TPG test pattern
 * generator module.
 *
 *
 * NOTES:
 * 02/04/14 by JAZ::Design created.
 *****************************************************************************/

// !!! UNDER CONSTRUCTION: Only use as a high-level structural reference !!!

#include <stdio.h>

// Located in: microblaze_0/include/xparameters.h

// VIVADO 2014 updated to xil_testmem.h  #include "xutil.h"
#include "xil_testmem.h"

#include "xil_io.h"
#include "xparameters.h"
#include "xtpg_app.h"
#define SinTableDepth 2048
#define ENABLE_ZONEPLATE_AND_SWEEPS 1
//#ifdef XPAR_AXI_TPG_0_BASEADDR
void xTPG_help(void);
void xTPG_config(int width, int height, int TPG_BASEADDR, int reset);

void xTPG_main (int width, int height, int TPG_BASEADDR)
{
   unsigned char inchar;
   Xuint32 Reg32Value;
   Xuint32 NewXHairHCoord, NewXHairVCoord, CompMaskValue, CrCbPolarity, NewBoxComponentVal, NewBoxSize, CurrentPattern, CurrentBoxColourPreset, CurrentSpeed, CurrentSPVal, CurrentNoiseGain, CurrentThrottleIndex;
   volatile unsigned long *address, *taddress;
   unsigned int data, mask, data_read;
   int x_offset, y_offset, fstore_address;
   int i, j, num, offset, pass;
//   int width = 1280;
//   int height = 720;

   CurrentBoxColourPreset =0;

   xTPG_config(width, height, TPG_BASEADDR, 1);

   // Default to Zone Plate
   Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
   XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFF0) | 10);
   ZPlate_config(width, height, TPG_BASEADDR);
   // Enable motion
   Reg32Value = XTPG_mReadSlaveReg1(TPG_BASEADDR, 0);
   if ((Reg32Value & 0x00000001)==0)
      XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, (Reg32Value | 0x00000001));
   else
      XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFFE));
   // Speed = 4
   Reg32Value = XTPG_mReadSlaveReg1(TPG_BASEADDR, 0);
   CurrentSpeed = (Reg32Value&0x000001fe)>>1;
   CurrentSpeed = 4;
   XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, (CurrentSpeed<<1)|(Reg32Value & 0x00000001));

   xTPG_help();
   while (1)
   {
      print(">");
      inchar = inbyte();
      xil_printf("%c\n\r",inchar);
      switch (inchar)
      {
         case 0x1b :
         case 'q' :
            xil_printf("- exit menu -\n\n\r"); return; break;
         case '?' :
         {
            xTPG_help();
            break;
         }
         case 'R' :
         {

            print("XTPG register testing and reset !");
            xTPG_help();
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, 0x00);   // PatternSel = 0 (passthru);
            xTPG_config(width, height, TPG_BASEADDR, 1);
            break;
         }
         case 'r' :
         {
            print("Read TPG registers: \r\n\n");

            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 0  (offset 0x0000) (PatternSel)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg1(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 1  (offset 0x0004) (Motion)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg2(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 2  (offset 0x0008) (XHairs)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg3(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 3  (offset 0x000C) (FrameSize)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg4(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 4  (offset 0x0010) (HSweep register)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg5(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 5  (offset 0x0014) (VSweep register)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg6(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 6  (offset 0x0018) (Box Size)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg7(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 7  (offset 0x001C) (Box Colour)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg8(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 8  (offset 0x0020) (Stuck-Pixel frequency index)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg9(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 9  (offset 0x0024) (Noise Gain)\n\r", Reg32Value);
            Reg32Value = XTPG_mReadSlaveReg10(TPG_BASEADDR, 0);
            xil_printf("   Read 0x%08x from register 10 (offset 0x0028) (Throttle level)\n\r", Reg32Value);
            break;
         }
         case '0':  // Video Passthrough
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFF0)|0x00000000);
            break;
         }
         case '1': // HRamp
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, ((Reg32Value & 0xFFFFFFF0)|0x00000001));
            break;
         }
         case '2': // VRamp
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, ((Reg32Value & 0xFFFFFFF0)|0x00000002));
            break;
         }
         case '3': // TRamp
         {
            xil_printf("   Enable motion (hit 'm') to make TRamp work\r\n");
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, ((Reg32Value & 0xFFFFFFF0)|0x00000003));
            break;
         }

         case '4': // Rotate full frame colours R G B Black White
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            CurrentPattern = Reg32Value&0x0000000F;
            if ((CurrentPattern<4) || (CurrentPattern>8))
            {
               CurrentPattern = 4;
            }
            else
            {
               if (CurrentPattern == 8)
                  CurrentPattern = 4;
               else
                  CurrentPattern = CurrentPattern+1;
            }
            xil_printf("   Current Pattern = %d\n\r", CurrentPattern);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFD0) | CurrentPattern);
            break;
         }

         case '5': // Rotate preset box colours White Yellow Cyan Green Magenta Red Blue Black
         {
            if (CurrentBoxColourPreset >= 7)
               CurrentBoxColourPreset = 0;
            else
               CurrentBoxColourPreset++;

            switch (CurrentBoxColourPreset) // Cr Cb Y
            {
               case 0:
               { // white
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x008080EB);
                  break;
               }
               case 1:
               { // yellow
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x000892D2);
                  break;
               }
               case 2:
               { // cyan
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x00A608AA);
                  break;
               }
               case 3:
               { // green
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x00362291);
                  break;
               }
               case 4:
               { // magenta
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x00CABE6A);
                  break;
               }
               case 5:
               { // red
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x005AF051);
                  break;
               }
               case 6:
               { // blue
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x00F06E29);
                  break;
               }
               case 7:
               { // black
                  XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x00808008);
                  break;
               }
            }
            break;
         }



         case '9': // Colour bars
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, ((Reg32Value & 0xFFFFFFD0)|0x00000009));
            break;
         }
         case 'T': // Tartan Bars
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, ((Reg32Value & 0xFFFFFFF0)|0x0000000b));
            break;
         }

         case 'X': // Crosshatch
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, ((Reg32Value & 0xFFFFFFF0)|0x0000000c));
            break;
         }

         case 'S': // Stuck-Pixel emulator
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value & 0x00004000)>>14==0)
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value | 0x00004000);
               xil_printf("Stuck-pixel emulator ON\r\n");
            }
            else
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value & 0xFFFFBFFF);
               xil_printf("Stuck-pixel emulator OFF\r\n");
            }
            break;
         }

         case '=' : // Enable AXI bypass
         {
            //*(unsigned long *)(TPG_BASEADDR+0) ^= 0x04;
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value & 0x00040000) == 0x00040000)
            {
               xil_printf("Disabling AXI Bypass.\n\r");
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value & 0xFFFBFFFF);
            }
            else
            {
               xil_printf("Enabling AXI Bypass.\n\r");
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value | 0x00040000);
            }
            break;
         }



         case 'A': // Decrease Stuck-pixel occurrence speed
         {
            Reg32Value = XTPG_mReadSlaveReg8(TPG_BASEADDR, 0);
            CurrentSPVal = Reg32Value;
            if (CurrentSPVal>15)
            	CurrentSPVal = CurrentSPVal-16;
            else
            	CurrentSPVal = 65520;
            XTPG_mWriteSlaveReg8(TPG_BASEADDR, 0, CurrentSPVal);
            xil_printf("   Current Stuck-Pixel frequency index = %d\r\n", CurrentSPVal);

            break;
         }

         case 'D': // Increase Stuck-pixel occurrence speed
         {
            Reg32Value = XTPG_mReadSlaveReg8(TPG_BASEADDR, 0);
            CurrentSPVal = Reg32Value;
            if (CurrentSPVal<65500)
            	CurrentSPVal = CurrentSPVal+16;
            else
            	CurrentSPVal = 0;
            XTPG_mWriteSlaveReg8(TPG_BASEADDR, 0, CurrentSPVal);
            xil_printf("   Current Stuck-Pixel frequency index = %d\r\n", CurrentSPVal);

            break;
         }
         case '+': // Disable/Enable AXI input interface
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value & 0x00010000)>>16==0)
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value | 0x00010000);
               xil_printf("Input AXI Interface is being IGNORED\r\n");
            }
            else
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value & 0xFFFEFFFF);
               xil_printf("Input AXI Interface is ACTIVE\r\n");
            }
            break;
         }

         case '}': // Data Throttle on/off
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value & 0x00020000)>>17==0)
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value | 0x00020000);
               xil_printf("Data Throttle ON\r\n");
            }
            else
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value & 0xFFFDFFFF);
               xil_printf("Data Throttle OFF\r\n");
            }
            break;
         }


         case '{': // Decrease Throttle Index
         {
            Reg32Value = XTPG_mReadSlaveReg10(TPG_BASEADDR, 0);
            CurrentThrottleIndex = Reg32Value;
            if (CurrentThrottleIndex > 63)
            	CurrentThrottleIndex = CurrentThrottleIndex-64;
            else
            	CurrentThrottleIndex = 65535;
            XTPG_mWriteSlaveReg10(TPG_BASEADDR, 0, CurrentThrottleIndex);
            xil_printf("   Current Throttle Index = %d\r\n", CurrentThrottleIndex);

            break;
         }

         case '|': // Increase Throttle Index
         {
            Reg32Value = XTPG_mReadSlaveReg10(TPG_BASEADDR, 0);
            CurrentThrottleIndex = Reg32Value;
            if (CurrentThrottleIndex<65472)
            	CurrentThrottleIndex = CurrentThrottleIndex+64;
            else
            	CurrentThrottleIndex = 0;
            XTPG_mWriteSlaveReg10(TPG_BASEADDR, 0, CurrentThrottleIndex);
            xil_printf("   Current Throttle Index = %d\r\n", CurrentThrottleIndex);

            break;
         }

         case 'N': // Noise generator
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value & 0x00008000)>>15==0)
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value | 0x00008000);
               xil_printf("Noise generator ON\r\n");
            }
            else
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value & 0xFFFF7FFF);
               xil_printf("Noise generator OFF\r\n");
            }
            break;
         }

         case 'J': // Decrease Noise Gain
         {
            Reg32Value = XTPG_mReadSlaveReg9(TPG_BASEADDR, 0);
            CurrentNoiseGain = Reg32Value;
            if (CurrentNoiseGain>0)
            	CurrentNoiseGain = CurrentNoiseGain-1;
            else
            	CurrentNoiseGain = 255;
            XTPG_mWriteSlaveReg9(TPG_BASEADDR, 0, CurrentNoiseGain);
            xil_printf("   Current Noise Gain Value = %d\r\n", CurrentNoiseGain);

            break;
         }

         case 'K': // Increase Noise Gain
         {
            Reg32Value = XTPG_mReadSlaveReg9(TPG_BASEADDR, 0);
            CurrentNoiseGain = Reg32Value;
            if (CurrentNoiseGain<255)
            	CurrentNoiseGain = CurrentNoiseGain+1;
            else
            	CurrentNoiseGain = 0;
            XTPG_mWriteSlaveReg9(TPG_BASEADDR, 0, CurrentNoiseGain);
            xil_printf("   Current Noise Gain Value = %d\r\n", CurrentNoiseGain);

            break;
         }


         case 'H':  // Invert HSync
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if (((Reg32Value & 0x00000800)>>11)==0)
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value | 0x00000800));
            else
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFF7FF));
            break;
         }
         
         case 'V':  // Invert VSync
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if (((Reg32Value & 0x00000400)>>10)==0)
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value | 0x00000400));
            else
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFBFF));
            break;
         }

         case 'm':
         {
            Reg32Value = XTPG_mReadSlaveReg1(TPG_BASEADDR, 0);
            if ((Reg32Value & 0x00000001)==0)
               XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, (Reg32Value | 0x00000001));
            else
               XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFFE));
            break;
         }

         case 'n': // Decrease Motion speed
         {
            Reg32Value = XTPG_mReadSlaveReg1(TPG_BASEADDR, 0);
            CurrentSpeed = (Reg32Value&0x000001fe)>>1;
            if (CurrentSpeed>0)
               CurrentSpeed--;
            XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, (CurrentSpeed<<1)|(Reg32Value & 0x00000001));
            xil_printf("   Current speed = %d\r\n", CurrentSpeed);

            break;
         }

         case ',': // Increase Motion speed
         {
            Reg32Value = XTPG_mReadSlaveReg1(TPG_BASEADDR, 0);
            CurrentSpeed = (Reg32Value&0x000001fe)>>1;
            if (CurrentSpeed<63)
               CurrentSpeed++;
            XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, (CurrentSpeed<<1)|(Reg32Value & 0x00000001));
            xil_printf("   Current speed = %d\r\n", CurrentSpeed);
            break;
         }


         case 'h': // Horizontal Sweep
         {
            if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
            {
               Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFF0) | 10);
               XTPG_mWriteSlaveReg4(TPG_BASEADDR, 0, 1);
               XTPG_mWriteSlaveReg5(TPG_BASEADDR, 0, 0); // Zoneplate with 0 vertical component
            }
            break;
         }

         case 'j': // Horizontal Sweep increase delta2
         {
            if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
            {
               Reg32Value = XTPG_mReadSlaveReg4(TPG_BASEADDR, 0);
               XTPG_mWriteSlaveReg4(TPG_BASEADDR, 0, Reg32Value+1);
            }
            break;
         }
         case 'g': // Horizontal Sweep decrease delta2
         {
            if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
            {
               Reg32Value = XTPG_mReadSlaveReg4(TPG_BASEADDR, 0);
               XTPG_mWriteSlaveReg4(TPG_BASEADDR, 0, Reg32Value-1);
            }
            break;
         }

         case 'v': // Vertical Sweep
         {
            if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
            {
               Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFF0) | 10);
               XTPG_mWriteSlaveReg4(TPG_BASEADDR, 0, 0); // Zoneplate with 0 horizontal component
               XTPG_mWriteSlaveReg5(TPG_BASEADDR, 0, 1);
            }
            break;
         }

         case 'b': // Horizontal Sweep increase delta2
         {
            if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
            {
               Reg32Value = XTPG_mReadSlaveReg5(TPG_BASEADDR, 0);
               XTPG_mWriteSlaveReg5(TPG_BASEADDR, 0, Reg32Value+1);
            }
            break;
         }
         case 'c': // Horizontal Sweep decrease delta2
         {
            if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
            {
               Reg32Value = XTPG_mReadSlaveReg5(TPG_BASEADDR, 0);
               XTPG_mWriteSlaveReg5(TPG_BASEADDR, 0, Reg32Value-1);
            }
            break;
         }

         case 'z': // Zone-Plate
         {
            if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
            {
               Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0xFFFFFFF0) | 10);
               ZPlate_config(width, height, TPG_BASEADDR);
//               XTPG_mWriteSlaveReg4(TPG_BASEADDR, 0, ((1<<5)*SinTableDepth)*2/width);
//               XTPG_mWriteSlaveReg5(TPG_BASEADDR, 0, ((1<<5)*SinTableDepth)*2/height);
            }
            break;
         }

         case 'x': // Crosshairs
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value&0x00000010)>>4==0)
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value|0x00000010);
            else
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value&0xFFFFFFEF);
            break;
         }
         case 'p': // Horizontal Crosshair - decrease V coordinate
         {
            Reg32Value = XTPG_mReadSlaveReg2(TPG_BASEADDR, 0);
            NewXHairVCoord = ((Reg32Value & 0x0000FFFF)-1) & 0x0000FFFF;
            XTPG_mWriteSlaveReg2(TPG_BASEADDR, 0, (Reg32Value & 0xFFFF0000) | NewXHairVCoord);
            xil_printf("New XHair: (%d, %d)\n\r", (Reg32Value>>16 & 0x0000FFFF), NewXHairVCoord);
            break;
         }
         case 'l': // Horizontal Crosshair - increase V coordinate
         {
            Reg32Value = XTPG_mReadSlaveReg2(TPG_BASEADDR, 0);
            NewXHairVCoord = ((Reg32Value & 0x0000FFFF)+1) & 0x0000FFFF;
            XTPG_mWriteSlaveReg2(TPG_BASEADDR, 0, (Reg32Value & 0xFFFF0000) | NewXHairVCoord);
            xil_printf("New XHair: (%d, %d)\n\r", (Reg32Value>>16 & 0x0000FFFF), NewXHairVCoord);
            break;
         }
         case 'a': // Vertical Crosshair - decrease H coordinate
         {
            Reg32Value = XTPG_mReadSlaveReg2(TPG_BASEADDR, 0);
            NewXHairHCoord = (((Reg32Value>>16) & 0x0000FFFF)-1) & 0x0000FFFF;
            XTPG_mWriteSlaveReg2(TPG_BASEADDR, 0, (Reg32Value & 0x0000FFFF) | NewXHairHCoord<<16);
            xil_printf("New XHair: (%d, %d)\n\r", NewXHairHCoord, (Reg32Value & 0x0000FFFF));
            break;
         }
         case 's': // Vertical Crosshair - increase H coordinate
         {
            Reg32Value = XTPG_mReadSlaveReg2(TPG_BASEADDR, 0);
            NewXHairHCoord = (((Reg32Value>>16) & 0x0000FFFF)+1) & 0x0000FFFF;
            XTPG_mWriteSlaveReg2(TPG_BASEADDR, 0, (Reg32Value & 0x0000FFFF) | NewXHairHCoord<<16);
            xil_printf("New XHair: (%d, %d)\n\r", NewXHairHCoord, (Reg32Value & 0x0000FFFF));
            break;
         }
         case 'i': // Invert CbCr
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value&0x00000020)>>5==0)
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value|0x00000020);
               CrCbPolarity = 1;
            }
            else
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value&0xFFFFFFDF);
               CrCbPolarity = 0;
            }
            xil_printf("New CbCrPolarity = %d\n\r", CrCbPolarity);
            break;

         }

         case '.': // Increment component mask
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            CompMaskValue = (Reg32Value&0x000001C0)>>6;
            if (CompMaskValue == 7)
               CompMaskValue = 0;
            else
               CompMaskValue++;
            xil_printf("   New Component Mask: 0x%x\r\n", CompMaskValue);
            XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (CompMaskValue<<6)|(Reg32Value&0xFFFFFE3F));
            break;
         }

         case 'B': // Box on/off
         {
            Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
            if ((Reg32Value&0x00000200)>>6==0)
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value|0x00000200);
               xil_printf("   Switching box on.\r\n");
            }
            else
            {
               XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, Reg32Value&0xfffffdff);
               xil_printf("   Switching box off.\r\n");
            }
            break;
         }

         case '~': // Box Colour - decrease B/Cr value
         {
            Reg32Value = XTPG_mReadSlaveReg7(TPG_BASEADDR, 0);
            NewBoxComponentVal = (((Reg32Value&0x00FF0000)>>16)-1)&0x000000FF;
            xil_printf("   New Box R/Y value: 0x%x\r\n", NewBoxComponentVal);
            XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, (Reg32Value&0xFF00FFFF)|(NewBoxComponentVal<<16));

            break;
         }
         case '!': // Box Colour - increase B/Cr value
         {
            Reg32Value = XTPG_mReadSlaveReg7(TPG_BASEADDR, 0);
            NewBoxComponentVal = (((Reg32Value&0x00FF0000)>>16)+1)&0x000000FF;
            xil_printf("   New Box R/Y value: 0x%x\r\n", NewBoxComponentVal);
            XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, (Reg32Value&0xFF00FFFF)|(NewBoxComponentVal<<16));
            break;
         }

         case '@': // Box Colour - decrease G/Cb value
         {
            Reg32Value = XTPG_mReadSlaveReg7(TPG_BASEADDR, 0);
            NewBoxComponentVal = (((Reg32Value&0x0000FF00)>>8)-1)&0x000000FF;
            xil_printf("   New Box G/Cb value: 0x%x\r\n", NewBoxComponentVal);
            XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, (Reg32Value&0xFFFF00FF)|(NewBoxComponentVal<<8));
            break;
         }
         case '#': // Box Colour - increase G/Cb value
         {
            Reg32Value = XTPG_mReadSlaveReg7(TPG_BASEADDR, 0);
            NewBoxComponentVal = (((Reg32Value&0x0000FF00)>>8)+1)&0x000000FF;
            xil_printf("   New Box G/Cb value: 0x%x\r\n", NewBoxComponentVal);
            XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, (Reg32Value&0xFFFF00FF)|(NewBoxComponentVal<<8));
            break;
         }

         case '$': // Box Colour - decrease R/Y value
         {
            Reg32Value = XTPG_mReadSlaveReg7(TPG_BASEADDR, 0);
            NewBoxComponentVal = ((Reg32Value&0x000000FF)-1)&0x000000FF;
            xil_printf("   New Box B/Cr value: 0x%x\r\n", NewBoxComponentVal);
            XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, (Reg32Value&0xFFFFFF00)|NewBoxComponentVal);
            break;
         }
         case '%': // Box Colour - increase R/Y value
         {
            Reg32Value = XTPG_mReadSlaveReg7(TPG_BASEADDR, 0);
            NewBoxComponentVal = ((Reg32Value&0x000000FF)+1)&0x000000FF;
            xil_printf("   New Box B/Cr value: 0x%x\r\n", NewBoxComponentVal);
            XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, (Reg32Value&0xFFFFFF00)|NewBoxComponentVal);
            break;
         }
         case '^': // Decrease Box Size
         {
            Reg32Value = XTPG_mReadSlaveReg6(TPG_BASEADDR, 0);
            NewBoxSize = ((Reg32Value&0x00000FFF)-1)&0x00000FFF;
            xil_printf("   New Box Size: 0x%x\r\n", NewBoxSize);
            XTPG_mWriteSlaveReg6(TPG_BASEADDR, 0, NewBoxSize);
            break;
         }
         case '&': // Increase Box Size
         {
            Reg32Value = XTPG_mReadSlaveReg6(TPG_BASEADDR, 0);
            NewBoxSize = ((Reg32Value&0x00000FFF)+1)&0x00000FFF;
            xil_printf("   New Box Size: 0x%x\r\n", NewBoxSize);
            XTPG_mWriteSlaveReg6(TPG_BASEADDR, 0, NewBoxSize);
            break;
         }
      }
   }
}

/************************** Function Definitions ***************************/

/**
 * Enable all possible interrupts from XTPG device.
 * @param   baseaddr_p is the base address of the XTPG device.
 * @return  None.
 * @note    None.
 */
//void XTPG_EnableInterrupt(void * baseaddr_p)
//{
//  Xuint32 baseaddr;
//  baseaddr = (Xuint32) baseaddr_p;
//
//  /*
//   * Enable all interrupt source from user logic.
//   */
//  XTPG_mWriteReg(baseaddr, XTPG_INTR_IPIER_OFFSET, 0x0000000F);
//
//  /*
//   * Enable all possible interrupt sources from device.
//   */
//  XTPG_mWriteReg(baseaddr, XTPG_INTR_DIER_OFFSET,
//    INTR_TERR_MASK
//    | INTR_DPTO_MASK
//    | INTR_IPIR_MASK
//    );
//
//  /*
//   * Set global interrupt enable.
//   */
//  XTPG_mWriteReg(baseaddr, XTPG_INTR_DGIER_OFFSET, INTR_GIE_MASK);
//}

void ZPlate_config(width, height, TPG_BASEADDR)
{

    XTPG_mWriteSlaveReg4(TPG_BASEADDR, 0, ((1<<5)*SinTableDepth)*2/width);
    XTPG_mWriteSlaveReg5(TPG_BASEADDR, 0, ((1<<5)*SinTableDepth)*2/height);
}

void xTPG_help(void)
{
  print("\n\r");
  print("-----------------------------------------------------\n\r");
  print("-- Test-Pattern Generator Menu                          --\n\r");
  print("-----------------------------------------------------\n\r");
  print("\n\r");
  print(" Select option\n\r");
  print(" 0 = TPG Off (pass video input through)\n\r");
  print(" 1 = Horizontal Ramp\n\r");
  print(" 2 = Vertical Ramp\n\r");
  print(" 3 = Temporal Ramp (please switch motion ('m') on for T-Ramp)\n\r");
  print(" 4 = Rotate Planes - Red, Green, Blue, Black, White\n\r");
  print(" 9 = Colour Bars\n\r");
  print(" T = Tartan Bars\n\r");
  print(" X = Crosshatch\n\r");
  if (ENABLE_ZONEPLATE_AND_SWEEPS == 1)
  {
     print(" h = Horizontal Frequency Sweep on/off\n\r");
     print("    g = Decrease H Frequency for H sweep/zone plate\n\r");
     print("    j = Increase H Frequency for H sweep/zone plate\n\r");
     print(" v = Vertical Frequency Sweep on/off\n\r");
     print("    c = Decrease V Frequency for V sweep/zone plate\n\r");
     print("    b = Increase V Frequency for V sweep/zone plate\n\r");
     print(" z = Zone Plate on/off\n\r");
     print("    H and V frequency controls for H/V sweeps apply to zone-plate. \n\r");
   }
  print(" m = Motion on/off\n\r");
  print("    n = Decrease Motion Speed\n\r");
  print("    , = Increase Motion Speed\n\r");
  print(" . = Increment Component Mask (0 = mask off; 7 = mask all 3 components)\n\r");
  print(" x = Cross-Hairs on/off\n\r");
  print("    p = Move H-Cross-Hair up\n\r");
  print("    l = Move H-Cross-Hair down\n\r");
  print("    a = Move V-Cross-Hair left\n\r");
  print("    s = Move V-Cross-Hair right\n\r");
  print(" S = Stuck-pixel emulation on/off\n\r");
  print("    A = Decrease number of stuck-pixel occurrences.\n\r");
  print("    D = Increase number of stuck-pixel occurrences.\n\r");
  print(" N = Random noise generation on/off\n\r");
  print("    J = Decrease noise gain.\n\r");
  print("    K = Increase noise gain.\n\r");
  print(" } = Data-throttle on/off (for AXI-Stream stress-testing)\n\r");
  print("    { = Decrease throttle index.\n\r");
  print("    | = Increase throttle index.\n\r");
  print(" + = Disable Input AXI interface (so TPG ignores TValid)\n\r");
  print(" B = Enable Box overlay.\n\r");
  print("    B = Enable Box overlay.\n\r");
  print("    &, ^ = Increase, decrease Box Size\n\r");
  print("    %, $ = Increase, decrease Box R, Y component value\n\r");
  print("    #, @ = Increase, decrease Box G, Cb component value\n\r");
  print("    !, ~ = Increase, decrease Box B, Cr component value\n\r");
  print("    5 = Rotate Box color around 8 presets.\n\r");
  print("\r\n = = Enable/Disable AXI Bypass       \n\r");
  print("\r\n r = Register print out \n\r");
  print(" R = Reset TPG\n\r");
  print("\n\r");
  print(" q = exit \n\r");
  print(" ?  = help \n\r");
  print("-----------------------------------------------------\n\r");
}


void xTPG_config(int width, int height, int TPG_BASEADDR, int reset)
{
   Xuint32 Reg32Value;
   Xuint32 FrameSize;


   FrameSize = height<<16 | width;
   XTPG_mWriteSlaveReg3(TPG_BASEADDR, 0, FrameSize); 	// Frame size set according to parameter
   ZPlate_config(width, height, TPG_BASEADDR);
   XTPG_mWriteSlaveReg6(TPG_BASEADDR, 0, 64);  			// Set box size to 64

   if (reset == 0)
   {
	   // Read current test-pattern setting.
	   Reg32Value = XTPG_mReadSlaveReg0(TPG_BASEADDR, 0);
	   XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, (Reg32Value & 0x0007FFFF));
   }
   else
   {
	   XTPG_mWriteSlaveReg0(TPG_BASEADDR, 0, 0x00); // PatternSel = 0 (passthru);
													// Box off (bit 9)
													// XHairs Off (bit 4)
													// CrCb Polarity = 0 (bit 5)
	   XTPG_mWriteSlaveReg1(TPG_BASEADDR, 0, 2);   	// Motion Off (bit 0)
													// Motion speed set to 1 (bits 8:1)
	   XTPG_mWriteSlaveReg2(TPG_BASEADDR, 0, (573<<16) | 360); // XHairs set at (573H, 360V)
	   XTPG_mWriteSlaveReg4(TPG_BASEADDR, 0, 1);   	// Set zone plate HDelta2 to 1 (bits 11:0)
                                                    // Set zone plate H sin-table address offset to 0 (bits 27:16)
	   XTPG_mWriteSlaveReg5(TPG_BASEADDR, 0, 1);   	// Set zone plate VDelta2 to 1 (bits 11:0)
                                                    // Set zone plate V sin-table address offset to 0 (bits 27:16)
	   XTPG_mWriteSlaveReg7(TPG_BASEADDR, 0, 0x005AF051);  // Set box colour to mid blue
	   XTPG_mWriteSlaveReg8(TPG_BASEADDR, 0, 0);  	// Set Stuck-Pixel frequency index to 0
	   XTPG_mWriteSlaveReg9(TPG_BASEADDR, 0, 0);  	// Set Noise level index to 0
	   XTPG_mWriteSlaveReg10(TPG_BASEADDR, 0, 0);  	// Set Throttle level to 0
    }

}


//#else
//void xTPG_main(void)
//{
//
//        print("XTPG Not Implemented in this Design!\r\n\r\n");
//}
//
//#endif
