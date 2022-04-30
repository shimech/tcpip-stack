package net

type Device interface {
	Next() *Device
	SetNext(d *Device)
	Index() int
	SetIndex(i int)
	Name() string
	SetName(n string)
	Type() DeviceType
	MTU() uint16
	Flags() uint16
	SetFlags(f uint16)
	Hlen() uint16
	Alen() uint16
	Addr() uint8
	Peer() uint8
	Broadcast() uint8
	IsUP() uint16
	State() string
	Open() error
	Close() error
	Transmit(ptype uint16, data []uint8, len int, dst *any) error
}

type DeviceType uint16

const (
	NET_DEVICE_FLAG_UP uint16 = 0x0001

	NET_DEVICE_TYPE_DUMMY    = 0x0000
	NET_DEVICE_TYPE_LOOPBACK = 0x0001
)

func IsUP(d Device) uint16 {
	return d.Flags() & NET_DEVICE_FLAG_UP
}

func State(d Device) string {
	if d.IsUP() > 0 {
		return "up"
	} else {
		return "down"
	}
}
