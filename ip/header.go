package ip

import (
	"bytes"
	"encoding/binary"

	"github.com/shimech/tcpip-stack/util/log"
)

type Header struct {
	VHL            uint8
	TypeOfService  uint8
	TotalLength    uint16
	ID             uint16
	FragmentOffset uint16
	TTL            uint8
	Protocol       ProtocolType
	Checksum       uint16
	Src            Address
	Dst            Address
}

const (
	IP_HEADER_SIZE_MIN = 20
	IP_HEADER_SIZE_MAX = 60
)

func decodeHeader(data []byte) (*Header, error) {
	h := &Header{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, h); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return h, nil
}

func (h *Header) version() uint8 {
	return (h.VHL & 0xf0) >> 4
}

func (h *Header) ihl() uint8 {
	return h.VHL & 0x0f
}

func (h *Header) len() int {
	return int(h.ihl() << 2)
}

func (h *Header) encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, h); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
