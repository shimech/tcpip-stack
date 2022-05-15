package test

import (
	"github.com/shimech/tcpip-stack/icmp"
	"github.com/shimech/tcpip-stack/ip"
	"github.com/shimech/tcpip-stack/platform/linux/intr"
	"github.com/shimech/tcpip-stack/util/log"
)

func Init() {
	if err := ip.Init(); err != nil {
		log.Errorf("ip.Init() failure")
	}
	if err := intr.Init(); err != nil {
		log.Errorf("intr.Init() failure")
	}
	if err := icmp.Init(); err != nil {
		log.Errorf("icmp.Init() failure")
	}
	log.Infof("initialized")
}
