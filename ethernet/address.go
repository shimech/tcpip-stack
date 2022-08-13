package ethernet

import (
	"fmt"
	"strconv"
	"strings"
)

type Address [ETHERNET_ADDRESS_LEN]byte

const ETHERNET_ADDRESS_LEN = 6

var (
	ETHERNET_ADDRESS_ANY       = Address{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	ETHERNET_ADDRESS_BROADCAST = Address{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

func ParseAddress(s string) (Address, error) {
	var a Address
	for i, part := range strings.Split(s, ":") {
		p, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return a, err
		}
		a[i] = byte(p)
	}
	return a, nil
}

func (a *Address) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", a[0], a[1], a[2], a[3], a[4], a[5])
}
