package axpert

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
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

func CVModeChargingTime(c Connector) (chargingTime uint8, err error) {
	time, err := sendRequest(c, "QCVT")
	if err != nil {
		return
	}
	b, err := strconv.ParseUint(time, 10, 8)
	if err != nil {
		return
	}
	chargingTime = uint8(b)
	return
}

func DeviceChargingStage(c Connector) (chargingStage ChargingStage, err error) {
	stage, err := sendRequest(c, "QCST")
	if err != nil {
		return
	}
	b, err := strconv.ParseUint(stage, 10, 8)
	if err != nil {
		return
	}
	chargingStage = ChargingStage(b)
	return
}

func DeviceOutputMode(c Connector) (outputMode OutputMode, err error) {
	mode, err := sendRequest(c, "QOPM")
	if err != nil {
		return
	}
	b, err := strconv.ParseUint(mode, 10, 8)
	if err != nil {
		return
	}
	outputMode = OutputMode(b)
	return
}

func DSPBootstrapped(c Connector) (hasBootstrap bool, err error) {
	bootstrap, err := sendRequest(c, "QBOOT")
	if err != nil {
		return
	}
	hasBootstrap = bootstrap != "0"
	return
}

func MaxSolarChargingCurrent(c Connector) (chargingCurrent string, err error) {
	chargingCurrent, err = sendRequest(c, "QMSCHGCR")
	return
}

func MaxUtilityChargingCurrent(c Connector) (chargingCurrent string, err error) {
	chargingCurrent, err = sendRequest(c, "QMUCHGCR")
	return
}

func MaxTotalChargingCurrent(c Connector) (chargingCurrent string, err error) {
	chargingCurrent, err = sendRequest(c, "QMCHGCR")
	return
}

func DefaultSettings(c Connector) (defaultSettings string, err error) {
	defaultSettings, err = sendRequest(c, "QDI")
	return
}

//go:generate enumer -type=BatteryType -json
type BatteryType uint8

const (
	AGM BatteryType = iota
	Flooded
	User
)

//go:generate enumer -type=VoltageRange -json
type VoltageRange uint8

const (
	Appliance VoltageRange = iota
	UPS
)

//go:generate enumer -type=OutputSourcePriority -json
type OutputSourcePriority uint8

const (
	OutputUtilityFirst OutputSourcePriority = iota
	OutputSolarFirst
	OutputSBUFirst
)

//go:generate enumer -type=ChargerSourcePriority -json
type ChargerSourcePriority uint8

const (
	ChargerUtilityFirst ChargerSourcePriority = iota
	ChargerSolarFirst
	ChargerSolarAndUtility
	ChargerSolarOnly
)

//go:generate enumer -type=MachineType -json
type MachineType uint8

const (
	GridTie          MachineType = 00
	OffGrid          MachineType = 01
	Hybrid           MachineType = 10
	OffGrid2Trackers MachineType = 11
	OffGrid3Trackers MachineType = 20
)

//go:generate enumer -type=Topology -json
type Topology uint8

const (
	Transfomerless Topology = iota
	Transformer
)

//go:generate enumer -type=OutputMode -json
type OutputMode uint8

const (
	SingleMachine OutputMode = iota
	Parallel
	Phase1
	Phase2
	Phase3
)

//go:generate enumer -type=ChargingStage -json
type ChargingStage uint8

const (
	Auto ChargingStage = iota
	TwoStage
	ThreeStage
)

//go:generate enumer -type=ParallelPVOK -json -text
type ParallelPVOK uint8

const (
	AnyInverterConnected ParallelPVOK = iota
	AllInvertersConnected
)

//go:generate enumer -type=PVPowerBalance -json -text
type PVPowerBalance uint8

const (
	InputCurrentIsChargedCurrent PVPowerBalance = iota
	InputPowerIsChargedPowerPlusLoadPower
)

type RatingInfo struct {
	GridRatingVoltage           float32
	GridRatingCurrent           float32
	ACOutputRatingVoltage       float32
	ACOutputRatingFrequency     float32
	ACOutputRatingCurrent       float32
	ACOutputRatingApparentPower int
	ACOutputRatingActivePower   int
	BatteryRatingVoltage        float32
	BatteryRechargeVoltage      float32
	BatteryUnderVoltage         float32
	BatteryBulkVoltage          float32
	BatteryFloatVoltage         float32
	BatteryType                 BatteryType
	MaxACChargingCurrent        int
	MaxChargingCurrent          int
	InputVoltageRange           VoltageRange
	OutputSourcePriority        OutputSourcePriority
	ChargerSourcePriority       ChargerSourcePriority
	ParallelMaxNumber           int
	MachineType                 MachineType
	Topology                    Topology
	OutputMode                  OutputMode
	BatteryRedischargeVoltage   float32
	ParallelPVOK                ParallelPVOK
	PVPowerBalance              PVPowerBalance
}

