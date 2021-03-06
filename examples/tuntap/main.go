package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/shimech/tcpip-stack/pkg/raw/tuntap"
)

func main() {
	name := flag.String("name", "", "device name")
	flag.Parse()
	tap, err := tuntap.NewTAP(*name)
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s, addr=%x\n", tap.Name(), tap.Address())
	buf := make([]byte, 4096)
	for {
		n, err := tap.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("--- [%s] incomming %d bytes data ---\n", tap.Name(), n)
		fmt.Printf("%s", hex.Dump(buf[:n]))
	}
}
