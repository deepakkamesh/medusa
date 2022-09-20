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
        Device            :  PIC16LF18325
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

// get/set DOOR aliases
#define DOOR_TRIS                 TRISAbits.TRISA0
#define DOOR_LAT                  LATAbits.LATA0
#define DOOR_PORT                 PORTAbits.RA0
#define DOOR_WPU                  WPUAbits.WPUA0
#define DOOR_OD                   ODCONAbits.ODCA0
#define DOOR_ANS                  ANSELAbits.ANSA0
#define DOOR_SetHigh()            do { LATAbits.LATA0 = 1; } while(0)
#define DOOR_SetLow()             do { LATAbits.LATA0 = 0; } while(0)
#define DOOR_Toggle()             do { LATAbits.LATA0 = ~LATAbits.LATA0; } while(0)
#define DOOR_GetValue()           PORTAbits.RA0
#define DOOR_SetDigitalInput()    do { TRISAbits.TRISA0 = 1; } while(0)
#define DOOR_SetDigitalOutput()   do { TRISAbits.TRISA0 = 0; } while(0)
#define DOOR_SetPullup()          do { WPUAbits.WPUA0 = 1; } while(0)
#define DOOR_ResetPullup()        do { WPUAbits.WPUA0 = 0; } while(0)
#define DOOR_SetPushPull()        do { ODCONAbits.ODCA0 = 0; } while(0)
#define DOOR_SetOpenDrain()       do { ODCONAbits.ODCA0 = 1; } while(0)
#define DOOR_SetAnalogMode()      do { ANSELAbits.ANSA0 = 1; } while(0)
#define DOOR_SetDigitalMode()     do { ANSELAbits.ANSA0 = 0; } while(0)

// get/set nRF24_CSN aliases
#define nRF24_CSN_TRIS                 TRISAbits.TRISA1
#define nRF24_CSN_LAT                  LATAbits.LATA1
#define nRF24_CSN_PORT                 PORTAbits.RA1
#define nRF24_CSN_WPU                  WPUAbits.WPUA1
#define nRF24_CSN_OD                   ODCONAbits.ODCA1
#define nRF24_CSN_ANS                  ANSELAbits.ANSA1
#define nRF24_CSN_SetHigh()            do { LATAbits.LATA1 = 1; } while(0)
#define nRF24_CSN_SetLow()             do { LATAbits.LATA1 = 0; } while(0)
#define nRF24_CSN_Toggle()             do { LATAbits.LATA1 = ~LATAbits.LATA1; } while(0)
#define nRF24_CSN_GetValue()           PORTAbits.RA1
#define nRF24_CSN_SetDigitalInput()    do { TRISAbits.TRISA1 = 1; } while(0)
#define nRF24_CSN_SetDigitalOutput()   do { TRISAbits.TRISA1 = 0; } while(0)
#define nRF24_CSN_SetPullup()          do { WPUAbits.WPUA1 = 1; } while(0)
#define nRF24_CSN_ResetPullup()        do { WPUAbits.WPUA1 = 0; } while(0)
#define nRF24_CSN_SetPushPull()        do { ODCONAbits.ODCA1 = 0; } while(0)
#define nRF24_CSN_SetOpenDrain()       do { ODCONAbits.ODCA1 = 1; } while(0)
#define nRF24_CSN_SetAnalogMode()      do { ANSELAbits.ANSA1 = 1; } while(0)
#define nRF24_CSN_SetDigitalMode()     do { ANSELAbits.ANSA1 = 0; } while(0)

// get/set nRF24_CE aliases
#define nRF24_CE_TRIS                 TRISAbits.TRISA2
#define nRF24_CE_LAT                  LATAbits.LATA2
#define nRF24_CE_PORT                 PORTAbits.RA2
#define nRF24_CE_WPU                  WPUAbits.WPUA2
#define nRF24_CE_OD                   ODCONAbits.ODCA2
#define nRF24_CE_ANS                  ANSELAbits.ANSA2
#define nRF24_CE_SetHigh()            do { LATAbits.LATA2 = 1; } while(0)
#define nRF24_CE_SetLow()             do { LATAbits.LATA2 = 0; } while(0)
#define nRF24_CE_Toggle()             do { LATAbits.LATA2 = ~LATAbits.LATA2; } while(0)
#define nRF24_CE_GetValue()           PORTAbits.RA2
#define nRF24_CE_SetDigitalInput()    do { TRISAbits.TRISA2 = 1; } while(0)
#define nRF24_CE_SetDigitalOutput()   do { TRISAbits.TRISA2 = 0; } while(0)
#define nRF24_CE_SetPullup()          do { WPUAbits.WPUA2 = 1; } while(0)
#define nRF24_CE_ResetPullup()        do { WPUAbits.WPUA2 = 0; } while(0)
#define nRF24_CE_SetPushPull()        do { ODCONAbits.ODCA2 = 0; } while(0)
#define nRF24_CE_SetOpenDrain()       do { ODCONAbits.ODCA2 = 1; } while(0)
#define nRF24_CE_SetAnalogMode()      do { ANSELAbits.ANSA2 = 1; } while(0)
#define nRF24_CE_SetDigitalMode()     do { ANSELAbits.ANSA2 = 0; } while(0)

