package icmp

import (
	"bytes"
	"encoding/binary"

	"github.com/shimech/tcpip-stack/util/log"
)

type Echo struct {
	Type           Type
	Code           uint8
	Checksum       uint16
	ID             uint16
	SequenceNumber uint16
}

func decodeEcho(data []byte) (*Echo, error) {
	e := &Echo{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, e); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return e, nil
}
