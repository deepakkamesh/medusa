/************* Utility Functions *********************/



/********** ParseConfigPkt() - Parses config response into the Config ***************/
uint8_t  ParseConfigPkt(uint8_t * bufferRX, uint8_t len) {
  // Valid packet?
  if (!(bufferRX[0] == PKT_TYPE_RELAY_CONFIG_ANS && len == RELAY_CONFIG_ANS_LEN)) {
    return 0;
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

  Config.pipe_addr[6][0] = bufferRX[15];
  SuperMemCpy(Config.pipe_addr[6], 1, bufferRX, 7, PIPE_ADDR_LEN - 1);

  /* Copy pipe address into config */
  Config.nrf24Channel = bufferRX[16];

  /* Copy vBoard address into config */
  SuperMemCpy(Config.vboard_addr, 0, bufferRX, 17, ADDR_LEN);

#ifdef DEBUG
  for (int i = 0; i < PIPE_ADDR_NUM; i++) {
    PrintPkt("Addr:", Config.pipe_addr[i], PIPE_ADDR_LEN);
  }
  Serial.printf("CHANNEL: %02X\n", Config.nrf24Channel);
  PrintPkt("vBoard:", Config.vboard_addr, ADDR_LEN);
#endif

  return 1;
}

IPAddress GetBroadcastIP() {
  // Find the broadcast IP.
  IPAddress mask = WiFi.subnetMask();
  IPAddress ip = WiFi.gatewayIP();
  IPAddress bcastIP;
  for (int i = 0; i < 4; i++) {
    mask[i] = ~mask[i];
    bcastIP[i] = ip[i] | mask[i];
  }
  return bcastIP;
}

int NetSend(uint8_t *buff, uint8_t sz) {
  return clientConn.write(buff, sz);
}

// NetSend sends the buff of size sz over the network.
int NetSendUDP(uint8_t *buff, uint8_t sz, IPAddress ip, uint16_t port) {
  int ok = Udp.beginPacket(ip, port);
  if (!ok) {
    return 0;
  }
  Udp.write(buff, sz);
  ok = Udp.endPacket();
  if (!ok) {
    return 0;
  }
  return 1;
}

void PrintPkt(char *str, uint8_t buff[], int len) {
#ifdef DEBUG
  Serial.print(str);
  for (int i = 0; i < len; i++) {
    Serial.printf(" %02X,", buff[i]);
  }
  Serial.println();
#endif
}

/* Handle Wifi Disconnects */
unsigned long previousMillis = 0;
unsigned long interval = 30000;

void WifiKeepAlive(void) {
  unsigned long currentMillis = millis();
  if (currentMillis - previousMillis >= interval) {
    if (WiFi.status() != WL_CONNECTED) {
      delay(1000);
      ESP.restart();
    }
    previousMillis = currentMillis;
  }
}

/* SuperMemCpy copies arrays with start index and size */
void SuperMemCpy(uint8_t *dest, uint8_t destStart, uint8_t *src, uint8_t srcStart, uint8_t sz) {
  for (uint8_t i = 0; i < sz; i++) {
    dest[i + destStart] = src[i + srcStart];
  }
}

/* FindPipeNum returns the pipe number associated with the address. -1 if none found */
int FindPipeNum(uint8_t *pipe_addr, uint8_t pipe_addr_sz) {
  for (uint8_t i = 0 ; i < PIPE_ADDR_NUM ; i++) {
    if (CompareArray(pipe_addr, Config.pipe_addr[i], pipe_addr_sz)) {
      return i;
    }
  }
  return -1;
}

/* CompareArray returns 1 if src == dst otherwise 0 */
uint8_t CompareArray(uint8_t *src, uint8_t *dst, uint8_t sz) {
  for (uint8_t i = 0; i < sz; i++) {
    if (src[i] != dst[i]) {
      return 0;
    }
  }
  return 1;
}

/***************************** Queuing Functions *****************************/
void initQueue(Queue *q) {
  q->readPtr = 0;
  q->writePtr = 0;
  q->overflow = 0;
  uint8_t i = 0;
  for (i = 0; i < MAX_TX_QUEUE_SZ; i++) {
    q->packets[i].size = 0;
    q->packets[i].pipeNum = 0;
  }
}

uint8_t enQueue(uint8_t *buf, uint8_t sz, uint8_t pipeNum, Queue *q) {
  if (q->overflow && q->writePtr == q->readPtr) {
    return 0;
  }

  memcpy(q->packets[q->writePtr].packet, buf, sz);
  q->packets[q->writePtr].size = sz;
  q->packets[q->writePtr].pipeNum = pipeNum;
  q->writePtr++;

  if (q->writePtr == MAX_TX_QUEUE_SZ) {
    q->writePtr = 0;
    q->overflow = 1;
  }
  return 1;
}

uint8_t deQueue(uint8_t *buff, uint8_t *pipeNum, Queue *q) {
  if (!q->overflow && q->readPtr == q->writePtr) {
    return 0;
  }
  if (q->overflow && q->readPtr < q->writePtr) {
    q->readPtr = q->writePtr;
  }

  memcpy(buff, q->packets[q->readPtr].packet, q->packets[q->readPtr].size);
  uint8_t sz = q->packets[q->readPtr].size;
  *pipeNum = q->packets[q->readPtr].pipeNum;
  q->readPtr++;

  if (q->readPtr == MAX_TX_QUEUE_SZ) {
    q->readPtr = 0;
    q->overflow = 0;
  }
  return sz;
}
