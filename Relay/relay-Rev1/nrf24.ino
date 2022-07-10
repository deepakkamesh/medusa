
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
  radio.setDataRate(RF24_2MBPS);
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


  radio.startListening();
  return 1;
}


void RadioLoop() {
  uint8_t pipeNum;
  if (!radio.available(&pipeNum)) {
    return;
  }

  int sz = radio.getDynamicPayloadSize();
  radio.read(&bufferRX, sizeof(bufferRX));

  SendRadioPacket(pipeNum, bufferRX, sz);

  PrintPkt("RadioPkt", bufferRX, sz);
}


void SendNetPacket(uint8_t pipeNum, uint8_t * data, uint8_t sz) {
  radio.writeAckPayload(pipeNum, data, sz);
}
