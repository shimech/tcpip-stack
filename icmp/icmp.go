package icmp

import (
	"github.com/shimech/tcpip-stack/ip"
	"github.com/shimech/tcpip-stack/util/log"
)

func Init() error {
	if err := ip.RegisterProtocol(ip.IP_PROTOCOL_ICMP, input); err != nil {
		return err
	}
	return nil
}

func input(data []byte, src ip.Address, dst ip.Address, i *ip.Iface) {
	log.Debugf("%s => %s, len=%d", src.String(), dst.String(), len(data))
	log.Debugdump(data)
}
