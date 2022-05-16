package icmp

import (
	"github.com/shimech/tcpip-stack/ip"
	"github.com/shimech/tcpip-stack/util/byteops"
	"github.com/shimech/tcpip-stack/util/checksum"
	"github.com/shimech/tcpip-stack/util/log"
)

func Init() error {
	if err := ip.RegisterProtocol(ip.IP_PROTOCOL_ICMP, input); err != nil {
		return err
	}
	return nil
}

func input(data []byte, src ip.Address, dst ip.Address, i *ip.Iface) {
	if len(data) < ICMP_HEADER_SIZE {
		log.Errorf("too short")
		return
	}
	h, err := decodeHeader(data)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	if checksum.Cksum16(data, 0) != 0 {
		log.Errorf("checksum error, sum=0x%04x, verify=0x%04x", byteops.NtoH16(h.Checksum), byteops.NtoH16(checksum.Cksum16(data, -uint32(h.Checksum))))
		return
	}
	log.Debugf("%s => %s, len=%d", src.String(), dst.String(), len(data))
	dump(data)
}
