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
    LoadAddrFromEE();
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
    nrf24_write_register(NRF24_MEM_RF_CH, DEFAULT_RF_CHANNEL);
    // RF_PWR=0bDm, RF_DR_HIGH=250kbps.
    nrf24_write_register(NRF24_MEM_RF_SETUP, 0b00100110);
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
    memcpy(config.Address, BoardAddress, ADDR_LEN);
    config.ARD = DEFAULT_ARD;

    // Initialize Transmit buffer.
    initQ(&TXQueue);
}
// DiscoverRFChannel iterates through the list of RF channels and attempts to 
// send a packet. If there is a response on the default pipe address it returns
// the RF channel.

uint8_t DiscoverRFChannel(void) {
    bufferTX[0] = PKT_NOOP;
    SuperMemCpy(bufferTX, 1, BoardAddress, 0, ADDR_LEN);

    for (uint8_t rf = 0; rf < 125; rf++) {
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
        return rf;
    }
    // Channel not found. retry.
    __delay_ms(5000);
    RESET();
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
    // becomes < FAILED_PERCENT.
    if (sentPktCnt == FAILURE_SAMPLE_RATE) {
        isRelayAvail = true;
        if ((float) failedPktCnt / (float) sentPktCnt >= FAILED_PERCENT) {
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
        case PKT_CFG_1:
            config.RFChannel = buffer[4];
            SuperMemCpy(config.PipeAddr1, 0, buffer, 5, PIPE_ADDR_LEN);
            config.ARD = buffer[10];
            break;
        case PKT_CFG_2:
            SuperMemCpy(config.Address, 0, buffer, 4, ADDR_LEN);
            config.PingInterval = buffer[7];
            break;
        default:
            SendError(ERR_UNKNOWN_PKT_TYPE);
    }
}

void ProcessActionRequest(uint8_t actionID, uint8_t * data) {
    uint8_t tmpHumidity[] = {0, 0};

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
            GetMockTempHumidity(tmpHumidity);
            SendData(ACTION_GET_TEMP_HUMIDITY, tmpHumidity, 2);
            break;
#else
            SendError(ERR_NOT_IMPL);
#endif
        case ACTION_RELOAD_CONFIG:
            ReloadConfig();
            break;
        case ACTION_RESET_DEVICE:
            RESET();
            break;
        case ACTION_TEST:
            TestFunc();
            break;
        default:
            SendError(ERR_NOT_IMPL);
    }
}

/* ReloadConfig loads the config in the config struct and makes it active */
void ReloadConfig(void) {
    config.IsConfigured = true;
    nrf24_write_register(NRF24_MEM_RF_CH, config.RFChannel);

    nrf24_write_buff(NRF24_MEM_TX_ADDR, config.PipeAddr1, PIPE_ADDR_LEN);
    nrf24_write_buff(NRF24_MEM_RX_ADDR_P0, config.PipeAddr1, PIPE_ADDR_LEN);

    uint8_t ard = (config.ARD << 4) | 0xF;
    nrf24_write_register(NRF24_MEM_SETUP_RETR, ard);

    memcpy(BoardAddress, config.Address, ADDR_LEN);
    WriteAddrToEE(); // Save the address to EEPROM.
    PingInterval = config.PingInterval;
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

void LoadAddrFromEE(void) {
    for (uint8_t i = 0; i < ADDR_LEN; i++) {
        BoardAddress[i] = DATAEE_ReadByte(EEPROM_ADDR + i);
    }
}

void WriteAddrToEE(void) {
    for (uint8_t i = 0; i < ADDR_LEN; i++) {
        DATAEE_WriteByte(EEPROM_ADDR + i, BoardAddress[i]);
    }
}

/***************************** Utility Functions *****************************/

void SuperMemCpy(uint8_t *dest, uint8_t destStart, uint8_t *src, uint8_t srcStart, uint8_t sz) {
    for (uint8_t i = 0; i < sz; i++) {
        dest[i + destStart] = src[i + srcStart];
    }
}

void TestFunc(void) {
    uint8_t buff[] = {0};
    SendData(ACTION_TEST, buff, 1);
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
