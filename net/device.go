package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/util/log"
)

type Device interface {
	Index() int
	SetIndex(i int)
	Name() string
	SetName(n string)
	Type() DeviceType
	SetType(t DeviceType)
	MTU() uint16
	SetMTU(mtu uint16)
	Flags() Flag
	SetFlags(f Flag)
	Hlen() uint16
	SetHlen(hlen uint16)
	Alen() uint16
	SetAlen(alen uint16)
	Addr() Address
	Peer() Address
	Broadcast() Address
	SetBroadcast(b Address)
	Ifaces() []Iface
	PrependIface(i Iface)
	Open() error
	Close() error
	Transmit(dtype uint16, data []byte, dst any) error
}

type Flag uint16

type DeviceType uint16

const (
	NET_DEVICE_FLAG_UP        Flag = 0x0001
	NET_DEVICE_FLAG_LOOPBACK  Flag = 0x0010
	NET_DEVICE_FLAG_BROADCAST Flag = 0x0020
	NET_DEVICE_FLAG_P2P       Flag = 0x0040
	NET_DEVICE_FLAG_NEED_ARP  Flag = 0x0100

	NET_DEVICE_TYPE_DUMMY    DeviceType = 0x0000
	NET_DEVICE_TYPE_LOOPBACK DeviceType = 0x0001
	NET_DEVICE_TYPE_ETHERNET DeviceType = 0x0002
)

var (
	devices []Device
	index   = 0
)

func RegisterDevice(d Device) {
	d.SetIndex(index)
	d.SetName(fmt.Sprintf("net%d", d.Index()))
	devices = append([]Device{d}, devices...)
	log.Infof("registered, dev=%s, type=0x%04x", d.Name(), d.Type())

	index += 1
}

func openDevice(d Device) error {
	if isDeviceUP(d) {
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
	if !isDeviceUP(d) {
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

func isDeviceUP(d Device) bool {
	return d.Flags()&NET_DEVICE_FLAG_UP > 0
}

func deviceState(d Device) string {
	if isDeviceUP(d) {
		return "up"
	} else {
		return "down"
	}
}
