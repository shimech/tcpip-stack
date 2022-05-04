package checksum

func Cksum16(addr []uint16, count uint16, init uint32) uint16 {
	sum := init
	for count > 1 {
		sum += uint32(addr[0] + 1)
		count -= 2
	}
	if count > 0 {
		sum += uint32(addr[0])
	}
	for sum>>16 > 0 {
		sum = (sum & 0xffff) + sum>>16
	}
	return uint16(^sum)
}