func DeviceRatingInfo(c Connector) (ratingInfo *RatingInfo, err error) {
	resp, err := sendRequest(c, "QPIRI")
	if err != nil {
		return
	}

	ratingInfo, err = parseRatingInfo(resp)
	return
}

//go:generate enumer -type=FlagStatus -json -text
type FlagStatus byte

const (
	FlagDisabled FlagStatus = iota
	FlagEnabled
)

func (s FlagStatus) char() byte {
	switch s {
	case FlagDisabled:
		return 'D'
	case FlagEnabled:
		return 'E'
	}
	return 0
}

//go:generate enumer -type=DeviceFlag -json -text
type DeviceFlag byte

const (
	Buzzer DeviceFlag = iota
	OverloadBypass
	PowerSaving
	DisplayTimeout
	OverloadRestart
	OverTemperatureRestart
	BacklightOn
	PrimarySourceInterruptAlarm
	FaultCodeRecord
	DataLogPopUp
)

func (f DeviceFlag) char() byte {
	switch f {
	case Buzzer:
		return 'A'
	case OverloadBypass:
		return 'B'
	case PowerSaving:
		return 'J'
	case DisplayTimeout:
		return 'K'
	case OverloadRestart:
		return 'U'
	case OverTemperatureRestart:
		return 'V'
	case BacklightOn:
		return 'X'
	case PrimarySourceInterruptAlarm:
		return 'Y'
	case FaultCodeRecord:
		return 'Z'
	case DataLogPopUp:
		return 'L'
	}
	return 0
}

func DeviceFlagStatus(c Connector) (flags map[DeviceFlag]FlagStatus, err error) {
	resp, err := sendRequest(c, "QFLAG")
	if err != nil {
		return
	}

	flags, err = parseDeviceFlags(resp)
	return
}

type DeviceStatusParams struct {
	GridVoltage                       float32
	GridFrequency                     float32
	ACOutputVoltage                   float32
	ACOutputFrequency                 float32
	ACOutputApparentPower             int
	ACOutputActivePower               int
	OutputLoadPercent                 int
	BusVoltage                        int
	BatteryVoltage                    float32
	BatteryChargingCurrent            int
	BatteryCapacity                   int
	HeatSinkTemperature               int
	PVInputCurrent1                   int
	PVInputVoltage1                   float32
	BatteryVoltageSCC1                float32
	PVInputCurrent2                   int
	PVInputVoltage2                   float32
	BatteryVoltageSCC2                float32
	PVInputCurrent3                   int
	PVInputVoltage3                   float32
	BatteryVoltageSCC3                float32
	BatteryDischargeCurrent           int
	AddSBUPriorityVersion             bool
	ConfigStatusChanged               bool
	SCCFirmwareVersionUpdated         bool
	LoadOn                            bool
	BatteryVoltageSteadyWhileCharging bool
	ChargingOn                        bool
	SCC1ChargingOn                    bool
	SCC2ChargingOn                    bool
	SCC3ChargingOn                    bool
	ACChargingOn                      bool
	FanBatteryVoltageOffset           int
	EEPROMVersion                     string
	PVChargingPower1                  int
	PVChargingPower2                  int
	PVChargingPower3                  int
	PVTotalChargingPower              int
	FloatingModeCharging              bool
	SwitchOn                          bool
	ACChargingCurrent                 int
	ACChargingPower                   int
}

func DeviceGeneralStatus(c Connector) (params *DeviceStatusParams, err error) {
	resp, err := sendRequest(c, "QPIGS")
	if err != nil {
		return
	}

	params, err = parseDeviceStatusParams(resp)

	return
}

func DeviceGeneralStatus2(c Connector, p *DeviceStatusParams) (params *DeviceStatusParams, err error) {
	resp, err := sendRequest(c, "QPIGS2")
	if err != nil {
		return
	}

	if p != nil {
		params = p
	}

	params, err = parseDeviceStatusParams2(resp, params)

	return
}

func DeviceMode(c Connector) (mode string, err error) {
	mode, err = sendRequest(c, "QMOD")
	return
}

