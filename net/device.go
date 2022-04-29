package net

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
	IsUP() uint16
	State() string
	Open() error
	Close() error
	Transmit(type_ uint16, data []uint8, len int, dst *any) error
}

const NET_DEVICE_FLAG_UP uint16 = 0x0001

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
