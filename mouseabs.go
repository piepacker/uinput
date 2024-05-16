package uinput

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// A MouseAbs is a device that will trigger an absolute change event.
// For details see: https://www.kernel.org/doc/Documentation/input/event-codes.txt
type MouseAbs interface {
	// MoveTo will move the cursor to the specified position on the screen
	MoveTo(x int32, y int32) error

	// LeftClick will issue a single left click.
	LeftClick() error

	// RightClick will issue a right click.
	RightClick() error

	// MiddleClick will issue a middle click.
	MiddleClick() error

	// LeftPress will simulate a press of the left mouse button. Note that the button will not be released until
	// LeftRelease is invoked.
	LeftPress() error

	// LeftRelease will simulate the release of the left mouse button.
	LeftRelease() error

	// RightPress will simulate the press of the right mouse button. Note that the button will not be released until
	// RightRelease is invoked.
	RightPress() error

	// RightRelease will simulate the release of the right mouse button.
	RightRelease() error

	// MiddlePress will simulate the press of the middle mouse button. Note that the button will not be released until
	// MiddleRelease is invoked.
	MiddlePress() error

	// MiddleRelease will simulate the release of the middle mouse button.
	MiddleRelease() error

	// Wheel will simulate a wheel movement.
	Wheel(horizontal bool, delta int32) error

	// FetchSysPath will return the syspath to the device file.
	FetchSyspath() (string, error)

	io.Closer
}

type vMouseAbs struct {
	name       []byte
	deviceFile *os.File
}

// CreateMouseAbs will create a new mouse input device. A mouseAbs is a device that allows absolute input.
func CreateMouseAbs(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (MouseAbs, error) {
	err := validateDevicePath(path)
	if err != nil {
		return nil, err
	}
	err = validateUinputName(name)
	if err != nil {
		return nil, err
	}

	fd, err := createMouseAbs(path, name, minX, maxX, minY, maxY)
	if err != nil {
		return nil, err
	}

	return vMouseAbs{name: name, deviceFile: fd}, nil
}

// MoveTo will move the cursor to the specified position on the screen
func (vAbs vMouseAbs) MoveTo(x int32, y int32) error {
	return vAbs.sendAbsEvent(x, y)
}

// LeftClick will issue a LeftClick.
func (vAbs vMouseAbs) LeftClick() error {
	err := sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnLeft}, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the LeftClick event: %v", err)
	}

	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnLeft}, btnStateReleased)
}

// RightClick will issue a RightClick
func (vAbs vMouseAbs) RightClick() error {
	err := sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnRight}, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the RightClick event: %v", err)
	}

	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnRight}, btnStateReleased)
}

// MiddleClick will issue a MiddleClick
func (vAbs vMouseAbs) MiddleClick() error {
	err := sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnMiddle}, btnStatePressed)
	if err != nil {
		return fmt.Errorf("Failed to issue the MiddleClick event: %v", err)
	}

	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnMiddle}, btnStateReleased)
}

// LeftPress will simulate a press of the left mouse button. Note that the button will not be released until
// LeftRelease is invoked.
func (vAbs vMouseAbs) LeftPress() error {
	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnLeft}, btnStatePressed)
}

// LeftRelease will simulate the release of the left mouse button.
func (vAbs vMouseAbs) LeftRelease() error {
	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnLeft}, btnStateReleased)
}

// RightPress will simulate the press of the right mouse button. Note that the button will not be released until
// RightRelease is invoked.
func (vAbs vMouseAbs) RightPress() error {
	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnRight}, btnStatePressed)
}

// RightRelease will simulate the release of the right mouse button.
func (vAbs vMouseAbs) RightRelease() error {
	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnRight}, btnStateReleased)
}

// MiddlePress will simulate the press of the middle mouse button. Note that the button will not be released until
// MiddleRelease is invoked.
func (vAbs vMouseAbs) MiddlePress() error {
	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnMiddle}, btnStatePressed)
}

// MiddleRelease will simulate the release of the middle mouse button.
func (vAbs vMouseAbs) MiddleRelease() error {
	return sendBtnEvent(vAbs.deviceFile, []int{evMouseBtnMiddle}, btnStateReleased)
}

