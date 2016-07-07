#include "Timer.h"
#include <RFM69.h>
#include <SPI.h>

// Timer setup
//
Timer t;

// Radio setup
//
RFM69 radio;
byte ackCount = 0;

#define NODE_ID        1
#define NETWORK_ID     100
#define FREQUENCY      RF69_433MHZ
#define ENCRYPT_KEY    "a7bd91kaxlchdk36"
#define ACK_TIME       30

// Serial setup
//
// NOTE: Serial output comes out on PIN 1
//       http://farm4.staticflickr.com/3818/10585364014_d028c66523_o.png
//
#define SERIAL_BAUD   115200

// LED setup
//
#define LED_PIN      9

void blink_led(byte pin, int delay_ms) {
  pinMode(pin, OUTPUT);
  digitalWrite(pin, HIGH);

  delay(delay_ms);
  digitalWrite(pin, LOW);
}

void radio_data_process() {
  if (radio.receiveDone()) {

    char data[128];

    for (byte i = 0; i < radio.DATALEN; i++) { data[i] = (char)radio.DATA[i]; }

    data[radio.DATALEN] = '\0';

    if (radio.ACK_REQUESTED) {
      byte node_id = radio.SENDERID;
      radio.sendACK();

      if (ackCount++%3==0) {
        delay(3); // need this when sending right after reception .. ?
        radio.sendWithRetry(node_id, "ACK TEST", 8, 0);
      }
    }

    Serial.println(data);
    free(data);
    blink_led(LED_PIN, 3);
  }
}

void blank_line() {
  Serial.println();
}

void setup_radio() {
  delay(10);
  radio.initialize(FREQUENCY, NODE_ID, NETWORK_ID);
  radio.encrypt(ENCRYPT_KEY);
}

void setup_serial() {
  Serial.begin(SERIAL_BAUD);
}

void setup() {
  setup_serial();
  setup_radio();
  t.every(1, radio_data_process);
}

void loop() {
  t.update();
}
