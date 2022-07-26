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
 * Comments: This header is included with the nRF24_lib to expose MCC functions.
 * Revision history: 
 */

// This is a guard condition so that contents of this file are not included
// more than once.  
#ifndef XC_HEADER_TEMPLATE_H
#define	XC_HEADER_TEMPLATE_H

#include <xc.h> // include processor files - each processor file is guarded.  
#include "mcc_generated_files/mcc.h"


#ifdef	__cplusplus
extern "C" {
#endif /* __cplusplus */

#ifdef	__cplusplus
}
#endif /* __cplusplus */

#endif	/* XC_HEADER_TEMPLATE_H */

#define HW_REV_0

#ifdef HW_REV_0
/******************************************************************************/
// SPI and GPIO Helper Function.
/******************************************************************************/

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

#define DEV_STATUS_LED 1
#define DEV_TEMP_HUMIDITY 1
#endif