//go:generate enumer -type=DeviceWarning -json
type DeviceWarning uint8

const (
	WarnReserved DeviceWarning = iota
	WarnInverterFault
	WarnBusOver
	WarnBusUnder
	WarnBusSoftFail
	WarnLineFail
	WarnOPVShort
	WarnInverterVoltageLow
	WarnInverterVoltageHigh
	WarnOverTemperature
	WarnFanLocked
	WarnBatteryVoltageHigh
	WarnBatteryLowAlarm
	WarnReservedOvercharge
	WarnBatteryShutdown
	WarnReservedBatteryDerating
	WarnOverload
	WarnEEPROMFault
	WarnInverterOverCurrent
	WarnInverterSoftFail
	WarnSelfTestFail
	WarnOPDCVoltageOver
	WarnBatteryOpen
	WarnCurrentSensorFail
	WarnBatteryShort
	WarnPowerLimit
	WarnPVVoltageHigh
	WarnMPPTOverloadFault
	WarnMPPTOverloadWarning
	WarnBatteryTooLowToCharge
	WarnPVVoltageHigh2
	WarnMPPTOverloadFault2
	WarnMPPTOverloadWarning2
	WarnBatteryTooLowToCharge2
	WarnPVVoltageHigh3
	WarnMPPTOverloadFault3
	WarnMPPTOverloadWarning3
	WarnBatteryTooLowToCharge3
)

func WarningStatus(c Connector) (warnings []DeviceWarning, err error) {
	status, err := sendRequest(c, "QPIWS")
	if err != nil {
		return
	}

	warnings, err = parseWarnings(status)
	return
}

func EnableDeviceFlags(c Connector, flags []DeviceFlag) error {
	command := formatDeviceFlags(flags, FlagEnabled)
	resp, err := sendRequest(c, command)
	if err != nil {
		return err
	}
	if resp == "NAK" {
		return fmt.Errorf("command not acknowledged, %v", command)
	}
	return nil
}

func DisableDeviceFlags(c Connector, flags []DeviceFlag) error {
	command := formatDeviceFlags(flags, FlagDisabled)
	return sendCommand(c, command)
}

func formatDeviceFlags(flags []DeviceFlag, status FlagStatus) string {
	cmdBuilder := new(strings.Builder)
	cmdBuilder.WriteByte('P')
	cmdBuilder.WriteByte(status.char())
	for _, flag := range flags {
		cmdBuilder.WriteByte(flag.char())
	}

	return cmdBuilder.String()
}

func SetOutputSourcePriority(c Connector, priority OutputSourcePriority) error {
	command := fmt.Sprintf("POP%02d", priority)
	return sendCommand(c, command)
}

func SetDefaultSettings(c Connector) error {
	command := "PF"
	return sendCommand(c, command)
}

func SetMaxTotalChargingCurrent(c Connector, current uint8, parallelNumber uint8) error {
	command := fmt.Sprintf("MCHGC%1d%03d", parallelNumber, current)
	return sendCommand(c, command)
}

func SetParallelMaxTotalChargingCurrent(c Connector, current uint8) error {
	command := fmt.Sprintf("MNCHGC%03d", current)
	return sendCommand(c, command)

}

func SetMaxUtilityChargingCurrent(c Connector, current uint8) error {
	command := fmt.Sprintf("MUCHGC%03d", current)
	return sendCommand(c, command)
}

func SetMaxSolarChargingCurrent(c Connector, current uint8) error {
	command := fmt.Sprintf("MSCHGC%03d", current)
	return sendCommand(c, command)
}

func SetOutputRatingFrequency(c Connector, frequency uint8) error {
	command := fmt.Sprintf("F%02d", frequency)
	return sendCommand(c, command)
}

func SetBatteryRechargeVoltage(c Connector, voltage float32) error {
	command := fmt.Sprintf("PBCV%.1f", voltage)
	return sendCommand(c, command)
}

func SetBatteryRedischargeVoltage(c Connector, voltage float32) error {
	command := fmt.Sprintf("PBDV%.1f", voltage)
	return sendCommand(c, command)
}

func SetChargerSourcePriority(c Connector, priority ChargerSourcePriority) error {
	command := fmt.Sprintf("PCP%02d", priority)
	return sendCommand(c, command)
}

func SetGridWorkingRange(c Connector, voltageRange VoltageRange) error {
	command := fmt.Sprintf("PGR%02d", voltageRange)
	return sendCommand(c, command)
}

