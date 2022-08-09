package intr

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shimech/tcpip-stack/util/log"
)

type Handler struct {
	SoftIRQ func() error
}

const (
	INTR_IRQ_SHARED  = 0x0001
	INTR_IRQ_BASE    = syscall.Signal(0x22) + 1
	INTR_IRQ_SOFTIRQ = syscall.SIGUSR1
)

var irqs *IRQEntry
var sigs = make(chan os.Signal, 1)
var raise = make(chan struct{}, 1)
var terminate = make(chan struct{}, 1)

func RequestIRQ(
	irq os.Signal,
	handler func(irq os.Signal, dev any) error,
	flags int,
	name string,
	dev any) error {
	log.Debugf("irq=%d, flags=%d, name=%s", irq, flags, name)
	for entry := irqs; entry != nil; entry = entry.next {
		if entry.irq == irq {
			if entry.flags^INTR_IRQ_SHARED != 0 || flags^INTR_IRQ_SHARED != 0 {
				err := fmt.Errorf("conflicts with already registered IRQs")
				log.Errorf(err.Error())
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
	log.Debugf("registered: irq=%d, name=%s", irq, name)
	return nil
}

func RaiseIRQ(irq os.Signal) {
	sigs <- irq
}

func Thread(h *Handler) {
	term := false

	log.Debugf("start...")
	raise <- struct{}{}
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGHUP:
			term = true
			break
		case syscall.SIGUSR1:
			h.SoftIRQ()
			break
		default:
			for entry := irqs; entry != nil; entry = entry.next {
				if entry.irq == sig {
					log.Debugf("irq=%d, name=%s", entry.irq, entry.name)
					entry.handler(entry.irq, entry.dev)
				}
			}
			break
		}
		if term {
			break
		}
	}

	log.Debugf("terminated")
	terminate <- struct{}{}
}

func Run(h *Handler) error {
	go Thread(h)
	<-raise
	return nil
}

func Shutdown() {
	sigs <- syscall.SIGHUP
	<-terminate
}

func Init() error {
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGUSR1)
	return nil
}
