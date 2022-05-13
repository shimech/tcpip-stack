package ip

import (
	"fmt"
	"strconv"
	"strings"
)

type Address [IPV4_SIZE]byte

var (
	IP_ADDR_ANY       = Address{0x00, 0x00, 0x00, 0x00}
	IP_ADDR_BROADCAST = Address{0xff, 0xff, 0xff, 0xff}
)

func ParseAddress(s string) (Address, error) {
	var a Address
	for i, part := range strings.Split(s, ".") {
		p, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return a, err
		}
		a[i] = byte(p)
	}
	return a, nil
}

func (a *Address) string() string {
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}

func networkAddress(address Address, mask Address) Address {
	b := Address{}
	for i := 0; i < len(b); i++ {
		b[i] = address[i] & mask[i]
	}
	return b
}
