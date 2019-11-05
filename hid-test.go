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
		fmt.Println(err)
	}
	fmt.Println("ProtocolId: ", protocolId)

	serialNo, err := axpert.SerialNo(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("SerialNo: ", serialNo)

	version, err := axpert.InverterFirmwareVersion(conn)
	if err != nil {
		fmt.Println(err)
	}
	jsonVersion, err := json.Marshal(version)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("FirmwareVersion: ", string(jsonVersion))

	chargingTime, err := axpert.CVModeChargingTime(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("CV Mode Charging Time: ", chargingTime)

	chargingStage, err := axpert.DeviceChargingStage(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Charging stage: ", chargingStage)

	outputMode, err := axpert.DeviceOutputMode(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Output Mode: ", outputMode)

	bootstrapped, err := axpert.DSPBootstrapped(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DSPBootstrapped: ", bootstrapped)
	/*##
	maxSolarChargingCurrent, err := axpert.MaxSolarChargingCurrent(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("MaxSolarChargingCurrent: ", maxSolarChargingCurrent)

	maxUtilityChargingCurrent, err := axpert.MaxUtilityChargingCurrent(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("MaxUtilityChargingCurrent: ", maxUtilityChargingCurrent)

	maxTotalChargingCurrent, err := axpert.MaxTotalChargingCurrent(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("MaxTotalChargingCurrent: ", maxTotalChargingCurrent)
	*/

	defaults, err := axpert.DefaultSettings(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Default Settings: ", defaults)

	ratingInfo, err := axpert.DeviceRatingInfo(conn)
	if err != nil {
		fmt.Println(err)
	}
	jsonRating, err := json.Marshal(ratingInfo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Device rating info ", string(jsonRating))

	flags, err := axpert.DeviceFlagStatus(conn)
	if err != nil {
		fmt.Println(err)
	}
	jsonFlags, err := json.Marshal(flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Device status flags ", string(jsonFlags))

	deviceInfo, err := axpert.ParallelDeviceInfo(conn, 0)
	if err != nil {
		fmt.Println(err)
	}
	jsonInfo, err := json.MarshalIndent(deviceInfo, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Parallel Device 0 Info ", string(jsonInfo))

	device2Info, err := axpert.ParallelDeviceInfo(conn, 1)
	if err != nil {
		fmt.Println(err)
	}
	json2Info, err := json.MarshalIndent(device2Info, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Parallel Device 1 Info ", string(json2Info))

	params, err := axpert.DeviceGeneralStatus(conn)
	if err != nil {
		fmt.Println(err)
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Device status params ", string(jsonParams))

	// Remove due to timeout error
	//params, err = axpert.DeviceGeneralStatus2(conn, params)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("All Device status params ", params)

	mode, err := axpert.DeviceMode(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Device mode ", mode)

	warnings, err := axpert.WarningStatus(conn)
	if err != nil {
		fmt.Println(err)
	}
	jsonWarnings, err := json.Marshal(warnings)
	fmt.Println("Warning status ", string(jsonWarnings))

	fmt.Println("Closing connection")
	conn.Close()
}
