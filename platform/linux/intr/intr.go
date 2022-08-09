package intr

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shimech/tcpip-stack/util"
)

const (
	INTR_IRQ_SHARED = 0x0001
	INTR_IRQ_BASE   = syscall.Signal(0x22) + 1
)

var irqs *IRQEntry
var sigs = make(chan os.Signal)
var terminate = make(chan bool)

func RequestIRQ(
	irq os.Signal,
	handler func(irq os.Signal, dev any) error,
	flags int,
	name string,
	dev any) error {
	util.Debugf("irq=%d, flags=%d, name=%s", irq, flags, name)
	for entry := irqs; entry != nil; entry = entry.next {
		if entry.irq == irq {
			if entry.flags^INTR_IRQ_SHARED != 0 || flags^INTR_IRQ_SHARED != 0 {
				err := fmt.Errorf("conflicts with already registered IRQs")
				util.Errorf(err.Error())
				return err
			}
		}
	}

	entry := &IRQEntry{
		next:    irqs,
		irq:     irq,
		handler: handler,
		flags:   flags,
		name:    name,
		dev:     dev,
	}
	irqs = entry
	signal.Notify(sigs, irq)
	util.Debugf("registered: irq=%d, name=%s", irq, name)
	return nil
}

func RaiseIRQ(irq os.Signal) {
	sigs <- irq
}

func Thread() {
	term := false

	util.Debugf("start...")
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGHUP:
			term = true
			break
		default:
			for entry := irqs; entry != nil; entry = entry.next {
				if entry.irq == sig {
					util.Debugf("irq=%d, name=%s", entry.irq, entry.name)
					entry.handler(entry.irq, entry.dev)
				}
			}
			break
		}
		if term {
			break
		}
	}

	util.Debugf("terminated")
	terminate <- term
}

func Run() error {
	go Thread()
	return nil
}

func Shutdown() {
	sigs <- syscall.SIGHUP
	<-terminate
}

func Init() error {
	signal.Notify(sigs, syscall.SIGHUP)
	return nil
}