// get/set MOTION aliases
#define MOTION_PORT                 PORTAbits.RA3
#define MOTION_WPU                  WPUAbits.WPUA3
#define MOTION_GetValue()           PORTAbits.RA3
#define MOTION_SetPullup()          do { WPUAbits.WPUA3 = 1; } while(0)
#define MOTION_ResetPullup()        do { WPUAbits.WPUA3 = 0; } while(0)

// get/set RELAY aliases
#define RELAY_TRIS                 TRISAbits.TRISA4
#define RELAY_LAT                  LATAbits.LATA4
#define RELAY_PORT                 PORTAbits.RA4
#define RELAY_WPU                  WPUAbits.WPUA4
#define RELAY_OD                   ODCONAbits.ODCA4
#define RELAY_ANS                  ANSELAbits.ANSA4
#define RELAY_SetHigh()            do { LATAbits.LATA4 = 1; } while(0)
#define RELAY_SetLow()             do { LATAbits.LATA4 = 0; } while(0)
#define RELAY_Toggle()             do { LATAbits.LATA4 = ~LATAbits.LATA4; } while(0)
#define RELAY_GetValue()           PORTAbits.RA4
#define RELAY_SetDigitalInput()    do { TRISAbits.TRISA4 = 1; } while(0)
#define RELAY_SetDigitalOutput()   do { TRISAbits.TRISA4 = 0; } while(0)
#define RELAY_SetPullup()          do { WPUAbits.WPUA4 = 1; } while(0)
#define RELAY_ResetPullup()        do { WPUAbits.WPUA4 = 0; } while(0)
#define RELAY_SetPushPull()        do { ODCONAbits.ODCA4 = 0; } while(0)
#define RELAY_SetOpenDrain()       do { ODCONAbits.ODCA4 = 1; } while(0)
#define RELAY_SetAnalogMode()      do { ANSELAbits.ANSA4 = 1; } while(0)
#define RELAY_SetDigitalMode()     do { ANSELAbits.ANSA4 = 0; } while(0)

// get/set LED aliases
#define LED_TRIS                 TRISAbits.TRISA5
#define LED_LAT                  LATAbits.LATA5
#define LED_PORT                 PORTAbits.RA5
#define LED_WPU                  WPUAbits.WPUA5
#define LED_OD                   ODCONAbits.ODCA5
#define LED_ANS                  ANSELAbits.ANSA5
#define LED_SetHigh()            do { LATAbits.LATA5 = 1; } while(0)
#define LED_SetLow()             do { LATAbits.LATA5 = 0; } while(0)
#define LED_Toggle()             do { LATAbits.LATA5 = ~LATAbits.LATA5; } while(0)
#define LED_GetValue()           PORTAbits.RA5
#define LED_SetDigitalInput()    do { TRISAbits.TRISA5 = 1; } while(0)
#define LED_SetDigitalOutput()   do { TRISAbits.TRISA5 = 0; } while(0)
#define LED_SetPullup()          do { WPUAbits.WPUA5 = 1; } while(0)
#define LED_ResetPullup()        do { WPUAbits.WPUA5 = 0; } while(0)
#define LED_SetPushPull()        do { ODCONAbits.ODCA5 = 0; } while(0)
#define LED_SetOpenDrain()       do { ODCONAbits.ODCA5 = 1; } while(0)
#define LED_SetAnalogMode()      do { ANSELAbits.ANSA5 = 1; } while(0)
#define LED_SetDigitalMode()     do { ANSELAbits.ANSA5 = 0; } while(0)

