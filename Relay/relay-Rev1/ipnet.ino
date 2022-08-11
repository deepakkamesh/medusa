#define UDP_LISTEN_PORT 6069 // default listen port.

#define RELAY_CONFIG_ANS_LEN 19 // Length of  Relay Get Answer
#define MAC_ADDR_LEN 6



WiFiUDP Udp;
uint8_t bufferRX[255];
uint8_t bufferTX[255];

// TODO: get from wifimanager.
const char* ssid     = "utopia";         // The SSID (name) of the Wi-Fi network you want to connect to
const char* password = "0d9f48a148";     // The password of the Wi-Fi network
//const char* ssid     = "Utopian";         // The SSID (name) of the Wi-Fi network you want to connect to
//const char* password = "moretti308!!";     // The password of the Wi-Fi network

char* controllerHost = "192.168.1.255";
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

    int ok = NetSend(bufferTX, 7);
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
    if (!ParseConfigPkt(bufferRX, len)) {
      continue;
    }
    return 1;
  }
}


/******* IpLoop parses UDP packets and sends it over Radio or locally for vboard **********/
void IpLoop(void) {

  int packetSize = Udp.parsePacket(); // check if there is a packet.
  if (!packetSize) {
    return;
  }

  int sz = Udp.read(bufferRX, 255);
  PrintPkt("Ctrl pkt:", bufferRX, sz);

  // Process packets and send to radio or process locally if pipe_addr is for virtual pipe #6.
  uint8_t pktType = bufferRX[0];
  switch (pktType) {
    case PKT_TYPE_BOARD_DATA_RELAY:
      ProcessBoardDataRelay(bufferRX, sz);
      break;
    case PKT_TYPE_RELAY_CONFIG_ANS:
      ParseConfigPkt(bufferRX, sz);
      RadioSetup();
      break;
    default:
      SendError(ERROR_UNKNOWN_PKT);
      break;
  }
  yield();
}

void ProcessBoardDataRelay(uint8_t *bufferRX, uint8_t sz) {
  uint8_t radioPkt[32];
  uint8_t pipe_addr[PIPE_ADDR_LEN];

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

  return NetSend(bufferTX,2);
}

/********** SendRadioPacket sends the Radio packet on UDP *****************/
int SendRadioPacket( uint8_t pipeNum, uint8_t  buff[], uint8_t sz) {
  uint8_t i = 0;
  bufferTX[i] = PKT_TYPE_BOARD_DATA_RELAY;
  i++;
  SuperMemCpy(bufferTX, i, Config.pipe_addr[pipeNum], 0, PIPE_ADDR_LEN);
  i += PIPE_ADDR_LEN;
  SuperMemCpy(bufferTX, i, buff, 0, sz);
  i += sz;

  return NetSend(bufferTX,i);
}
