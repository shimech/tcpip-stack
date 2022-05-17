package icmp

import (
	"bytes"
	"encoding/binary"
)

type Datagram[H any] struct {
	Header *H
	Data   []byte
}

func (d *Datagram[H]) encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, d.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, d.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
