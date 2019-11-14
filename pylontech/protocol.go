package pylontech

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/tmthrgd/go-hex"

	"github.com/dbld-org/energia/internal/connector"
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
	Version     = 0x20
	start       = 0x7E
	end         = 0x0D
	batteryData = 0x46
)

func ProtocolVersion(c connector.Connector) (string, error) {
	encoded, err := encodeProtocolVersion()
	response, err := sendRequest(c, encoded)
	if err != nil {
		return "", err
	}

	return string(response), err
}

func sendRequest(c connector.Connector, encoded []byte) ([]byte, error) {
	err := c.Write(encoded)
	if err != nil {
		return nil, err
	}
	readBytes, err := c.ReadUntilCR()
	if err != nil {
		return nil, err
	}
	return readBytes, nil
}

func encodeProtocolVersion() ([]byte, error) {
	f := newFrame(1, GetProtocolVersion, nil)

	encode, err := f.encode()
	return encode, err
}

type frame struct {
	ver  byte
	adr  byte
	cid1 byte
	cid2 Command
	info []byte
}

func newFrame(address byte, command Command, info []byte) *frame {
	return &frame{
		ver:  Version,
		adr:  address,
		cid1: batteryData,
		cid2: command,
		info: info,
	}
}

func (f *frame) encode() ([]byte, error) {
	buf := bytes.Buffer{}
	info := hex.EncodeUpperToString(f.info)
	length, err := lengthChecksum(len(info))
	if err != nil {
		return nil, err
	}
	buf.WriteByte(start)
	data := fmt.Sprintf("%02X%02X%02X%02X%04X%s", f.ver, f.adr, f.cid1, byte(f.cid2), length, info)
	buf.WriteString(data)
	checksum, err := frameChecksum(data)
	if err != nil {
		return nil, err
	}
	buf.WriteString(fmt.Sprintf("%04X", checksum))
	buf.WriteByte(end)
	log.Println("Encoded: ", buf.String())
	return buf.Bytes(), nil
}

func lengthChecksum(len int) (uint16, error) {
	if len < 0 {
		return 0, fmt.Errorf("invalid length, must be >= 0")
	}
	if len > 0x0FFF {
		return 0, fmt.Errorf("invalid length, must be <= %d", 0xFFF)
	}

	if len == 0 {
		return 0, nil
	}

	ulen := uint16(len)
	length := (ulen & 0x000F) + ((ulen >> 4) & 0x000F) + ((ulen >> 8) & 0x000F)

	length = (^(length%0x10)+1)<<12 + ulen

	return length, nil
}

func frameChecksum(frameData string) (uint16, error) {
	bs := []byte(strings.ToUpper(frameData))
	var sum uint16
	for _, b := range bs {
		sum += uint16(b)
	}

	sum = ^uint16(uint32(sum)%0x10000) + 1

	return sum, nil
}
