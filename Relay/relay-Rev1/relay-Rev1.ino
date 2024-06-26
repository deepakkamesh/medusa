#include <SPI.h>
#include <RF24.h>
#include <ESP8266WiFi.h>
#include <WiFiUdp.h>
#include <ArduinoOTA.h>
#include <stdarg.h>
#include "DHT.h"
#include <Adafruit_Sensor.h>
#include "Adafruit_BME680.h"

#define DEBUG 1
#define ADDR_LEN 3
#define PIPE_ADDR_LEN 5
#define PIPE_ADDR_NUM 7  // number of address. 6 + 1 virtual pipe addr.
#define VIRT_PIPE 6      // number of virtual pipe.
#define PING_INT 50000   // ping interval in ms.
#define SENSOR_INT 120   // Sensor poll interval in sec.

#define PKT_TYPE_RELAY_GET_CONFIG 0xAA  // Packet type for Relay Get Config
#define PKT_TYPE_RELAY_CONFIG_ANS 0xAB  // Packet type for Relay Get Answer
#define PKT_TYPE_RELAY_ERROR 0xAC       // Packet type for Relay  error.
#define PKT_TYPE_BOARD_DATA_RELAY 0xAD  // Packet type for Relay  data packet.
#define PKT_TYPE_PING 0x02
#define PKT_TYPE_ACTION 0x10
#define PKT_TYPE_DATA 0x01

#define ERROR_RELAY_RADIO_INIT_FAILED 0x02
#define ERROR_RELAY_ACK_PAYLOAD_LOAD 0x03
#define ERROR_RELAY_NOT_IMPLEMENTED 0x04
#define ERROR_UNKNOWN_PKT 0x05
#define ERROR_PIPE_ADDR_404 0x06
#define ERROR_SENSOR 0x07

#define ACTION_MOTION 0x01
#define ACTION_TEMP 0x02
#define ACTION_VERSION 0x06
#define ACTION_GAS 0x07
#define ACTION_PRESSURE 0x08
#define ACTION_ALTITUDE 0x09
#define ACTION_BUZZER 0x10
#define ACTION_STATUS_LED 0x13
#define ACTION_RESET_DEVICE 0x14
#define ACTION_FLUSH_TX_FIFO 0x17

struct RelayConfig {
  uint8_t pipe_addr[PIPE_ADDR_NUM][PIPE_ADDR_LEN];
  bool isConfigured;
  uint8_t nrf24Channel;
  uint8_t vboard_addr[3];
};

struct RelayConfig Config = {
  .isConfigured = false,
  .nrf24Channel = 115,
  .vboard_addr = { 0xA, 0xB, 0xC },
};

/*************** START CONFIGURE HERE *************************/
// Sea pressure in millibar
#define SEALEVELPRESSURE_HPA (1020)

// NOTE: ENSURE PINS DONT CONFLICT.
// Board pin connectivity and configuration.
#define BOARD_TYPE_WEMOS
//#define BOARD_TYPE_NODEMCU


#ifdef BOARD_TYPE_NODEMCU
#define LED_ONBOARD 16
#define DHTTYPE DHT11
#define DHTPIN D4
#define MOTIONPIN D3
RF24 radio(5, 4);  // CE, CSN.
#endif

#ifdef BOARD_TYPE_WEMOS
#define LED_ONBOARD LED_BUILTIN  // LED is on D4.
#define MOTIONPIN D1
#define DHTTYPE DHT11
#define DHTPIN D2
#define BUZZERPIN D3
RF24 radio(D0, D8);  // CE, CSN.
#endif

//  Sensors onboard. AVOID PIN CONFLICT.
//#define DHT11SENSOR
//#define RCWL516SENSOR
//#define BUZZERDEV
#define BME680  // Uses I2C pins D1, D2.

/*************** END CONFIGURE HERE *************************/

void setup() {
  int ok;

  // Setup Basics.
  Serial.begin(9600);
  pinMode(LED_ONBOARD, OUTPUT);
#ifdef BUZZERDEV
  pinMode(BUZZERPIN, OUTPUT);
#endif
#ifdef RCWL516
  pinMode(MOTIONPIN, INPUT);
#endif

  // Find the strongest Wifi
  WifiAPFinder();
  // Setup Wifi and get configs.
  WifiConnect();
  OTAInit();
  ok = RelaySetup();
  if (!ok) {
    DbgPrint("%s\n", "Relay Setup Failed");
    digitalWrite(LED_ONBOARD, false);  // false turns it on?.
    delay(3000);
    ESP.restart();
  }

  // Setup Radio.
  ok = RadioSetup();
  if (!ok) {
    DbgPrint("%s\n", "Radio Setup Failed");
    SendError(ERROR_RELAY_RADIO_INIT_FAILED);
    digitalWrite(LED_ONBOARD, false);  // false turns it on?.
    delay(3000);
    ESP.restart();
  }

  // Setup Onboard Sensors.
#ifdef DHT11SENSOR
  dhtstart();
#endif
#ifdef BME680
  ok = BME680Setup();
  if (!ok) {
    DbgPrint("%s\n", "BME680 init failed");
    SendError(ERROR_SENSOR);
  }
#endif
}


void loop() {
  ArduinoOTA.handle();
  RadioRcvLoop();
  IpLoop();
  RadioSendLoop();
  WifiKeepAlive();
  PingLoop();
  SensorDataLoop();
#ifdef RCWL516SENSOR
  HandleMotionSensorLoop();
#endif
#ifdef BUZZERDEV
  HandleBuzzerLoop();
#endif
}
