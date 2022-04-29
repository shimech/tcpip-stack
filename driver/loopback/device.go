package loopback

import (
	"fmt"
	"math"
	"os"
	"sync"
	"syscall"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/platform/linux/intr"
	"github.com/shimech/tcpip-stack/util/log"
	"github.com/shimech/tcpip-stack/util/queue"
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
	irq       syscall.Signal
	mu        sync.Mutex
	q         *queue.Queue
}

type QueueEntry struct {
	type_ uint16
	len   int
	data  []uint8
}

const (
	LOOPBACK_MTU         = math.MaxUint16
	LOOPBACK_QUEUE_LIMIT = 16
	LOOPBACK_IRQ         = intr.INTR_IRQ_BASE + 1
)

func NewDevice() *Device {
	d := &Device{
		type_: net.NET_DEVICE_TYPE_LOOPBACK,
		mtu:   LOOPBACK_MTU,
		hlen:  0,
		alen:  0,
		flags: net.NET_DEVICE_FLAG_LOOPBACK,
		irq:   LOOPBACK_IRQ,
		q:     queue.NewQueue(),
	}
	net.Register(d)
	intr.RequestIRQ(LOOPBACK_IRQ, LoopbackISR, intr.INTR_IRQ_SHARED, d.name, d)
	log.Debugf("initialized, dev=%s", d.name)
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
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.q.Size() >= LOOPBACK_QUEUE_LIMIT {
		err := fmt.Errorf("queue is full")
		log.Errorf(err.Error())
		return err
	}

	entry := &QueueEntry{
		type_: type_,
		len:   len,
		data:  data,
	}
	d.q.Push(entry)
	log.Debugf("queue pushed (size:%d), dev=%s, type=0x%04x, len=%d", d.q.Size(), d.name, type_, len)
	intr.RaiseIRQ(d.irq)
	return nil
}

func LoopbackISR(irq os.Signal, id any) error {
	d, ok := id.(*Device)
	if !ok {
		return fmt.Errorf("fail cast")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	for {
		e := d.q.Pop()
		if e == nil {
			break
		}
		entry, ok := e.(*QueueEntry)
		if !ok {
			return fmt.Errorf("fail cast")
		}

		log.Debugf("queue popped (num:%d), dev=%s, type=0x%04x, len=%d", d.q.Size(), d.name, entry.type_, entry.len)
		log.Debugdump(entry.data, entry.len)
		net.InputHandler(d, entry.type_, entry.data, entry.len)
	}
	return nil
}
