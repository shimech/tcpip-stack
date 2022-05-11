package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/util/log"
)

type Device interface {
	Next() *Device
	SetNext(d *Device)
	Index() int
	SetIndex(i int)
	Name() string
	SetName(n string)
	Type() uint16
	MTU() uint16
	Flags() uint16
	SetFlags(f uint16)
	Hlen() uint16
	Alen() uint16
	Addr() uint8
	Peer() uint8
	Broadcast() uint8
	Open() error
	Close() error
	Transmit(dtype uint16, data []byte, dst any) error
}

const (
	NET_DEVICE_FLAG_UP uint16 = 0x0001

	NET_DEVICE_TYPE_DUMMY    = 0x0000
	NET_DEVICE_TYPE_LOOPBACK = 0x0001

	NET_DEVICE_FLAG_LOOPBACK = 0x0010
)

var devices *Device
var index = 0

func HeadDevice() *Device {
	return devices
}

func RegisterDevice(d Device) {
	d.SetIndex(index)
	d.SetName(fmt.Sprintf("net%d", d.Index()))
	pushNewDevice(d)
	log.Infof("registered, dev=%s, type=0x%04x", d.Name(), d.Type())

	index += 1
}

func openDevice(d Device) error {
	if isDeviceUP(d) > 0 {
		err := fmt.Errorf("already opened, dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	if err := d.Open(); err != nil {
		err := fmt.Errorf("failure dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	d.SetFlags(d.Flags() | NET_DEVICE_FLAG_UP)
	log.Infof("dev=%s, state=%s", d.Name(), deviceState(d))
	return nil
}

func closeDevice(d Device) error {
	if isDeviceUP(d) == 0 {
		err := fmt.Errorf("not opened, dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	if err := d.Close(); err != nil {
		err := fmt.Errorf("failure dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	d.SetFlags(d.Flags() & ^NET_DEVICE_FLAG_UP)
	log.Infof("dev=%s, state=%s", d.Name(), deviceState(d))
	return nil
}

func isDeviceUP(d Device) uint16 {
	return d.Flags() & NET_DEVICE_FLAG_UP
}

func deviceState(d Device) string {
	if isDeviceUP(d) > 0 {
		return "up"
	} else {
		return "down"
	}
}

func pushNewDevice(d Device) {
	d.SetNext(devices)
	devices = &d
}
