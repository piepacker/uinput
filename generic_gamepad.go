package uinput

import (
	"fmt"
	"os"
)

// CreateGamepad will create a new gamepad using the given uinput
// device path of the uinput device.
func CreateGenericGamepad(path string, bustype uint16, name []byte, vendor uint16, product uint16, version uint16, keys []uint16, absEvents []uint16, mscEvents []uint16) (Gamepad, error) {
	err := validateDevicePath(path)
	if err != nil {
		return nil, err
	}
	err = validateUinputName(name)
	if err != nil {
		return nil, err
	}

	fd, err := createVGenericGamepadDevice(path, bustype, name, vendor, product, version, keys, absEvents, mscEvents)
	if err != nil {
		return nil, err
	}

	return vGamepad{name: name, deviceFile: fd}, nil
}

func createVGenericGamepadDevice(path string, bustype uint16, name []byte, vendor uint16, product uint16, version uint16, keys []uint16, absEvents []uint16, mscEvents []uint16) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual gamepad device: %v", err)
	}

	// register button events
	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to register virtual gamepad device: %v", err)
	}

	for _, code := range keys {
		fmt.Printf("registering uiSetKeyBit 0x%x\n", code)
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(code))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register key number %d: %v", code, err)
		}
	}

	// register absolute events
	err = registerDevice(deviceFile, uintptr(evAbs))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute event input device: %v", err)
	}

	for _, event := range absEvents {
		fmt.Printf("registering uiSetAbsBit 0x%x\n", event)
		err = ioctl(deviceFile, uiSetAbsBit, uintptr(event))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register absolute event %v: %v", event, err)
		}
	}

	// misc event
	if len(mscEvents) > 0 {
		err = registerDevice(deviceFile, uintptr(evMsc))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register misc event input device: %v", err)
		}

		for _, event := range mscEvents {
			fmt.Printf("registering uiSetMscBit 0x%x\n", event)
			err = ioctl(deviceFile, uiSetMscBit, uintptr(event))
			if err != nil {
				_ = deviceFile.Close()
				return nil, fmt.Errorf("failed to register misc event %v: %v", event, err)
			}
		}

	}

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: bustype,
				Vendor:  vendor,
				Product: product,
				Version: version}})
}