func SetBatteryType(c Connector, batteryType BatteryType) error {
	command := fmt.Sprintf("PBT%02d", batteryType)
	return sendCommand(c, command)
}

func SetDeviceOutputMode(c Connector, mode OutputMode) error {
	command := fmt.Sprintf("POPM%02d", mode)
	return sendCommand(c, command)
}

func SetParallelChargerSourcePriority(c Connector, priority ChargerSourcePriority, parallelNumber uint8) error {
	command := fmt.Sprintf("PPCP%1d%02d", parallelNumber, priority)
	return sendCommand(c, command)
}

func SetBatteryCutoffVoltage(c Connector, voltage float32) error {
	command := fmt.Sprintf("PSDV%.1f", voltage)
	return sendCommand(c, command)
}

func SetCVModeChargingVoltage(c Connector, voltage float32) error {
	command := fmt.Sprintf("PCVV%.1f", voltage)
	return sendCommand(c, command)
}

func SetFloatChargingVoltage(c Connector, voltage float32) error {
	command := fmt.Sprintf("PBFT%.1f", voltage)
	return sendCommand(c, command)
}

func SetDeviceChargingStage(c Connector, mode OutputMode) error {
	command := fmt.Sprintf("PCST%02d", mode)
	return sendCommand(c, command)
}

func SetCVModeChargingTime(c Connector, chargingTime uint8) error {
	command := fmt.Sprintf("PCVT%03d", chargingTime)
	return sendCommand(c, command)
}

func SetParallelPVOK(c Connector, pvok ParallelPVOK) error {
	command := fmt.Sprintf("PPVOKC%1d", pvok)
	return sendCommand(c, command)
}

func SetPVPowerBalance(c Connector, balance PVPowerBalance) error {
	command := fmt.Sprintf("PSPB%1d", balance)
	return sendCommand(c, command)
}

func sendCommand(c Connector, command string) error {
	resp, err := sendRequest(c, command)
	if err != nil {
		return err
	}
	if resp == "NAK" {
		return fmt.Errorf("command not acknowledged, %v", command)
	}
	return nil
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
	i := crc16.Checksum(data, crc16.CCITTFalseTable)
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

func parseRatingInfo(resp string) (*RatingInfo, error) {
	parts := strings.Split(resp, " ")
	if len(parts) < 25 {
		return nil, fmt.Errorf("invalid response %s : not enough fields", resp)
	}

	info := RatingInfo{}

	f, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return nil, err
	}
	info.GridRatingVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return nil, err
	}
	info.GridRatingCurrent = float32(f)

	f, err = strconv.ParseFloat(parts[2], 32)
	if err != nil {
		return nil, err
	}
	info.ACOutputRatingVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[3], 32)
	if err != nil {
		return nil, err
	}
	info.ACOutputRatingFrequency = float32(f)

	f, err = strconv.ParseFloat(parts[4], 32)
	if err != nil {
		return nil, err
	}
	info.ACOutputRatingCurrent = float32(f)

	i, err := strconv.Atoi(parts[5])
	if err != nil {
		return nil, err
	}
	info.ACOutputRatingApparentPower = i

	i, err = strconv.Atoi(parts[6])
	if err != nil {
		return nil, err
	}
	info.ACOutputRatingActivePower = i

	f, err = strconv.ParseFloat(parts[7], 32)
	if err != nil {
		return nil, err
	}
	info.BatteryRatingVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[8], 32)
	if err != nil {
		return nil, err
	}
	info.BatteryRechargeVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[9], 32)
	if err != nil {
		return nil, err
	}
	info.BatteryUnderVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[10], 32)
	if err != nil {
		return nil, err
	}
	info.BatteryBulkVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[11], 32)
	if err != nil {
		return nil, err
	}
	info.BatteryFloatVoltage = float32(f)

	b, err := strconv.ParseUint(parts[12], 10, 8)
	if err != nil {
		return nil, err
	}
	info.BatteryType = BatteryType(b)

	i, err = strconv.Atoi(parts[13])
	if err != nil {
		return nil, err
	}
	info.MaxACChargingCurrent = i

	i, err = strconv.Atoi(parts[14])
	if err != nil {
		return nil, err
	}
	info.MaxChargingCurrent = i

	b, err = strconv.ParseUint(parts[15], 10, 8)
	if err != nil {
		return nil, err
	}
	info.InputVoltageRange = VoltageRange(b)

	b, err = strconv.ParseUint(parts[16], 10, 8)
	if err != nil {
		return nil, err
	}
	info.OutputSourcePriority = OutputSourcePriority(b)

	b, err = strconv.ParseUint(parts[17], 10, 8)
	if err != nil {
		return nil, err
	}
	info.ChargerSourcePriority = ChargerSourcePriority(b)

	i, err = strconv.Atoi(parts[18])
	if err != nil {
		return nil, err
	}
	info.ParallelMaxNumber = i

	b, err = strconv.ParseUint(parts[19], 10, 8)
	if err != nil {
		return nil, err
	}
	info.MachineType = MachineType(b)

	b, err = strconv.ParseUint(parts[20], 10, 8)
	if err != nil {
		return nil, err
	}
	info.Topology = Topology(b)

	b, err = strconv.ParseUint(parts[21], 10, 8)
	if err != nil {
		return nil, err
	}
	info.OutputMode = OutputMode(b)

	f, err = strconv.ParseFloat(parts[22], 32)
	if err != nil {
		return nil, err
	}
	info.BatteryRedischargeVoltage = float32(f)

	b, err = strconv.ParseUint(parts[23], 10, 8)
	if err != nil {
		return nil, err
	}
	info.ParallelPVOK = ParallelPVOK(b)

	b, err = strconv.ParseUint(parts[24], 10, 8)
	if err != nil {
		return nil, err
	}
	info.PVPowerBalance = PVPowerBalance(b)

	return &info, nil
}

