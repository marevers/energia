package axpert

import (
	"bytes"
	"fmt"
	"github.com/howeyc/crc16"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func ProtocolId(c Connector) (id int, err error) {
	resp, err := sendRequest(c, "QPI")
	if err != nil {
		return
	}
	id, err = strconv.Atoi(resp)
	return
}

func SerialNo(c Connector) (serialNo string, err error) {
	serialNo, err = sendRequest(c, "QID")
	return
}

type FirmwareVersion struct {
	Series  string
	Version string
}

func InverterFirmwareVersion(c Connector) (version *FirmwareVersion, err error) {
	resp, err := sendRequest(c, "QVFW")
	if err != nil {
		return
	}

	version, err = parseFirmwareVersion(resp)
	return
}

func SCC1FirmwareVersion(c Connector) (version *FirmwareVersion, err error) {
	resp, err := sendRequest(c, "QVFW2")
	if err != nil {
		return
	}

	version, err = parseFirmwareVersion(resp)
	return
}

func SCC2FirmwareVersion(c Connector) (version *FirmwareVersion, err error) {
	resp, err := sendRequest(c, "QVFW3")
	if err != nil {
		return
	}

	version, err = parseFirmwareVersion(resp)
	return
}

func SCC3FirmwareVersion(c Connector) (version *FirmwareVersion, err error) {
	resp, err := sendRequest(c, "QVFW4")
	if err != nil {
		return
	}

	version, err = parseFirmwareVersion(resp)
	return
}

func sendRequest(c Connector, req string) (resp string, err error) {
	reqBytes := []byte(req)
	reqBytes = append(reqBytes, crc(reqBytes)...)
	reqBytes = append(reqBytes, 0x0d)
	log.Println("Sending ", reqBytes)
	err = c.Write(reqBytes)
	if err != nil {
		return
	}

	readBytes, err := c.ReadUntilCR()
	if err != nil {
		return
	}

	log.Println("Received ", readBytes)
	err = validateResponse(readBytes)
	if err != nil {
		return
	}

	resp = string(readBytes[1 : len(readBytes)-2])
	return
}

func validateResponse(read []byte) error {
	readLen := len(read)
	if read[0] != 0x28 {
		return fmt.Errorf("invalid response start %x", read[0])
	}
	if read[readLen-1] != 0x0d {
		return fmt.Errorf("invalid response end %x", read[readLen-1])
	}
	readCrc := read[readLen-3 : readLen-1]
	calcCrc := crc(read[:readLen-3])
	if !bytes.Equal(readCrc, calcCrc) {
		return fmt.Errorf("CRC error, received %v, expected %v", readCrc, calcCrc)
	}

	return nil
}

func crc(data []byte) []byte {
	i := crc16.Checksum(data, crc16.MakeBitsReversedTable(crc16.CCITTFalse))
	bs := []byte{uint8(i >> 8), uint8(i & 0xff)}
	for i := range bs {
		if bs[i] == 0x0a || bs[i] == 0x0d || bs[i] == 0x28 {
			bs[i] += 1
		}
	}
	return bs
}

func parseFirmwareVersion(s string) (*FirmwareVersion, error) {
	const fwPrefix = "VERFW:"
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid response %s", s)
	}
	if parts[0] != fwPrefix {
		return nil, fmt.Errorf("invalid prefix %s", parts[0])
	}

	version := strings.Split(parts[1], ".")
	if len(version) != 2 {
		return nil, fmt.Errorf("invalid version %s", parts[1])
	}
	if !isNumeric(version[0]) {
		return nil, fmt.Errorf("invalid firmware series %s", version[1])
	}
	if !isNumeric(version[1]) {
		return nil, fmt.Errorf("invalid firmware version number %s", version[1])
	}

	return &FirmwareVersion{version[0], version[1]}, nil
}

func isNumeric(s string) bool {
	numeric, _ := regexp.MatchString("[0-9]+", s)
	return numeric
}
