/*
 * File:   handler.c
 * Author: dkg
 *
 * Created on July 5, 2022, 8:38 PM
 */


#include "mcc_generated_files/mcc.h"
#include "../lib/nrf24_lib.h"
#include "../lib/dht11_lib.h"

#include "handler.h"


uint32_t Ticks = 0; // Ticks of timer.
struct Config config; // Board config.

uint8_t bufferTX[32];
uint8_t bufferRX[32];

void InitRadio(void) {
    nrf24_rf_init();

    nrf24_write_buff(NRF24_MEM_TX_ADDR, DEFAULT_PIPE_ADDR, 5);
    nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, DEFAULT_PIPE_ADDR, 5);
    // EN_CRC,CRCO=1byte,PWR_UP,PTX.   
    nrf24_write_register(NRF24_MEM_CONFIG, 0b1010);
    // ENAA_P0. 
    nrf24_write_register(NRF24_MEM_EN_AA, 0b1);
    // ERX_P0.
    nrf24_write_register(NRF24_MEM_EN_RXADDR, 0b1);
    // AW=5byte.
    nrf24_write_register(NRF24_MEM_SETUP_AW, 0b11);
    // Retry Settings. 
    nrf24_write_register(NRF24_MEM_SETUP_RETR, 0b10101010);
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
}

void TimerInterruptHandler(void) {
    Ticks++;

    // Send ping packets. 
    if (Ticks % config.PingInterval != 0) {
        return;
    }
    unsigned char sz = MakePingPkt(bufferTX);
    nrf24_send_rf_data(bufferTX, sz);

    // Wait for successful transmission or MAX_RT assertion.
    unsigned char status = 0;
    while (1) {
        status = nrf24_read_register(NRF24_MEM_STATUSS);
        if ((status & 0x20) || (status & 0x10)) {
            break;
        }
    }
    // Clear status register.
    nrf24_write_register(NRF24_MEM_STATUSS, 0x70);

    // MAX_RT exceeded.
    if (status & 0x10) {
        return;
        // TODO: Update primary address to another pipe address.
    }

    // Check for ack payload. 
    if (status & 0x40) {
        uint8_t sz = nrf24_read_dynamic_payload_length();
        nrf24_read_rf_data(bufferRX, sz);
        ProcessAckPayload(bufferRX, sz);
    }

}

void ProcessAckPayload(unsigned char * buffer, uint8_t sz) {
    LED_Toggle();
}

unsigned char MakePingPkt(unsigned char *buffer) {
    buffer[0] = PKT_PING;
    for (char i = 0; i < ADDR_LEN; i++) {
        buffer[i + 1] = config.Address[i];
    }
    return ADDR_LEN + 1;
}
