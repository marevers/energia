package main

import (
	"fmt"

	"github.com/dbld-org/energia/src/axpert"
	"github.com/kristoiv/hid"
)

func main() {
	devs, err := hid.Devices()
	if err != nil {
		panic(err)
	}

	var di *hid.DeviceInfo
	for i, dev := range devs {
		fmt.Println(i, dev)
		di = dev
	}

	conn, err := axpert.NewUSBConnector(di.Path)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to ", conn.Path())
	err = conn.Open()
	if err != nil {
		panic(err)
	}

	protocolId, err := axpert.ProtocolId(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProtocolId: ", protocolId)

	serialNo, err := axpert.SerialNo(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("SerialNo: ", serialNo)

	version, err := axpert.InverterFirmwareVersion(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("FirmwareVersion: ", version)

	chargingTime, err := axpert.CVModeChargingTime(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("CV Mode Charging Time: ", chargingTime)

	chargingStage, err := axpert.ChargingStage(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Charging stage: ", chargingStage)

	outputMode, err := axpert.DeviceOutputMode(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Output Mode: ", outputMode)

	bootstraped, err := axpert.DSPBootstraped(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("DSPBootstraped: ", bootstraped)
	/*##
	maxSolarChargingCurrent, err := axpert.MaxSolarChargingCurrent(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("MaxSolarChargingCurrent: ", maxSolarChargingCurrent)

	maxUtilityChargingCurrent, err := axpert.MaxUtilityChargingCurrent(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("MaxUtilityChargingCurrent: ", maxUtilityChargingCurrent)

	maxTotalChargingCurrent, err := axpert.MaxTotalChargingCurrent(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("MaxTotalChargingCurrent: ", maxTotalChargingCurrent)
	*/

	defaults, err := axpert.DefaultSettings(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Default Settings: ", defaults)

	ratingInfo, err := axpert.DeviceRatingInfo(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Device rating info ", ratingInfo)

	fmt.Println("Closing connection")
	conn.Close()
}
