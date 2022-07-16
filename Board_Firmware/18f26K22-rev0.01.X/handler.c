/*
 * File:   handler.c
 * Author: dkg
 *
 * Created on July 5, 2022, 8:38 PM
 */

#include <stdlib.h>

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

    nrf24_write_buff(NRF24_MEM_TX_ADDR, DEFAULT_PIPE_ADDR, 5);
    nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, DEFAULT_PIPE_ADDR, 5);
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


    config.IsConfigured = false;
    config.PingInterval = PING_INT;

    // Init Transmit buffer.
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
    for (uint8_t j = 0; j < sz; j++) {
        packetsTX[i].packet[j] = buffer[j];
    }
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
    if (Ticks % config.PingInterval != 0) {
        return;
    }
    SendPing();
}

/*
void TimerInterruptHandlerOld(void) {
    Ticks++;

    // Send ping packets. 
    if (Ticks % config.PingInterval != 0) {
        return;
    }
    uint8_t sz = SendPing();
    nrf24_send_rf_data(bufferTX, sz);

    // Wait for successful transmission or MAX_RT assertion.
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
        LED_SetHigh();
        nrf24_flush_tx_rx();
        return;
        // TODO: Update primary address to another pipe address.
    }

    // Check for ack payload. 
    if (status & 0x40) {
        uint8_t sz = nrf24_read_dynamic_payload_length();
        nrf24_read_rf_data(bufferRX, sz);
        if (!VerifyBoardAddress(bufferRX)) { // Address does not match.
            return;
        }
        ProcessAckPayload(bufferRX, sz);
    }
}*/

bool VerifyBoardAddress(uint8_t *bufferRX) {
    for (int i = 0; i < 3; i++) {
        if (config.Address[i] != bufferRX[i + 1]) {
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
            for (int i = 0; i < sz - 5; i++) {
                data[i] = buffer[i + 5];
            }
            ProcessActionRequest(actionID, data);
            break;
        case PKT_CFG:
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
        default:
            SendError(ERR_NOT_IMPL);
    }
}

uint8_t SendError(uint8_t errorCode) {
    uint8_t i = 0;
    bufferTX[i] = PKT_DATA;
    for (i = 1; i <= ADDR_LEN; i++) {
        bufferTX[i] = config.Address[i - 1];
    }
    bufferTX[i] = 0; // ActionID.
    bufferTX[++i] = errorCode;
    return QueueTXPacket(bufferTX, (i+1));
}

uint8_t SendPing() {
    bufferTX[0] = PKT_PING;
    for (char i = 0; i < ADDR_LEN; i++) {
        bufferTX[i + 1] = config.Address[i];
    }
    return QueueTXPacket(bufferTX, (ADDR_LEN + 1));
}
