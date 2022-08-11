package dummy

import (
	"os"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/platform/linux/intr"
	"github.com/shimech/tcpip-stack/util/log"
)

type Device struct {
	index     int
	name      string
	dtype     uint16
	mtu       uint16
	flags     uint16
	hlen      uint16
	alen      uint16
	addr      uint8
	peer      uint8
	broadcast uint8
	ifaces    []net.Iface
}

const (
	DUMMY_MTU = 0xffff
	DUMMY_IRQ = intr.INTR_IRQ_BASE
)

func NewDevice() *Device {
	d := &Device{
		dtype: net.NET_DEVICE_TYPE_DUMMY,
		mtu:   DUMMY_MTU,
		hlen:  0,
		alen:  0,
	}
	net.RegisterDevice(d)
	intr.RequestIRQ(DUMMY_IRQ, dummyISR, intr.INTR_IRQ_SHARED, d.name, d)
	log.Debugf("initialized, dev=%s", d.name)
	return d
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
	return d.dtype
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
	return d.addr
}

func (d *Device) Peer() uint8 {
	return d.peer
}

func (d *Device) Broadcast() uint8 {
	return d.broadcast
}

func (d *Device) Ifaces() []net.Iface {
	return d.ifaces
}

func (d *Device) PrependIface(i net.Iface) {
	d.ifaces = append([]net.Iface{i}, d.ifaces...)
}

func (d *Device) Open() error {
	return nil
}

func (d *Device) Close() error {
	return nil
}

func (d *Device) Transmit(dtype uint16, data []byte, dst any) error {
	log.Debugf("dev=%s, type=0x%04x, len=%d", d.name, dtype, len(data))
	log.Debugdump(data)
	// drop data
	intr.RaiseIRQ(DUMMY_IRQ)
	return nil
}

func dummyISR(irq os.Signal, id any) error {
	log.Debugf("irq=%d, dev=%s", irq, id.(net.Device).Name())
	return nil
}
