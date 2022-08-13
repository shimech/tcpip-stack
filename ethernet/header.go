package ethernet

import (
	"bytes"
	"encoding/binary"

	"github.com/shimech/tcpip-stack/util/log"
)

type Header struct {
	Src  Address
	Dst  Address
	Type Type
}

const ETHERNET_HEADER_SIZE = 14

func decodeHeader(data []byte) (*Header, error) {
	h := &Header{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, h); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return h, nil
}
