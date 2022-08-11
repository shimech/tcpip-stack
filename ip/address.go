package ip

import (
	"fmt"
	"strconv"
	"strings"
)

type Address [IPV4_SIZE]byte

var IP_ADDR_BROADCAST = Address{0xff, 0xff, 0xff, 0xff}

func newAddress(s string) (Address, error) {
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
