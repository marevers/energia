package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/goburrow/serial"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/dbld-org/energia/axpert"
	"github.com/dbld-org/energia/internal/connector"
	"github.com/dbld-org/energia/pylontech"
)

var timerInterval int

var mqttServer string
var mqttPort int
var mqttClientId string

var inverterPath string
var inverterCount int
var inverterTopic string

var batteryPath string
var batteryBaud int
var batteryTopic string

type messageData struct {
	Timestamp   time.Time
	MessageType string
	Data        interface{}
}

func main() {

	err := initConfig()
	if err != nil {
		panic(err)
	}

	uc, err := connector.NewUSBConnector(inverterPath)
	if err != nil {
		panic(err)
	}
	err = uc.Open()
	if err != nil {
		panic(err)
	}
	defer uc.Close()
	fmt.Println("connected to ", inverterPath)

	serialConfig := serial.Config{
		Address:  batteryPath,
		BaudRate: batteryBaud,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  30 * time.Second,
	}

	sc := connector.NewSerialConnector(serialConfig)
	err = sc.Open()
	if err != nil {
		log.Panic(err)
	}
	defer sc.Close()

	clientOpts := mqtt.NewClientOptions()
	clientOpts.AddBroker("tcp://" + mqttServer + ":" + strconv.Itoa(mqttPort))
	clientOpts.SetAutoReconnect(true)
	clientOpts.SetStore(mqtt.NewFileStore("/tmp/mqtt"))
	clientOpts.SetCleanSession(false)
	clientOpts.SetClientID(mqttClientId)
	clientOpts.SetOnConnectHandler(logConnect)
	clientOpts.SetConnectionLostHandler(logConnectionLost)

	client := mqtt.NewClient(clientOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)
	fmt.Println("Connected to mqtt")

	ticker := time.NewTicker(time.Duration(timerInterval) * time.Second)

	go func() {
		for t := range ticker.C {
			mode, err := axpert.DeviceMode(uc)
			if err != nil {
				panic(err)
			}
			m := map[string]string{"Mode": mode}
			msgData := messageData{Timestamp: t, MessageType: "Mode", Data: m}
			err = sendInverterMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			for inv := 0; inv < inverterCount; inv++ {
				deviceInfo, err := axpert.ParallelDeviceInfo(uc, inv)
				if err != nil {
					panic(err)
				}
				msgData = messageData{Timestamp: t, MessageType: "DeviceInfo", Data: deviceInfo}
				err = sendInverterMessage(msgData, client)
				if err != nil {
					panic(err)
				}
			}

			status, err := axpert.DeviceGeneralStatus(uc)
			if err != nil {
				panic(err)
			}
			msgData = messageData{Timestamp: t, MessageType: "Status", Data: status}
			err = sendInverterMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			warnings, err := axpert.WarningStatus(uc)
			if err != nil {
				panic(err)
			}
			msgData = messageData{Timestamp: t, MessageType: "Warnings", Data: warnings}
			err = sendInverterMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			flags, err := axpert.DeviceFlagStatus(uc)
			msgData = messageData{Timestamp: t, MessageType: "Flags", Data: flags}
			err = sendInverterMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			ratingInfo, err := axpert.DeviceRatingInfo(uc)
			msgData = messageData{Timestamp: t, MessageType: "RatingInfo", Data: ratingInfo}
			err = sendInverterMessage(msgData, client)
			if err != nil {
				panic(err)
			}

			batteryStatus, err := pylontech.GetBatteryStatus(sc)
			msgData = messageData{Timestamp: t, MessageType: "BatteryStatus", Data: batteryStatus}
			err = sendBatteryMessage(msgData, client)
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

func sendInverterMessage(data messageData, client mqtt.Client) error {
	return sendMessage(data, inverterTopic, client)
}

func sendBatteryMessage(data messageData, client mqtt.Client) error {
	return sendMessage(data, batteryTopic, client)
}

func sendMessage(data messageData, topic string, client mqtt.Client) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	token := client.Publish(topic, 1, true, msg)
	token.Wait()
	return nil
}

func logConnect(_ mqtt.Client) {
	fmt.Println("Connected to broker")
}

func logConnectionLost(_ mqtt.Client, err error) {
	fmt.Println("Connection lost:", err)
}

func initConfig() error {
	var configPath string
	pflag.StringVarP(&configPath, "config-path", "c", ".", "Path to config file (datalogd-conf.yaml)")
	pflag.Parse()

	viper.SetDefault("mqtt.server", "localhost")
	viper.SetDefault("mqtt.port", 1883)
	viper.SetDefault("mqtt.clientid", "datalogd")
	viper.SetDefault("timer.interval", 30)
	viper.SetDefault("inverter.count", 1)
	viper.SetDefault("inverter.topic", "datalogd/inverter")
	viper.SetDefault("battery.baud", 1200)
	viper.SetDefault("battery.topic", "datalogd/battery")

	viper.SetEnvPrefix("dlog")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("datalog-conf")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Println("Config file not found, relying on defaults/ENV")
	} else {
		return err
	}

	timerInterval = viper.GetInt("timer.interval")
	mqttServer = viper.GetString("mqtt.server")
	mqttPort = viper.GetInt("mqtt.port")
	mqttClientId = viper.GetString("mqtt.clientId")
	inverterPath = viper.GetString("inverter.path")
	inverterCount = viper.GetInt("inverter.count")
	inverterTopic = viper.GetString("inverter.topic")
	batteryPath = viper.GetString("battery.path")
	batteryBaud = viper.GetInt("battery.baud")
	batteryTopic = viper.GetString("battery.topic")

	return nil
}
