package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/goburrow/serial"

	"github.com/dbld-org/energia/internal/connector"
	"github.com/dbld-org/energia/pylontech"
)

var (
	address  string
	baudrate int
	databits int
	stopbits int
	parity   string
)

func main() {
	flag.StringVar(&address, "a", "/dev/ttyUSB0", "address")
	flag.IntVar(&baudrate, "b", 1200, "baud rate")
	flag.IntVar(&databits, "d", 8, "data bits")
	flag.IntVar(&stopbits, "s", 1, "stop bits")
	flag.StringVar(&parity, "p", "N", "parity (N/E/O)")
	flag.Parse()

	config := serial.Config{
		Address:  address,
		BaudRate: baudrate,
		DataBits: databits,
		StopBits: stopbits,
		Parity:   parity,
		Timeout:  30 * time.Second,
	}

	log.Printf("connecting %+v", config)

	sc := connector.NewSerialConnector(config)
	err := sc.Open()
	if err != nil {
		log.Panic(err)
	}
	defer sc.Close()

	version, err := pylontech.GetProtocolVersion(sc)

	if err != nil {
		log.Panic(err)
	}

	log.Println("Protocol version:", version)

	manufacturerInfo, err := pylontech.GetManufacturerInfo(sc)

	bytes, err := json.MarshalIndent(manufacturerInfo, "", "   ")
	log.Println("Manufacturer info:", string(bytes))

	batteryStatus, err := pylontech.GetBatteryStatus(sc)

	bytes, err = json.MarshalIndent(batteryStatus, "", "   ")
	log.Println("Battery status:", string(bytes))

}
