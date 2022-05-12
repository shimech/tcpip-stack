package intr

import (
	"os"
)

type IRQEntry struct {
	irq     os.Signal
	handler func(irq os.Signal, device any) error
	flags   int
	name    string
	device  any
}
