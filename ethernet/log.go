package ethernet

import (
	"fmt"
	"os"

	"github.com/shimech/tcpip-stack/util/byteops"
	"github.com/shimech/tcpip-stack/util/log"
)

func dump(frame []byte) {
	h, err := decodeHeader(frame)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	fmt.Fprintf(os.Stderr, "        src: %s\n", h.Src.String())
	fmt.Fprintf(os.Stderr, "        dst: %s\n", h.Dst.String())
	fmt.Fprintf(os.Stderr, "       type: 0x%04x\n", byteops.NtoH16(uint16(h.Type)))
}
