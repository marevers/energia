package axpert

import (
	"fmt"
	"testing"
)

func TestParseFirmwareVersion(t *testing.T) {
	resp := string([]byte{86, 69, 82, 70, 87, 58, 48, 48, 48, 55, 50, 46, 55, 48})

	expectedFv := FirmwareVersion{"00072", "70"}

	fv, err := parseFirmwareVersion(resp, "VERFW")
	fmt.Println(fv)

	if err != nil {
		t.Error("expected no error, got", err)
	}
	if fv == nil {
		t.Error("expected result, got nil")
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
	if info == nil {
		t.Error("expected result, got nil")
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
	fmt.Println(flags)

	if err != nil {
		t.Error("expected no error, got", err)
	}
	if flags == nil {
		t.Error("expected result, got nil")
	}
	if expectedFlags != *flags {
		t.Error("expected ", expectedFlags, " got ", *flags)
	}

}

func TestParseDeviceFlagsWithLowercaseResponse(t *testing.T) {
	resp := "EabjklDuvxyz"

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
	fmt.Println(flags)

	if err != nil {
		t.Error("expected no error, got", err)
	}
	if flags == nil {
		t.Error("expected result, got nil")
	}
	if expectedFlags != *flags {
		t.Error("expected ", expectedFlags, " got ", *flags)
	}

}

func TestParseDeviceStatusParams(t *testing.T) {
	resp := "230.0 50.0 231.0 49.9 0300 0250 010 460 57.50 012 100 0069 0014 103.8 57.45 00000 00110110 00 07 00856 010"

	expectedParams := DeviceStatusParams{
		GridVoltage:                       230.0,
		GridFrequency:                     50.0,
		ACOutputVoltage:                   231.0,
		ACOutputFrequency:                 49.9,
		ACOutputApparentPower:             300,
		ACOutputActivePower:               250,
		OutputLoadPercent:                 10,
		BusVoltage:                        460,
		BatteryVoltage:                    57.5,
		BatteryChargingCurrent:            12,
		BatteryCapacity:                   100,
		HeatSinkTemperature:               69,
		PVInputCurrent1:                   14,
		PVInputVoltage1:                   103.8,
		BatteryVoltageSCC1:                57.45,
		BatteryDischargeCurrent:           0,
		AddSBUPriorityVersion:             false,
		ConfigStatusChanged:               false,
		SCCFirmwareVersionUpdated:         true,
		LoadOn:                            true,
		BatteryVoltageSteadyWhileCharging: false,
		ChargingOn:                        true,
		SCC1ChargingOn:                    true,
		ACChargingOn:                      false,
		FanBatteryVoltageOffset:           0,
		EEPROMVersion:                     "07",
		PVChargingPower1:                  856,
		FloatingModeCharging:              false,
		SwitchOn:                          true,
	}

	params, err := parseDeviceStatusParams(resp)
	fmt.Println(params)

	if err != nil {
		t.Error("expected no error, got", err)
	}
	if params == nil {
		t.Error("expected result, got nil")
	}
	if expectedParams != *params {
		t.Error("expected ", expectedParams, " got ", *params)
	}
}
