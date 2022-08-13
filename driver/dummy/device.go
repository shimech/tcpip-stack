package dummy

import (
	"os"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/platform/intr"
	"github.com/shimech/tcpip-stack/util/log"
)

type Device struct {
	index     int
	name      string
	dtype     net.DeviceType
	mtu       uint16
	flags     net.Flag
	hlen      uint16
	alen      uint16
	addr      net.Address
	peer      net.Address
	broadcast net.Address
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

func (d *Device) Type() net.DeviceType {
	return d.dtype
}

func (d *Device) SetType(t net.DeviceType) {
	d.dtype = t
}

func (d *Device) MTU() uint16 {
	return d.mtu
}

func (d *Device) SetMTU(mtu uint16) {
	d.mtu = mtu
}

func (d *Device) Flags() net.Flag {
	return d.flags
}

func (d *Device) SetFlags(f net.Flag) {
	d.flags = f
}

func (d *Device) Hlen() uint16 {
	return d.hlen
}

func (d *Device) SetHlen(hlen uint16) {
	d.hlen = hlen
}

func (d *Device) Alen() uint16 {
	return d.alen
}

func (d *Device) SetAlen(alen uint16) {
	d.alen = alen
}

func (d *Device) Addr() net.Address {
	return d.addr
}

func (d *Device) Peer() net.Address {
	return d.peer
}

func (d *Device) Broadcast() net.Address {
	return d.broadcast
}

func (d *Device) SetBroadcast(b net.Address) {
	d.broadcast = b
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