func parseDeviceFlags(resp string) (map[DeviceFlag]FlagStatus, error) {
	flags := make(map[DeviceFlag]FlagStatus)

	if len(resp) < 2 {
		return nil, fmt.Errorf("response too short: %s", resp)
	}
	if strings.HasPrefix(resp, "E") {
		value := FlagEnabled
		for i := 1; i < len(resp); i++ {
			switch resp[i] {
			case 'A', 'a':
				flags[Buzzer] = value
			case 'B', 'b':
				flags[OverloadBypass] = value
			case 'J', 'j':
				flags[PowerSaving] = value
			case 'K', 'k':
				flags[DisplayTimeout] = value
			case 'L', 'l':
				flags[DataLogPopUp] = value
			case 'U', 'u':
				flags[OverloadRestart] = value
			case 'V', 'v':
				flags[OverTemperatureRestart] = value
			case 'X', 'x':
				flags[BacklightOn] = value
			case 'Y', 'y':
				flags[PrimarySourceInterruptAlarm] = value
			case 'Z', 'z':
				flags[FaultCodeRecord] = value
			case 'D':
				value = FlagDisabled
			default:
				return nil, fmt.Errorf("unknown flag %c", resp[i])
			}
		}
	}
	return flags, nil
}

func parseDeviceStatusParams(resp string) (*DeviceStatusParams, error) {
	parts := strings.Split(resp, " ")
	if len(parts) < 21 {
		return nil, fmt.Errorf("response too short: %s", resp)
	}

	params := DeviceStatusParams{}

	f, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return nil, err
	}
	params.GridVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return nil, err
	}
	params.GridFrequency = float32(f)

	f, err = strconv.ParseFloat(parts[2], 32)
	if err != nil {
		return nil, err
	}
	params.ACOutputVoltage = float32(f)

	f, err = strconv.ParseFloat(parts[3], 32)
	if err != nil {
		return nil, err
	}
	params.ACOutputFrequency = float32(f)

	i, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, err
	}
	params.ACOutputApparentPower = i

	i, err = strconv.Atoi(parts[5])
	if err != nil {
		return nil, err
	}
	params.ACOutputActivePower = i

	i, err = strconv.Atoi(parts[6])
	if err != nil {
		return nil, err
	}
	params.OutputLoadPercent = i

	i, err = strconv.Atoi(parts[7])
	if err != nil {
		return nil, err
	}
	params.BusVoltage = i

	f, err = strconv.ParseFloat(parts[8], 32)
	if err != nil {
		return nil, err
	}
	params.BatteryVoltage = float32(f)

	i, err = strconv.Atoi(parts[9])
	if err != nil {
		return nil, err
	}
	params.BatteryChargingCurrent = i

	i, err = strconv.Atoi(parts[10])
	if err != nil {
		return nil, err
	}
	params.BatteryCapacity = i

	i, err = strconv.Atoi(parts[11])
	if err != nil {
		return nil, err
	}
	params.HeatSinkTemperature = i

	i, err = strconv.Atoi(parts[12])
	if err != nil {
		return nil, err
	}
	params.PVInputCurrent1 = i

	f, err = strconv.ParseFloat(parts[13], 32)
	if err != nil {
		return nil, err
	}
	params.PVInputVoltage1 = float32(f)

	f, err = strconv.ParseFloat(parts[14], 32)
	if err != nil {
		return nil, err
	}
	params.BatteryVoltageSCC1 = float32(f)

	i, err = strconv.Atoi(parts[15])
	if err != nil {
		return nil, err
	}
	params.BatteryDischargeCurrent = i

	b, err := strconv.ParseUint(parts[16], 2, 8)
	if err != nil {
		return nil, err
	}

	sflags := uint8(b)
	params.AddSBUPriorityVersion = sflags&0x80 == 0x80
	params.ConfigStatusChanged = sflags&0x40 == 0x40
	params.SCCFirmwareVersionUpdated = sflags&0x20 == 0x20
	params.LoadOn = sflags&0x10 == 0x10
	params.BatteryVoltageSteadyWhileCharging = sflags&0x08 == 0x08
	params.ChargingOn = sflags&0x04 == 0x04
	params.SCC1ChargingOn = sflags&0x02 == 0x02
	params.ACChargingOn = sflags&0x01 == 0x01

	i, err = strconv.Atoi(parts[17])
	if err != nil {
		return nil, err
	}
	params.FanBatteryVoltageOffset = i

	params.EEPROMVersion = parts[18]

	i, err = strconv.Atoi(parts[19])
	if err != nil {
		return nil, err
	}
	params.PVChargingPower1 = i

	if len(parts[20]) < 3 {
		return nil, fmt.Errorf("invalid status %s", parts[20])
	}
	params.FloatingModeCharging = parts[20][0] == '1'
	params.SwitchOn = parts[20][1] == '1'

	return &params, nil
}

