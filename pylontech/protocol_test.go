package pylontech

import "testing"

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

func Test_infoChecksum(t *testing.T) {
	type args struct {
		info []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{name: "Doc example", args: args{info: []byte{0x12, 0x03, 0x40, 0x04, 0x56, 0xAB, 0xCE, 0xFE}}, want: 0xFC71, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := infoChecksum(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("infoChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("infoChecksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_infoStrChecksum(t *testing.T) {
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
			got, err := infoStrChecksum(tt.args.info)
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
