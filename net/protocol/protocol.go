package protocol

import (
	"fmt"

	"github.com/shimech/tcpip-stack/net/device"
	"github.com/shimech/tcpip-stack/util/log"
	"github.com/shimech/tcpip-stack/util/queue"
)

type Protocol struct {
	Next    *Protocol
	Type    uint16
	Queue   *queue.Queue
	handler func(data []uint8, len int, d device.Device)
}

type QueueEntry struct {
	d    *device.Device
	len  int
	data []uint8
}

const (
	NET_PROTOCOL_TYPE_IP   = 0x0800
	NET_PROTOCOL_TYPE_ARP  = 0x0806
	NET_PROTOCOL_TYPE_IPV6 = 0x86dd
)

var Protocols *Protocol

func NewQueueEntry(d *device.Device, len int, data []uint8) *QueueEntry {
	return &QueueEntry{
		d:    d,
		len:  len,
		data: data,
	}
}

func Register(ptype uint16, handler func(data []uint8, len int, d device.Device)) error {
	for p := Protocols; p != nil; p = p.Next {
		if ptype == p.Type {
			err := fmt.Errorf("already registered, type=0x%04x", ptype)
			log.Errorf(err.Error())
			return err
		}
	}
	p := &Protocol{
		Next:    Protocols,
		Type:    ptype,
		Queue:   queue.NewQueue(),
		handler: handler,
	}
	Protocols = p
	log.Infof("registered, type=0x%04x", ptype)
	return nil
}
