package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/util/log"
)

type Iface interface {
	Device() Device
	SetDevice(d Device)
	Family() IfaceFamily
}

type IfaceFamily int

const (
	NET_IFACE_FAMILY_IP   IfaceFamily = 1
	NET_IFACE_FAMILY_IPV6 IfaceFamily = 2
)

func AddIfaceToDevice(d Device, i Iface) error {
	for _, e := range d.Ifaces() {
		if e.Family() == i.Family() {
			err := fmt.Errorf("already exists, dev=%s, family=%d", d.Name(), e.Family())
			log.Errorf(err.Error())
			return err
		}
	}
	i.SetDevice(d)
	d.PrependIface(i)
	return nil
}

func GetIfaceInDevice(d Device, f IfaceFamily) Iface {
	for _, e := range d.Ifaces() {
		if e.Family() == f {
			return e
		}
	}
	return nil
}
