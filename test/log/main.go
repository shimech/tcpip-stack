package main

import (
	"github.com/shimech/tcpip-stack/test"
	"github.com/shimech/tcpip-stack/util/log"
)

func main() {
	log.Debugf("Hello, World!")
	log.Debugdump(test.TestData)
}
