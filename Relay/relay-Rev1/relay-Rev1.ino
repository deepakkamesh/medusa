#include <SPI.h>
#include <RF24.h>
#include <ESP8266WiFi.h>
#include <WiFiUdp.h>
#include <ArduinoOTA.h>

#define DEBUG 1
#define LED_ONBOARD 16
#define ADDR_LEN 3
#define PIPE_ADDR_LEN 5
#define PIPE_ADDR_NUM 7 // number of address. 6 + 1 virtual pipe addr. 
#define VIRT_PIPE 6 // number of virtual pipe.
#define PING_INT 5000 // ping interval in ms.

#define PKT_TYPE_RELAY_GET_CONFIG 0xAA // Packet type for Relay Get Config
#define PKT_TYPE_RELAY_CONFIG_ANS 0xAB // Packet type for Relay Get Answer
#define PKT_TYPE_RELAY_ERROR 0xAC // Packet type for Relay  error. 
#define PKT_TYPE_BOARD_DATA_RELAY 0xAD // Packet type for Relay  data packet. 
#define PKT_TYPE_PING 0x02
#define PKT_TYPE_ACTION 0x10

#define ERROR_RELAY_RADIO_INIT_FAILED 0x02
#define ERROR_RELAY_ACK_PAYLOAD_LOAD 0x03
#define ERROR_RELAY_NOT_IMPLEMENTED 0x04
#define ERROR_PIPE_ADDR_404 0x06
#define ERROR_UNKNOWN_PKT 0x05

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
  .vboard_addr = {0xA, 0xB, 0xC},
};

RF24 radio(5, 4); // CE, CSN.


void setup() {
  Serial.begin(9600);
  pinMode(LED_ONBOARD, OUTPUT);
  int ok ;
  WifiConnect();
  OTAInit();
  ok = RelaySetup();
  if (!ok) {
    digitalWrite(LED_ONBOARD, false); // false turns it on?.
    delay(3000);
    ESP.restart();
  }
  ok = RadioSetup();
  if (!ok) {
    SendError(ERROR_RELAY_RADIO_INIT_FAILED);
    digitalWrite(LED_ONBOARD, false); // false turns it on?.
    delay(3000);
    ESP.restart();
  }
}

void loop() {
  ArduinoOTA.handle();
  RadioRcvLoop();
  IpLoop();
  RadioSendLoop();
  WifiKeepAlive();
  PingLoop();
}
