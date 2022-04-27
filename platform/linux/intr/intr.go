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
var sigmask = make(chan os.Signal)

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
	signal.Notify(sigmask, irq)
	util.Debugf("registered: irq=%d, name=%s", irq, name)
	return nil
}

func RaiseIRQ(irq os.Signal) {
	sigmask <- irq
}

func Thread() {
	terminate := false

	util.Debugf("start...")
	for {
		sig := <-sigmask
		switch sig {
		case syscall.SIGHUP:
			terminate = true
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
		if terminate {
			break
		}
	}

	util.Debugf("terminated")
}

func Run() error {
	go Thread()
	return nil
}

func Shutdown() {
	sigmask <- syscall.SIGHUP
}

func Init() error {
	signal.Notify(sigmask, syscall.SIGHUP)
	return nil
}
