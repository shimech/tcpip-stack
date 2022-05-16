package icmp

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
	fmt.Fprintf(os.Stderr, "       type: %d (%s)\n", h.Type, h.Type.String())
	fmt.Fprintf(os.Stderr, "       code: %d\n", h.Code)
	fmt.Fprintf(os.Stderr, "        sum: 0x%04x\n", byteops.NtoH16(h.Checksum))
	switch h.Type {
	case ICMP_TYPE_ECHOREPLY:
	case ICMP_TYPE_ECHO:
		e, err := decodeEcho(data)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		fmt.Fprintf(os.Stderr, "         id: %d\n", byteops.NtoH16(e.ID))
		fmt.Fprintf(os.Stderr, "        seq: %d\n", byteops.NtoH16(e.SequenceNumber))
		break
	default:
		fmt.Fprintf(os.Stderr, "     values: 0x%08x\n", byteops.NtoH32(h.Value))
	}
	log.Debugdump(data)
}
