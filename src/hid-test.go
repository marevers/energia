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

	fmt.Println("Closing connection")
	conn.Close()
}
