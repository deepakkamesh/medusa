/**
  @Generated Pin Manager Header File

  @Company:
    Microchip Technology Inc.

  @File Name:
    pin_manager.h

  @Summary:
    This is the Pin Manager file generated using PIC10 / PIC12 / PIC16 / PIC18 MCUs

  @Description
    This header file provides APIs for driver for .
    Generation Information :
        Product Revision  :  PIC10 / PIC12 / PIC16 / PIC18 MCUs - 1.81.7
        Device            :  PIC18F26K22
        Driver Version    :  2.11
    The generated drivers are tested against the following:
        Compiler          :  XC8 2.31 and above
        MPLAB 	          :  MPLAB X 5.45	
*/

/*
    (c) 2018 Microchip Technology Inc. and its subsidiaries. 
    
    Subject to your compliance with these terms, you may use Microchip software and any 
    derivatives exclusively with Microchip products. It is your responsibility to comply with third party 
    license terms applicable to your use of third party software (including open source software) that 
    may accompany Microchip software.
    
    THIS SOFTWARE IS SUPPLIED BY MICROCHIP "AS IS". NO WARRANTIES, WHETHER 
    EXPRESS, IMPLIED OR STATUTORY, APPLY TO THIS SOFTWARE, INCLUDING ANY 
    IMPLIED WARRANTIES OF NON-INFRINGEMENT, MERCHANTABILITY, AND FITNESS 
    FOR A PARTICULAR PURPOSE.
    
    IN NO EVENT WILL MICROCHIP BE LIABLE FOR ANY INDIRECT, SPECIAL, PUNITIVE, 
    INCIDENTAL OR CONSEQUENTIAL LOSS, DAMAGE, COST OR EXPENSE OF ANY KIND 
    WHATSOEVER RELATED TO THE SOFTWARE, HOWEVER CAUSED, EVEN IF MICROCHIP 
    HAS BEEN ADVISED OF THE POSSIBILITY OR THE DAMAGES ARE FORESEEABLE. TO 
    THE FULLEST EXTENT ALLOWED BY LAW, MICROCHIP'S TOTAL LIABILITY ON ALL 
    CLAIMS IN ANY WAY RELATED TO THIS SOFTWARE WILL NOT EXCEED THE AMOUNT 
    OF FEES, IF ANY, THAT YOU HAVE PAID DIRECTLY TO MICROCHIP FOR THIS 
    SOFTWARE.
*/

#ifndef PIN_MANAGER_H
#define PIN_MANAGER_H

/**
  Section: Included Files
*/

#include <xc.h>

#define INPUT   1
#define OUTPUT  0

#define HIGH    1
#define LOW     0

#define ANALOG      1
#define DIGITAL     0

#define PULL_UP_ENABLED      1
#define PULL_UP_DISABLED     0

// get/set LED aliases
#define LED_TRIS                 TRISAbits.TRISA1
#define LED_LAT                  LATAbits.LATA1
#define LED_PORT                 PORTAbits.RA1
#define LED_ANS                  ANSELAbits.ANSA1
#define LED_SetHigh()            do { LATAbits.LATA1 = 1; } while(0)
#define LED_SetLow()             do { LATAbits.LATA1 = 0; } while(0)
#define LED_Toggle()             do { LATAbits.LATA1 = ~LATAbits.LATA1; } while(0)
#define LED_GetValue()           PORTAbits.RA1
#define LED_SetDigitalInput()    do { TRISAbits.TRISA1 = 1; } while(0)
#define LED_SetDigitalOutput()   do { TRISAbits.TRISA1 = 0; } while(0)
#define LED_SetAnalogMode()      do { ANSELAbits.ANSA1 = 1; } while(0)
#define LED_SetDigitalMode()     do { ANSELAbits.ANSA1 = 0; } while(0)

// get/set nRF24_CSN aliases
#define nRF24_CSN_TRIS                 TRISCbits.TRISC1
#define nRF24_CSN_LAT                  LATCbits.LATC1
#define nRF24_CSN_PORT                 PORTCbits.RC1
#define nRF24_CSN_SetHigh()            do { LATCbits.LATC1 = 1; } while(0)
#define nRF24_CSN_SetLow()             do { LATCbits.LATC1 = 0; } while(0)
#define nRF24_CSN_Toggle()             do { LATCbits.LATC1 = ~LATCbits.LATC1; } while(0)
#define nRF24_CSN_GetValue()           PORTCbits.RC1
#define nRF24_CSN_SetDigitalInput()    do { TRISCbits.TRISC1 = 1; } while(0)
#define nRF24_CSN_SetDigitalOutput()   do { TRISCbits.TRISC1 = 0; } while(0)

