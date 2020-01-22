package main

import (
  "log"
  "time"

  "github.com/influxdata/influxdb1-client/v2"
  "gobot.io/x/gobot/drivers/i2c"
  "gobot.io/x/gobot/platforms/raspi"
)

func main() {
  r := raspi.NewAdaptor()

  bme280 := i2c.NewBME280Driver(r, i2c.WithBus(1), i2c.WithAddress(0x76))
  err := bme280.Start()
  if err != nil {
    log.Fatal(err)
  }
  t, err := bme280.Temperature()
  if err != nil {
    log.Fatal(err)
  }
  p, err := bme280.Pressure()
  if err != nil {
    log.Fatal(err)
  }
  p = p / 100
  h, err := bme280.Humidity()
  if err != nil {
    log.Fatal(err)
  }
  bme280.Halt()

  tsl2561 := i2c.NewTSL2561Driver(r, i2c.WithBus(1), i2c.WithAddress(0x39))
  err = tsl2561.Start()
  if err != nil {
    log.Fatal(err)
  }
  bb, ir, err := tsl2561.GetLuminocity()
  if err != nil {
    log.Fatal(err)
  }
  tsl2561.Halt()

  c, err := client.NewHTTPClient(client.HTTPConfig{
    Addr: "http://192.168.0.45:8086",
    Username: "env_writer",
    Password: "det34fgdgser453dfg",
  })
  if err != nil {
    log.Fatal(err)
  }
  defer c.Close()

  bpconf := client.BatchPointsConfig{
    Database: "environment",
  }
  bp, err := client.NewBatchPoints(bpconf)
  if err != nil {
    log.Fatal(err)
  }

  tags := map[string]string{"location":"front_lounge"}
  tph_fields := map[string]interface{}{
    "temperature": t,
    "pressure": p,
    "humidity": h,
  }
  tph, err := client.NewPoint("temp_pressure_humidity", tags, tph_fields, time.Now())
  if err != nil {
    log.Fatal(err)
  }
  bp.AddPoint(tph)
  lux_fields := map[string]interface{}{
    "broadband": bb,
    "infrared": ir,
  }
  lux, err := client.NewPoint("luminocity", tags, lux_fields, time.Now())
  if err != nil {
    log.Fatal(err)
  }
  bp.AddPoint(lux)
  err = c.Write(bp)
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("Temperature: %v*C", t)
  log.Printf("Pressure: %v hPa", p)
  log.Printf("Humidity: %v %%RH", h)
  log.Printf("Broadband Luminocity: %v lux", bb)
  log.Printf("InfraRed Luminocity: %v lux", ir)
}
