package ip

import (
	"fmt"
	"os"

	"github.com/shimech/tcpip-stack/util/byteops"
	"github.com/shimech/tcpip-stack/util/log"
)

func dump(data []byte) {
	h, err := decodeHeader(data)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	d := &Datagram{
		Header: *h,
		Data:   data[h.len():],
	}
	fmt.Fprintf(os.Stderr, "        vhl: 0x%02x [v: %d, hl: %d (%d)]\n", d.VHL, d.version(), d.ihl(), d.len())
	fmt.Fprintf(os.Stderr, "        tos: 0x%02x\n", d.TypeOfService)
	tl := byteops.NtoH16(d.TotalLength)
	fmt.Fprintf(os.Stderr, "      total: %d (payload: %d)\n", tl, tl-uint16(d.len()))
	fmt.Fprintf(os.Stderr, "         id: %d\n", byteops.NtoH16(d.ID))
	fo := byteops.NtoH16(d.FragmentOffset)
	fmt.Fprintf(os.Stderr, "     offset: 0x%04x [flags=%x, offset=%d]\n", fo, (fo&0xe000)>>13, fo&0x1fff)
	fmt.Fprintf(os.Stderr, "        ttl: %d\n", d.TTL)
	fmt.Fprintf(os.Stderr, "   protocol: %d\n", d.Protocol)
	fmt.Fprintf(os.Stderr, "        sum: 0x%04x\n", byteops.NtoH16(d.Checksum))
	fmt.Fprintf(os.Stderr, "        src: %s\n", d.Src.String())
	fmt.Fprintf(os.Stderr, "        dst: %s\n", d.Dst.String())
	log.Debugdump(d.Data)
}
