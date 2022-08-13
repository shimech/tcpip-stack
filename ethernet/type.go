package ethernet

type Type uint16

const (
	ETHERNET_TYPE_IP   Type = 0x0800
	ETHERNET_TYPE_ARP  Type = 0x0806
	ETHERNET_TYPE_IPV6 Type = 0x86dd
)
