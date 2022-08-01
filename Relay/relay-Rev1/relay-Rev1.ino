#include <SPI.h>
#include <RF24.h>
#include <ESP8266WiFi.h>
#include <WiFiUdp.h>
#include <ArduinoOTA.h>

#define DEBUG 1
#define LED_ONBOARD 16

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