// get/set RC0 procedures
#define RC0_SetHigh()            do { LATCbits.LATC0 = 1; } while(0)
#define RC0_SetLow()             do { LATCbits.LATC0 = 0; } while(0)
#define RC0_Toggle()             do { LATCbits.LATC0 = ~LATCbits.LATC0; } while(0)
#define RC0_GetValue()              PORTCbits.RC0
#define RC0_SetDigitalInput()    do { TRISCbits.TRISC0 = 1; } while(0)
#define RC0_SetDigitalOutput()   do { TRISCbits.TRISC0 = 0; } while(0)
#define RC0_SetPullup()             do { WPUCbits.WPUC0 = 1; } while(0)
#define RC0_ResetPullup()           do { WPUCbits.WPUC0 = 0; } while(0)
#define RC0_SetAnalogMode()         do { ANSELCbits.ANSC0 = 1; } while(0)
#define RC0_SetDigitalMode()        do { ANSELCbits.ANSC0 = 0; } while(0)

// get/set RC1 procedures
#define RC1_SetHigh()            do { LATCbits.LATC1 = 1; } while(0)
#define RC1_SetLow()             do { LATCbits.LATC1 = 0; } while(0)
#define RC1_Toggle()             do { LATCbits.LATC1 = ~LATCbits.LATC1; } while(0)
#define RC1_GetValue()              PORTCbits.RC1
#define RC1_SetDigitalInput()    do { TRISCbits.TRISC1 = 1; } while(0)
#define RC1_SetDigitalOutput()   do { TRISCbits.TRISC1 = 0; } while(0)
#define RC1_SetPullup()             do { WPUCbits.WPUC1 = 1; } while(0)
#define RC1_ResetPullup()           do { WPUCbits.WPUC1 = 0; } while(0)
#define RC1_SetAnalogMode()         do { ANSELCbits.ANSC1 = 1; } while(0)
#define RC1_SetDigitalMode()        do { ANSELCbits.ANSC1 = 0; } while(0)

// get/set RC2 procedures
#define RC2_SetHigh()            do { LATCbits.LATC2 = 1; } while(0)
#define RC2_SetLow()             do { LATCbits.LATC2 = 0; } while(0)
#define RC2_Toggle()             do { LATCbits.LATC2 = ~LATCbits.LATC2; } while(0)
#define RC2_GetValue()              PORTCbits.RC2
#define RC2_SetDigitalInput()    do { TRISCbits.TRISC2 = 1; } while(0)
#define RC2_SetDigitalOutput()   do { TRISCbits.TRISC2 = 0; } while(0)
#define RC2_SetPullup()             do { WPUCbits.WPUC2 = 1; } while(0)
#define RC2_ResetPullup()           do { WPUCbits.WPUC2 = 0; } while(0)
#define RC2_SetAnalogMode()         do { ANSELCbits.ANSC2 = 1; } while(0)
#define RC2_SetDigitalMode()        do { ANSELCbits.ANSC2 = 0; } while(0)

// get/set ADC_LIGHT aliases
#define ADC_LIGHT_TRIS                 TRISCbits.TRISC3
#define ADC_LIGHT_LAT                  LATCbits.LATC3
#define ADC_LIGHT_PORT                 PORTCbits.RC3
#define ADC_LIGHT_WPU                  WPUCbits.WPUC3
#define ADC_LIGHT_OD                   ODCONCbits.ODCC3
#define ADC_LIGHT_ANS                  ANSELCbits.ANSC3
#define ADC_LIGHT_SetHigh()            do { LATCbits.LATC3 = 1; } while(0)
#define ADC_LIGHT_SetLow()             do { LATCbits.LATC3 = 0; } while(0)
#define ADC_LIGHT_Toggle()             do { LATCbits.LATC3 = ~LATCbits.LATC3; } while(0)
#define ADC_LIGHT_GetValue()           PORTCbits.RC3
#define ADC_LIGHT_SetDigitalInput()    do { TRISCbits.TRISC3 = 1; } while(0)
#define ADC_LIGHT_SetDigitalOutput()   do { TRISCbits.TRISC3 = 0; } while(0)
#define ADC_LIGHT_SetPullup()          do { WPUCbits.WPUC3 = 1; } while(0)
#define ADC_LIGHT_ResetPullup()        do { WPUCbits.WPUC3 = 0; } while(0)
#define ADC_LIGHT_SetPushPull()        do { ODCONCbits.ODCC3 = 0; } while(0)
#define ADC_LIGHT_SetOpenDrain()       do { ODCONCbits.ODCC3 = 1; } while(0)
#define ADC_LIGHT_SetAnalogMode()      do { ANSELCbits.ANSC3 = 1; } while(0)
#define ADC_LIGHT_SetDigitalMode()     do { ANSELCbits.ANSC3 = 0; } while(0)

