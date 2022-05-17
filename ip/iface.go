package ip

import (
	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/util/log"
)

type Iface struct {
	device    net.Device
	family    net.IfaceFamily
	Unicast   Address
	Newmask   Address
	Broadcast Address
}

var ifaces []*Iface

func (i *Iface) Device() net.Device {
	return i.device
}

func (i *Iface) SetDevice(d net.Device) {
	i.device = d
}

func (i *Iface) Family() net.IfaceFamily {
	return i.family
}

func (i *Iface) calcBroadcast() {
	n := networkAddress(i.Unicast, i.Newmask)
	for j := 0; j < len(n); j++ {
		n[j] |= ^i.Newmask[j]
	}
	i.Broadcast = n
}

func NewIface(unicast string, netmask string) (*Iface, error) {
	i := &Iface{
		family: net.NET_IFACE_FAMILY_IP,
	}

	u, err := ParseAddress(unicast)
	if err != nil {
		return i, err
	}
	i.Unicast = u

	n, err := ParseAddress(netmask)
	if err != nil {
		return i, err
	}
	i.Newmask = n

	i.calcBroadcast()

	return i, nil
}

func RegisterIface(d net.Device, i *Iface) error {
	err := net.AddIfaceToDevice(d, i)
	if err != nil {
		return err
	}
	ifaces = append([]*Iface{i}, ifaces...)
	log.Infof("registered: dev=%s, unicast=%s, netmask=%s, broadcast=%s", d.Name(), i.Unicast.String(), i.Newmask.String(), i.Broadcast.String())
	return nil
}

func SelectIface(a Address) *Iface {
	for _, i := range ifaces {
		if i.Unicast == a {
			return i
		}
	}
	return nil
}
