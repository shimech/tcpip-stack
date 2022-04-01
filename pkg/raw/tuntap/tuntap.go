package tuntap

import (
	"io"
)

type TAP struct {
	io.ReadWriteCloser
	name string
}

func NewTAP(name string) (*TAP, error) {
	n, f, err := openTAP(name)
	if err != nil {
		return nil, err
	}
	return &TAP{
		ReadWriteCloser: f,
		name:            n,
	}, nil
}

func (t TAP) Address() []byte {
	addr, _ := getAddress(t.name)
	return addr[:6]
}

func (t TAP) Name() string {
	return t.name
}
