package keyboard

import (
	"syscall"
	"unsafe"
)

const (
	// EvKey describes the state changes of keyboards, buttons or other key-like devices
	EvKey EventType = 0x01
)

//
// EventType ...
//
type EventType uint16

var eventsize = int(unsafe.Sizeof(InputEvent{}))

//
// InputEvent ...
//
type InputEvent struct {
	Time  syscall.Timeval
	Type  EventType
	Code  uint16
	Value uint32
}

//
// KeyPress ...
//
func (i *InputEvent) KeyPress() bool {
	return i.Value == 1
}

//
// KeyRelease ...
//
func (i *InputEvent) KeyRelease() bool {
	return i.Value == 0
}
