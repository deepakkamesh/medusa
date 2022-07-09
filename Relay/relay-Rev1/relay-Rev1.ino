#include <SPI.h>
//#include <nRF24L01.h>
#include <RF24.h>
#include <ESP8266WiFi.h>
#include <WiFiUdp.h>


RF24 nrf24(5, 4); // CE, CSN

struct RelayConfig {
  uint8_t pipe_addr_p0[5];
  uint8_t pipe_addr_p1[5];
  uint8_t pipe_addr_p2;
  uint8_t pipe_addr_p3;
  uint8_t pipe_addr_p4;
  uint8_t pipe_addr_p5;
  bool isConfigured;
  uint8_t nrf24Channel;
};

struct RelayConfig Config = {
  .pipe_addr_p0 = {'h', 'e', 'l', 'l', 'o'},
  .pipe_addr_p1 = {'w', 'o', 'r', 'l', 'd'},
  .pipe_addr_p2 = '1',
  .pipe_addr_p3 = '2',
  .pipe_addr_p4 = '3',
  .pipe_addr_p5 = '4',
  .isConfigured = false,
  .nrf24Channel = 115,
};



void setup() {
  Serial.begin(9600);

  //RadioSetup(nrf24);

  // Wifi Setup.
  WifiConnect();
  RelaySetup();
}

void loop() {
  // RadioLoop(nrf24);
  IpLoop();
}
