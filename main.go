package main

import (
  "io"
  "bufio"
  "log"
  "fmt"
  "os"
  "regexp"
  "strconv"
  "strings"
  "time"

  // "github.com/davecgh/go-spew/spew"
  "github.com/tarm/goserial"
  "github.com/urfave/cli"
  "menteslibres.net/gosexy/redis"
)

type SerialReader struct {
  Device    string
  BaudRate  uint
}

type RedisConnection struct {
  Host  string
  Port  uint
}

func (r RedisConnection) client() *redis.Client {
  rc := redis.New()
  err := rc.ConnectNonBlock(r.Host, r.Port)
  if err != nil { log.Fatalf("FATAL RedisConnection.client(): %s", err) }

  return rc
}

func (s SerialReader) newSerialConnection() io.ReadWriteCloser {
  serialConfig := &serial.Config{ Name: s.Device, Baud: int(s.BaudRate) }
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
      log.Printf("processSerialData() #1: %s", err)
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

      value := fmt.Sprintf("%s,%d,%s", metric_key_and_value[0], uint64(time.Now().Unix()), metric_key_and_value[1])

      _, err := redisConnection.LPush(channel, value)
      if err != nil { log.Printf("processSerialData() #2: %s", err) }

      _, err = redisConnection.Publish(channel, value)
      if err != nil { log.Printf("processSerialData() #3: %s", err) }
    }

    _, err = redisConnection.LTrim(channel, 0, 10)
    if err != nil { log.Printf("processSerialData() #4: %s", err) }
  }
}

func validateSerialDevice(device string) (bool) {
  if _, err := os.Stat(device); os.IsNotExist(err) {
    log.Fatalf("FATAL main() serial device '%s' is invalid: %s", device, err)
  }

  return false
}

func main() {
  app := cli.NewApp()
  app.Name = "moteino-collector"
  app.Usage = "Moteino Collector"
  app.Version = "0.1.0"

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "redis-host",
      Usage: "redis host:port to consume from",
      Value: "127.0.0.1:6379",
      EnvVar: "REDIS_HOST",
    },
    cli.StringFlag{
      Name: "serial-device",
      Usage: "serial device to connect to",
      EnvVar: "SERIAL_DEVICE",
    },
    cli.UintFlag{
      Name: "serial-speed",
      Usage: "speed to use with serial port",
      Value: 115200,
      EnvVar: "SERIAL_SPEED",
    },
  }

  app.Action = func(c *cli.Context) {
    log.Print("Waiting for data...")

    // /dev/tty.usbserial-A50285BI

    s := strings.Split(c.String("redis-host"), ":")
    redisHost := s[0]
    redisPort, _ := strconv.ParseUint(s[1], 10, 64)

    serialDevice := c.String("serial-device")
    serialSpeed := c.Uint("serial-speed")

    validateSerialDevice(serialDevice)

    redisConnection := RedisConnection{ Host: redisHost, Port: uint(redisPort) }
    serialReader    := SerialReader{ Device: serialDevice, BaudRate: serialSpeed }

    serialReader.processSerialData(redisConnection.client())
  }

  app.Run(os.Args)
}