// get/set RC4 procedures
#define RC4_SetHigh()            do { LATCbits.LATC4 = 1; } while(0)
#define RC4_SetLow()             do { LATCbits.LATC4 = 0; } while(0)
#define RC4_Toggle()             do { LATCbits.LATC4 = ~LATCbits.LATC4; } while(0)
#define RC4_GetValue()              PORTCbits.RC4
#define RC4_SetDigitalInput()    do { TRISCbits.TRISC4 = 1; } while(0)
#define RC4_SetDigitalOutput()   do { TRISCbits.TRISC4 = 0; } while(0)
#define RC4_SetPullup()             do { WPUCbits.WPUC4 = 1; } while(0)
#define RC4_ResetPullup()           do { WPUCbits.WPUC4 = 0; } while(0)
#define RC4_SetAnalogMode()         do { ANSELCbits.ANSC4 = 1; } while(0)
#define RC4_SetDigitalMode()        do { ANSELCbits.ANSC4 = 0; } while(0)

// get/set RC5 procedures
#define RC5_SetHigh()            do { LATCbits.LATC5 = 1; } while(0)
#define RC5_SetLow()             do { LATCbits.LATC5 = 0; } while(0)
#define RC5_Toggle()             do { LATCbits.LATC5 = ~LATCbits.LATC5; } while(0)
#define RC5_GetValue()              PORTCbits.RC5
#define RC5_SetDigitalInput()    do { TRISCbits.TRISC5 = 1; } while(0)
#define RC5_SetDigitalOutput()   do { TRISCbits.TRISC5 = 0; } while(0)
#define RC5_SetPullup()             do { WPUCbits.WPUC5 = 1; } while(0)
#define RC5_ResetPullup()           do { WPUCbits.WPUC5 = 0; } while(0)
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


/**
 * @Param
    none
 * @Returns
    none
 * @Description
    Interrupt on Change Handler for the IOCAF0 pin functionality
 * @Example
    IOCAF0_ISR();
 */
void IOCAF0_ISR(void);

/**
  @Summary
    Interrupt Handler Setter for IOCAF0 pin interrupt-on-change functionality

  @Description
    Allows selecting an interrupt handler for IOCAF0 at application runtime
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    InterruptHandler function pointer.

  @Example
    PIN_MANAGER_Initialize();
    IOCAF0_SetInterruptHandler(MyInterruptHandler);

*/
void IOCAF0_SetInterruptHandler(void (* InterruptHandler)(void));

/**
  @Summary
    Dynamic Interrupt Handler for IOCAF0 pin

  @Description
    This is a dynamic interrupt handler to be used together with the IOCAF0_SetInterruptHandler() method.
    This handler is called every time the IOCAF0 ISR is executed and allows any function to be registered at runtime.
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCAF0_SetInterruptHandler(IOCAF0_InterruptHandler);

*/
extern void (*IOCAF0_InterruptHandler)(void);

/**
  @Summary
    Default Interrupt Handler for IOCAF0 pin

  @Description
    This is a predefined interrupt handler to be used together with the IOCAF0_SetInterruptHandler() method.
    This handler is called every time the IOCAF0 ISR is executed. 
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCAF0_SetInterruptHandler(IOCAF0_DefaultInterruptHandler);

*/
void IOCAF0_DefaultInterruptHandler(void);


/**
 * @Param
    none
 * @Returns
    none
 * @Description
    Interrupt on Change Handler for the IOCAF3 pin functionality
 * @Example
    IOCAF3_ISR();
 */
void IOCAF3_ISR(void);

/**
  @Summary
    Interrupt Handler Setter for IOCAF3 pin interrupt-on-change functionality

  @Description
    Allows selecting an interrupt handler for IOCAF3 at application runtime
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    InterruptHandler function pointer.

  @Example
    PIN_MANAGER_Initialize();
    IOCAF3_SetInterruptHandler(MyInterruptHandler);

*/
void IOCAF3_SetInterruptHandler(void (* InterruptHandler)(void));

/**
  @Summary
    Dynamic Interrupt Handler for IOCAF3 pin

  @Description
    This is a dynamic interrupt handler to be used together with the IOCAF3_SetInterruptHandler() method.
    This handler is called every time the IOCAF3 ISR is executed and allows any function to be registered at runtime.
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCAF3_SetInterruptHandler(IOCAF3_InterruptHandler);

*/
extern void (*IOCAF3_InterruptHandler)(void);

/**
  @Summary
    Default Interrupt Handler for IOCAF3 pin

  @Description
    This is a predefined interrupt handler to be used together with the IOCAF3_SetInterruptHandler() method.
    This handler is called every time the IOCAF3 ISR is executed. 
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCAF3_SetInterruptHandler(IOCAF3_DefaultInterruptHandler);

*/
void IOCAF3_DefaultInterruptHandler(void);



#endif // PIN_MANAGER_H
/**
 End of File
*/