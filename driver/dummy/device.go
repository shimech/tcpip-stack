package dummy

import (
	"os"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/platform/linux/intr"
	"github.com/shimech/tcpip-stack/util"
)

const (
	NET_DEVICE_TYPE_DUMMY = 0x0000
	DUMMY_MTU             = 0xffff
	DUMMY_IRQ             = intr.INTR_IRQ_BASE
)

type Device struct {
	next      *net.Device
	index     int
	name      string
	type_     uint16
	mtu       uint16
	flags     uint16
	hlen      uint16
	alen      uint16
	adr       uint8
	peer      uint8
	broadcast uint8
	priv      *any
}

func NewDevice() *Device {
	d := &Device{
		type_: NET_DEVICE_TYPE_DUMMY,
		mtu:   DUMMY_MTU,
		hlen:  0,
		alen:  0,
	}
	net.Register(d)
	intr.RequestIRQ(DUMMY_IRQ, DummyISR, intr.INTR_IRQ_SHARED, d.Name(), d)
	util.Debugf("initialized, dev=%s", d.name)
	return d
}

func (d *Device) Next() *net.Device {
	return d.next
}

func (d *Device) SetNext(_d *net.Device) {
	d.next = _d
}

func (d *Device) Index() int {
	return d.index
}

func (d *Device) SetIndex(i int) {
	d.index = i
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) SetName(n string) {
	d.name = n
}

func (d *Device) Type() uint16 {
	return d.type_
}

func (d *Device) MTU() uint16 {
	return d.mtu
}

func (d *Device) Flags() uint16 {
	return d.flags
}

func (d *Device) SetFlags(f uint16) {
	d.flags = f
}

func (d *Device) Hlen() uint16 {
	return d.hlen
}

func (d *Device) Alen() uint16 {
	return d.alen
}

func (d *Device) Addr() uint8 {
	return d.adr
}

func (d *Device) Peer() uint8 {
	return d.peer
}

func (d *Device) Broadcast() uint8 {
	return d.broadcast
}

func (d *Device) Priv() *any {
	return d.priv
}

func (d *Device) IsUP() uint16 {
	return net.IsUP(d)
}

func (d *Device) State() string {
	return net.State(d)
}

func (d *Device) Open() error {
	return nil
}

func (d *Device) Close() error {
	return nil
}

func (d *Device) Transmit(type_ uint16, data []uint8, len int, dst *any) error {
	util.Debugf("dev=%s, type=0x%04x, len=%d", d.name, type_, len)
	util.Debugdump(data, len)
	// drop data
	intr.RaiseIRQ(DUMMY_IRQ)
	return nil
}

func DummyISR(irq os.Signal, id any) error {
	util.Debugf("irq=%d, dev=%s", irq, id.(net.Device).Name())
	return nil
}
