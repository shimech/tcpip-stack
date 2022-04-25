package main

import (
	"github.com/shimech/tcpip-stack/test"
	"github.com/shimech/tcpip-stack/util"
)

func main() {
	util.Debugf("Hello, World!")
	util.Debugdump(test.TestData(), len(test.TestData()))
}
