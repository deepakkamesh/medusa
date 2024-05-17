/* Handles the virtual board configuration */
uint16_t buzzInt = 500;
bool buzzerOn = false;

unsigned long prevTicks = 0;

void PingLoop() {
  unsigned long currTicks = millis();

  if (currTicks - prevTicks >= PING_INT) {
    SendPing();
    prevTicks = currTicks;
  }
}

unsigned long prevSensorDataTicks = 0;
void SensorDataLoop() {
  unsigned long currTicks = millis();

  if (currTicks - prevSensorDataTicks >= SENSOR_INT*1000) {

#ifdef DHT11SENSOR
    TempHumidity();
#endif
#ifdef BME680
    BME680TempHumidity();
    BME680Gas();
    BME680Pressure();
    BME680Altitude();
#endif
    prevSensorDataTicks = currTicks;
  }
}

void ProcessVBoardPacket(uint8_t *pkt, uint8_t sz) {

  uint8_t data[32];
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

void ProcessAction(uint8_t actionID, uint8_t *data) {

  switch (actionID) {
    case ACTION_RESET_DEVICE:
      ESP.restart();
      break;

    case ACTION_STATUS_LED:
      digitalWrite(LED_ONBOARD, !data[0]);  // For some reason 0 turns on LED.
      break;

    case ACTION_FLUSH_TX_FIFO:
      radio.flush_tx();
      break;

    case ACTION_TEMP:
#ifdef DHT11SENSOR
      TempHumidity();
      break;
#endif
#ifdef BME680
      BME680TempHumidity();
      break;
#endif
      SendError(ERROR_RELAY_NOT_IMPLEMENTED);
      break;

    case ACTION_GAS:
#ifndef BME680
      SendError(ERROR_RELAY_NOT_IMPLEMENTED);
      break;
#endif
      BME680Gas();
      break;

    case ACTION_PRESSURE:
#ifndef BME680
      SendError(ERROR_RELAY_NOT_IMPLEMENTED);
      break;
#endif
      BME680Pressure();
      break;

    case ACTION_ALTITUDE:
#ifndef BME680
      SendError(ERROR_RELAY_NOT_IMPLEMENTED);
      break;
#endif
      BME680Altitude();
      break;

    case ACTION_BUZZER:
#ifndef BUZZERDEV
      SendError(ERROR_RELAY_NOT_IMPLEMENTED);
      break;
#endif
      buzzerOn = data[0];
      buzzInt = data[1];
      buzzInt = (buzzInt << 8) | data[2];
      break;

    default:
      SendError(ERROR_RELAY_NOT_IMPLEMENTED);
  }
}

unsigned long prevBuzzTicks = 0;
void HandleBuzzerLoop() {

  unsigned long currTicks = millis();

  if (currTicks - prevBuzzTicks >= buzzInt) {
    prevBuzzTicks = currTicks;

    if (!buzzerOn) {
      digitalWrite(BUZZERPIN, LOW);
      return;
    }
    digitalWrite(BUZZERPIN, !digitalRead(BUZZERPIN));
  }
}


bool motionFlag = false;
unsigned long startTicks = 0;

void HandleMotionSensorLoop() {
  uint8_t buff[1];

  // Wait 2s to prevent flapping. The sensor is super sensitive
  // and senses continuously.
  if (millis() - startTicks < 2000) {
    return;
  }

  buff[0] = digitalRead(MOTIONPIN);

  if (buff[0] && !motionFlag) {
    SendData(ACTION_MOTION, buff, 1);
    motionFlag = true;
    startTicks = millis();
    return;
  }

  if (!buff[0] && motionFlag) {
    SendData(ACTION_MOTION, buff, 1);
    motionFlag = false;
    startTicks = millis();
    return;
  }
}

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
  return NetSend(bufferTX, i);
}


// Sends the data for the particular action.
bool SendData(uint8_t Action, uint8_t *data, int8_t sz) {
  uint8_t i = 0;
  bufferTX[i] = PKT_TYPE_BOARD_DATA_RELAY;
  i++;
  SuperMemCpy(bufferTX, i, Config.pipe_addr[VIRT_PIPE], 0, PIPE_ADDR_LEN);
  i += PIPE_ADDR_LEN;
  bufferTX[i] = PKT_TYPE_DATA;
  i++;
  SuperMemCpy(bufferTX, i, Config.vboard_addr, 0, ADDR_LEN);
  i += ADDR_LEN;
  bufferTX[i] = Action;
  i++;
  bufferTX[i] = 0x00;
  i++;
  SuperMemCpy(bufferTX, i, data, 0, sz);
  i += sz;
  return NetSend(bufferTX, i);
}

