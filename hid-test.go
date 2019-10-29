package main

import (
	"encoding/json"
	"fmt"

	"github.com/kristoiv/hid"

	"github.com/dbld-org/energia/axpert"
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
	jsonVersion, err := json.Marshal(version)
	if err != nil {
		panic(err)
	}
	fmt.Println("FirmwareVersion: ", string(jsonVersion))

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

	bootstrapped, err := axpert.DIPBootstrapped(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("DIPBootstrapped: ", bootstrapped)
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
	jsonRating, err := json.Marshal(ratingInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Device rating info ", string(jsonRating))

	flags, err := axpert.DeviceFlagStatus(conn)
	if err != nil {
		panic(err)
	}
	jsonFlags, err := json.Marshal(flags)
	if err != nil {
		panic(err)
	}
	fmt.Println("Device status flags ", string(jsonFlags))

	params, err := axpert.DeviceGeneralStatus(conn)
	if err != nil {
		panic(err)
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Device status params ", string(jsonParams))

	// Remove due to timeout error
	//params, err = axpert.DeviceGeneralStatus2(conn, params)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("All Device status params ", params)
	//

	mode, err := axpert.DeviceMode(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Device mode ", mode)

	warnings, err := axpert.WarningStatus(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Warning status ", warnings)

	fmt.Println("Closing connection")
	conn.Close()
}
