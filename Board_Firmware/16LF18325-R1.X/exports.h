/* Microchip Technology Inc. and its subsidiaries.  You may use this software 
 * and any derivatives exclusively with Microchip products. 
 * 
 * THIS SOFTWARE IS SUPPLIED BY MICROCHIP "AS IS".  NO WARRANTIES, WHETHER 
 * EXPRESS, IMPLIED OR STATUTORY, APPLY TO THIS SOFTWARE, INCLUDING ANY IMPLIED 
 * WARRANTIES OF NON-INFRINGEMENT, MERCHANTABILITY, AND FITNESS FOR A 
 * PARTICULAR PURPOSE, OR ITS INTERACTION WITH MICROCHIP PRODUCTS, COMBINATION 
 * WITH ANY OTHER PRODUCTS, OR USE IN ANY APPLICATION. 
 *
 * IN NO EVENT WILL MICROCHIP BE LIABLE FOR ANY INDIRECT, SPECIAL, PUNITIVE, 
 * INCIDENTAL OR CONSEQUENTIAL LOSS, DAMAGE, COST OR EXPENSE OF ANY KIND 
 * WHATSOEVER RELATED TO THE SOFTWARE, HOWEVER CAUSED, EVEN IF MICROCHIP HAS 
 * BEEN ADVISED OF THE POSSIBILITY OR THE DAMAGES ARE FORESEEABLE.  TO THE 
 * FULLEST EXTENT ALLOWED BY LAW, MICROCHIP'S TOTAL LIABILITY ON ALL CLAIMS 
 * IN ANY WAY RELATED TO THIS SOFTWARE WILL NOT EXCEED THE AMOUNT OF FEES, IF 
 * ANY, THAT YOU HAVE PAID DIRECTLY TO MICROCHIP FOR THIS SOFTWARE.
 *
 * MICROCHIP PROVIDES THIS SOFTWARE CONDITIONALLY UPON YOUR ACCEPTANCE OF THESE 
 * TERMS. 
 */

/* 
 * File:   
 * Author: 
 * Comments:
 * Revision history: 
 */

// This is a guard condition so that contents of this file are not included
// more than once.  
#ifndef XC_HEADER_TEMPLATE_H
#define	XC_HEADER_TEMPLATE_H

#include <xc.h> // include processor files - each processor file is guarded.  
#include "mcc_generated_files/mcc.h"

#endif	/* XC_HEADER_TEMPLATE_H */


/******************************************************************************/
// SPI and GPIO Helper Mappers
/******************************************************************************/
// NRF24 Functions. 
#define NRF24L01_CSN_H()            nRF24_CSN_SetHigh()
#define NRF24L01_CSN_L()            nRF24_CSN_SetLow()
#define NRF24L01_CSN_SetOutput()    nRF24_CSN_SetDigitalOutput()
#define NRF24L01_CSN_SetInput()     nRF24_CSN_SetDigitalInput()

#define NRF24L01_CE_H()             nRF24_CE_SetHigh()
#define NRF24L01_CE_L()             nRF24_CE_SetLow()
#define NRF24L01_CE_SetOutput()     nRF24_CE_SetDigitalOutput()
#define NRF24L01_CE_SetInput()      nRF24_CE_SetDigitalInput()

#define SPI_WRITE_BYTE(dt)          SPI1_ExchangeByte(dt)
#define SPI_READ_BYTE(dt)           SPI1_ExchangeByte(dt)
#define SPI_INIT()                  SPI1_Open(SPI1_DEFAULT)

// Motion Functions. 
#define Motion_SetInterruptHandler(ih) IOCAF3_SetInterruptHandler(ih)
#define MOTIONGetValue()    MOTION_GetValue()

// Door Functions.
#define Door_SetInterruptHandler(ih) IOCAF0_SetInterruptHandler(ih)
#define zDOOR_GetValue() DOOR_GetValue()

// ADC Functions.
#define zADC_GetConversion(a) ADC_GetConversion(a)

// I2C Functions.
#define i2cRead1bReg(a,b)    i2c_read1ByteRegister(a,b)
#define i2cWriteBytes(a,b,c) i2c_writeNBytes(a,b,c)
#define i2cReadBytes(a,b,c)  i2c_readNBytes(a,b,c)
#define i2cAddr i2c2_address_t

// LED Functions.
#define zLED_Toggle() LED_Toggle()
#define zLED_SetHigh() LED_SetHigh()
#define zLED_SetLow() LED_SetLow()

// Relay Functions. 
#define zRELAY_Toggle() RELAY_Toggle()
#define zRELAY_SetHigh() RELAY_SetHigh()
#define zRELAY_SetLow() RELAY_SetLow()

#// EEPROM functions.
#define zDATAEE_ReadByte(a) DATAEE_ReadByte(a)
#define zDATAEE_WriteByte(a,b) DATAEE_WriteByte(a,b)

// Version:
// TODO: Update with every major change. 
#define VER_LOW_BYTE 0x01 // Software version updates. 
#define VER_HIGH_BYTE 0x1 // HW Version Major number

// Set Hardware revision here. 
#define HW_REV_1_3

#ifdef HW_REV_1_1 // Available Actions for Rev 1_1
#define DEV_STATUS_LED 
#define DEV_VOLTS 
#define DEV_LIGHT 
#define DEV_MOTION 
#define DEV_TEMP_HUMIDITY 
#define DEV_RELAY 
#define DEV_DOOR 
#endif

#ifdef HW_REV_1_2 // Available Actions for Rev 1_2
#define DEV_STATUS_LED 
#define DEV_VOLTS 
#define DEV_LIGHT 
#define DEV_MOTION 
#endif

#ifdef HW_REV_1_3 // Available Actions for Rev 1_2
#define DEV_STATUS_LED 
#define DEV_VOLTS 
#define DEV_LIGHT 
#endif