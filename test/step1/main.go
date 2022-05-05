package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimech/tcpip-stack/driver/dummy"
	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/net/device"
	"github.com/shimech/tcpip-stack/test"
	"github.com/shimech/tcpip-stack/util/log"
)

func main() {
	if err := net.Init(); err != nil {
		log.Errorf("net.Init() failure")
		return
	}

	d := dummy.NewDevice()

	if err := net.Run(); err != nil {
		log.Errorf("device.Run() failure")
		return
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)

	go func() {
		for {
			if err := device.Output(d, 0x0800, test.TestData, len(test.TestData), nil); err != nil {
				log.Errorf("device.Output() failure")
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
