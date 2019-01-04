package axpert

import (
	"bytes"
	"testing"
)

func TestCrc(t *testing.T) {
	data := "(NAK"
	expectedCrc := []byte{0x73, 0x73}
	crc := crc([]byte(data))

	if !bytes.Equal(expectedCrc, crc) {
		t.Error("Expected ", expectedCrc, "got ", crc)
	}
}

func TestValidateResponse(t *testing.T) {
	data := "(NAKss\r"
	err := validateResponse([]byte(data))

	if err != nil {
		t.Error("Expected no error", "got ", err)
	}

}
