package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/platform/linux/intr"
	"github.com/shimech/tcpip-stack/util/log"
)

func InputHandler(ptype uint16, data []byte, d Device) error {
	for _, p := range protocols {
		if uint16(p.Type) == ptype {
			e := NewQueueEntry(d, data)
			p.Queue.Push(e)
			log.Debugf("queue pushed (num:%d), dev=%s, type=0x%04x, len=%d", p.Queue.Size(), d.Name(), ptype, len(data))
			log.Debugdump(data, len(data))
			intr.RaiseIRQ(intr.INTR_IRQ_SOFTIRQ)
			return nil
		}
	}
	return nil
}

func SoftIRQHandler() error {
	for _, p := range protocols {
		for {
			x := p.Queue.Pop()
			if x == nil {
				break
			}
			e, ok := x.(*QueueEntry)
			if !ok {
				return fmt.Errorf("fail cast")
			}
			log.Debugf("queue popped (num:%d), dev=%s, type=0x%04x, len=%d", p.Queue.Size(), e.Device.Name(), p.Type, len(e.Data))
			log.Debugdump(e.Data, len(e.Data))
			p.Handler(e.Data, e.Device)
		}
	}
	return nil
}

func Run() error {
	h := &intr.Handler{
		SoftIRQ: SoftIRQHandler,
	}
	if err := intr.Run(h); err != nil {
		err := fmt.Errorf("intr.Run() failure")
		log.Errorf(err.Error())
		return err
	}
	log.Debugf("open all devices...")

	for _, d := range devices {
		if err := openDevice(d); err != nil {
			return err
		}
	}
	log.Debugf("running...")
	return nil
}

func Output(d Device, dtype uint16, data []byte, len int, dst any) error {
	if !isDeviceUP(d) {
		err := fmt.Errorf("not opened, dev=%s", d.Name())
		log.Errorf(err.Error())
		return err
	}

	if len > int(d.MTU()) {
		err := fmt.Errorf("too long, dev=%s, mtu=%x, len=%d", d.Name(), d.MTU(), len)
		log.Errorf(err.Error())
		return err
	}

	log.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), dtype, len)
	log.Debugdump(data, len)
	if err := d.Transmit(dtype, data, dst); err != nil {
		err := fmt.Errorf("device transmit failure, dev=%s, len=%d", d.Name(), len)
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func Shutdown() error {
	log.Debugf("close all devices...")

	for _, d := range devices {
		if err := closeDevice(d); err != nil {
			return err
		}
	}
	intr.Shutdown()
	log.Debugf("shutting down")
	return nil
}
