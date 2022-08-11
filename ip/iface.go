package ip

import (
	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/util/log"
)

type Iface struct {
	device    net.Device
	family    net.IfaceFamily
	unicast   Address
	netmask   Address
	broadcast Address
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
	b := Address{}
	for j := 0; j < len(b); j++ {
		n := i.unicast[j] & i.netmask[j]
		b[j] = n | ^i.netmask[j]
	}
	i.broadcast = b
}

func NewIface(unicast string, netmask string) (*Iface, error) {
	i := &Iface{
		family: net.NET_IFACE_FAMILY_IP,
	}

	u, err := newAddress(unicast)
	if err != nil {
		return i, err
	}
	i.unicast = u

	n, err := newAddress(netmask)
	if err != nil {
		return i, err
	}
	i.netmask = n

	i.calcBroadcast()

	return i, nil
}

func RegisterIface(d net.Device, i *Iface) error {
	err := net.AddIfaceToDevice(d, i)
	if err != nil {
		return err
	}
	ifaces = append([]*Iface{i}, ifaces...)
	log.Infof("registered: dev=%s, unicast=%s, netmask=%s, broadcast=%s", d.Name(), i.unicast.string(), i.netmask.string(), i.broadcast.string())
	return nil
}

func SelectIface(a Address) *Iface {
	for _, i := range ifaces {
		if i.unicast == a {
			return i
		}
	}
	return nil
}
