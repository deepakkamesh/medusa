#include "master_exports.h"

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
#define ERR_NA 0x00
#define ERR_NOT_IMPL 0x04

// Action Types.
#define ACTION_STATUS_LED 0x13
#define ACTION_RESET_DEVICE 0x14
#define ACTION_RELOAD_CONFIG 0x15
#define ACTION_GET_TEMP_HUMIDITY 0x02
#define ACTION_GET_LIGHT 0x03

#define PIPE_ADDR_LEN 5 
#define DEFAULT_RF_CHANNEL 115
#define DEFAULT_ARD 0xA // default ARD setting. (val*250 +250)
uint8_t DEFAULT_PIPE_ADDR[] = "hello"; // Default pipe address to bootstrap.
uint8_t PingInterval = 1; // Default ping interval.
uint8_t BoardAddress[3] = {0xFF, 0xFF, 0xFF}; // Default board address.

void TimerInterruptHandler(void);
void InitRadio(void);
void ProcessAckPayload(uint8_t * buffer, uint8_t sz);
void ProcessActionRequest(uint8_t actionID, uint8_t * data);
bool VerifyBoardAddress(uint8_t *bufferRX);
void HandlePacketLoop(void);
void SendError(uint8_t errorCode);
void SendPing(void);
void SuperMemCpy(uint8_t *dest, uint8_t destStart, uint8_t *src, uint8_t srcStart, uint8_t sz);
void ReloadConfig(void);
void SendData(uint8_t actionID, uint8_t *data, uint8_t dataSz);

typedef struct {
    uint8_t packet[MAX_PKT_SZ];
    bool free;
    uint8_t size; // Size of packet.
    uint32_t tmpstmp; // Time stamp of when packet was queued.
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