#define UDP_LISTEN_PORT 6069 // default listen port.

#define RELAY_CONFIG_ANS_LEN 19 // Length of  Relay Get Answer
#define MAC_ADDR_LEN 6



WiFiUDP Udp;
uint8_t bufferRX[255];
uint8_t bufferTX[255];

// TODO: get from wifimanager.
const char* ssid     = "utopia";         // The SSID (name) of the Wi-Fi network you want to connect to
const char* password = "0d9f48a148";     // The password of the Wi-Fi network
char* controllerHost = "192.168.1.108";
uint16_t controllerPort = 6000;

/* WifiConnect establishes connection to the specified access point*/
void WifiConnect(void) {
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  while (WiFi.waitForConnectResult() != WL_CONNECTED) {
    Serial.println("Connection Failed! Rebooting...");
    digitalWrite(LED_ONBOARD, !digitalRead(LED_ONBOARD));
    delay(3000);
    ESP.restart();
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

/* RelaySetup retrieves the configuration for the relay */
int RelaySetup(void) {
  if (Config.isConfigured) {
    return 1;
  }

  // Build Relay Get Config Packet.
  bufferTX[0] = PKT_TYPE_RELAY_GET_CONFIG;
  uint8_t mac[6];
  WiFi.macAddress(mac);
  SuperMemCpy(bufferTX, 1, mac, 0, MAC_ADDR_LEN);
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

    // See if we got a response.
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
    /* Copy pipe address into config */
    SuperMemCpy(Config.pipe_addr[0], 0, bufferRX, 1, PIPE_ADDR_LEN);
    SuperMemCpy(Config.pipe_addr[1], 0, bufferRX, 6, PIPE_ADDR_LEN);

    Config.pipe_addr[2][0] = bufferRX[11];
    SuperMemCpy(Config.pipe_addr[2], 1, bufferRX, 7, PIPE_ADDR_LEN - 1);

    Config.pipe_addr[3][0] = bufferRX[12];
    SuperMemCpy(Config.pipe_addr[3], 1, bufferRX, 7, PIPE_ADDR_LEN - 1);

    Config.pipe_addr[4][0] = bufferRX[13];
    SuperMemCpy(Config.pipe_addr[4], 1, bufferRX, 7, PIPE_ADDR_LEN - 1);

    Config.pipe_addr[5][0] = bufferRX[14];
    SuperMemCpy(Config.pipe_addr[5], 1, bufferRX, 7, PIPE_ADDR_LEN - 1);

    SuperMemCpy(Config.pipe_addr[6], 0, pipe_addr_6, 0, PIPE_ADDR_LEN); // Virtual pipe address #6.

    /* Copy pipe address into config */
    Config.nrf24Channel = bufferRX[15];

    /* Copy vBoard address into config */
    SuperMemCpy(Config.vboard_addr, 0, bufferRX, 16, PIPE_ADDR_LEN);

#ifdef DEBUG
    Serial.printf("Mac: %02X:%02X:%02X:%02X:%02X:%02X\n", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
    for (int i = 0; i < PIPE_ADDR_NUM; i++) {
      PrintPkt("Addr:", Config.pipe_addr[i], PIPE_ADDR_LEN);
    }
    Serial.printf("CHANNEL: %02X\n", Config.nrf24Channel);
    PrintPkt("vBoard:", Config.vboard_addr, ADDR_LEN);
#endif
    return 1;
  }
}

/* IpLoop parses UDP packets and sends it over Radio or locally for vboard */
void IpLoop(void) {

  int packetSize = Udp.parsePacket(); // check if there is a packet.
  if (!packetSize) {
    return;
  }

  int sz = Udp.read(bufferRX, 255);
  PrintPkt("Ctrl pkt:", bufferRX, sz);

  uint8_t radioPkt[32];
  uint8_t pipe_addr[PIPE_ADDR_LEN];

  // Process packets and send to radio or process locally if pipe_addr is for virtual pipe #6.
  uint8_t pktType = bufferRX[0];
  if (pktType !=  PKT_TYPE_BOARD_DATA_RELAY)
  {
    SendError(ERROR_UNKNOWN_PKT);
    return;
  }

  SuperMemCpy(pipe_addr, 0, bufferRX, 1, PIPE_ADDR_LEN);
  uint8_t radioPktSz = sz - 1 - PIPE_ADDR_LEN;
  SuperMemCpy(radioPkt, 0, bufferRX, PIPE_ADDR_LEN + 1 , radioPktSz);

  // pipe number determines if this is sent to radio or processed locally.
  int pipeNum = FindPipeNum(pipe_addr, PIPE_ADDR_LEN);
  switch (pipeNum) {
    case -1:
      SendError(ERROR_PIPE_ADDR_404);
      break;
    case VIRT_PIPE:
      ProcessVBoardPacket(radioPkt, radioPktSz);
      break;
    default:
      int ok = SendNetPacket(pipeNum, radioPkt, radioPktSz) ;
      if (!ok) {
        SendError(ERROR_RELAY_ACK_PAYLOAD_LOAD);
      }
      break;
  }
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

/* SendRadioPacket sends the Radio packet on UDP */
int SendRadioPacket( uint8_t pipeNum, uint8_t  buff[], uint8_t sz) {
  uint8_t i = 0;
  bufferTX[i] = PKT_TYPE_BOARD_DATA_RELAY;
  i++;
  SuperMemCpy(bufferTX, i, Config.pipe_addr[pipeNum], 0, PIPE_ADDR_LEN);
  i += PIPE_ADDR_LEN;
  SuperMemCpy(bufferTX, i, buff, 0, sz);
  i += sz;
  int ok = Udp.beginPacket(controllerHost, controllerPort);
  if (!ok) {
    return 0;
  }
  Udp.write(bufferTX, i);
  ok = Udp.endPacket();
  if (!ok) {
    return 0;
  }
  return 1;
}
