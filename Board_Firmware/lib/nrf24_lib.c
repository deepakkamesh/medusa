/*
 * File:   nrf24_lib.c
 * Author: Noyel Seth (noyelseth@gamil.com)
 */

#include "nrf24_lib.h"

void nrf24_write_register(unsigned char mnemonic_addr, unsigned char value) {
    NRF24L01_CSN_L();
    if (mnemonic_addr < W_REGISTER) {
        // This is a register access
        SPI_WRITE_BYTE(W_REGISTER | (mnemonic_addr & NRF24_MEM_REGISTER_MASK));
        SPI_WRITE_BYTE(value);
    } else {
        // This is a single byte command or future command/register
        SPI_WRITE_BYTE(mnemonic_addr);
        if ((mnemonic_addr != FLUSH_TX) && (mnemonic_addr != FLUSH_RX) && \
				(mnemonic_addr != REUSE_TX_PL) && (mnemonic_addr != NRF24_MEM_CMD_NOP)) {
            // Send register value
            SPI_WRITE_BYTE(value);
        }
    }
    __delay_us(10);
    NRF24L01_CSN_H();
}

unsigned char nrf24_read_register(unsigned char mnemonic_addr) {
    unsigned char byte0;
    NRF24L01_CSN_L();
    SPI_WRITE_BYTE(R_REGISTER | (mnemonic_addr & NRF24_MEM_REGISTER_MASK));
    byte0 = SPI_READ_BYTE(NRF24_MEM_CMD_NOP);
    NRF24L01_CSN_H();
    return byte0;
}

uint8_t nrf24_read_dynamic_payload_length(void) {
    uint8_t byte0;
    NRF24L01_CSN_L();
    SPI_WRITE_BYTE(R_RX_PL_WID);
    byte0 = SPI_READ_BYTE(NRF24_MEM_CMD_NOP);
    NRF24L01_CSN_H();
    return byte0;
}

void nrf24_write_buff(unsigned char mnemonic_addr, unsigned char *buffer, unsigned char bytes) {
    unsigned char i;
    NRF24L01_CSN_L();
    SPI_WRITE_BYTE(W_REGISTER | mnemonic_addr);
    for (i = 0; i < bytes; i++) {
        SPI_WRITE_BYTE(*buffer);
        buffer++;
        __delay_us(10);
    }
    NRF24L01_CSN_H();
}

void nrf24_read_buff(unsigned char mnemonic_addr, unsigned char *buffer, unsigned char bytes) {
    unsigned char i;
    NRF24L01_CSN_L();
    SPI_WRITE_BYTE(R_REGISTER | mnemonic_addr);
    for (i = 0; i < bytes; i++) {
        *buffer = SPI_READ_BYTE(NRF24_MEM_CMD_NOP);
        buffer++;
    }
    *buffer = (unsigned char) NULL;
    NRF24L01_CSN_H();
}

void nrf24_rf_init() {
    SPI_INIT();
    NRF24L01_CSN_SetOutput();
    NRF24L01_CE_SetOutput();
    NRF24L01_CSN_H();
    NRF24L01_CE_L();
}


void nrf24_send_rf_data(unsigned char *buffer,unsigned char sz) {
    nrf24_write_buff(W_TX_PAYLOAD, buffer, sz);
    NRF24L01_CE_H();
    __delay_ms(1);
    NRF24L01_CE_L();
}

unsigned char nrf24_is_rf_data_available(void) {
    if ((nrf24_read_register(NRF24_MEM_STATUSS) & 0x40) == 0x40) {
        return 0;
    }
    return 1;
}

void nrf24_read_rf_data(unsigned char *buffer,unsigned char sz) {
    nrf24_read_buff(R_RX_PAYLOAD, buffer, sz);
    nrf24_write_register(NRF24_MEM_STATUSS, 0x70); // Clear STATUS.
    nrf24_flush_tx_rx();
}

void nrf24_flush_tx_rx(void) {
    NRF24L01_CSN_L();
    nrf24_write_register(NRF24_MEM_STATUSS, 0x70);
    __delay_ms(10);
    NRF24L01_CSN_H();

    NRF24L01_CSN_L();
    SPI_WRITE_BYTE(FLUSH_TX);
    __delay_ms(10);
    NRF24L01_CSN_H();

    NRF24L01_CSN_L();
    SPI_WRITE_BYTE(FLUSH_RX);
    __delay_ms(10);
    NRF24L01_CSN_H();
}
