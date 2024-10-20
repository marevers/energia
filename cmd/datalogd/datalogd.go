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

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/goburrow/serial"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/marevers/energia/pkg/axpert"
	"github.com/marevers/energia/pkg/connector"
	"github.com/marevers/energia/pkg/pylontech"
)

var timerInterval int

var mqttServer string
var mqttPort int
var mqttClientId string
var mqttUsername string
var mqttPassword string

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

type queryFunc func(chan connector.Connector, mqtt.Client, time.Time) error

type query struct {
	f        queryFunc
	cc       chan connector.Connector
	interval time.Duration
}

var ucc chan connector.Connector

func main() {
	fmt.Println("initializing config ")

	err := initConfig()
	fmt.Println("initialized config ", viper.AllSettings())
	if err != nil {
		panic(err)
	}

	fmt.Println("connecting to ", inverterPath)
	uc, err := connector.NewUSBConnector(inverterPath)
	if err != nil {
		panic(err)
	}
	err = uc.Open()
	if err != nil {
		panic(err)
	}
	defer uc.Close()

	ucc = make(chan connector.Connector, 1)
	ucc <- uc
	fmt.Println("connected to ", inverterPath)

	var sc connector.Connector
	var scc chan connector.Connector

	if viper.IsSet("battery.path") {

		serialConfig := serial.Config{
			Address:  batteryPath,
			BaudRate: batteryBaud,
			DataBits: 8,
			StopBits: 1,
			Parity:   "N",
			Timeout:  30 * time.Second,
		}

		sc = connector.NewSerialConnector(serialConfig)
		err = sc.Open()
		if err != nil {
			log.Panic(err)
		}
		defer sc.Close()

		scc = make(chan connector.Connector, 1)
		scc <- sc
	}

	clientOpts := mqtt.NewClientOptions()
	clientOpts.AddBroker("tcp://" + mqttServer + ":" + strconv.Itoa(mqttPort))
	clientOpts.SetAutoReconnect(true)
	clientOpts.SetStore(mqtt.NewFileStore("/tmp/mqtt"))
	clientOpts.SetCleanSession(false)
	clientOpts.SetClientID(mqttClientId)
	clientOpts.SetOnConnectHandler(logConnect)
	clientOpts.SetConnectionLostHandler(logConnectionLost)
	clientOpts.SetUsername(mqttUsername)
	clientOpts.SetPassword(mqttPassword)

	client := mqtt.NewClient(clientOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)
	fmt.Println("Connected to mqtt")

	queries := []query{
		{deviceMode, ucc, 30 * time.Second},
		{parallelDeviceInfo, ucc, 30 * time.Second},
		{deviceGeneralStatus, ucc, 10 * time.Second},
		{deviceFlagStatus, ucc, 30 * time.Second},
		{warningStatus, ucc, 30 * time.Second},
		{deviceRating, ucc, 30 * time.Second},
	}

	if viper.IsSet("battery.path") {
		queries = append(queries, query{batteryStatus, scc, 10 * time.Second})
	}

	ts := make([]*time.Ticker, len(queries))

	for i, q := range queries {
		ts[i] = schedule(q.f, q.interval, q.cc, client)
	}

	client.Subscribe("inverter/cmd/setOutputSourcePriority", 1, messageReceiver)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	close(sigChan)
	fmt.Println(sig, " stopping tickers")

	for _, t := range ts {
		t.Stop()
	}

	fmt.Println("exiting")
}

func schedule(f queryFunc, interval time.Duration, ucc chan connector.Connector, client mqtt.Client) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for t := range ticker.C {
			f(ucc, client, t)
		}
	}()
	return ticker
}

func deviceGeneralStatus(ucc chan connector.Connector, client mqtt.Client, t time.Time) error {

	uc := <-ucc

	defer func() { ucc <- uc }()

	status, err := axpert.DeviceGeneralStatus(uc)
	if err != nil {
		return err
	}
	msgData := messageData{Timestamp: t, MessageType: "Status", Data: status}
	err = sendInverterMessage(msgData, client)
	if err != nil {
		return err
	}

	return nil

}

