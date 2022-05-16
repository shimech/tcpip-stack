package checksum

import (
	"encoding/binary"
)

func Cksum16(data []byte, init uint32) uint16 {
	len := len(data)
	sum := init
	var count int
	for count = len; count > 1; count -= 2 {
		d := data[len-count : len-count+2]
		sum += uint32(binary.BigEndian.Uint16(d))
	}
	if count > 0 {
		sum += uint32(data[len-count])
	}
	for sum>>16 > 0 {
		sum = (sum & 0xffff) + sum>>16
	}
	return uint16(^sum)
}