// Wheel will simulate a wheel movement.
func (vAbs vMouseAbs) Wheel(horizontal bool, delta int32) error {
	w := relWheel
	if horizontal {
		w = relHWheel
	}
	return vAbs.sendRelEvent(uint16(w), delta)
}

// Close closes the device and releases the device.
func (vAbs vMouseAbs) Close() error {
	return closeDevice(vAbs.deviceFile)
}

func createMouseAbs(path string, name []byte, minX int32, maxX int32, minY int32, maxY int32) (fd *os.File, err error) {
	deviceFile, err := createDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not create absolute axis input device: %v", err)
	}

	err = registerDevice(deviceFile, uintptr(evKey))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register key device: %v", err)
	}

	// register button events (in order to enable left, right and middle click)
	for _, event := range []int{evMouseBtnLeft, evMouseBtnRight, evMouseBtnMiddle} {
		err = ioctl(deviceFile, uiSetKeyBit, uintptr(event))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register click event %v: %v", event, err)
		}
	}

	err = registerDevice(deviceFile, uintptr(evRel))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register relative axis input device: %v", err)
	}

	// register relative events
	for _, event := range []int{relWheel, relHWheel} {
		err = ioctl(deviceFile, uiSetRelBit, uintptr(event))
		if err != nil {
			deviceFile.Close()
			return nil, fmt.Errorf("failed to register relative event %v: %v", event, err)
		}
	}

	err = registerDevice(deviceFile, uintptr(evAbs))
	if err != nil {
		deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute axis input device: %v", err)
	}

	// register x and y-axis events
	for _, event := range []int{absX, absY} {
		err = ioctl(deviceFile, uiSetAbsBit, uintptr(event))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register absolute axis event %v: %v", event, err)
		}
	}

	var absMin [absSize]int32
	absMin[absX] = minX
	absMin[absY] = minY

	var absMax [absSize]int32
	absMax[absX] = maxX
	absMax[absY] = maxY

	return createUsbDevice(deviceFile,
		uinputUserDev{
			Name: toUinputName(name),
			ID: inputID{
				Bustype: busUsb,
				Vendor:  0x4711,
				Product: 0x0816,
				Version: 1},
			Absmin: absMin,
			Absmax: absMax})
}

func (vAbs vMouseAbs) sendAbsEvent(xPos int32, yPos int32) error { // TODO: Perhaps move this to a more generic function? This conflicts with the gamepad ABS events which only have one value.
	var ev [2]inputEvent
	ev[0].Type = evAbs
	ev[0].Code = absX
	ev[0].Value = xPos

	// Various tests (using evtest) have shown that positioning on x=0;y=0 doesn't trigger any event and will not move
	// the cursor as expected. Setting at least one of the coordinates to -1 will however have the desired effect of
	// moving the cursor to the upper left corner. Interestingly, the same is true for equivalent code in C, which rules
	// out issues related to Go's data type representation or the like. This will need to be investigated further...
	if xPos == 0 && yPos == 0 {
		yPos--
	}

	ev[1].Type = evAbs
	ev[1].Code = absY
	ev[1].Value = yPos

	for _, iev := range ev {
		buf, err := inputEventToBuffer(iev)
		if err != nil {
			return fmt.Errorf("writing abs event failed: %v", err)
		}

		_, err = vAbs.deviceFile.Write(buf)
		if err != nil {
			return fmt.Errorf("failed to write abs event to device file: %v", err)
		}
	}

	return syncEvents(vAbs.deviceFile)
}

func (vAbs vMouseAbs) sendRelEvent(eventCode uint16, pixel int32) error {
	iev := inputEvent{
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
		Type:  evRel,
		Code:  eventCode,
		Value: pixel}

	buf, err := inputEventToBuffer(iev)
	if err != nil {
		return fmt.Errorf("writing abs event failed: %v", err)
	}

	_, err = vAbs.deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write rel event to device file: %v", err)
	}

	return syncEvents(vAbs.deviceFile)
}

func (vAbs vMouseAbs) FetchSyspath() (string, error) {
	return fetchSyspath(vAbs.deviceFile)
}
