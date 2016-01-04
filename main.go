package main

import (
  "io"
  "bufio"
  "log"
  "fmt"
  "regexp"
  "strings"
  "time"

  "github.com/tarm/goserial"
  "menteslibres.net/gosexy/redis"
)

var redisHost = "127.0.0.1"
var redisPort = uint(6379)

var serialPort = "/dev/ttyAMA0"
var serialBaudRate = 115200

type SerialReader struct {
  DeviceName  string
  BaudRate    int
}

type RedisConnection struct {
  Host    string
  Port    uint
}

func (r RedisConnection) client() *redis.Client {
  rc := redis.New()
  rc.ConnectNonBlock(redisHost, redisPort)

  return rc
}

func (s SerialReader) newSerialConnection() io.ReadWriteCloser {
  serialConfig := &serial.Config{ Name: s.DeviceName, Baud: s.BaudRate }
  serialPort, err := serial.OpenPort(serialConfig)
  if err != nil { log.Fatalf("FATAL newSerialConnection(): %s", err) }

  return serialPort
}

func (s SerialReader) processSerialData(redisConnection *redis.Client) {
  serialConnection := s.newSerialConnection()
  reader := bufio.NewReader(serialConnection)
  re := regexp.MustCompile("^([^:]+):(.+)$")

  for {
    reply, err := reader.ReadBytes('\n')

    if err != nil {
      log.Printf("FATAL processSerialData() #1: %s", err)
      continue
    }

    data := strings.TrimSpace(string(reply))

    if len(data) == 0 { continue }

    log.Printf("%s", data)
    raw_data := re.FindStringSubmatch(data)

    if len(raw_data) == 0 { continue }

    metrics := strings.Split(raw_data[2], ",")
    channel := fmt.Sprintf("rpi-moteino-collector:%s", raw_data[1])

    for i := 0; i < len(metrics); i++ {
      metric_key_and_value := strings.Split(metrics[i], ":")

      if len(metric_key_and_value) != 2 { continue }

      value := fmt.Sprintf("%s,%d,%s", metric_key_and_value[0], int64(time.Now().Unix()), metric_key_and_value[1])

      _, err := redisConnection.LPush(channel, value)
      if err != nil { log.Printf("FATAL processSerialData() #2: %s", err) }

      _, err = redisConnection.Publish(channel, value)
      if err != nil { log.Printf("FATAL processSerialData() #3: %s", err) }
    }

    _, err = redisConnection.LTrim(channel, 0, 10)
    if err != nil { log.Printf("FATAL processSerialData() #4: %s", err) }
  }
}

func main() {
  log.Print("Waiting for data...")

  redisConnection := RedisConnection{ Host: redisHost, Port: redisPort }
  serialReader := SerialReader{ DeviceName: serialPort, BaudRate: serialBaudRate }

  serialReader.processSerialData(redisConnection.client())
}

