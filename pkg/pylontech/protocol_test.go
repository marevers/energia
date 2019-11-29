package pylontech

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func Test_lengthChecksum(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{name: "Negative", args: args{len: -1}, want: 0, wantErr: true},
		{name: "Too big", args: args{len: 5000}, want: 0, wantErr: true},
		{name: "Zero", args: args{len: 0}, want: 0, wantErr: false},
		{name: "Given example", args: args{len: 18}, want: 0xD012, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lengthChecksum(tt.args.len)
			if (err != nil) != tt.wantErr {
				t.Errorf("lengthChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lengthChecksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_frameChecksum(t *testing.T) {
	type args struct {
		info string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{name: "Doc example (lower)", args: args{info: "1203400456abcefe"}, want: 0xFC71, wantErr: false},
		{name: "Doc example (UPPER)", args: args{info: "1203400456ABCEFE"}, want: 0xFC71, wantErr: false},
		{name: "Version query", args: args{info: "2001464F0000"}, want: 0xFD99, wantErr: false},
		{name: "Version response", args: args{info: "200146000000"}, want: 0xFDB3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := frameChecksum(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("infoStrChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("infoStrChecksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeProtocolVersion(t *testing.T) {
	want := "~2001464F0000FD99\r"
	got, err := encodeProtocolVersion()
	fmt.Println(string(got))

	if err != nil {
		t.Errorf("encodeProtocolVersion() error = %v", err)
		return
	}

	if string(got) != want {
		t.Errorf("encodeProtocolVersion() got = %v, want %v", got, want)
	}

}

func Test_parseProtocolVersionResponse(t *testing.T) {
	want := "~200146000000FDB3\r"
	got, err := parseResponse([]byte(want))
	fmt.Println(got)

	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	bytes, err := got.encode()
	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	if string(bytes) != want {
		t.Errorf("parseResponse() got = %v, want %v", string(bytes), want)
	}
}

func Test_parseManufacturerInfoResponse(t *testing.T) {
	want := "~20014600C0405553324B42504C000000020150796C6F6E2D2D2D2D2D2D2D2D2D2D2D2D2D2D2DEF9B\r"
	got, err := parseResponse([]byte(want))
	fmt.Println(got)

	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	bytes, err := got.encode()
	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	if string(bytes) != want {
		t.Errorf("parseResponse() got = %v, want %v", string(bytes), want)
	}

}

func Test_parseManufacturerInfo(t *testing.T) {
	want := ManufacturerInfo{"US2KBPL", "21", "Pylon---------------"}

	resp := "~20014600C0405553324B42504C000000020150796C6F6E2D2D2D2D2D2D2D2D2D2D2D2D2D2D2DEF9B\r"

	f, err := parseResponse([]byte(resp))
	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	got, err := parseManufacturerInfo(f.info)
	if err != nil {
		t.Errorf("parseManufacturerInfo() error = %v", err)
		return
	}

	if got.DeviceName != want.DeviceName {
		t.Errorf("parseManufacturerInfo() got = [%v], want [%v]", got.DeviceName, want.DeviceName)
		return
	}

	if got.SoftwareVersion != want.SoftwareVersion {
		t.Errorf("parseManufacturerInfo() got = [%v], want [%v]", got.SoftwareVersion, want.SoftwareVersion)
		return
	}

	if got.ManufacturerName != want.ManufacturerName {
		t.Errorf("parseManufacturerInfo() got = [%v], want [%v]", got.ManufacturerName, want.ManufacturerName)
		return
	}

}

func Test_parseBatteryGroupStatusUS2000B(t *testing.T) {

	resp :=
		"~20014600B0D811020F0D6F0D6F0D6D0D6F0D6C0D6E0D6F0D6E0D760D780D760D780D770D780D76050BAF0B7D0B7D0B7D0B7D0000C9B2C35002C35000050F0DEE0DF80DF50DF20DF00DEE0DF60DF60E040E020E030E030E030E040E04050BB90B7D0B7D0B7D0B7D0000D1AEC35002C3500011CD77\r"

	f, err := parseResponse([]byte(resp))
	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	got, err := parseBatteryGroupStatus(f.info)
	bytes, err := json.Marshal(got)
	fmt.Println(string(bytes))

	if err != nil {
		t.Errorf("parseManufacturerInfo() error = %v", err)
		return
	}

}

func Test_parseBatteryGroupStatusUS3000A(t *testing.T) {

	resp :=
		"~2001460010F011020F0D1A0D220D220D200D1D0D210D1D0D190D1A0D1E0D210D1F0D1C0D1A0D1C050BB90BB90BB90BC30BB900BEC4BCFFFF04FFFF010A00BEC80121100F0D220D230D1F0D1C0D1C0D1C0D1C0D1A0D1C0D1D0D1D0D1C0D1C0D1C0D1D050BC30BB90BB90BB90BB900BDC4B5FFFF04FFFF010600B900012110C7D3\r"

	f, err := parseResponse([]byte(resp))
	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	got, err := parseBatteryGroupStatus(f.info)
	bytes, err := json.Marshal(got)
	fmt.Println(string(bytes))

	if err != nil {
		t.Errorf("parseManufacturerInfo() error = %v", err)
		return
	}

	if math.Abs(float64(74.0-got.Status[0].TotalCapacity)) > 0.01 {
		t.Error("parseBatteryGroupStatus(), total capacity should be 74 ")
	}

}

func Test_parseBatteryGroupStatusNegativeCurrent(t *testing.T) {

	resp :=
		"~2001460010F011020F0D6F0D780D770D770D740D740D740D720D7B0D790D7A0D7A0D790D7A0D77050BCD0BC30BC30BCD0BC3FFFEC9F5FFFF04FFFF010A0126D80121100F0D770D780D790D790D780D780D780D690D7A0D790D760D780D780D790D78050BCD0BC30BC30BCD0BC3FFFEC9FCFFFF04FFFF01060126D8012110C74C\r"

	f, err := parseResponse([]byte(resp))
	if err != nil {
		t.Errorf("parseResponse() error = %v", err)
		return
	}

	got, err := parseBatteryGroupStatus(f.info)
	bytes, err := json.Marshal(got)
	fmt.Println(string(bytes))

	if err != nil {
		t.Errorf("parseBatteryGroupStatus() error = %v", err)
		return
	}

	if got.Status[0].Current > 0 {
		t.Error("parseBatteryGroupStatus(), current should be negative ")
	}

}
