package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimech/tcpip-stack/driver"
	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/test"
	"github.com/shimech/tcpip-stack/util"
)

func main() {
	if err := net.Init(); err != nil {
		util.Errorf("net.Init() failure")
		return
	}

	dd := driver.NewDummyDevice()

	if err := net.Run(); err != nil {
		util.Errorf("net.Run() failure")
		return
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)

	go func() {
		for {
			if err := net.Output(dd, 0x0800, test.TestData(), len(test.TestData()), nil); err != nil {
				util.Errorf("net.Output() failure")
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	<-sig

	if err := net.Shutdown(); err != nil {
		return
	}
}
