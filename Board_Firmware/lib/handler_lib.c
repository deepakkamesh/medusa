/*
 * File:   handler.c
 * Author: dkg
 *
 * Created on July 5, 2022, 8:38 PM
 */

#include <stdlib.h>
#include <string.h>
#include "handler_lib.h"

#include "nrf24_lib.h"
#include "dht11_lib.h"

uint32_t prevTicks = 0;
uint32_t Ticks = 0; // Ticks of timer.
struct Config config; // Board config.

uint8_t bufferTX[32];
uint8_t bufferRX[32];
uint8_t sentPktCnt = 0;
uint8_t failedPktCnt = 0;
Queue TXQueue; //Transmit Queue;

void InitHandlerLib(void) {
    LoadConfigFromEE();
    InitRadio();
    TMR1_SetInterruptHandler(TimerInterruptHandler);
    uint8_t rfChan = DiscoverRFChannel(); // Roughly 10sec delay to discover channel.
    config.RFChannel = rfChan;
}

void HandlerLoop(void) {
    HandlePacketLoop();
    HandleTimeLoop();
    NOP();
    CLRWDT();
}

void InitRadio(void) {
    nrf24_rf_init();

    nrf24_write_buff(NRF24_MEM_TX_ADDR, config.PipeAddr1, PIPE_ADDR_LEN);
    nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, config.PipeAddr1, PIPE_ADDR_LEN);
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
    nrf24_write_register(NRF24_MEM_RF_CH, DEFAULT_RF_CHANNEL);
    // RF_PWR=0bDm, RF_DR_HIGH=250kbps.
    nrf24_write_register(NRF24_MEM_RF_SETUP, 0b00100110);
    // EN_DPL, EN_ACK_PAY.
    nrf24_write_register(NRF24_MEM_FEATURE, 0b110); // Enable Dynamic payload, ack payload.
    // DPL_P0.
    nrf24_write_register(NRF24_MEM_DYNPD, 0b1); // Dynamic payload on Pipe 0.
    // Setup ARD.
    uint8_t ard = (config.ARD << 4) | 0xF;
    nrf24_write_register(NRF24_MEM_SETUP_RETR, ard);
    __delay_us(10);

    PingInterval = config.PingInterval;
    memcpy(BoardAddress, config.Address, ADDR_LEN);

    // Initialize Transmit buffer.
    initQ(&TXQueue);
}
// DiscoverRFChannel iterates through the list of RF channels and attempts to 
// send a packet. If no response after set retries it tries the default pipe address.

uint8_t DiscoverRFChannel(void) {
    bufferTX[0] = PKT_NOOP;
    SuperMemCpy(bufferTX, 1, BoardAddress, 0, ADDR_LEN);

    // If board not configured no point in trying to flip addresses. 
    if (config.IsConfigured) {
        FlipPipeAddress();
    }

    // Cycle through available channels. 
    for (uint8_t rf = 0; rf < 125; rf++) {
#ifdef DEV_STATUS_LED
        LED_Toggle();
#endif
        nrf24_write_register(NRF24_MEM_RF_CH, rf);
        nrf24_send_rf_data(bufferTX, ADDR_LEN + 1);
        __delay_us(10);

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
            continue;
        }
#ifdef DEV_STATUS_LED
        LED_SetLow();
#endif
        // Found the channel. Reset Flip counter. 
        ResetFlipCounter();
        return rf;
    }
    // Channel not found. retry.
    __delay_ms(5000);
    RESET();
    return 0;
}

// FlipPipeAddress flips the pipe address between whats in memory to the default
// pipe address after MAX_CONNECT_RETRIES. 

void FlipPipeAddress(void) {
    uint8_t n = DATAEE_ReadByte(EEPROM_ADDR + EE_RETRY_OFFSET);
    DATAEE_WriteByte(EEPROM_ADDR + EE_RETRY_OFFSET, n + 1);

    if (n < MAX_CONNECT_RETRIES) {
        return;
    }

    if (n % 2 == 0) {
        nrf24_write_buff(NRF24_MEM_TX_ADDR, config.PipeAddr1, PIPE_ADDR_LEN);
        nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, config.PipeAddr1, PIPE_ADDR_LEN);
    } else {
        nrf24_write_buff(NRF24_MEM_TX_ADDR, DEFAULT_PIPE_ADDR, PIPE_ADDR_LEN);
        nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, DEFAULT_PIPE_ADDR, PIPE_ADDR_LEN);
    }
}

void ResetFlipCounter(void) {
    DATAEE_WriteByte(EEPROM_ADDR + EE_RETRY_OFFSET, 0);
}

