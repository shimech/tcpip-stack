package ethernet

import (
	"bytes"
	"encoding/binary"

	"github.com/shimech/tcpip-stack/pkg/net"
)

type Header struct {
	Dst  Address
	Src  Address
	Type net.EthernetType
}

type Frame struct {
	Header
	payload []byte
}

func parse(data []byte) (*Frame, error) {
	frame := Frame{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &frame.Header); err != nil {
		return nil, err
	}
	frame.payload = buf.Bytes()
	return &frame, nil
}
