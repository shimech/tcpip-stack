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

	src, err := ip.ParseAddress(test.LOOPBACK_IP_ADDR)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	dst := src

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)

	payload := test.TestData[ip.IP_HDR_SIZE_MIN:]
	go func() {
		for {
			if err := ip.Output(ip.IP_PROTOCOL_ICMP, payload, src, dst); err != nil {
				log.Errorf("ip.Output() failure")
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
