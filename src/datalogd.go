package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.mqtt.golang"

	"github.com/dbld-org/energia/src/axpert"
)

// TODO: These should be configurable
const path = "/dev/hidraw0"
const interval = 30 * time.Second

const mqttServer = "192.168.24.39"
const mqttPort = "1883"
const mqttTopic = "axpert/data"

type messageData struct {
	Timestamp    time.Time
	MesssageType string
	Data         interface{}
}

func main() {

	uc, err := axpert.NewUSBConnector(path)
	if err != nil {
		panic(err)
	}
	err = uc.Open()
	if err != nil {
		panic(err)
	}
	defer uc.Close()
	fmt.Println("connected to ", path)

	clientOpts := mqtt.NewClientOptions()
	clientOpts.AddBroker("tcp://" + mqttServer + ":" + mqttPort)
	clientOpts.SetAutoReconnect(true)
	clientOpts.SetStore(mqtt.NewFileStore("/tmp/mqtt"))
	clientOpts.SetCleanSession(false)
	clientOpts.SetClientID("Axpert")

	client := mqtt.NewClient(clientOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)
	fmt.Println("Connected to mqtt")

	ticker := time.NewTicker(interval)

	go func() {
		for t := range ticker.C {
			mode, err := axpert.DeviceMode(uc)
			if err != nil {
				panic(err)
			}
			m := map[string]string{"Mode": mode}
			msgData := messageData{Timestamp: t, MesssageType: "Mode", Data: m}
			err = sendMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			status, err := axpert.DeviceGeneralStatus(uc)
			if err != nil {
				panic(err)
			}
			msgData = messageData{Timestamp: t, MesssageType: "Status", Data: status}
			err = sendMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			warnings, err := axpert.WarningStatus(uc)
			if err != nil {
				panic(err)
			}
			msgData = messageData{Timestamp: t, MesssageType: "Warnings", Data: warnings}
			err = sendMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			flags, err := axpert.DeviceFlagStatus(uc)
			msgData = messageData{Timestamp: t, MesssageType: "Flags", Data: flags}
			err = sendMessage(msgData, client)
			if err != nil {
				panic(err)
			}

		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	close(sigChan)
	fmt.Println(sig, " stopping ticker")
	ticker.Stop()

	fmt.Println("exiting")
}

func sendMessage(data messageData, client mqtt.Client) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	token := client.Publish(mqttTopic, 1, false, msg)
	token.Wait()
	return nil
}
