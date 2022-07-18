/*
 * File:   handler.c
 * Author: dkg
 *
 * Created on July 5, 2022, 8:38 PM
 */

#include <stdlib.h>
#include <string.h>

#include "mcc_generated_files/mcc.h"
#include "../lib/nrf24_lib.h"
#include "../lib/dht11_lib.h"

#include "handler.h"


uint32_t Ticks = 0; // Ticks of timer.
struct Config config; // Board config.

uint8_t bufferTX[32];
uint8_t bufferRX[32];

Packet packetsTX[MAX_TX_QUEUE_SZ]; // Transmit buffer.

void InitRadio(void) {
    nrf24_rf_init();

    nrf24_write_buff(NRF24_MEM_TX_ADDR, DEFAULT_PIPE_ADDR, PIPE_ADDR_LEN);
    nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, DEFAULT_PIPE_ADDR, PIPE_ADDR_LEN);
    // Mask all interrupts, EN_CRC,CRCO=1byte,PWR_UP,PTX, 
    nrf24_write_register(NRF24_MEM_CONFIG, 0b1111010);
    // ENAA_P0. 
    nrf24_write_register(NRF24_MEM_EN_AA, 0b1);
    // ERX_P0.
    nrf24_write_register(NRF24_MEM_EN_RXADDR, 0b1);
    // AW=5byte.
    nrf24_write_register(NRF24_MEM_SETUP_AW, 0b11);
    // Retry Settings. 
    nrf24_write_register(NRF24_MEM_SETUP_RETR, 0b11011111);
    // RF channel.
    nrf24_write_register(NRF24_MEM_RF_CH, 115);
    // RF_PWR=0bDm, RF_DR_HIGH=2Mbps.
    nrf24_write_register(NRF24_MEM_RF_SETUP, 0b1110);
    // EN_DPL, EN_ACK_PAY.
    nrf24_write_register(NRF24_MEM_FEATURE, 0b110); // Enable Dynamic payload, ack payload.
    // DPL_P0.
    nrf24_write_register(NRF24_MEM_DYNPD, 0b1); // Dynamic payload on Pipe 0.
    __delay_us(10);

    // Setup defaults in Config.
    config.IsConfigured = false;
    config.PingInterval = PingInterval;
    config.RFChannel = DEFAULT_RF_CHANNEL;
    memcpy(config.PipeAddr1, DEFAULT_PIPE_ADDR, PIPE_ADDR_LEN);
    memcpy(config.PipeAddr2, DEFAULT_PIPE_ADDR, PIPE_ADDR_LEN);
    memcpy(config.Address, BoardAddress, ADDR_LEN);
    config.ARD = DEFAULT_ARD;

    // Initialize Transmit buffer.
    for (uint8_t i = 0; i < MAX_TX_QUEUE_SZ; i++) {
        packetsTX[i].free = true;
        packetsTX[i].size = 0;
    }
}

// QueueTXPacket queues a packet to be transmitted. 
// returns 1 if success or 0 if no free slot.

uint8_t QueueTXPacket(uint8_t *buffer, uint8_t sz) {
    uint8_t i;
    for (i = 0; i < MAX_TX_QUEUE_SZ; i++) {
        if (packetsTX[i].free) {
            break;
        }
    }
    if (i == MAX_TX_QUEUE_SZ) {
        return 0;
    }
    packetsTX[i].free = false;
    packetsTX[i].size = sz;
    memcpy(packetsTX[i].packet, buffer, sz);
    return 1;
}

// Sends packets 1 at a time.

