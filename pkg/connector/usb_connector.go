package connector

import (
	"errors"
	"fmt"
	"time"

	"github.com/sstallion/go-hid"
)

type USBConnector struct {
	deviceInfo *hid.DeviceInfo
	device     *hid.Device
}

func GetUSBPaths() (paths []string, err error) {
	paths = make([]string, 0)

	err = hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		paths = append(paths, info.Path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func NewUSBConnector(path string) (uc *USBConnector, err error) {
	device, err := hid.OpenPath(path)
	if err != nil {
		return nil, err
	}

	deviceInfo, err := device.GetDeviceInfo()
	if err != nil {
		return nil, err
	}

	return &USBConnector{deviceInfo: deviceInfo, device: device}, nil
}

func (uc *USBConnector) DeviceInfo() *hid.DeviceInfo {
	return uc.deviceInfo
}

func (uc *USBConnector) Path() string {
	return uc.deviceInfo.Path
}

func (uc *USBConnector) Open() error {
	// Device is already open - no need to open
	return nil
}

func (uc *USBConnector) Close() {
	uc.device.Close()
	uc.device = nil
}

func (uc *USBConnector) ReadUntilCR() ([]byte, error) {
	return uc.Read(0x0d)
}

// TODO This should take timout as argument or set by config
func (uc *USBConnector) Read(terminator byte) ([]byte, error) {
	bytesRead := make([]byte, 0, 8)
	buffer := make([]byte, 64) // HID report buffer
	timeout := 5 * time.Second

	// Create channels for communication with the read goroutine
	type readResult struct {
		data []byte
		n    int
		err  error
	}
	resultCh := make(chan readResult, 1)

	// Start a goroutine to perform the blocking read
	go func() {
		n, err := uc.device.Read(buffer)
		resultCh <- readResult{data: buffer, n: n, err: err}
	}()

	for {
		select {
		case result := <-resultCh:
			if result.err != nil {
				return nil, fmt.Errorf("failed to read from HID device: %w", result.err)
			}

			// Process the read bytes
			for i := 0; i < result.n; i++ {
				b := result.data[i]
				if b > 0 {
					bytesRead = append(bytesRead, b)
				}
				if b == terminator {
					return bytesRead, nil
				}
			}

			// If we haven't found the terminator, start another read
			go func() {
				n, err := uc.device.Read(buffer)
				resultCh <- readResult{data: buffer, n: n, err: err}
			}()

		case <-time.After(timeout):
			return nil, errors.New("timeout reading HID")
		}
	}
}

func (uc *USBConnector) Write(bytes []byte) error {
	_, err := uc.device.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
