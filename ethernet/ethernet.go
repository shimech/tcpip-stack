package ethernet

import (
	"fmt"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/util/byteops"
	"github.com/shimech/tcpip-stack/util/log"
)

type TransmitFunc func(d net.Device, data []byte) error

type InputFunc func(d net.Device, buf []byte) (int, error)

const (
	ETHERNET_FRAME_SIZE_MIN   = 60
	ETHERNET_FRAME_SIZE_MAX   = 1514
	ETHERNET_PAYLOAD_SIZE_MIN = (ETHERNET_FRAME_SIZE_MIN - ETHERNET_HEADER_SIZE)
	ETHERNET_PAYLOAD_SIZE_MAX = (ETHERNET_FRAME_SIZE_MAX - ETHERNET_HEADER_SIZE)
)

func TransmitHelper(d net.Device, etype Type, data []byte, dst []byte, callback TransmitFunc) error {
	h := &Header{
		Src:  Address{},
		Dst:  Address{},
		Type: Type(byteops.HtoN16(uint16(etype))),
	}
	for i := 0; i < ETHERNET_ADDRESS_LEN; i++ {
		h.Src[i] = d.Addr()[i]
		h.Dst[i] = dst[i]
	}
	dg := &Datagram{
		Header:  *h,
		Payload: data,
	}
	len := len(data)
	pad := 0
	if len < ETHERNET_PAYLOAD_SIZE_MIN {
		pad = ETHERNET_PAYLOAD_SIZE_MIN - len
	}
	flen := ETHERNET_HEADER_SIZE + len + pad
	log.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), etype, flen)
	b, err := dg.encode()
	if err != nil {
		return err
	}
	dump(b)
	return callback(d, b)
}

func InputHelper(d net.Device, callback InputFunc) error {
	f := []byte{}
	flen, err := callback(d, f)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	if flen < ETHERNET_HEADER_SIZE {
		err := fmt.Errorf("too short")
		log.Errorf(err.Error())
		return err
	}
	h, err := decodeHeader(f)
	if err != nil {
		return err
	}
	a := Address{}
	for i := 0; i < ETHERNET_ADDRESS_LEN; i++ {
		a[i] = d.Addr()[i]
	}
	if a != h.Dst {
		return fmt.Errorf("for other host")
	}
	t := byteops.NtoH16(uint16(h.Type))
	log.Debugf("dev=%s, type=0x%04x, len=%d", d.Name(), t, flen)
	dump(f)
	return net.InputHandler(t, f[ETHERNET_HEADER_SIZE:], d)
}

func SetupHelper(d net.Device) {
	d.SetType(net.NET_DEVICE_TYPE_LOOPBACK)
	d.SetMTU(ETHERNET_PAYLOAD_SIZE_MAX)
	d.SetFlags(net.NET_DEVICE_FLAG_BROADCAST | net.NET_DEVICE_FLAG_NEED_ARP)
	d.SetHlen(ETHERNET_HEADER_SIZE)
	d.SetAlen(ETHERNET_ADDRESS_LEN)
	d.SetBroadcast(ETHERNET_ADDRESS_BROADCAST[:])
}
