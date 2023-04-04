#define UDP_LISTEN_PORT 6069 // default listen port.

#define RELAY_CONFIG_ANS_LEN 20 // Length of  Relay Get Answer
#define MAC_ADDR_LEN 6


WiFiClient clientConn;
WiFiUDP Udp;
uint8_t bufferRX[255];
uint8_t bufferTX[255];

// TODO: get from wifimanager.
char ssid[20];
const char* ssidPrefix     = "utopia";         // The SSID (name) of the Wi-Fi network you want to connect to
const char* password = "0d9f48a148";     // The password of the Wi-Fi network

uint16_t ctrPort = 3345;
IPAddress ctrIP;

/* Find the strongest signal and connect */
void WifiAPFinder() {
  WiFi.mode(WIFI_STA);
  WiFi.disconnect();
  int n = WiFi.scanNetworks();
  if (n == 0) {
    Serial.println("no networks found");
    delay(2000);
    ESP.restart();
  }

  int strength = -80; // Assume some initial strength.
  bool sel = false;

  for (int i = 0; i < n; ++i) {
    if (!strstr(WiFi.SSID(i).c_str(), ssidPrefix)) {
      continue;
    }
    if (WiFi.RSSI(i) >  strength) {
      strength = WiFi.RSSI(i);
      strcpy(ssid, WiFi.SSID(i).c_str());
      sel = true;
    }
  }

  if (!sel) {
    Serial.printf("no network with prefix found %s or strength below -80", ssidPrefix);
    delay(2000);
    ESP.restart();
  }
  Serial.printf("\nWifi selected %s (%i)\n", ssid, strength);

}

/* WifiConnect establishes connection to the specified access point*/
void WifiConnect(void) {
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  int retryCnt = 0;

  while (WiFi.status() != WL_CONNECTED) {
    Serial.printf(".");
    digitalWrite(LED_ONBOARD, !digitalRead(LED_ONBOARD));
    delay(1000);
    if (retryCnt > 15) {
      ESP.restart();
    }
    retryCnt ++;
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
  IPAddress bcastIP = GetBroadcastIP();

  // Build Relay Get Config Packet.
  bufferTX[0] = PKT_TYPE_RELAY_GET_CONFIG;
  uint8_t mac[6];
  WiFi.macAddress(mac);
  SuperMemCpy(bufferTX, 1, mac, 0, MAC_ADDR_LEN);
  while (1) {
    digitalWrite(LED_ONBOARD, !digitalRead(LED_ONBOARD));

    int ok = NetSendUDP(bufferTX, 7, bcastIP, ctrPort);
    if (!ok) {
      return 0;
    }
    delay(200);

    // See if we got a response.
    int packetSize = Udp.parsePacket();
    if (!packetSize) {
      continue;
    }

    int len = Udp.read(bufferRX, 255);
    if (!ParseConfigPkt(bufferRX, len)) {
      continue;
    }
    digitalWrite(LED_ONBOARD, true); // For some reason true turns this off.

    // Get the controller IP.
    ctrIP = Udp.remoteIP();
    return 1;
  }
}


/******* IpLoop parses UDP packets and sends it over Radio or locally for vboard **********/
void IpLoop(void)  {
  // reconnect if connection is broken.
  if (!clientConn.connected()) {
    clientConn.stop();
    if (!clientConn.connect(ctrIP, ctrPort)) {
      // Restart if unable to connect to controller.
#ifdef DEBUG
      Serial.println("Failed trying to reconnect to controller.");
#endif
      delay(3000);
      ESP.restart();
    }
  }

  if (!clientConn.available()) {
    return;
  }

  uint8_t  sz = clientConn.read(bufferRX, 255);
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

  return NetSend(bufferTX, 2);
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

  return NetSend(bufferTX, i);
}
