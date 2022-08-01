#include <SPI.h>
#include <RF24.h>
#include <ESP8266WiFi.h>
#include <WiFiUdp.h>
#include <ArduinoOTA.h>

#define DEBUG 1
#define LED_ONBOARD 16
#define PKT_TYPE_RELAY_GET_CONFIG 0xAA // Packet type for Relay Get Config
#define PKT_TYPE_RELAY_CONFIG_ANS 0xAB // Packet type for Relay Get Answer
#define PKT_TYPE_RELAY_ERROR 0xAC // Packet type for Relay  error. 
#define PKT_TYPE_RELAY_DATA 0xAD // Packet type for Relay  data packet. 
#define PKT_TYPE_RELAY_ACTION 0xAE // Packet type for relay action.

#define ERROR_RELAY_RADIO_INIT_FAILED 0x02
#define ERROR_RELAY_ACK_PAYLOAD_LOAD 0x03
#define ERROR_RELAY_NOT_IMPLEMENTED 0x04

#define ACTION_STATUS_LED 0x13
#define ACTION_RESET_DEVICE 0x14
#define ACTION_FLUSH_TX_FIFO 0x17
struct RelayConfig {
  uint8_t pipe_addr_p0[5];
  uint8_t pipe_addr_p1[5];
  uint8_t pipe_addr_p2[1];
  uint8_t pipe_addr_p3[1];
  uint8_t pipe_addr_p4[1];
  uint8_t pipe_addr_p5[1];
  bool isConfigured;
  uint8_t nrf24Channel;
};

struct RelayConfig Config = {
  .pipe_addr_p0 = {'h', 'e', 'l', 'l', 'o'},
  .pipe_addr_p1 = {'w', 'o', 'r', 'l', 'd'},
  .pipe_addr_p2 = {'1'},
  .pipe_addr_p3 = {'2'},
  .pipe_addr_p4 = {'3'},
  .pipe_addr_p5 = {'4'},
  .isConfigured = false,
  .nrf24Channel = 115,
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
    delay(1000);
    ESP.restart();
  }
  ok = RadioSetup();
  if (!ok) {
    SendError(ERROR_RELAY_RADIO_INIT_FAILED);
    digitalWrite(LED_ONBOARD, false); // false turns it on?.
    delay(1000);
    ESP.restart();
  }
}

void loop() {
  ArduinoOTA.handle();
  RadioRcvLoop();
  IpLoop();
  RadioSendLoop();
  WifiKeepAlive();

}
