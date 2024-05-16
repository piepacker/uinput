package uinput

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// This test confirms that all basic mouse moves are working as expected.
func TestBasicMouseAbsMoves(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	defer func(relDev MouseAbs) {
		err := relDev.Close()
		if err != nil {
			t.Fatalf("failed to close virtual mouse: %v", err)
		}
	}(relDev)

	err = relDev.MoveTo(100, 200)
	if err != nil {
		t.Fatalf("Failed to move mouse left. Last error was: %s\n", err)
	}
}

func TestMouseAbsButtonPresses(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	defer func(relDev MouseAbs) {
		err := relDev.Close()
		if err != nil {
			t.Fatalf("failed to close virtual mouse: %v", err)
		}
	}(relDev)

	err = relDev.LeftPress()
	if err != nil {
		t.Fatalf("Failed to perform left key press. Last error was: %s\n", err)
	}

	err = relDev.LeftRelease()
	if err != nil {
		t.Fatalf("Failed to perform left key release. Last error was: %s\n", err)
	}

	err = relDev.RightPress()
	if err != nil {
		t.Fatalf("Failed to perform right key press. Last error was: %s\n", err)
	}

	err = relDev.RightRelease()
	if err != nil {
		t.Fatalf("Failed to perform right key release. Last error was: %s\n", err)
	}

	err = relDev.MiddlePress()
	if err != nil {
		t.Fatalf("Failed to perform middle key press. Last error was: %s\n", err)
	}

	err = relDev.MiddleRelease()
	if err != nil {
		t.Fatalf("Failed to perform middle key release. Last error was: %s\n", err)
	}
}

func TestVMouseAbs_Wheel(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	defer func(relDev MouseAbs) {
		err := relDev.Close()
		if err != nil {
			t.Fatalf("failed to close virtual mouse: %v", err)
		}
	}(relDev)

	err = relDev.Wheel(false, 1)
	if err != nil {
		t.Fatalf("Failed to perform wheel movement. Last error was: %s\n", err)
	}

	err = relDev.Wheel(true, 1)
	if err != nil {
		t.Fatalf("Failed to perform horizontal wheel movement. Last error was: %s\n", err)
	}
}

func TestMouseAbsClicks(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	defer func(relDev MouseAbs) {
		err := relDev.Close()
		if err != nil {
			t.Fatalf("failed to close virtual mouse: %v", err)
		}
	}(relDev)

	err = relDev.RightClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

	err = relDev.LeftClick()
	if err != nil {
		t.Fatalf("Failed to perform right click. Last error was: %s\n", err)
	}

	err = relDev.MiddleClick()
	if err != nil {
		t.Fatalf("Failed to perform middle click. Last error was: %s\n", err)
	}

}

func TestMouseAbsCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	_, err := CreateMouseAbs("", []byte("MouseAbsDevice"), 0, 1900, 0, 1080)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMouseAbsCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	_, err := CreateMouseAbs(path, []byte("MouseAbsDevice"), 0, 1900, 0, 1080)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestMouseAbsCreationFailsOnWrongPathName(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "uinput-mouse-test-")
	if err != nil {
		t.Fatalf("Failed to setup test. Unable to create tempfile: %v", err)
	}
	defer file.Close()

	expected := "failed to register key device: failed to close device: inappropriate ioctl for device"
	_, err = CreateMouseAbs(file.Name(), []byte("DialDevice"), 0, 1900, 0, 1080)
	if err == nil || !(expected == err.Error()) {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMouseAbsCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	_, err := CreateMouseAbs("/dev/uinput", []byte(name), 0, 1900, 0, 1080)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestMouseAbsLeftClickFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.LeftClick()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsLeftPressFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.LeftPress()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsLeftReleaseFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.LeftRelease()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsRightClickFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.RightClick()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsRightPressFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.RightPress()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestVMouseAbs_RightReleaseFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.RightRelease()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsMiddleClickFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MiddleClick()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsMiddlePressFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MiddlePress()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestVMouseAbs_MiddleReleaseFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MiddleRelease()
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsMoveToFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.MoveTo(1, 1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsWheelFailsIfDeviceIsClosed(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}
	relDev.Close()

	err = relDev.Wheel(false, 1)
	if err == nil {
		t.Fatalf("Expected error due to closed device, but no error was returned.")
	}
}

func TestMouseAbsSyspath(t *testing.T) {
	relDev, err := CreateMouseAbs("/dev/uinput", []byte("Test Basic MouseAbs"), 0, 1900, 0, 1080)
	if err != nil {
		t.Fatalf("Failed to create the virtual mouse. Last error was: %s\n", err)
	}

	sysPath, err := relDev.FetchSyspath()
	if err != nil {
		t.Fatalf("Failed to fetch syspath. Last error was: %s\n", err)
	}

	if sysPath[:32] != "/sys/devices/virtual/input/input" {
		t.Fatalf("Expected syspath to start with /sys/devices/virtual/input/input, but got %s", sysPath)
	}
	t.Logf("Syspath: %s", sysPath)
}
