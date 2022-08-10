package ip

import (
	"fmt"
)

type Address [IPV4_SIZE]byte

func (a *Address) string() string {
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}
