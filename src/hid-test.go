package main

import (
	"fmt"
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


	fmt.Println("Opening USB device ", di.Path)
	d, err := di.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading from device" )
	ch := d.ReadCh()
	bytes := <- ch
	fmt.Println(bytes)

}
