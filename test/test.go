package test

const (
	LOOPBACK_IP_ADDR = "127.0.0.1"
	LOOPBACK_NETMASK = "255.0.0.0"

	ETHER_TAP_NAME    = "tap0"
	ETHER_TAP_HW_ADDR = "00:00:5e:00:53:01"
	ETHER_TAP_IP_ADDR = "192.0.2.2"
	ETHER_TAP_NETMASK = "255.255.255.0"

	DEFAULT_GATEWAY = "192.0.2.1"
)

var TestData = []byte{
	0x45, 0x00, 0x00, 0x30,
	0x00, 0x80, 0x00, 0x00,
	0xff, 0x01, 0xbd, 0x4a,
	0x7f, 0x00, 0x00, 0x01,
	0x7f, 0x00, 0x00, 0x01,
	0x08, 0x00, 0x35, 0x64,
	0x00, 0x80, 0x00, 0x01,
	0x31, 0x32, 0x33, 0x34,
	0x35, 0x36, 0x37, 0x38,
	0x39, 0x30, 0x21, 0x40,
	0x23, 0x24, 0x25, 0x5e,
	0x26, 0x2a, 0x28, 0x29,
}
