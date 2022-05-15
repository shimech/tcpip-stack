package ip

import (
	"fmt"

	"github.com/shimech/tcpip-stack/util/log"
)

type Protocol struct {
	Type    ProtocolType
	Handler func(data []byte, src Address, dst Address, i *Iface)
}

type ProtocolType uint8

const (
	IP_PROTOCOL_ICMP ProtocolType = 1
	IP_PROTOCOL_TCP  ProtocolType = 6
	IP_PROTOCOL_UDP  ProtocolType = 17
)

var protocols []*Protocol

func RegisterProtocol(ptype ProtocolType, handler func(data []byte, src Address, dst Address, i *Iface)) error {
	for _, p := range protocols {
		if ptype == p.Type {
			return fmt.Errorf("type: %d is already registered", ptype)
		}
	}

	p := &Protocol{
		Type:    ptype,
		Handler: handler,
	}
	protocols = append([]*Protocol{p}, protocols...)
	log.Infof("registered, type=%d", p.Type)
	return nil
}
