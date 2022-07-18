/*
 * File:   nrf24_lib.h
 * Author: Deepak Guruswamy (deepak.kamesh@gmail.com)
 * Author: [Adapted from Noyel Seth (noyelseth@gmail.com) 
 */

#ifndef NRF24_LIB_H
#define	NRF24_LIB_H
#include "master_exports.h"
// NRF24L01 Operation Modes
typedef enum {
    RX_MODE = 1,
    TX_MODE = 2
}NRF24_OPERATION_MODE;


/******************************************************************************/
// Register Map.
/******************************************************************************/

#define R_REGISTER          0x00  
#define W_REGISTER          0x20  
#define R_RX_PAYLOAD        0x61  
#define W_TX_PAYLOAD        0xA0  
#define FLUSH_TX            0xE1  
#define FLUSH_RX            0xE2  
#define REUSE_TX_PL         0xE3  
#define R_RX_PL_WID         0x60
#define W_TX_PAYLOAD_NOACK  0xB0

#define NRF24_MEM_CONFIG              0x00  
#define NRF24_MEM_EN_AA               0x01  
#define NRF24_MEM_EN_RXADDR           0x02  
#define NRF24_MEM_SETUP_AW            0x03  
#define NRF24_MEM_SETUP_RETR          0x04  
#define NRF24_MEM_RF_CH               0x05  
#define NRF24_MEM_RF_SETUP            0x06  
#define NRF24_MEM_STATUSS             0x07  
#define NRF24_MEM_OBSERVE_TX          0x08  
#define NRF24_MEM_CD                  0x09  
#define NRF24_MEM_RX_ADDR_P0          0x0A  
#define NRF24_MEM_RX_ADDR_P1          0x0B  
#define NRF24_MEM_RX_ADDR_P2          0x0C  
#define NRF24_MEM_RX_ADDR_P3          0x0D  
#define NRF24_MEM_RX_ADDR_P4          0x0E  
#define NRF24_MEM_RX_ADDR_P5          0x0F  
#define NRF24_MEM_TX_ADDR             0x10  
#define NRF24_MEM_RX_PW_P0            0x11  
#define NRF24_MEM_RX_PW_P1            0x12  
#define NRF24_MEM_RX_PW_P2            0x13  
#define NRF24_MEM_RX_PW_P3            0x14  
#define NRF24_MEM_RX_PW_P4            0x15  
#define NRF24_MEM_RX_PW_P5            0x16  
#define NRF24_MEM_FIFO_STATUS         0x17
#define NRF24_MEM_CMD_NOP             0xFF // No operation (used for reading status register)
#define NRF24_MEM_DYNPD               0x1C
#define NRF24_MEM_FEATURE             0x1D
#define NRF24_MEM_REGISTER_MASK       0x1F


/******************************************************************************/
// NRF24L01 Functions.
/******************************************************************************/

/**
 * @brief  Write or update value into the Memory Address of NRF24L01
 *
 * @param[in]	mnemonic_addr NRF24L01 Memory Address
 * @param[in]	value NRF24L01 Memory Address's value
 *
 */
void nrf24_write_register(uint8_t mnemonic_addr, uint8_t value);

/**
 * @brief  Read value from the Memory Address of NRF24L01
 *
 * @param[in]	mnemonic_addr NRF24L01 Memory Address
 * 
 * @return      current value of the read NRF24L01 Memory Address
 *
 */
uint8_t nrf24_read_register(uint8_t mnemonic_addr);

/**
 * @brief  Write or update buffer into the Memory Address of NRF24L01
 *
 * @param[in]	mnemonic_addr NRF24L01 Memory Address
 * @param[in]	buffer write buffer data for NRF24L01 Memory Address's value
 * @param[in]	bytes size of the write buffer
 *
 */
void nrf24_write_buff(uint8_t mnemonic_addr, uint8_t *buffer, uint8_t bytes);

/**
 * @brief  Read buffer data from the Memory Address of NRF24L01
 *
 * @param[in]	mnemonic_addr NRF24L01 Memory Address
 * @param[out]	buffer read buffer for read NRF24L01 Memory Address's data
 * @param[in]	bytes size of the read buffer
 *
 */
void nrf24_read_buff(uint8_t mnemonic_addr, uint8_t *buffer, uint8_t bytes);

/**
 * @brief  Initialize NRF24L01 to setup CSN, CE and SPI
 *
 *
 */
void nrf24_rf_init(void);

/**
 * @brief  Set the NRF24L01 Operation Mode Tx/Rx
 *
 * @param[in]	mode NRF24L01 Tx/Rx Operation Mode
 *
 */
void nrf24_set_rf_mode(NRF24_OPERATION_MODE mode);

/**
 * @brief  Send Payload to NRF24L01
 *
 * @param[in]	buffer NRF24L01 Send Payload buffer pointer
 * @param[in]   sz is the size of buffer
 *
 */
void nrf24_send_rf_data(uint8_t *buffer, uint8_t sz);

/**
 * @brief  Receive data availability check
 * 
 * @return  1 if data present
 *          0 if data not present
 *
 */
uint8_t nrf24_is_rf_data_available(void);

/**
 * @brief  Read payload from NRF24L01
 *
 * @param[in]	buffer pointer buffer for receive payload
 * @param[in]	sz number of bytes to read
 *  
 */
void nrf24_read_rf_data(uint8_t *buffer, uint8_t sz);

/**
 * @brief  Set NRF24L01 RF Channel Frequency
 *
 * @param[in]	rf_channel NRF24L01 RF channel - radio frequency channel, value from 0 to 127
 * 
 * @Note: NRF24L01 frequency will be (2400 + rf_channel)GHz
 *
 */
void nrf24_set_channel_frq(uint8_t rf_channel);

/**
 * @brief  Get NRF24L01 RF Channel Frequency
 *
 * @return rf_channel NRF24L01 RF channel - radio frequency channel, value from 0 to 127
 * 
 * @Note: For actual NRF24L01 frequency will be (2400 + rf_channel)GHz
 *
 */
uint8_t nrf24_get_channel_frq(void);

/**
 * @brief  Set NRF24L01 into StandBy-I
 * 
 */
void nrf24_standby_I(void);

/**
 * @brief  Flush the Tx an Rx
 *
 */
void nrf24_flush_tx_rx(void);

/**
 * @brief read dynamic payload length of the received packet. 
 *
 */
uint8_t nrf24_read_dynamic_payload_length(void) ;

#endif	/* NRF24_LIB_H */

