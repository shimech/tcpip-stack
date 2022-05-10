package ip

import (
	"fmt"
	"os"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/util/byteops"
	"github.com/shimech/tcpip-stack/util/checksum"
	"github.com/shimech/tcpip-stack/util/log"
)

const (
	IPV4_SIZE             = 4
	IP_VERSION_IPV4 uint8 = 4

	IP_HDR_SIZE_MIN = 20
	IP_HDR_SIZE_MAX = 60
)

func Init() error {
	if err := net.RegisterProtocol(net.NET_PROTOCOL_TYPE_IP, input); err != nil {
		err := fmt.Errorf("protocol.Register() failure")
		return err
	}
	return nil
}

func input(data []byte, d net.Device) {
	size := len(data)
	if size < IP_HDR_SIZE_MIN {
		log.Errorf("too short")
		return
	}

	h, err := newHeader(data)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	if h.version() != IP_VERSION_IPV4 {
		log.Errorf("illegal version")
		return
	}

	if int(h.ihl()) > size {
		log.Errorf("ihl > size")
		return
	}

	tl := byteops.NtoH16(h.TotalLength)
	if int(tl) > size {
		log.Errorf("total length > size")
		return
	}

	if checksum.Cksum16(data, 0) != 0 {
		log.Errorf("checksum error")
		return
	}

	fo := byteops.NtoH16(h.FragmentOffset)
	if fo&0x2000 > 0 || fo&0x1fff > 0 {
		log.Errorf("fragments does not support")
		return
	}

	log.Debugf("dev=%s, protocol=%d, total=%d", d.Name(), h.Protocol, tl)
	dump(data, int(tl))
}

func dump(data []uint8, len int) {
	h, err := newHeader(data)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	v := h.version()
	ihl := h.ihl()
	hlen := ihl << 2
	fmt.Fprintf(os.Stderr, "        vhl: 0x%02x [v: %d, hl: %d (%d)]\n", h.VHL, v, ihl, hlen)
	fmt.Fprintf(os.Stderr, "        tos: 0x%02x\n", h.TypeOfService)
	tl := byteops.NtoH16(h.TotalLength)
	fmt.Fprintf(os.Stderr, "      total: %d (payload: %d)\n", tl, tl-uint16(hlen))
	fmt.Fprintf(os.Stderr, "         id: %d\n", byteops.NtoH16(h.ID))
	fo := byteops.NtoH16(h.FragmentOffset)
	fmt.Fprintf(os.Stderr, "     offset: 0x%04x [flags=%x, offset=%d]\n", fo, (fo&0xe000)>>13, fo&0x1fff)
	fmt.Fprintf(os.Stderr, "        ttl: %d\n", h.TTL)
	fmt.Fprintf(os.Stderr, "   protocol: %d\n", h.Protocol)
	fmt.Fprintf(os.Stderr, "        sum: 0x%04x\n", byteops.NtoH16(h.CheckSum))
	fmt.Fprintf(os.Stderr, "        src: %s\n", h.Src.string())
	fmt.Fprintf(os.Stderr, "        dst: %s\n", h.Dst.string())
}
