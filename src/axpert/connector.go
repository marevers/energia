package axpert

import "github.com/kristoiv/hid"

type Connector interface {
	Open() error
	Close()
	ReadUntilCR() ([]byte, error)
	Read(terminator byte) ([]byte, error)
	Write(bytes []byte) error
}

type USBConnector struct {
	path string
	deviceInfo *hid.DeviceInfo
	device hid.Device
}

func NewUSBConnector(path string) *USBConnector {
	return &USBConnector{path: path}
}

func (uc *USBConnector) DeviceInfo() *hid.DeviceInfo {
	return uc.deviceInfo
}

func (uc *USBConnector) Path() string {
	return uc.path
}

func (uc *USBConnector) Open() error {
	deviceInfo, err := hid.ByPath(uc.Path())
	if err != nil {
		return err
	}
	uc.deviceInfo = deviceInfo

	device, err := uc.deviceInfo.Open()
	if err != nil {
		return err
	}
	uc.device = device
	return nil
}

func (uc *USBConnector) Close()  {
	uc.device.Close()
}

func (uc *USBConnector) ReadUntilCR() ([]byte, error) {
	return uc.Read(0x0d)
}

// TODO This should have a timeout
func (uc *USBConnector) Read(terminator byte) ([]byte, error) {
	ch := uc.device.ReadCh()
	bytesRead := make([]byte, 0, 8)
	reading := true
	for reading {
		bs := <-ch
		for _, b := range bs {
			if b > 0 {
				bytesRead = append(bytesRead, b)
			}
			if b == terminator {
				reading = false
			}
		}
	}

	return bytesRead, nil
}

func (uc *USBConnector) Write(bytes []byte) error {
	return uc.device.Write(bytes)
}
