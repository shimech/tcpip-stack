package icmp

import (
	"bytes"
	"encoding/binary"

	"github.com/shimech/tcpip-stack/util/log"
)

type Header struct {
	Type     Type
	Code     uint8
	Checksum uint16
	Value    uint32
}

const ICMP_HEADER_SIZE = 8

func decodeHeader(data []byte) (*Header, error) {
	h := &Header{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, h); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return h, nil
}
