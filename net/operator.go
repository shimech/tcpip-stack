package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/util"
)

var devices *Device
var index = 0

func Register(d Device) {
	d.SetIndex(index)
	d.SetName(fmt.Sprintf("net%d", d.Index()))
	d.SetNext(devices)
	devices = &d
	util.Infof("registered, dev=%s, type=0x%04x", d.Name(), d.Type())

	index += 1
}

func Open(d Device) error {
	if d.IsUP() > 0 {
		err := fmt.Errorf("already opened, dev=%s", d.Name())
		util.Errorf(err.Error())
		return err
	}

	if err := d.Open(); err != nil {
		err := fmt.Errorf("failure dev=%s", d.Name())
		util.Errorf(err.Error())
		return err
	}

	d.SetFlags(d.Flags() | NET_DEVICE_FLAG_UP)
	util.Infof("dev=%s, state=%s", d.Name(), d.State())
	return nil
}

func Close(d Device) error {
	if d.IsUP() == 0 {
		err := fmt.Errorf("not opened, dev=%s", d.Name())
		util.Errorf(err.Error())
		return err
	}

	if err := d.Close(); err != nil {
		err := fmt.Errorf("failurem dev=%s", d.Name())
		util.Errorf(err.Error())
		return err
	}

	d.SetFlags(d.Flags() & ^NET_DEVICE_FLAG_UP)
	util.Infof("dev=%s, state=%s", d.Name(), d.State())
	return nil
}

func Output(d Device, type_ uint16, data []uint8, len int, dst *any) error {
	if d.IsUP() == 0 {
		err := fmt.Errorf("not opened, dev=%s", d.Name())
		util.Errorf(err.Error())
		return err
	}

	if len > int(d.MTU()) {
		err := fmt.Errorf("too long, dev=%s, mtu=%x, len=%d", d.Name(), d.MTU(), len)
		util.Errorf(err.Error())
		return err
	}

	util.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), type_, len)
	util.Debugdump(data, len)
	if err := d.Transmit(type_, data, len, dst); err != nil {
		err := fmt.Errorf("device transmit failure, dev=%s, len=%d", d.Name(), len)
		util.Errorf(err.Error())
		return err
	}
	return nil
}

func InputHandler(d Device, type_ uint16, data []uint8, len int) error {
	util.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), type_, len)
	util.Debugdump(data, len)
	return nil
}

func Run() error {
	util.Debugf("open all devices...")
	for d := devices; d != nil; d = (*d).Next() {
		if err := Open(*d); err != nil {
			return err
		}
	}
	util.Debugf("running...")
	return nil
}

func Shutdown() error {
	util.Debugf("close all devices...")
	for d := devices; d != nil; d = (*d).Next() {
		if err := Close(*d); err != nil {
			return err
		}
	}
	util.Debugf("shutting down")
	return nil
}

func Init() error {
	util.Infof("initialized")
	return nil
}
