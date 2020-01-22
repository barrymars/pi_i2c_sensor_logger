package main

import (
  "log"
  "time"

  "github.com/d2r2/go-i2c"
  "github.com/d2r2/go-bsbmp"
  "github.com/d2r2/go-logger"
  "github.com/influxdata/influxdb1-client/v2"
)

func main() {
  i2c_bme280, err := i2c.NewI2C(0x76, 1)
  if err != nil {
    log.Fatal(err)
  }
  defer i2c_bme280.Close()
  logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

  bme280, err := bsbmp.NewBMP(bsbmp.BME280, i2c_bme280)
  if err != nil {
    log.Fatal(err)
  }
  logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)

  t, err := bme280.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
  if err != nil {
    log.Fatal(err)
  }
  p, err := bme280.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
  if err != nil {
    log.Fatal(err)
  }
  p = p / 100
  _, h, err := bme280.ReadHumidityRH(bsbmp.ACCURACY_STANDARD)
  if err != nil {
    log.Fatal(err)
  }

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
  fields := map[string]interface{}{
    "temperature": t,
    "pressure": p,
    "humidity": h,
  }
  tp, err := client.NewPoint("temp_pressure_humidity", tags, fields, time.Now())
  if err != nil {
    log.Fatal(err)
  }
  bp.AddPoint(tp)
  err = c.Write(bp)
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Temperature: %v*C", t)
  log.Printf("Pressure: %v hPa", p)
  log.Printf("Humidity: %v %%RH", h)

  i2c_tsl2561, err := i2c.NewI2C(0x39, 1)
  if err != nil {
    log.Fatal(err)
  }
  defer i2c_tsl2561.Close()

  _, err := NewTSL2561()
  if err != nil {
    log.Fatal(err)
  }

}
