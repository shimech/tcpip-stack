package byteops

const (
	BIG_ENDIAN    = 4321
	LITTLE_ENDIAN = 1234
)

var endian int

func NtoH16(n uint16) uint16 {
	if endian == 0 {
		endian = byteorder()
	}

	if endian == LITTLE_ENDIAN {
		n = byteswap16(n)
	}
	return n
}

func byteorder() int {
	x := [4]uint8{0x00, 0x00, 0x00, 0x01}
	if x[0] == 0 {
		return BIG_ENDIAN
	} else {
		return LITTLE_ENDIAN
	}
}

func byteswap16(v uint16) uint16 {
	return (v&0x00ff)<<8 | (v&0xff00)>>8
}