void HandleTimeLoop(void) {

    uint32_t currTicks = Ticks;
    if (currTicks == prevTicks) {
        return;
    }
    prevTicks = currTicks;

    // Perform any timed activity here.
    if (currTicks % PingInterval == 0) {
        SendPing();
    }
}

// Sends packets 1 at a time.

void HandlePacketLoop(void) {
    uint8_t TXPacket[MAX_PKT_SZ];
    uint8_t TXPktSz = 0;

    TXPktSz = deQueue(TXPacket, &TXQueue);

    // Check queue; if nothing sleep.
    if (TXPktSz == 0) {
        SLEEP();
        return;
    }

    // If in the last FAILURE_SAMPLE_RATE packet the failure exceeds FAILED_PERCENT
    // Relay availability is marked down. Only ping packets are sent until it 
    // becomes < FAILED_PERCENT. If fail rate is 1.0 reset.  
    if (sentPktCnt == FAILURE_SAMPLE_RATE) {
        isRelayAvail = true;
        float failedRate = (float) failedPktCnt / (float) sentPktCnt;
        if (failedRate == 1.0) {
            RESET();
        } else if (failedRate >= FAILED_PERCENT) {
            isRelayAvail = false;
        }
        failedPktCnt = 0;
        sentPktCnt = 0;
    }

    nrf24_send_rf_data(TXPacket, TXPktSz);
    sentPktCnt++;
    __delay_us(10);
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
        failedPktCnt++;
        // Only retry packets if the packet failure rate is within limits to avoid
        // continuous running loop preventing sleep.
        if (isRelayAvail) {
            enQueue(TXPacket, TXPktSz, &TXQueue); // Send failed so enqueue packet.
        }
        return;
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
}

void TimerInterruptHandler(void) {
    Ticks++;
}