// Utility struct for pack and send.
union val {
  float f;
  uint8_t uc[4];
};

union valUint32 {
  uint32_t f;
  uint8_t uc[4];
};


/***************** DHT Sensor *******************/
DHT dht(DHTPIN, DHTTYPE);
void dhtstart() {
  dht.begin();
}

void TempHumidity() {

  union val temp, humidity;
  uint8_t buff[10];
  // Reading temperature or humidity takes about 250 milliseconds!
  // Sensor readings may also be up to 2 seconds 'old' (its a very slow sensor)
  humidity.f = dht.readHumidity();
  // Read temperature as Celsius (the default)
  temp.f = dht.readTemperature();

  buff[0] = temp.uc[0];
  buff[1] = temp.uc[1];
  buff[2] = temp.uc[2];
  buff[3] = temp.uc[3];
  buff[4] = humidity.uc[0];
  buff[5] = humidity.uc[1];
  buff[6] = humidity.uc[2];
  buff[7] = humidity.uc[3];

  // Check if any reads failed and exit early (to try again).
  if (isnan(humidity.f) || isnan(temp.f)) {
    Serial.println(F("Failed to read from DHT sensor!"));
    return;
  }

  SendData(ACTION_TEMP, buff, 8);
}
/***************** BME680 Sensor *******************/
// comp_gas = log(R_gas[ohm]) + 0.04 log(Ohm)/%rh * hum[%rh] IAQ
// https://forums.pimoroni.com/t/bme680-observed-gas-ohms-readings/6608/20
// https://esphome.io/components/sensor/bme680.html#bme680-oversampling


Adafruit_BME680 bme;
int BME680Setup() {
  if (!bme.begin()) {
    return 0;
  }
  // Set up oversampling and filter initialization
  bme.setTemperatureOversampling(BME680_OS_8X);
  bme.setHumidityOversampling(BME680_OS_2X);
  bme.setPressureOversampling(BME680_OS_4X);
  bme.setIIRFilterSize(BME680_FILTER_SIZE_3);
  bme.setGasHeater(320, 150);  // 320*C for 150 ms
}

void BME680TempHumidity() {
  union val temp, humidity;
  uint8_t buff[10];
  humidity.f = bme.readHumidity();
  temp.f = bme.readTemperature();

  buff[0] = temp.uc[0];
  buff[1] = temp.uc[1];
  buff[2] = temp.uc[2];
  buff[3] = temp.uc[3];
  buff[4] = humidity.uc[0];
  buff[5] = humidity.uc[1];
  buff[6] = humidity.uc[2];
  buff[7] = humidity.uc[3];

  // Check if any reads failed and exit early (to try again).
  if (isnan(humidity.f) || isnan(temp.f)) {
    Serial.println(F("Failed to read from BME sensor!"));
    return;
  }

  SendData(ACTION_TEMP, buff, 8);
}

void BME680Pressure() {
  union val pressure;
  uint8_t buff[5];
  pressure.f = bme.readPressure();

  buff[0] = pressure.uc[0];
  buff[1] = pressure.uc[1];
  buff[2] = pressure.uc[2];
  buff[3] = pressure.uc[3];

  // Check if any reads failed and exit early (to try again).
  if (isnan(pressure.f)) {
    Serial.println(F("Failed to read from BME sensor!"));
    return;
  }
  SendData(ACTION_PRESSURE, buff, 4);
}

void BME680Gas() {
  union valUint32 gas;
  uint8_t buff[5];
  gas.f = bme.readGas();

  buff[0] = gas.uc[0];
  buff[1] = gas.uc[1];
  buff[2] = gas.uc[2];
  buff[3] = gas.uc[3];

  // Check if any reads failed and exit early (to try again).
  if (isnan(gas.f)) {
    Serial.println(F("Failed to read from BME sensor!"));
    return;
  }
  SendData(ACTION_GAS, buff, 4);
}

void BME680Altitude() {
  union val alt;
  uint8_t buff[5];
  alt.f = bme.readAltitude(SEALEVELPRESSURE_HPA);

  buff[0] = alt.uc[0];
  buff[1] = alt.uc[1];
  buff[2] = alt.uc[2];
  buff[3] = alt.uc[3];

  // Check if any reads failed and exit early (to try again).
  if (isnan(alt.f)) {
    Serial.println(F("Failed to read from BME sensor!"));
    return;
  }
  SendData(ACTION_ALTITUDE, buff, 4);
}