void HandlePacketLoop(void) {
    uint8_t i;

    // Check if there are any packets to send.
    for (i = 0; i < MAX_TX_QUEUE_SZ; i++) {
        if (!packetsTX[i].free) {
            break;
        }
    }
    if (i == MAX_TX_QUEUE_SZ) {
        return; // Nothing to do.
    }

    nrf24_send_rf_data(packetsTX[i].packet, packetsTX[i].size);

    // Wait for successful TX or MAX_RT assertion.
    uint8_t status = 0;
    while (1) {
        status = nrf24_read_register(NRF24_MEM_STATUSS);
        if ((status & 0x20) || (status & 0x10)) {
            break;
        }
        __delay_us(10);
    }
    // Clear status register.
    nrf24_write_register(NRF24_MEM_STATUSS, 0x70);

    // MAX_RT exceeded. 
    if (status & 0x10) {
        nrf24_flush_tx_rx();
        return;
        // TODO: Update primary address to another pipe address.
        // after x retries shutdown to conserve battery.
    }

    // Free up the slot since packet was transmitted successfully.
    packetsTX[i].free = true;

    // Check for ack payload. 
    if (status & 0x40) {
        uint8_t sz = nrf24_read_dynamic_payload_length();
        nrf24_read_rf_data(bufferRX, sz);
        if (!VerifyBoardAddress(bufferRX)) { // Address does not match.
            return;
        }
        ProcessAckPayload(bufferRX, sz);
    }
}

void TimerInterruptHandler(void) {
    Ticks++;

    // Send ping packets. 
    if (Ticks % PingInterval != 0) {
        return;
    }
    SendPing();
}

bool VerifyBoardAddress(uint8_t *bufferRX) {
    for (int i = 0; i < ADDR_LEN; i++) {
        if (BoardAddress[i] != bufferRX[i + 1]) {
            return false;
        }
    }
    return true;
}

void ProcessAckPayload(uint8_t * buffer, uint8_t sz) {
    uint8_t data[32];
    uint8_t actionID;

    uint8_t pktType = buffer[0];
    switch (pktType) {
        case PKT_ACTION:
            actionID = buffer[4];
            SuperMemCpy(data, 0, buffer, 5, sz - 5);
            ProcessActionRequest(actionID, data);
            break;
        case PKT_CFG_1:
            config.RFChannel = buffer[4];
            SuperMemCpy(config.PipeAddr1, 0, buffer, 5, PIPE_ADDR_LEN);
            SuperMemCpy(config.PipeAddr2, 0, buffer, 10, PIPE_ADDR_LEN);
            config.ARD = buffer[15];
            break;
        case PKT_CFG_2:
            SuperMemCpy(config.Address, 0, buffer, 4, ADDR_LEN);
            config.PingInterval = buffer[7];
            break;
        default:
            SendError(ERR_NOT_IMPL);
    }

}

void ProcessActionRequest(uint8_t actionID, uint8_t * data) {

    switch (actionID) {
        case ACTION_STATUS_LED:
            LED_SetLow();
            if (data[0]) {
                LED_SetHigh();
            }
            break;
        case ACTION_RELOAD_CONFIG:
            ReloadConfig();
            break;
        default:
            SendError(ERR_NOT_IMPL);
    }
}

/* ReloadConfig loads the config in the config struct and makes tit active */
void ReloadConfig(void) {
    config.IsConfigured = true;
    nrf24_write_register(NRF24_MEM_RF_CH, config.RFChannel);

    nrf24_write_buff(NRF24_MEM_TX_ADDR, config.PipeAddr1, PIPE_ADDR_LEN);
    nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, config.PipeAddr1, PIPE_ADDR_LEN);

    uint8_t ard = (config.ARD << 4) | 0xF;
    nrf24_write_register(NRF24_MEM_SETUP_RETR, ard);

    memcpy(BoardAddress, config.Address, ADDR_LEN);
    PingInterval = config.PingInterval;
}

uint8_t SendError(uint8_t errorCode) {
    uint8_t i = 0;
    bufferTX[i] = PKT_DATA;
    SuperMemCpy(bufferTX, 1, BoardAddress, 0, ADDR_LEN);
    i += ADDR_LEN;
    bufferTX[++i] = 0; // ActionID.
    bufferTX[++i] = errorCode;
    return QueueTXPacket(bufferTX, (i + 1));
}

uint8_t SendPing() {
    bufferTX[0] = PKT_PING;
    SuperMemCpy(bufferTX, 1, BoardAddress, 0, ADDR_LEN);
    return QueueTXPacket(bufferTX, (ADDR_LEN + 1));
}

void SuperMemCpy(uint8_t *dest, uint8_t destStart, uint8_t *src, uint8_t srcStart, uint8_t sz) {
    for (uint8_t i = 0; i < sz; i++) {
        dest[i + destStart] = src[i + srcStart];
    }
}