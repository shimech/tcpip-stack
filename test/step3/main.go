package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimech/tcpip-stack/driver/loopback"
	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/test"
	"github.com/shimech/tcpip-stack/util/log"
)

func init() {
	test.Init()
}

func main() {
	d := loopback.NewDevice()

	if err := net.Run(); err != nil {
		log.Errorf("net.Run() failure")
		return
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)

	go func() {
		for {
			if err := net.Output(d, 0x0800, test.TestData, nil); err != nil {
				log.Errorf("net.Output() failure")
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
