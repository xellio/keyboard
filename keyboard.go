package keyboard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

//
// Device ...
//
type Device struct {
	fd *os.File
}

//
// New ...
//
func New(devicePath string) (*Device, error) {
	d := &Device{}
	if !d.IsRoot() {
		return nil, errors.New("Must be run as root")
	}

	fd, err := os.Open(devicePath)
	d.fd = fd
	return d, err
}

//
// FindDevice ...
//
func FindDevice(query string) (string, error) {

	device := ""
	path := "/sys/class/input/event%d/device/name"
	resolved := "/dev/input/event%d"

	for i := 0; i < 255; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(path, i))
		if err != nil {
			return device, err
		}

		if strings.Contains(strings.ToLower(string(buff)), query) {
			device = fmt.Sprintf(resolved, i)
			return device, nil
		}
	}

	return device, nil
}

//
// IsRoot ...
//
func (d *Device) IsRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
}

func (d *Device) Read() chan InputEvent {
	event := make(chan InputEvent)

	go func(event chan InputEvent) {
		for {
			e, err := d.read()
			if err != nil {
				close(event)
				break
			}

			if e != nil {
				event <- *e
			}
		}
	}(event)

	return event
}

func (d *Device) read() (*InputEvent, error) {
	buffer := make([]byte, eventsize)
	n, err := d.fd.Read(buffer)
	if err != nil {
		return nil, err
	}

	if n <= 0 {
		return nil, nil
	}

	return d.eventFromBuffer(buffer)
}

func (d *Device) eventFromBuffer(buffer []byte) (*InputEvent, error) {
	event := &InputEvent{}
	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
	return event, err
}

//
// Close ...
//
func (d *Device) Close() error {
	if d.fd == nil {
		return nil
	}
	return d.fd.Close()
}