func warningStatus(ucc chan connector.Connector, client mqtt.Client, t time.Time) error {

	uc := <-ucc

	defer func() { ucc <- uc }()

	warnings, err := axpert.WarningStatus(uc)
	if err != nil {
		return err
	}
	msgData := messageData{Timestamp: t, MessageType: "Warnings", Data: warnings}
	err = sendInverterMessage(msgData, client)
	if err != nil {
		return err
	}
	return nil

}

func deviceFlagStatus(ucc chan connector.Connector, client mqtt.Client, t time.Time) error {

	uc := <-ucc

	defer func() { ucc <- uc }()

	flags, err := axpert.DeviceFlagStatus(uc)
	msgData := messageData{Timestamp: t, MessageType: "Flags", Data: flags}
	err = sendInverterMessage(msgData, client)
	if err != nil {
		return err
	}
	return nil

}

func deviceRating(ucc chan connector.Connector, client mqtt.Client, t time.Time) error {

	uc := <-ucc

	defer func() { ucc <- uc }()

	ratingInfo, err := axpert.DeviceRatingInfo(uc)
	msgData := messageData{Timestamp: t, MessageType: "RatingInfo", Data: ratingInfo}
	err = sendInverterMessage(msgData, client)
	if err != nil {
		return err
	}

	return nil

}

func batteryStatus(ucc chan connector.Connector, client mqtt.Client, t time.Time) error {

	uc := <-ucc

	defer func() { ucc <- uc }()

	batteryStatus, err := pylontech.GetBatteryStatus(uc)
	msgData := messageData{Timestamp: t, MessageType: "BatteryStatus", Data: batteryStatus}
	err = sendBatteryMessage(msgData, client)
	if err != nil {
		return err
	}

	return nil

}

func parallelDeviceInfo(ucc chan connector.Connector, client mqtt.Client, t time.Time) error {

	uc := <-ucc

	defer func() { ucc <- uc }()

	for inv := 0; inv < inverterCount; inv++ {
		deviceInfo, err := axpert.ParallelDeviceInfo(uc, inv)
		if err != nil {
			return err
		}
		msgData := messageData{Timestamp: t, MessageType: "DeviceInfo", Data: deviceInfo}
		err = sendInverterMessage(msgData, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func deviceMode(ucc chan connector.Connector, client mqtt.Client, t time.Time) error {

	uc := <-ucc

	defer func() { ucc <- uc }()

	mode, err := axpert.DeviceMode(uc)
	if err != nil {
		return err
	}
	m := map[string]string{"Mode": mode}
	msgData := messageData{Timestamp: t, MessageType: "Mode", Data: m}
	err = sendInverterMessage(msgData, client)
	if err != nil {
		return err
	}

	return nil
}

func sendInverterMessage(data messageData, client mqtt.Client) error {
	return sendMessage(data, inverterTopic+"/"+data.MessageType, client)
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

func messageReceiver(client mqtt.Client, msg mqtt.Message) {

	go func() {
		switch msg.Topic() {
		case "inverter/cmd/setOutputSourcePriority":
			fmt.Printf("%s\n", msg.Topic())
			msg.Topic()
			msg.Payload()
			uc := <-ucc

			defer func() { ucc <- uc }()
			priority, err := strconv.Atoi(string(msg.Payload()))
			if err != nil {
				fmt.Println("Value conversion error", err)
				return
			}

			err = axpert.SetOutputSourcePriority(uc, axpert.OutputSourcePriority(priority))
			if err != nil {
				fmt.Println("Failed sending command ", err)
				return
			}

		default:
			fmt.Printf("%s", msg.Topic())
		}
	}()
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

	viper.SetConfigName("datalogd-conf")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, relying on defaults/ENV")
		} else {
			return err
		}
	}

	fmt.Println("config: ", viper.AllSettings())
	timerInterval = viper.GetInt("timer.interval")
	mqttServer = viper.GetString("mqtt.server")
	mqttPort = viper.GetInt("mqtt.port")
	mqttUsername = viper.GetString("mqtt.username")
	mqttPassword = viper.GetString("mqtt.password")
	mqttClientId = viper.GetString("mqtt.clientId")
	inverterPath = viper.GetString("inverter.path")
	inverterCount = viper.GetInt("inverter.count")
	inverterTopic = viper.GetString("inverter.topic")
	batteryPath = viper.GetString("battery.path")
	batteryBaud = viper.GetInt("battery.baud")
	batteryTopic = viper.GetString("battery.topic")

	return nil
}
