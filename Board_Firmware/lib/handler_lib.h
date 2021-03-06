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
#define ERR_UNKNOWN_PKT_TYPE 0x05

// Action Types.
#define ACTION_STATUS_LED 0x13
#define ACTION_RESET_DEVICE 0x14
#define ACTION_RELOAD_CONFIG 0x15
#define ACTION_TEST 0x16
#define ACTION_GET_TEMP_HUMIDITY 0x02
#define ACTION_GET_LIGHT 0x03

#define PIPE_ADDR_LEN 5 
#define DEFAULT_RF_CHANNEL 115
#define DEFAULT_ARD 0xA // default ARD setting. (val*250 +250)
#define FAILURE_SAMPLE_RATE 10 // Number of packets to count.
#define FAILED_PERCENT 0.80 // Percent of failed packets. 
uint8_t DEFAULT_PIPE_ADDR[] = "hello"; // Default pipe address to bootstrap.
uint8_t PingInterval = 2; // Default ping interval.
uint8_t BoardAddress[3] = {0xFF, 0xFF, 0xFF}; // Default board address.
bool isRelayAvail = true; // Relay comms. If false dont sent any but ping packets.
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
void TestFunc(void);
void HandleTimeLoop(void);
void InitHandlerLib(void); // Main Library init routine. To be called in setup.
void HandlerLoop(void); // Main loop. To be called in loop.

// EEPROM stuff.
#define EEPROM_ADDR 0x10
void LoadAddrFromEE(void);
void WriteAddrToEE(void);

// Queue Functions.
typedef struct {
    uint8_t packet[MAX_PKT_SZ];
    uint8_t size; // Size of packet.
} Packet;

typedef struct {
    Packet packets[MAX_TX_QUEUE_SZ];
    uint8_t readPtr;
    uint8_t writePtr;
    uint8_t overflow;
} Queue;

void initQ(Queue *q);
void enQueue(uint8_t *buf, uint8_t sz, Queue *q);
uint8_t deQueue(uint8_t *buff, Queue *q);

// Config stores the configuration of the board.
struct Config {
    bool IsConfigured; // True if board is configured.
    uint8_t Address[3]; // Address of the board.
    uint8_t PingInterval; // Ping interval in seconds
    uint8_t RFChannel; // Frequency channel.
    uint8_t PipeAddr1[5];
    uint8_t ARD; // Auto Retry Duration. 
};