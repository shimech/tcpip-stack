package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/ip"
	"github.com/shimech/tcpip-stack/net/device"
	"github.com/shimech/tcpip-stack/net/protocol"
	"github.com/shimech/tcpip-stack/platform/linux/intr"
	"github.com/shimech/tcpip-stack/util/log"
)

func Init() error {
	if err := ip.Init(); err != nil {
		err := fmt.Errorf("ip.Init() failure")
		log.Errorf(err.Error())
		return err
	}
	if err := intr.Init(); err != nil {
		err := fmt.Errorf("intr.Init() failure")
		log.Errorf(err.Error())
		return err
	}
	log.Infof("initialized")
	return nil
}

func InputHandler(ptype uint16, data []uint8, len int, d device.Device) error {
	for p := protocol.Head(); p != nil; p = p.Next {
		if uint16(p.Type) == ptype {
			pqe := protocol.NewQueueEntry(d, len, data)
			p.Queue.Push(pqe)
			log.Debugf("queue pushed (num:%d), dev=%s, type=0x%04x, len=%d", p.Queue.Size(), d.Name(), ptype, len)
			log.Debugdump(data, len)
			intr.RaiseIRQ(intr.INTR_IRQ_SOFTIRQ)
			return nil
		}
	}
	return nil
}

func SoftIRQHandler() error {
	for p := protocol.Head(); p != nil; p = p.Next {
		for {
			e := p.Queue.Pop()
			if e == nil {
				break
			}
			entry, ok := e.(*protocol.QueueEntry)
			if !ok {
				return fmt.Errorf("fail cast")
			}
			log.Debugf("queue popped (num:%d), dev=%s, type=0x%04x, len=%d", p.Queue.Size(), entry.Device.Name(), p.Type, entry.Len)
			log.Debugdump(entry.Data, entry.Len)
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
	for d := device.Head(); d != nil; d = (*d).Next() {
		if err := device.Open(*d); err != nil {
			return err
		}
	}
	log.Debugf("running...")
	return nil
}

func Shutdown() error {
	log.Debugf("close all devices...")
	for d := device.Head(); d != nil; d = (*d).Next() {
		if err := device.Close(*d); err != nil {
			return err
		}
	}
	intr.Shutdown()
	log.Debugf("shutting down")
	return nil
}
