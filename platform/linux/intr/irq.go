package intr

import (
	"os"
)

type IRQEntry struct {
	next    *IRQEntry
	irq     os.Signal
	handler func(irq os.Signal, dev any) error
	flags   int
	name    string
	dev     any
}
