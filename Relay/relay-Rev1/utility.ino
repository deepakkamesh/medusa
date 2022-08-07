/************* Utility Functions *********************/


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
