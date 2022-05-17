package main

import (
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimech/tcpip-stack/driver/loopback"
	"github.com/shimech/tcpip-stack/icmp"
	"github.com/shimech/tcpip-stack/ip"
	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/test"
	"github.com/shimech/tcpip-stack/util/byteops"
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

	payload := test.TestData[ip.IP_HEADER_SIZE_MIN+icmp.ICMP_HEADER_SIZE:]
	id := os.Getpid() % math.MaxUint16
	seq := 0
	go func() {
		for {
			seq += 1
			if err := icmp.Output(icmp.ICMP_TYPE_ECHO, 0, byteops.HtoN32(uint32(id<<16|seq)), payload, src, dst); err != nil {
				log.Errorf("icmp.Output() failure")
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