// get/set nRF24_CE aliases
#define nRF24_CE_TRIS                 TRISCbits.TRISC2
#define nRF24_CE_LAT                  LATCbits.LATC2
#define nRF24_CE_PORT                 PORTCbits.RC2
#define nRF24_CE_ANS                  ANSELCbits.ANSC2
#define nRF24_CE_SetHigh()            do { LATCbits.LATC2 = 1; } while(0)
#define nRF24_CE_SetLow()             do { LATCbits.LATC2 = 0; } while(0)
#define nRF24_CE_Toggle()             do { LATCbits.LATC2 = ~LATCbits.LATC2; } while(0)
#define nRF24_CE_GetValue()           PORTCbits.RC2
#define nRF24_CE_SetDigitalInput()    do { TRISCbits.TRISC2 = 1; } while(0)
#define nRF24_CE_SetDigitalOutput()   do { TRISCbits.TRISC2 = 0; } while(0)
#define nRF24_CE_SetAnalogMode()      do { ANSELCbits.ANSC2 = 1; } while(0)
#define nRF24_CE_SetDigitalMode()     do { ANSELCbits.ANSC2 = 0; } while(0)

// get/set RC3 procedures
#define RC3_SetHigh()            do { LATCbits.LATC3 = 1; } while(0)
#define RC3_SetLow()             do { LATCbits.LATC3 = 0; } while(0)
#define RC3_Toggle()             do { LATCbits.LATC3 = ~LATCbits.LATC3; } while(0)
#define RC3_GetValue()              PORTCbits.RC3
#define RC3_SetDigitalInput()    do { TRISCbits.TRISC3 = 1; } while(0)
#define RC3_SetDigitalOutput()   do { TRISCbits.TRISC3 = 0; } while(0)
#define RC3_SetAnalogMode()         do { ANSELCbits.ANSC3 = 1; } while(0)
#define RC3_SetDigitalMode()        do { ANSELCbits.ANSC3 = 0; } while(0)

// get/set RC4 procedures
#define RC4_SetHigh()            do { LATCbits.LATC4 = 1; } while(0)
#define RC4_SetLow()             do { LATCbits.LATC4 = 0; } while(0)
#define RC4_Toggle()             do { LATCbits.LATC4 = ~LATCbits.LATC4; } while(0)
#define RC4_GetValue()              PORTCbits.RC4
#define RC4_SetDigitalInput()    do { TRISCbits.TRISC4 = 1; } while(0)
#define RC4_SetDigitalOutput()   do { TRISCbits.TRISC4 = 0; } while(0)
#define RC4_SetAnalogMode()         do { ANSELCbits.ANSC4 = 1; } while(0)
#define RC4_SetDigitalMode()        do { ANSELCbits.ANSC4 = 0; } while(0)

// get/set RC5 procedures
#define RC5_SetHigh()            do { LATCbits.LATC5 = 1; } while(0)
#define RC5_SetLow()             do { LATCbits.LATC5 = 0; } while(0)
#define RC5_Toggle()             do { LATCbits.LATC5 = ~LATCbits.LATC5; } while(0)
#define RC5_GetValue()              PORTCbits.RC5
#define RC5_SetDigitalInput()    do { TRISCbits.TRISC5 = 1; } while(0)
#define RC5_SetDigitalOutput()   do { TRISCbits.TRISC5 = 0; } while(0)
#define RC5_SetAnalogMode()         do { ANSELCbits.ANSC5 = 1; } while(0)
#define RC5_SetDigitalMode()        do { ANSELCbits.ANSC5 = 0; } while(0)

/**
   @Param
    none
   @Returns
    none
   @Description
    GPIO and peripheral I/O initialization
   @Example
    PIN_MANAGER_Initialize();
 */
void PIN_MANAGER_Initialize (void);

/**
 * @Param
    none
 * @Returns
    none
 * @Description
    Interrupt on Change Handling routine
 * @Example
    PIN_MANAGER_IOC();
 */
void PIN_MANAGER_IOC(void);



#endif // PIN_MANAGER_H
/**
 End of File
*/