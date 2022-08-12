package ip

import (
	"fmt"
	"os"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/util/byteops"
	"github.com/shimech/tcpip-stack/util/checksum"
	"github.com/shimech/tcpip-stack/util/log"
)

const (
	IPV4_SIZE             = 4
	IP_VERSION_IPV4 uint8 = 4

	IP_HDR_SIZE_MIN = 20
	IP_HDR_SIZE_MAX = 60
)

func Init() error {
	if err := net.RegisterProtocol(net.NET_PROTOCOL_TYPE_IP, input); err != nil {
		err := fmt.Errorf("net.RegisterProtocol() failure")
		return err
	}
	return nil
}

func input(data []byte, d net.Device) {
	len := len(data)
	if len < IP_HDR_SIZE_MIN {
		log.Errorf("too short")
		return
	}

	h, err := newHeader(data)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	if h.version() != IP_VERSION_IPV4 {
		log.Errorf("illegal version")
		return
	}

	if int(h.ihl()) > len {
		log.Errorf("ihl > size")
		return
	}

	tl := byteops.NtoH16(h.TotalLength)
	if int(tl) > len {
		log.Errorf("total length > size")
		return
	}

	if checksum.Cksum16(data, len, 0) != 0 {
		log.Errorf("checksum error")
		return
	}

	fo := byteops.NtoH16(h.FragmentOffset)
	if fo&0x2000 > 0 || fo&0x1fff > 0 {
		log.Errorf("fragments does not support")
		return
	}

	i, ok := net.GetIfaceInDevice(d, net.NET_IFACE_FAMILY_IP).(*Iface)
	if !ok || i == nil {
		log.Errorf("interface is not found")
		return
	}

	if h.Dst != i.unicast && h.Dst != i.broadcast && h.Dst != IP_ADDR_BROADCAST {
		log.Errorf("fort other host")
		return
	}

	log.Debugf("dev=%s, iface=%s, protocol=%d, total=%d", d.Name(), i.unicast.string(), h.Protocol, tl)
	dump(data)
}

func Output(protocol byte, data []byte, src Address, dst Address) error {
	len := len(data)
	if src == IP_ADDR_ANY {
		return fmt.Errorf("ip routing does not implement")
	}
	i := SelectIface(src)
	if i == nil {
		return fmt.Errorf("interface is not found")
	}
	n := networkAddress(i.unicast, i.netmask)
	if n != networkAddress(dst, i.netmask) && n != IP_ADDR_BROADCAST {
		return fmt.Errorf("illegal destination address")
	}
	if int(i.device.MTU()) < IP_HDR_SIZE_MIN+len {
		err := fmt.Errorf("too long, dev=%s, mtu=%d < %d", i.device.Name(), i.device.MTU(), IP_HDR_SIZE_MIN+len)
		log.Errorf(err.Error())
		return err
	}
	id := generateID()
	if err := outputCore(i, protocol, data, src, dst, id, 0); err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func outputCore(i *Iface, protocol uint8, data []byte, src Address, dst Address, id uint16, offset uint16) error {
	len := len(data)
	hlen := IP_HDR_SIZE_MIN
	tl := hlen + len
	d := &Datagram{
		Header: Header{
			VHL:            IP_VERSION_IPV4<<4 | uint8(hlen)>>2,
			TypeOfService:  0,
			TotalLength:    byteops.HtoN16(uint16(tl)),
			ID:             byteops.HtoN16(id),
			FragmentOffset: byteops.HtoN16(offset),
			TTL:            255,
			Protocol:       protocol,
			Checksum:       0,
			Src:            src,
			Dst:            dst,
		},
		Data: data,
	}
	hb, err := d.Header.encode()
	if err != nil {
		return err
	}
	d.Checksum = checksum.Cksum16(hb, hlen, 0)
	log.Debugf("dev=%s, dst=%s, protocol=%d, len=%d", i.device.Name(), dst.string(), protocol, tl)
	db, err := d.encode()
	if err != nil {
		return err
	}
	dump(db)
	return outputDevice(i, db, dst)
}

func outputDevice(i *Iface, data []byte, dst Address) error {
	var hwaddr uint8
	if (i.device.Flags() & net.NET_DEVICE_FLAG_NEED_ARP) > 0 {
		if dst == i.broadcast || dst == IP_ADDR_BROADCAST {
			hwaddr = i.device.Broadcast()
		} else {
			err := fmt.Errorf("arp does not implement")
			return err
		}
	}
	return net.Output(i.device, net.NET_PROTOCOL_TYPE_IP, data, hwaddr)
}

func dump(data []byte) {
	h, err := newHeader(data)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	v := h.version()
	ihl := h.ihl()
	hlen := ihl << 2
	fmt.Fprintf(os.Stderr, "        vhl: 0x%02x [v: %d, hl: %d (%d)]\n", h.VHL, v, ihl, hlen)
	fmt.Fprintf(os.Stderr, "        tos: 0x%02x\n", h.TypeOfService)
	tl := byteops.NtoH16(h.TotalLength)
	fmt.Fprintf(os.Stderr, "      total: %d (payload: %d)\n", tl, tl-uint16(hlen))
	fmt.Fprintf(os.Stderr, "         id: %d\n", byteops.NtoH16(h.ID))
	fo := byteops.NtoH16(h.FragmentOffset)
	fmt.Fprintf(os.Stderr, "     offset: 0x%04x [flags=%x, offset=%d]\n", fo, (fo&0xe000)>>13, fo&0x1fff)
	fmt.Fprintf(os.Stderr, "        ttl: %d\n", h.TTL)
	fmt.Fprintf(os.Stderr, "   protocol: %d\n", h.Protocol)
	fmt.Fprintf(os.Stderr, "        sum: 0x%04x\n", byteops.NtoH16(h.Checksum))
	fmt.Fprintf(os.Stderr, "        src: %s\n", h.Src.string())
	fmt.Fprintf(os.Stderr, "        dst: %s\n", h.Dst.string())
}
