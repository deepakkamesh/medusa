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


int RadioSetup() {

  if (!radio.begin()) {
    return 0;
  }

  if (!radio.isChipConnected()) {
    return 0;
  }

  // Set default radio params.
  radio.setPALevel(RF24_PA_MAX);
  radio.setAddressWidth(5);
  radio.setDataRate(RF24_250KBPS);
  radio.setCRCLength(RF24_CRC_8);

  // Other Enhanced Shockburt.
  radio.setAutoAck(true);
  radio.enableDynamicPayloads();
  radio.enableAckPayload();

  // Set params from config struct.
  radio.setChannel(Config.nrf24Channel);
  radio.openReadingPipe(0, Config.pipe_addr_p0);
  radio.openReadingPipe(1, Config.pipe_addr_p1);
  radio.openReadingPipe(2, Config.pipe_addr_p2);
  radio.openReadingPipe(3, Config.pipe_addr_p3);
  radio.openReadingPipe(4, Config.pipe_addr_p4);
  radio.openReadingPipe(5, Config.pipe_addr_p5);

  radio.flush_tx();
  radio.flush_rx();

  radio.startListening();

  // Init. AckPacketQ.
  initQueue(&qTX);

  return 1;
}

void RadioRcvLoop() {
  uint8_t pipeNum;
  if (!radio.available(&pipeNum)) {
    return;
  }

  int sz = radio.getDynamicPayloadSize();
  radio.read(&bufferRX, sizeof(bufferRX));

  bool okSig = radio.testRPD(); // returns true is strength > -64bdM.
  int ok =  SendRadioPacket(okSig, pipeNum, bufferRX, sz);
  // Try restart if network comms is broken.
  if (!ok) {
    delay(1000);
    ESP.restart();
  }

  PrintPkt("RadioPkt", bufferRX, sz);
}

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
  return enQueue(data, sz, pipeNum, &qTX);
}