func parseDeviceStatusParams2(resp string, params *DeviceStatusParams) (*DeviceStatusParams, error) {
	parts := strings.Split(resp, " ")
	if len(parts) < 12 {
		return params, fmt.Errorf("response too short: %s", resp)
	}

	i, err := strconv.Atoi(parts[0])
	if err != nil {
		return params, err
	}
	params.PVInputCurrent2 = i

	f, err := strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return params, err
	}
	params.PVInputVoltage2 = float32(f)

	f, err = strconv.ParseFloat(parts[2], 32)
	if err != nil {
		return params, err
	}
	params.BatteryVoltageSCC2 = float32(f)

	i, err = strconv.Atoi(parts[3])
	if err != nil {
		return params, err
	}
	params.PVChargingPower2 = i

	b, err := strconv.ParseUint(parts[4], 2, 8)
	if err != nil {
		return params, err
	}

	sflags := uint8(b)
	params.SCC2ChargingOn = sflags&0x80 == 0x80
	params.SCC3ChargingOn = sflags&0x40 == 0x40

	i, err = strconv.Atoi(parts[5])
	if err != nil {
		return params, err
	}
	params.ACChargingCurrent = i

	i, err = strconv.Atoi(parts[6])
	if err != nil {
		return params, err
	}
	params.ACChargingPower = i

	i, err = strconv.Atoi(parts[7])
	if err != nil {
		return params, err
	}
	params.PVInputCurrent3 = i

	f, err = strconv.ParseFloat(parts[8], 32)
	if err != nil {
		return params, err
	}
	params.PVInputVoltage3 = float32(f)

	f, err = strconv.ParseFloat(parts[9], 32)
	if err != nil {
		return params, err
	}
	params.BatteryVoltageSCC3 = float32(f)

	i, err = strconv.Atoi(parts[10])
	if err != nil {
		return params, err
	}
	params.PVChargingPower3 = i

	i, err = strconv.Atoi(parts[11])
	if err != nil {
		return params, err
	}
	params.PVTotalChargingPower = i

	return params, nil
}

func parseWarnings(status string) ([]DeviceWarning, error) {
	if len(status) < 32 {
		return nil, fmt.Errorf("not enough status flags, %d", len(status))
	}

	if len(status) > 38 {
		return nil, fmt.Errorf("too many status flags, %d", len(status))
	}

	warnings := make([]DeviceWarning, 0)
	for i, c := range status {
		switch c {
		case '1':
			warnings = append(warnings, DeviceWarning(i))
		default:
			continue
		}
	}

	return warnings, nil
}
