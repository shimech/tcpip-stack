package ip

import (
	"fmt"

	"github.com/shimech/tcpip-stack/net/device"
	"github.com/shimech/tcpip-stack/net/protocol"
	"github.com/shimech/tcpip-stack/util/log"
)

func Init() error {
	if err := protocol.Register(protocol.NET_PROTOCOL_TYPE_IP, input); err != nil {
		err := fmt.Errorf("protocol.Register() failure")
		return err
	}
	return nil
}

func input(data []uint8, len int, d device.Device) {
	log.Debugf("dev=%s, len=%d", d.Name(), len)
	log.Debugdump(data, len)
}
