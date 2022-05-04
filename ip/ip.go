package ip

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shimech/tcpip-stack/net/device"
	"github.com/shimech/tcpip-stack/net/protocol"
	"github.com/shimech/tcpip-stack/util/byte"
	"github.com/shimech/tcpip-stack/util/log"
)

type Address [4]uint8

type Header struct {
	vhl      uint8
	tos      uint8
	total    uint16
	id       uint16
	offset   uint16
	ttl      uint8
	protocol uint8
	sum      uint16
	src      Address
	dst      Address
	options  []uint8
}

const (
	IP_VERSION_IPV4 uint8 = 4

	IP_HDR_SIZE_MIN = 20
	IP_HDR_SIZE_MAX = 60
)

func (h *Header) v() uint8 {
	return (h.vhl & 0xf0) >> 4
}

func (h *Header) ihl() uint8 {
	return h.vhl & 0x0f
}

func Init() error {
	if err := protocol.Register(protocol.NET_PROTOCOL_TYPE_IP, input); err != nil {
		err := fmt.Errorf("protocol.Register() failure")
		return err
	}
	return nil
}

func newHeader(data []uint8) *Header {
	return &Header{
		vhl:      data[0],
		tos:      data[1],
		total:    uint16(data[2]*16 + data[3]),
		id:       uint16(data[4]*16 + data[5]),
		offset:   uint16(data[6]*16 + data[7]),
		ttl:      data[8],
		protocol: data[9],
		sum:      uint16(data[10]*16 + data[11]),
		src:      Address{data[12], data[13], data[14], data[15]},
		dst:      Address{data[16], data[17], data[18], data[19]},
		options:  data[20:]}
}

func addrPtoN(p string) (Address, error) {
	parts := strings.Split(p, ".")
	n := Address{}
	for i, part := range parts {
		u8, err := strconv.ParseInt(part, 10, 8)
		if err != nil {
			return Address{}, err
		}
		n[i] = uint8(u8)
	}
	return n, nil
}

func addrNtoP(n Address) string {
	return fmt.Sprintf("%d.%d.%d.%d", n[0], n[1], n[2], n[3])
}

func dump(data []uint8, len int) {
	h := newHeader(data)
	v := h.v()
	hl := h.ihl()
	hlen := hl << 2
	fmt.Fprintf(os.Stderr, "        vhl: 0x%02x [v: %d, hl: %d (%d)]\n", h.vhl, v, hl, hlen)
	fmt.Fprintf(os.Stderr, "        tos: 0x%02x\n", h.tos)
	total := byte.NtoH16(h.total)
	fmt.Fprintf(os.Stderr, "      total: %d (payload: %d)\n", total, total-uint16(hlen))
	fmt.Fprintf(os.Stderr, "         id: %d\n", byte.NtoH16(h.id))
	offset := byte.NtoH16(h.offset)
	fmt.Fprintf(os.Stderr, "     offset: 0x%04x [flags=%x, offset=%d]\n", offset, (offset&0xe000)>>13, offset&0x1fff)
	fmt.Fprintf(os.Stderr, "        ttl: %d\n", h.ttl)
	fmt.Fprintf(os.Stderr, "   protocol: %d\n", h.protocol)
	fmt.Fprintf(os.Stderr, "        sum: 0x%04x\n", byte.NtoH16(h.sum))
	fmt.Fprintf(os.Stderr, "        src: %s\n", addrNtoP(h.src))
	fmt.Fprintf(os.Stderr, "        dst: %s\n", addrNtoP(h.dst))
}

func input(data []uint8, len int, d device.Device) {
	if len < IP_HDR_SIZE_MIN {
		log.Errorf("too short")
		return
	}

	h := newHeader(data)

	if h.v() != IP_VERSION_IPV4 {
		log.Errorf("illegal version")
		return
	}

	if int(h.ihl()) > len {
		log.Errorf("ihl > len")
		return
	}

	total := byte.NtoH16(h.total)
	if int(total) > len {
		log.Errorf("tl > len")
		return
	}

	offset := byte.NtoH16(h.offset)
	if offset&0x2000 > 0 || offset&0x1fff > 0 {
		log.Errorf("fragments does not support")
		return
	}

	log.Debugf("dev=%s, protocol=%d, total=%d", d.Name(), h.protocol, total)
	dump(data, int(total))
}
