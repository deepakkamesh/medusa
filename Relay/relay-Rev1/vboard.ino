/* Handles the virtual board configuration */


int SendPing() {
  
  uint8_t i = 0;

  bufferTX[i] = PKT_TYPE_BOARD_DATA_RELAY;
  i++;
  SuperMemCpy(bufferTX, i, Config.pipe_addr[VIRT_PIPE], 0, PIPE_ADDR_LEN);
  i += PIPE_ADDR_LEN;
  bufferTX[i] = PKT_TYPE_PING;
  i++;
  SuperMemCpy(bufferTX, i, Config.vboard_addr, 0, ADDR_LEN);
  i += ADDR_LEN;
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

unsigned long prevTicks = 0;

void PingLoop() {
  unsigned long currTicks = millis();

  if (currTicks - prevTicks >= PING_INT) {
    SendPing();
    prevTicks = currTicks;
  }
}

void ProcessVBoardPacket(uint8_t *pkt, uint8_t sz) {

  uint8_t data [32];
  uint8_t actionID;

  uint8_t pktType = pkt[0];
  switch (pktType) {
    case PKT_TYPE_ACTION:
      actionID = pkt[4];
      SuperMemCpy(data, 0, pkt, 5, sz - 5);
      ProcessAction(actionID, data);
      break;

    default:
      SendError(ERROR_UNKNOWN_PKT);
      break;
  }
}

void ProcessAction(uint8_t actionID, uint8_t * data) {
#ifdef DEBUG
  Serial.printf("Got Relay Action Request:%d\n", actionID);
#endif

  switch (actionID) {
    case ACTION_RESET_DEVICE:
      ESP.restart();
      break;
    case ACTION_STATUS_LED:
      digitalWrite(LED_ONBOARD, data[0]);
      break;
    case ACTION_FLUSH_TX_FIFO:
      radio.flush_tx();
      break;
    default:
      SendError(ERROR_RELAY_NOT_IMPLEMENTED);
  }
}
