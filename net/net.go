package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/platform/linux/intr"
	"github.com/shimech/tcpip-stack/util/log"
)

const (
	NET_DEVICE_FLAG_LOOPBACK = 0x0010
)

var devices *Device
var index = 0

func Register(d Device) {
	d.SetIndex(index)
	d.SetName(fmt.Sprintf("net%d", d.Index()))
	d.SetNext(devices)
	devices = &d
	log.Infof("registered, dev=%s, type=0x%04x", d.Name(), d.Type())

	index += 1
}

func Open(d Device) error {
	if d.IsUP() > 0 {
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
	log.Infof("dev=%s, state=%s", d.Name(), d.State())
	return nil
}

func Close(d Device) error {
	if d.IsUP() == 0 {
		err := fmt.Errorf("not opened, dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	if err := d.Close(); err != nil {
		err := fmt.Errorf("failurem dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	d.SetFlags(d.Flags() & ^NET_DEVICE_FLAG_UP)
	log.Infof("dev=%s, state=%s", d.Name(), d.State())
	return nil
}

func Output(d Device, ptype uint16, data []uint8, len int, dst *any) error {
	if d.IsUP() == 0 {
		err := fmt.Errorf("not opened, dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	if len > int(d.MTU()) {
		err := fmt.Errorf("too long, dev=%s, mtu=%x, len=%d", d.Name(), d.MTU(), len)
		log.Errorf(err.Error())
		return err
	}

	log.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), ptype, len)
	log.Debugdump(data, len)
	if err := d.Transmit(ptype, data, len, dst); err != nil {
		err := fmt.Errorf("device transmit failure, dev=%s, len=%d", d.Name(), len)
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func InputHandler(d Device, ptype uint16, data []uint8, len int) error {
	log.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), ptype, len)
	log.Debugdump(data, len)
	return nil
}

func Run() error {
	if err := intr.Run(); err != nil {
		err := fmt.Errorf("intr.Run() failure")
		log.Errorf(err.Error())
		return err
	}
	log.Debugf("open all devices...")
	for d := devices; d != nil; d = (*d).Next() {
		if err := Open(*d); err != nil {
			return err
		}
	}
	log.Debugf("running...")
	return nil
}

func Shutdown() error {
	log.Debugf("close all devices...")
	for d := devices; d != nil; d = (*d).Next() {
		if err := Close(*d); err != nil {
			return err
		}
	}
	intr.Shutdown()
	log.Debugf("shutting down")
	return nil
}

func Init() error {
	if err := intr.Init(); err != nil {
		err := fmt.Errorf("intr.Init() failure")
		log.Errorf(err.Error())
		return err
	}
	log.Infof("initialized")
	return nil
}
