package driver

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

type DummyDevice struct {
	next      *net.Device
	index     int
	name      string
	type_     uint16
	mtu       uint16
	flags     uint16
	hlen      uint16
	alen      uint16
	addr      uint8
	peer      uint8
	broadcast uint8
	priv      *any
}

var _ net.Device = &DummyDevice{}

func NewDummyDevice() *DummyDevice {
	dd := &DummyDevice{
		type_: NET_DEVICE_TYPE_DUMMY,
		mtu:   DUMMY_MTU,
		hlen:  0,
		alen:  0,
	}
	net.Register(dd)
	intr.RequestIRQ(DUMMY_IRQ, DummyISR, intr.INTR_IRQ_SHARED, dd.Name(), dd)
	util.Debugf("initialized, dev=%s", dd.name)
	return dd
}

func (dd *DummyDevice) Next() *net.Device {
	return dd.next
}

func (dd *DummyDevice) SetNext(d *net.Device) {
	dd.next = d
}

func (dd *DummyDevice) Index() int {
	return dd.index
}

func (dd *DummyDevice) SetIndex(i int) {
	dd.index = i
}

func (dd *DummyDevice) Name() string {
	return dd.name
}

func (dd *DummyDevice) SetName(n string) {
	dd.name = n
}

func (dd *DummyDevice) Type() uint16 {
	return dd.type_
}

func (dd *DummyDevice) MTU() uint16 {
	return dd.mtu
}

func (dd *DummyDevice) Flags() uint16 {
	return dd.flags
}

func (dd *DummyDevice) SetFlags(f uint16) {
	dd.flags = f
}

func (dd *DummyDevice) Hlen() uint16 {
	return dd.hlen
}

func (dd *DummyDevice) Alen() uint16 {
	return dd.alen
}

func (dd *DummyDevice) Addr() uint8 {
	return dd.addr
}

func (dd *DummyDevice) Peer() uint8 {
	return dd.peer
}

func (dd *DummyDevice) Broadcast() uint8 {
	return dd.broadcast
}

func (dd *DummyDevice) Priv() *any {
	return dd.priv
}

func (dd *DummyDevice) IsUP() uint16 {
	return net.IsUP(dd)
}

func (dd *DummyDevice) State() string {
	return net.State(dd)
}

func (dd *DummyDevice) Open() error {
	return nil
}

func (dd *DummyDevice) Close() error {
	return nil
}

func (dd *DummyDevice) Transmit(type_ uint16, data []uint8, len int, dst *any) error {
	util.Debugf("dev=%s, type=0x%04x, len=%d", dd.name, type_, len)
	util.Debugdump(data, len)
	// drop data
	intr.RaiseIRQ(DUMMY_IRQ)
	return nil
}

func DummyISR(irq os.Signal, id any) error {
	util.Debugf("irq=%d, dev=%s", irq, id.(net.Device).Name())
	return nil
}
