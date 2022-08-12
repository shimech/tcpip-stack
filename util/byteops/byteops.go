package byteops

import (
	"encoding/binary"
	"unsafe"
)

var endian = byteorder()

func NtoH16(n uint16) uint16 {
	b := make([]byte, 2)
	endian.PutUint16(b, n)
	return endian.Uint16(b)
}

func HtoN16(h uint16) uint16 {
	b := make([]byte, 2)
	endian.PutUint16(b, h)
	return endian.Uint16(b)
}

func byteorder() binary.ByteOrder {
	x := 0x0001
	ptr := unsafe.Pointer(&x)
	if *(*byte)(ptr) == 0x00 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}
