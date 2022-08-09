package device

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
	Transmit(dtype uint16, data []uint8, len int, dst any) error
}

const (
	NET_DEVICE_FLAG_UP uint16 = 0x0001

	NET_DEVICE_TYPE_DUMMY    = 0x0000
	NET_DEVICE_TYPE_LOOPBACK = 0x0001

	NET_DEVICE_FLAG_LOOPBACK = 0x0010
)

var devices *Device
var index = 0

func Head() *Device {
	return devices
}

func Register(d Device) {
	d.SetIndex(index)
	d.SetName(fmt.Sprintf("net%d", d.Index()))
	push(d)
	log.Infof("registered, dev=%s, type=0x%04x", d.Name(), d.Type())

	index += 1
}

func Open(d Device) error {
	if isUP(d) > 0 {
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
	log.Infof("dev=%s, state=%s", d.Name(), state(d))
	return nil
}

func Close(d Device) error {
	if isUP(d) == 0 {
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
	log.Infof("dev=%s, state=%s", d.Name(), state(d))
	return nil
}

func Output(d Device, dtype uint16, data []uint8, len int, dst any) error {
	if isUP(d) == 0 {
		err := fmt.Errorf("not opened, dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	if len > int(d.MTU()) {
		err := fmt.Errorf("too long, dev=%s, mtu=%x, len=%d", d.Name(), d.MTU(), len)
		log.Errorf(err.Error())
		return err
	}

	log.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), dtype, len)
	log.Debugdump(data, len)
	if err := d.Transmit(dtype, data, len, dst); err != nil {
		err := fmt.Errorf("device transmit failure, dev=%s, len=%d", d.Name(), len)
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func isUP(d Device) uint16 {
	return d.Flags() & NET_DEVICE_FLAG_UP
}

func state(d Device) string {
	if isUP(d) > 0 {
		return "up"
	} else {
		return "down"
	}
}

func push(d Device) {
	d.SetNext(devices)
	devices = &d
}
