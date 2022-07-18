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

#ifdef	__cplusplus
extern "C" {
#endif /* __cplusplus */

    // TODO If C++ is being used, regular C code needs function names to have C 
    // linkage so the functions can be used by the c code. 

#ifdef	__cplusplus
}
#endif /* __cplusplus */

#endif	/* XC_HEADER_TEMPLATE_H */

// Protocol Stuff
#define ADDR_LEN 3
#define MIN_PKT_SZ 4
#define MAX_PKT_SZ 32
#define MAX_TX_QUEUE_SZ 8 

// Packet Types.
#define PKT_DATA 0x01
#define PKT_PING 0x02
#define PKT_CFG_1 0x03
#define PKT_CFG_2 0x04
#define PKT_ACTION 0x10

// Error Types.
#define ERR_NOT_IMPL 0x04

#define ACTION_STATUS_LED 0x13
#define ACTION_RELOAD_CONFIG 0x15

#define PIPE_ADDR_LEN 5 
#define DEFAULT_RF_CHANNEL 115
#define DEFAULT_ARD 0xA // default ARD setting. (val*250 +250)
uint8_t DEFAULT_PIPE_ADDR[] = "hello"; // Default pipe address to bootstrap.
uint8_t PingInterval = 1; // Default ping interval.
uint8_t BoardAddress[3] = {0xFF,0xFF,0xFF}; // Default board address.

void TimerInterruptHandler(void);
void InitRadio(void);
void ProcessAckPayload(uint8_t * buffer, uint8_t sz);
void ProcessActionRequest(uint8_t actionID, uint8_t * data);
bool VerifyBoardAddress(uint8_t *bufferRX);
void HandlePacketLoop(void);
uint8_t SendError(uint8_t errorCode);
uint8_t SendPing();
void SuperMemCpy(uint8_t *dest, uint8_t destStart, uint8_t *src, uint8_t srcStart, uint8_t sz);
void ReloadConfig(void);

typedef struct {
    uint8_t packet[MAX_PKT_SZ];
    bool free;
    uint8_t size; // Size of packet.
} Packet;

// Config stores the configuration of the board.

struct Config {
    bool IsConfigured; // True if board is configured.
    uint8_t Address[3]; // Address of the board.
    uint8_t PingInterval; // Ping interval in seconds
    uint8_t RFChannel; // Frequency channel.
    uint8_t PipeAddr1[5];
    uint8_t PipeAddr2[5]; // Backup Pipe Address.
    uint8_t ARD; // Auto Retry Duration. 
};