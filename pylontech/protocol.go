package pylontech

import (
	"encoding/hex"
	"fmt"
	"strings"
)

//go:generate enumer -type=Command -json
type Command uint8

const (
	GetAnalogValue          Command = 0x42
	GetAlarmData                    = 0x44
	GetSystemParameter              = 0x47
	GetProtocolVersion              = 0x4F
	GetManufacturerInfo             = 0x51
	GetChargeManagementInfo         = 0x92
	GetSeriesNumber                 = 0x93
	SetChargeManagementInfo         = 0x94
	TurnOff                         = 0x95
)

const (
	Version = 2.8
	start   = 0x7E
	end     = 0x0D
)

type frame struct {
	ver    byte
	adr    byte
	cid1   byte
	cid2   byte
	length uint16 // 4 bit length checksum + 12 bit info length
	info   []byte
	chksum uint16
}

func lengthChecksum(len int) (uint16, error) {
	if len < 0 {
		return 0, fmt.Errorf("invalid length, must be >= 0")
	}
	if len > 0x0FFF {
		return 0, fmt.Errorf("invalid length, must be <= %d", 0xFFF)
	}

	ulen := uint16(len)
	length := (ulen & 0x000F) + ((ulen >> 4) & 0x000F) + ((ulen >> 8) & 0x000F)

	length = (^(length%0x10)+1)<<12 + ulen

	return length, nil
}

func infoChecksum(info []byte) (uint16, error) {
	enc := hex.EncodeToString(info)
	return infoStrChecksum(enc)
}

func infoStrChecksum(info string) (uint16, error) {
	bs := []byte(strings.ToUpper(info))

	var sum uint16
	for _, b := range bs {
		sum += uint16(b)
	}

	sum = ^uint16(uint32(sum)%0x10000) + 1

	return sum, nil
}
