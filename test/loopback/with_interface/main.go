package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimech/tcpip-stack/driver/loopback"
	"github.com/shimech/tcpip-stack/ip"
	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/test"
	"github.com/shimech/tcpip-stack/util/log"
)

func init() {
	test.Init()
}

func main() {
	d := loopback.NewDevice()

	i, err := ip.NewIface(test.LOOPBACK_IP_ADDR, test.LOOPBACK_NETMASK)
	if err != nil {
		log.Errorf("ip.NewIface() failure")
		return
	}
	if err := ip.RegisterIface(d, i); err != nil {
		log.Errorf("ip.RegisterIface() failue")
		return
	}

	if err := net.Run(); err != nil {
		log.Errorf("net.Run() failure")
		return
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)

	go func() {
		for {
			if err := net.Output(d, net.NET_PROTOCOL_TYPE_IP, test.TestData, nil); err != nil {
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
