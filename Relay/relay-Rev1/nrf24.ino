
const byte address[6] = "hello";
//const byte address[6] = {0xe1, 0xe1, 0xe1, 0xe1, 0xe1};
char ack[10] = "hell";
int cnt = 0;

void RadioSetup(RF24 radio) {

  if (!radio.begin()) {
    Serial.println("Failed");
    while (1) {};
  }

  Serial.println("success");
  radio.setAutoAck(true);

  if (radio.isChipConnected()) {
    Serial.println("yay");
  } else {
    Serial.println("nop");
  }

  radio.openReadingPipe(0, address);
  radio.setPALevel(RF24_PA_MIN);

  //PIC STUFF
  radio.setDataRate(RF24_2MBPS);
  radio.setChannel(115);
  radio.setCRCLength(RF24_CRC_8);

  radio.enableDynamicPayloads();
  radio.enableAckPayload();
  radio.writeAckPayload(0, ack, sizeof(ack));

  radio.startListening();
}


void RadioLoop(RF24 radio) {
  if (!radio.available()) {
    return;
  }
  char text[32] = "";
  radio.read(&text, sizeof(text));
  sprintf(ack, "Ack %d", cnt);
  radio.writeAckPayload(0, ack, sizeof(ack));

  for (int i = 0; i < strlen(text); i++) {
    Serial.print(text[i], HEX);
    Serial.print(" ");
  }
  Serial.println("");
  cnt++;

}
