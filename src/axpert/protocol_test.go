package axpert

import (
	"fmt"
	"testing"
)

func TestParseFirmwareVersion(t *testing.T) {
	resp := string([]byte{86, 69, 82, 70, 87, 58, 48, 48, 48, 55, 50, 46, 55, 48})

	expectedFv := FirmwareVersion{"00072", "70"}

	fv, err := parseFirmwareVersion(resp, "VERFW")
	if err != nil {
		t.Error("expected no error, got", err)
	}

	if expectedFv != *fv {
		t.Error("expected ", expectedFv, " got ", *fv)
	}
}

func TestParseRatingInfo(t *testing.T) {
	resp := string([]byte{50, 51, 48, 46, 48, 32, 50, 49, 46, 55, 32, 50, 51, 48, 46, 48, 32, 53, 48, 46, 48, 32,
		50, 49, 46, 55, 32, 53, 48, 48, 48, 32, 52, 48, 48, 48, 32, 52, 56, 46, 48, 32, 52, 56, 46, 48, 32,
		52, 55, 46, 53, 32, 53, 51, 46, 50, 32, 53, 49, 46, 57, 32, 50, 32, 51, 48, 32, 49, 50, 48, 32, 48, 32,
		48, 32, 49, 32, 57, 32, 48, 49, 32, 48, 32, 48, 32, 53, 49, 46, 48, 32, 48, 32, 49, 32, 48, 48, 48})

	expectedInfo := RatingInfo{GridRatingVoltage: 230, GridRatingCurrent: 21.7, ACOutputRatingVoltage: 230,
		ACOutputRatingFrequency: 50, ACOutputRatingCurrent: 21.7, ACOutputRatingApparentPower: 5000,
		ACOutputRatingActivePower: 4000, BatteryRatingVoltage: 48, BatteryRechargeVoltage: 48,
		BatteryUnderVoltage: 47.5, BatteryBulkVoltage: 53.2, BatteryFloatVoltage: 51.9, BatteryType: User,
		MaxACChargingCurrent: 30, MaxChargingCurrent: 120, InputVoltageRange: Appliance,
		OutputSourcePriority: OutputUtilityFirst, ChargerSourcePriority: ChargerSolarFirst, ParallelMaxNumber: 9,
		MachineType: OffGrid, Topology: Transfomerless, OutputMode: SingleMachine, BatteryRedischargeVoltage: 51,
		ParallelPVOK: AnyInverterConnected, PVPowerBalance: InputPowerIsChargedPowerPlusLoadPower}

	info, err := parseRatingInfo(resp)
	fmt.Println(info)
	if err != nil {
		t.Error("expected no error, got", err)
	}
	if expectedInfo != *info {
		t.Error("expected ", expectedInfo, " got ", *info)
	}

}

func TestParseDeviceFlags(t *testing.T) {
	resp := "EABJKLDUVXYZ"

	expectedFlags := DeviceFlags{
		Buzzer:                      FlagEnabled,
		OverloadBypass:              FlagEnabled,
		PowerSaving:                 FlagEnabled,
		DisplayTimeout:              FlagEnabled,
		OverloadRestart:             FlagDisabled,
		OverTemperatureRestart:      FlagDisabled,
		BacklightOn:                 FlagDisabled,
		PrimarySourceInterruptAlarm: FlagDisabled,
		FaultCodeRecord:             FlagDisabled,
		DataLogPopUp:                FlagEnabled,
	}

	flags, err := parseDeviceFlags(resp)
	if err != nil {
		t.Error("expected no error, got", err)
	}
	if expectedFlags != *flags {
		t.Error("expected ", expectedFlags, " got ", *flags)
	}

}