bool VerifyBoardAddress(uint8_t * buffer) {
    for (int i = 0; i < ADDR_LEN; i++) {
        if (BoardAddress[i] != buffer[i + 1]) {
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
        case PKT_CFG:
            config.ARD = buffer[4];
            config.PingInterval = buffer[5];
            SuperMemCpy(config.PipeAddr1, 0, buffer, 6, PIPE_ADDR_LEN);
            SuperMemCpy(config.Address, 0, buffer, 11, ADDR_LEN);
            WriteConfigToEE();
            break;
        default:
            SendError(ERR_UNKNOWN_PKT_TYPE);
    }
}

void ProcessActionRequest(uint8_t actionID, uint8_t * data) {
    uint8_t buff[2] = {0, 0};
    adc_result_t volts;

    switch (actionID) {
        case ACTION_STATUS_LED:
#ifdef DEV_STATUS_LED
            LED_SetLow();
            if (data[0]) {
                LED_SetHigh();
            }
            break;
#else
            SendError(ERR_NOT_IMPL);
#endif

        case ACTION_GET_TEMP_HUMIDITY:
#ifdef DEV_TEMP_HUMIDITY
            GetMockTempHumidity(buff);
            SendData(ACTION_GET_TEMP_HUMIDITY, buff, 2);
            break;
#else
            SendError(ERR_NOT_IMPL);
#endif

        case ACTION_RESET_DEVICE:
            RESET();
            break;

        case ACTION_GET_VOLTS:
            volts = ADC_GetConversion(channel_FVR);
            buff[0] = volts & 0x00FF;
            buff[1] = volts >> 8;
            SendData(ACTION_GET_VOLTS, buff, 2);
            break;

        case ACTION_TEST:
            TestFunc();
            break;

        default:
            SendError(ERR_NOT_IMPL);
    }
}

void SendError(uint8_t errorCode) {
    uint8_t i = 0;
    bufferTX[i] = PKT_DATA;
    SuperMemCpy(bufferTX, 1, BoardAddress, 0, ADDR_LEN);
    i += ADDR_LEN;
    bufferTX[++i] = 0; // ActionID.
    bufferTX[++i] = errorCode;
    enQueue(bufferTX, (i + 1), &TXQueue);
}

void SendData(uint8_t actionID, uint8_t *data, uint8_t dataSz) {
    uint8_t i = 0;
    bufferTX[i] = PKT_DATA;
    SuperMemCpy(bufferTX, 1, BoardAddress, 0, ADDR_LEN);
    i += ADDR_LEN;
    bufferTX[++i] = actionID;
    bufferTX[++i] = ERR_NA;
    SuperMemCpy(bufferTX, i + 1, data, 0, dataSz);
    i += dataSz;
    enQueue(bufferTX, (i + 1), &TXQueue);
}

void SendPing(void) {
    bufferTX[0] = PKT_PING;
    SuperMemCpy(bufferTX, 1, BoardAddress, 0, ADDR_LEN);
    enQueue(bufferTX, (ADDR_LEN + 1), &TXQueue);
}

void ResetEE(void) {
    unsigned int idx = EEPROM_ADDR + EE_CONFIG_OFFSET;
    DATAEE_WriteByte(idx, 0xFF);
}

/* LoadConfigFromEE loads the configuration from memory. If none found default loaded*/
void LoadConfigFromEE(void) {
    unsigned int idx = EEPROM_ADDR + EE_CONFIG_OFFSET;
    uint8_t isConfigured = DATAEE_ReadByte(idx);
    if (isConfigured != IS_CONFIGURED) {
        config.ARD = DEFAULT_ARD;
        config.PingInterval = PingInterval;
        memcpy(config.PipeAddr1, DEFAULT_PIPE_ADDR, PIPE_ADDR_LEN);
        memcpy(config.Address, BoardAddress, ADDR_LEN);
        config.IsConfigured = false;
        return;
    }
    config.IsConfigured = true;
    idx++;
    config.ARD = DATAEE_ReadByte(idx);
    idx++;
    config.PingInterval = DATAEE_ReadByte(idx);
    idx++;
    for (uint8_t i = 0; i < PIPE_ADDR_LEN; i++) {
        config.PipeAddr1[i] = DATAEE_ReadByte(idx + i);
    }
    idx += PIPE_ADDR_LEN;
    for (uint8_t i = 0; i < ADDR_LEN; i++) {
        config.Address[i] = DATAEE_ReadByte(idx + i);
    }
}

void WriteConfigToEE(void) {
    unsigned int idx = EEPROM_ADDR + EE_CONFIG_OFFSET;
    DATAEE_WriteByte(idx, IS_CONFIGURED);
    idx++;
    DATAEE_WriteByte(idx, config.ARD);
    idx++;
    DATAEE_WriteByte(idx, config.PingInterval);
    idx++;
    for (uint8_t i = 0; i < PIPE_ADDR_LEN; i++) {
        DATAEE_WriteByte(i + idx, config.PipeAddr1[i]);
    }
    idx += PIPE_ADDR_LEN;
    for (uint8_t i = 0; i < ADDR_LEN; i++) {
        DATAEE_WriteByte(i + idx, config.Address[i]);
    }
}

/***************************** Utility Functions *****************************/

void SuperMemCpy(uint8_t *dest, uint8_t destStart, uint8_t *src, uint8_t srcStart, uint8_t sz) {
    for (uint8_t i = 0; i < sz; i++) {
        dest[i + destStart] = src[i + srcStart];
    }
}

void TestFunc(void) {
    //LoadConfigFromEE();
    //memcpy(DEFAULT_PIPE_ADDR, BoardAddress, ADDR_LEN);
    uint8_t buff[2] = {0, 0};
    adc_result_t volts = ADC_GetConversion(channel_FVR);
    buff[0] = volts & 0x00FF;
    buff[1] = volts >> 8;
    SendData(ACTION_TEST, buff, 2);
}

/***************************** Queuing Functions *****************************/

void initQ(Queue * q) {
    q->readPtr = 0;
    q->writePtr = 0;
    q->overflow = 0;
    uint8_t i = 0;
    for (i = 0; i < MAX_TX_QUEUE_SZ; i++) {
        q->packets[i].size = 0;
    }
}

void enQueue(uint8_t *buf, uint8_t sz, Queue * q) {
    memcpy(q->packets[q->writePtr].packet, buf, sz);
    q->packets[q->writePtr].size = sz;
    q->writePtr++;
    if (q->writePtr == MAX_TX_QUEUE_SZ) {
        q->writePtr = 0;
        q->overflow = 1;
    }
}

uint8_t deQueue(uint8_t *buff, Queue * q) {
    if (!q->overflow && q->readPtr == q->writePtr) {
        return 0;
    }
    if (q->overflow && q->readPtr < q->writePtr) {
        q->readPtr = q->writePtr;
    }

    memcpy(buff, q->packets[q->readPtr].packet, q->packets[q->readPtr].size);
    uint8_t sz = q->packets[q->readPtr].size;
    q->readPtr++;

    if (q->readPtr == MAX_TX_QUEUE_SZ) {
        q->readPtr = 0;
        q->overflow = 0;
    }
    return sz;
}
