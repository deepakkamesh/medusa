#define UDP_LISTEN_PORT 6069 // default listen port.
#define PKT_TYPE_RELAY_GET_CONFIG 0xAA // Packet type for Relay Get Config
#define PKT_TYPE_RELAY_CONFIG_ANS 0xAB // Packet type for Relay Get Answer
#define PKT_TYPE_RELAY_ERROR 0xAC // Packet type for Relay  error. 
#define PKT_TYPE_RELAY_DATA 0xAD // Packet type for Relay  data packet. 

#define RELAY_CONFIG_ANS_LEN 16 // Length of  Relay Get Answer

#define ERROR_RELAY_RADIO_INIT_FAILED 0x02

WiFiUDP Udp;
uint8_t bufferRX[255];
uint8_t bufferTX[255];

// TODO: get from wifimanager.
const char* ssid     = "utopia";         // The SSID (name) of the Wi-Fi network you want to connect to
const char* password = "0d9f48a148";     // The password of the Wi-Fi network
const char* controllerHost = "192.168.1.116";
const uint16_t controllerPort = 6000;

void WifiConnect(void) {
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) { // Wait for the Wi-Fi to connect
    delay(1000);
    digitalWrite(LED_ONBOARD, !digitalRead(LED_ONBOARD));
  }
  digitalWrite(LED_ONBOARD, true); // For some reason true turns this off.
  WiFi.setAutoReconnect(true);
  WiFi.persistent(true);

#ifdef DEBUG
  Serial.println("Connection established!");
  Serial.print("IP address:\t");
  Serial.println(WiFi.localIP());
#endif

  Udp.begin(UDP_LISTEN_PORT);
}

int RelaySetup(void) {
  if (Config.isConfigured) {
    return 1;
  }

  // Build Relay Get Config Packet.
  bufferTX[0] = PKT_TYPE_RELAY_GET_CONFIG;
  uint8_t mac[6];
  WiFi.macAddress(mac);
  for (int i = 0; i < 6; i++) {
    bufferTX[i + 1] = mac[i];
  }

  while (1) {
    delay(1000);

    // Send config packet request.
    int ok = Udp.beginPacket(controllerHost, controllerPort);
    if (!ok) {
      return 0;
    }
    Udp.write(bufferTX, 7);
    ok = Udp.endPacket();
    if (!ok) {
      return 0;
    }
    delay(10);

    // See if we got a return.
    int packetSize = Udp.parsePacket();
    if (!packetSize) {
      continue;
    }

    int len = Udp.read(bufferRX, 255);
    // Valid packet?
    if (!(bufferRX[0] == PKT_TYPE_RELAY_CONFIG_ANS && len == RELAY_CONFIG_ANS_LEN)) {
      continue;
    }

    // Parse config packet and store in Config struct.
    for (int i = 0; i < 5; i++) {
      Config.pipe_addr_p0[i] = bufferRX[i + 1];
    }
    for (int i = 0; i < 5; i++) {
      Config.pipe_addr_p1[i] = bufferRX[i + 6];
    }
    Config.pipe_addr_p2[0] = bufferRX[11];
    Config.pipe_addr_p3[0] = bufferRX[12];
    Config.pipe_addr_p4[0] = bufferRX[13];
    Config.pipe_addr_p5[0] = bufferRX[14];
    Config.nrf24Channel = bufferRX[15];

#ifdef DEBUG
    Serial.printf("Mac: %02X:%02X:%02X:%02X:%02X:%02X\n", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
    PrintPkt(Config.pipe_addr_p0, 5);
    PrintPkt(Config.pipe_addr_p1, 5);
    Serial.printf("%02X\n", Config.pipe_addr_p2[0]);
    Serial.printf("%02X\n", Config.pipe_addr_p3[0]);
    Serial.printf("%02X\n", Config.pipe_addr_p4[0]);
    Serial.printf("%02X\n", Config.pipe_addr_p5[0]);
    Serial.printf("%02X\n", Config.nrf24Channel);
#endif
    return 1;
  }

}



void IpLoop(void) {


  int packetSize = Udp.parsePacket(); // check if there is a packet.
  if (!packetSize) {
    return;
  }

  // receive incoming UDP packets
  Serial.printf("Received %d bytes from %s, port %d\n", packetSize, Udp.remoteIP().toString().c_str(), Udp.remotePort());
  int len = Udp.read(bufferRX, 255);
  if (len > 0)
  {
    bufferRX[len] = 0;
  }
  Serial.printf("UDP packet contents: %s\n", bufferRX);


  /*
    char reply[] = {1, 2, 0xFE, 4};

    // Send return packet
    UDP.beginPacket("192.168.1.116", 6000);
    UDP.write(reply);
    UDP.endPacket();
    delay(1000);
  */
}


void PrintPkt(uint8_t buff[], int len) {
  for (int i = 0; i < len; i++) {
    Serial.printf("%02X,", buff[i]);
  }
  Serial.println();
}

int SendError(uint8_t errorCode) {
  // Send config packet request.
  bufferTX[0] = PKT_TYPE_RELAY_ERROR;
  bufferTX[1] = errorCode;
  int ok = Udp.beginPacket(controllerHost, controllerPort);
  if (!ok) {
    return 0;
  }
  Udp.write(bufferTX, 2);
  ok = Udp.endPacket();
  if (!ok) {
    return 0;
  }
  return 1;
}

int SendRadioPacket(uint8_t pipeNum, uint8_t * buff, uint8_t sz) {
  bufferTX[0] = PKT_TYPE_RELAY_DATA;
  bufferTX[1] = pipeNum;
  for (int i = 0; i < sz; i++) {
    bufferTX[i + 2] = buff[i];
  }

  int ok = Udp.beginPacket(controllerHost, controllerPort);
  if (!ok) {
    return 0;
  }
  Udp.write(bufferTX, sz + 2);
  ok = Udp.endPacket();
  if (!ok) {
    return 0;
  }
  return 1;
}
