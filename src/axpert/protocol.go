package axpert

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/howeyc/crc16"
)

const (
	cr        byte = 0x0d
	lf        byte = 0x0a
	leftParen byte = 0x28
)

func ProtocolId(c Connector) (id string, err error) {
	id, err = sendRequest(c, "QPI")
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

	version, err = parseFirmwareVersion(resp, "VERFW")
	return
}

func SCC1FirmwareVersion(c Connector) (version *FirmwareVersion, err error) {
	resp, err := sendRequest(c, "QVFW2")
	if err != nil {
		return
	}

	version, err = parseFirmwareVersion(resp, "VERFW2")
	return
}

func SCC2FirmwareVersion(c Connector) (version *FirmwareVersion, err error) {
	resp, err := sendRequest(c, "QVFW3")
	if err != nil {
		return
	}

	version, err = parseFirmwareVersion(resp, "VERFW3")
	return
}

func SCC3FirmwareVersion(c Connector) (version *FirmwareVersion, err error) {
	resp, err := sendRequest(c, "QVFW4")
	if err != nil {
		return
	}

	version, err = parseFirmwareVersion(resp, "VERFW4")
	return
}

func CVModeChargingTime(c Connector) (chargingTime string, err error) {
	chargingTime, err = sendRequest(c, "QCVT")
	return
}

func ChargingStage(c Connector) (chargingStage string, err error) {
	chargingStage, err = sendRequest(c, "QCST")
	return
}

func OutputMode(c Connector) (otputMode string, err error) {
	otputMode, err = sendRequest(c, "QOPM")
	return
}

func DSPBootstraped(c Connector) (hasBootstrap string, err error) {
	hasBootstrap, err = sendRequest(c, "QBOOT")
	return
}

func MaxSolarChargingCurrent(c Connector) (charginCurrent string, err error) {
	charginCurrent, err = sendRequest(c, "QMSCHGCR")
	return
}

func MaxUtilityChargingCurrent(c Connector) (charginCurrent string, err error) {
	charginCurrent, err = sendRequest(c, "QMUCHGCR")
	return
}

func MaxTotalChargingCurrent(c Connector) (charginCurrent string, err error) {
	charginCurrent, err = sendRequest(c, "QMCHGCR")
	return
}

func sendRequest(c Connector, req string) (resp string, err error) {
	reqBytes := []byte(req)
	reqBytes = append(reqBytes, crc(reqBytes)...)
	reqBytes = append(reqBytes, cr)
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

	resp = string(readBytes[1 : len(readBytes)-3])
	return
}

func validateResponse(read []byte) error {
	readLen := len(read)
	if read[0] != leftParen {
		return fmt.Errorf("invalid response start %x", read[0])
	}
	if read[readLen-1] != cr {
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
		if bs[i] == lf || bs[i] == cr || bs[i] == leftParen {
			bs[i] += 1
		}
	}
	return bs
}

func parseFirmwareVersion(resp string, fwPrefix string) (*FirmwareVersion, error) {
	parts := strings.Split(resp, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid response %s", resp)
	}
	if parts[0] != fwPrefix {
		return nil, fmt.Errorf("invalid prefix %s", parts[0])
	}

	version := strings.Split(parts[1], ".")
	if len(version) != 2 {
		return nil, fmt.Errorf("invalid version %s", parts[1])
	}

	return &FirmwareVersion{version[0], version[1]}, nil
}
