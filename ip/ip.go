package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shimech/tcpip-stack/net/device"
	"github.com/shimech/tcpip-stack/net/protocol"
	"github.com/shimech/tcpip-stack/util/byteops"
	"github.com/shimech/tcpip-stack/util/log"
)

type Address [IPV4_SIZE]byte

type Header struct {
	VHL            uint8
	TypeOfService  uint8
	TotalLength    uint16
	ID             uint16
	FragmentOffset uint16
	TTL            uint8
	Protocol       uint8
	CheckSum       uint16
	Src            Address
	Dst            Address
}

const (
	IPV4_SIZE             = 4
	IP_VERSION_IPV4 uint8 = 4

	IP_HDR_SIZE_MIN = 20
	IP_HDR_SIZE_MAX = 60
)

func (h *Header) Version() uint8 {
	return (h.VHL & 0xf0) >> 4
}

func (h *Header) IHL() uint8 {
	return h.VHL & 0x0f
}

func Init() error {
	if err := protocol.Register(protocol.NET_PROTOCOL_TYPE_IP, input); err != nil {
		err := fmt.Errorf("protocol.Register() failure")
		return err
	}
	return nil
}

func newHeader(data []byte) (*Header, error) {
	h := &Header{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, h); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return h, nil
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
	h, err := newHeader(data)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	v := h.Version()
	hl := h.IHL()
	hlen := hl << 2
	fmt.Fprintf(os.Stderr, "        vhl: 0x%02x [v: %d, hl: %d (%d)]\n", h.VHL, v, hl, hlen)
	fmt.Fprintf(os.Stderr, "        tos: 0x%02x\n", h.TypeOfService)
	total := byteops.NtoH16(h.TotalLength)
	fmt.Fprintf(os.Stderr, "      total: %d (payload: %d)\n", total, total-uint16(hlen))
	fmt.Fprintf(os.Stderr, "         id: %d\n", byteops.NtoH16(h.ID))
	offset := byteops.NtoH16(h.FragmentOffset)
	fmt.Fprintf(os.Stderr, "     offset: 0x%04x [flags=%x, offset=%d]\n", offset, (offset&0xe000)>>13, offset&0x1fff)
	fmt.Fprintf(os.Stderr, "        ttl: %d\n", h.TTL)
	fmt.Fprintf(os.Stderr, "   protocol: %d\n", h.Protocol)
	fmt.Fprintf(os.Stderr, "        sum: 0x%04x\n", byteops.NtoH16(h.CheckSum))
	fmt.Fprintf(os.Stderr, "        src: %s\n", addrNtoP(h.Src))
	fmt.Fprintf(os.Stderr, "        dst: %s\n", addrNtoP(h.Dst))
}

func input(data []byte, d device.Device) {
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

	if h.Version() != IP_VERSION_IPV4 {
		log.Errorf("illegal version")
		return
	}

	if int(h.IHL()) > size {
		log.Errorf("ihl > len")
		return
	}

	total := byteops.NtoH16(h.TotalLength)
	if int(total) > size {
		log.Errorf("tl > len")
		return
	}

	offset := byteops.NtoH16(h.FragmentOffset)
	if offset&0x2000 > 0 || offset&0x1fff > 0 {
		log.Errorf("fragments does not support")
		return
	}

	log.Debugf("dev=%s, protocol=%d, total=%d", d.Name(), h.Protocol, total)
	dump(data, int(total))
}
