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

var (
	irqs      []*IRQEntry
	sigs      = make(chan os.Signal, 1)
	raise     = make(chan struct{}, 1)
	terminate = make(chan struct{}, 1)
)

func RequestIRQ(
	irq os.Signal,
	handler func(irq os.Signal, device any) error,
	flags int,
	name string,
	device any) error {
	log.Debugf("irq=%d, flags=%d, name=%s", irq, flags, name)
	for _, e := range irqs {
		if e.irq == irq {
			if e.flags^INTR_IRQ_SHARED != 0 || flags^INTR_IRQ_SHARED != 0 {
				err := fmt.Errorf("conflicts with already registered IRQs")
				log.Errorf(err.Error())
				return err
			}
		}
	}

	e := &IRQEntry{
		irq:     irq,
		handler: handler,
		flags:   flags,
		name:    name,
		device:  device,
	}
	irqs = append([]*IRQEntry{e}, irqs...)
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
	for !term {
		sig := <-sigs
		switch sig {
		case syscall.SIGHUP:
			term = true
		case syscall.SIGUSR1:
			h.SoftIRQ()
		default:
			for _, e := range irqs {
				if e.irq == sig {
					log.Debugf("irq=%d, name=%s", e.irq, e.name)
					e.handler(e.irq, e.device)
				}
			}
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
