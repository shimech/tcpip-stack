package net

import (
	"fmt"

	"github.com/shimech/tcpip-stack/util/log"
	"github.com/shimech/tcpip-stack/util/queue"
)

type Protocol struct {
	Next    *Protocol
	Type    uint16
	Queue   *queue.Queue
	Handler func(data []byte, d Device)
}

const (
	NET_PROTOCOL_TYPE_IP   = 0x0800
	NET_PROTOCOL_TYPE_ARP  = 0x0806
	NET_PROTOCOL_TYPE_IPV6 = 0x86dd
)

var protocols *Protocol

func RegisterProtocol(ptype uint16, handler func(data []byte, d Device)) error {
	for p := protocols; p != nil; p = p.Next {
		if ptype == p.Type {
			err := fmt.Errorf("already registered, type=0x%04x", ptype)
			log.Errorf(err.Error())
			return err
		}
	}
	p := &Protocol{
		Type:    ptype,
		Queue:   queue.NewQueue(),
		Handler: handler,
	}
	pushNewProtocol(p)
	log.Infof("registered, type=0x%04x", ptype)
	return nil
}

func pushNewProtocol(p *Protocol) {
	p.Next = protocols
	protocols = p
}
