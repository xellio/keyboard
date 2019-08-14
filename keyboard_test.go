package keyboard

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestFindDevices(t *testing.T) {
	devices := FindDevices()
	if len(devices) <= 0 {
		t.Error("No devices found.")
	}

	for _, device := range devices {
		num, err := strconv.Atoi(strings.Replace(device, "/dev/input/event", "", -1))
		if err != nil {
			t.Error("Unexpected device number.")
		}
		buff, err := ioutil.ReadFile(fmt.Sprintf(path, num))
		if err != nil {
			t.Error("Unable to open path.")
		}
		if !strings.Contains(strings.ToLower(string(buff)), "keyboard") {
			t.Error("Not a keyboard device.")
		}
	}
}
