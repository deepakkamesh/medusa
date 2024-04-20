#define MAX_RF_PKT_SZ 32
#define MAX_TX_QUEUE_SZ 8

typedef struct {
  uint8_t packet[MAX_RF_PKT_SZ];
  uint8_t pipeNum;
  uint8_t size; // Size of packet.
} RFPacket;

typedef struct {
  RFPacket packets[MAX_TX_QUEUE_SZ];
  uint8_t readPtr;
  uint8_t writePtr;
  uint8_t overflow;
} Queue;

Queue qTX;

void initQueue(Queue *q) ;
uint8_t enQueue(uint8_t *buf, uint8_t sz, uint8_t pipeNum, Queue *q);
uint8_t deQueue(uint8_t *buff, uint8_t *pipeNum, Queue *q) ;

/************ RadioSetup *****************/
int RadioSetup() {

  if (!radio.begin()) {
    return 0;
  }


  if (!radio.isChipConnected()) {
    return 0;
  }
  radio.stopListening();

  // Set default radio params.
  radio.setPALevel(RF24_PA_MAX,true);
  radio.setAddressWidth(5);
  radio.setDataRate(RF24_250KBPS);
  radio.setCRCLength(RF24_CRC_8);

  // Other: Enhanced Shockburt.
  radio.setAutoAck(true);
  radio.enableDynamicPayloads();
  radio.enableAckPayload();

  // Set params from config struct.
  radio.setChannel(Config.nrf24Channel);

  for (uint8_t ch = 0; ch < 6 ; ch++) {
    radio.openReadingPipe(ch, Config.pipe_addr[ch]);
  }

  radio.flush_tx();
  radio.flush_rx();

  radio.startListening();

  // Init. AckPacketQ.
  initQueue(&qTX);

  return 1;
}

/***************** RadioRcvLoop() *************/
void RadioRcvLoop() {
  uint8_t pipeNum;
  if (!radio.available(&pipeNum)) {
    return;
  }

  int sz = radio.getDynamicPayloadSize();
  radio.read(&bufferRX, sizeof(bufferRX));

  int ok =  SendRadioPacket(pipeNum, bufferRX, sz);
  // Try restart if network comms is broken.
  if (!ok) {
    delay(1000);
    ESP.restart();
  }
  PrintPkt("RadioPkt", bufferRX, sz);
}

/********** RadioSendLoop() ****************/
void RadioSendLoop() {
  uint8_t AckPkt[MAX_RF_PKT_SZ];
  uint8_t AckPktSz = 0;
  uint8_t pipeNum = 0;

  AckPktSz = deQueue(AckPkt, &pipeNum, &qTX);

  // Nothing to do.
  if (!AckPktSz) {
    return;
  }

  // Queue up ack payload to TX fifo.
  bool ok = radio.writeAckPayload(pipeNum, AckPkt, AckPktSz);
  if (!ok) {
    enQueue(AckPkt, AckPktSz, pipeNum, &qTX);
  }
  delay(10); // This delay is needed to prevent packet corruption!!.
}

int SendNetPacket(uint8_t pipeNum, uint8_t * data, uint8_t sz) {
  PrintPkt("Sending to Radio", data, sz);
  return enQueue(data, sz, pipeNum, &qTX);
}
