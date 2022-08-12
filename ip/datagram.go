package ip

import (
	"bytes"
	"encoding/binary"
)

type Datagram struct {
	Header
	Data []byte
}

func (d *Datagram) encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, d.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, d.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
