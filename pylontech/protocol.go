package pylontech

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tmthrgd/go-hex"

	"github.com/dbld-org/energia/internal/connector"
)

//go:generate enumer -type=command -json
type command uint8

const (
	getBatteryStatus        command = 0x42
	getAlarmData                    = 0x44
	getSystemParameter              = 0x47
	getProtocolVersion              = 0x4F
	getManufacturerInfo             = 0x51
	getChargeManagementInfo         = 0x92
	getSeriesNumber                 = 0x93
	setChargeManagementInfo         = 0x94
	turnOff                         = 0x95
)

const (
	AllBatteries   = 0xFF
	defaultVersion = 0x20
	start          = 0x7E
	end            = 0x0D
	batteryData    = 0x46
)

func GetProtocolVersion(c connector.Connector) (string, error) {
	encoded, err := encodeProtocolVersion()
	if err != nil {
		return "", err
	}

	response, err := sendRequest(c, encoded)
	if err != nil {
		return "", err
	}

	decoded, err := parseResponse(response)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%02X", decoded.ver), err
}

type ManufacturerInfo struct {
	DeviceName       string
	SoftwareVersion  string
	ManufacturerName string
}

func GetManufacturerInfo(c connector.Connector) (*ManufacturerInfo, error) {
	encoded, err := encodeManufacturerInfo()
	if err != nil {
		return nil, err
	}

	response, err := sendRequest(c, encoded)
	if err != nil {
		return nil, err
	}

	decoded, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return parseManufacturerInfo(decoded.info)
}

func parseManufacturerInfo(info []byte) (*ManufacturerInfo, error) {
	man := &ManufacturerInfo{
		DeviceName: strings.TrimFunc(string(info[0:10]), func(r rune) bool {
			return r < 32
		}),
		SoftwareVersion:  fmt.Sprintf("%d%d", info[10], info[11]),
		ManufacturerName: string(info[12:]),
	}
	return man, nil
}

func GetBatteryStatus(c connector.Connector) (string, error) {
	encoded, err := encodeBatteryStatus(1, AllBatteries)
	if err != nil {
		return "", err
	}

	response, err := sendRequest(c, encoded)
	if err != nil {
		return "", err
	}

	decoded, err := parseResponse(response)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", decoded.info), err
}

func encodeBatteryStatus(address byte, batteryNumber byte) ([]byte, error) {
	if batteryNumber == 0 {
		batteryNumber = AllBatteries
	}

	f := newFrame(address, getBatteryStatus, []byte{batteryNumber})

	encode, err := f.encode()
	return encode, err
}

func encodeManufacturerInfo() ([]byte, error) {
	f := newFrame(1, getManufacturerInfo, nil)

	encode, err := f.encode()
	return encode, err
}

func encodeProtocolVersion() ([]byte, error) {
	f := newFrame(1, getProtocolVersion, nil)

	encode, err := f.encode()
	return encode, err
}

func parseResponse(response []byte) (*frame, error) {
	log.Printf("received response: [%s]", string(response))
	respData, err := validateResponse(response)
	if err != nil {
		return nil, err
	}

	f := &frame{}
	f.ver = hex2Byte(respData[0:2])
	f.adr = hex2Byte(respData[2:4])
	f.cid1 = hex2Byte(respData[4:6])
	f.cid2 = command(hex2Byte(respData[6:8]))

	infoLen := uint16(hex2Byte(respData[8:10]))<<8 | uint16(hex2Byte(respData[10:12]))
	var info []byte
	if len(respData) > 12 {
		info = respData[12:]
	}

	lenCheck, err := lengthChecksum(len(info))
	if err != nil {
		return nil, err
	}
	if lenCheck != infoLen {
		return nil, fmt.Errorf("invalid length, received %v, calculated %v", infoLen, lenCheck)
	}
	f.info = hex2Bytes(info)

	return f, nil
}

func hex2Bytes(hexBytes []byte) []byte {
	hexLen := len(hexBytes)
	if hexLen%2 != 0 {
		return nil
	}

	bs := make([]byte, 0, hexLen/2)
	for i := 0; i < hexLen; i += 2 {
		bs = append(bs, hex2Byte(hexBytes[i:i+2]))
	}
	return bs
}

func hex2Byte(hexBytes []byte) byte {
	if len(hexBytes) > 2 {
		return 0
	}
	parsed, err := strconv.ParseUint(string(hexBytes), 16, 16)
	if err != nil {
		return 0
	}

	return byte(parsed)
}

func validateResponse(response []byte) ([]byte, error) {
	rlen := len(response)
	if rlen == 0 {
		return nil, fmt.Errorf("response is empty")
	}
	if response[0] != start {
		return nil, fmt.Errorf("invalid response start %v", response[0])
	}
	if response[rlen-1] != end {
		return nil, fmt.Errorf("invalid response end %v", response[rlen-1])
	}
	checkStart := rlen - 5
	respData := response[1:checkStart]
	respCheck := string(response[checkStart : rlen-1])
	dataSum, err := frameChecksum(string(respData))
	if err != nil {
		return nil, err
	}
	checkSum, err := strconv.ParseUint(respCheck, 16, 16)
	if err != nil {
		return nil, err
	}
	if uint16(checkSum) != dataSum {
		return nil, fmt.Errorf("invalid checksum, received: %v, calculated: %v", checkSum, dataSum)
	}

	return respData, nil
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

type frame struct {
	ver  byte
	adr  byte
	cid1 byte
	cid2 command
	info []byte
}

func newFrame(address byte, command command, info []byte) *frame {
	return &frame{
		ver:  defaultVersion,
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
	log.Printf("Encoded: %s", buf.String())
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
