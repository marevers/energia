package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbld-org/energia/src/axpert"
)

// TODO: These should be configurable
const path = "/dev/hidraw0"
const interval = 30 * time.Second

//const mqttServer = "192.168.24.39"
//const mqttPort = 1883

func main() {

	ticker := time.NewTicker(interval)

	go func() {
		uc, err := axpert.NewUSBConnector(path)
		if err != nil {
			panic(err)
		}
		defer uc.Close()
		fmt.Println("connected to ", path)

		for t := range ticker.C {
			mode, err := axpert.DeviceMode(uc)
			if err != nil {
				panic(err)
			}
			status, err := axpert.DeviceGeneralStatus(uc)
			if err != nil {
				panic(err)
			}
			// TODO more calls, send to mqtt
			fmt.Println(t.Format(time.RFC3339), " mode: ", mode, " status: ", status)
		}
		fmt.Println("closing connection")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Println(sig, " stopping ticker")
	ticker.Stop()

	fmt.Println("exiting")
}
