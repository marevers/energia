package pylontech

import (
	"fmt"
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

func Test_parseResponse(t *testing.T) {
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
