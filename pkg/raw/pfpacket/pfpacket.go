package pfpacket

import "syscall"

type PFPacket struct {
	fd   int
	name string
}

func NewPFPacket(name string) (*PFPacket, error) {
	fd, err := openPFPacket(name)
	if err != nil {
		return nil, err
	}
	return &PFPacket{
		fd:   fd,
		name: name,
	}, nil
}

func (p PFPacket) Name() string {
	return p.name
}

func (p PFPacket) Address() []byte {
	addr, _ := getAddress(p.name)
	return addr[:6]
}

func (p PFPacket) Read(b []byte) (int, error) {
	return syscall.Read(p.fd, b)
}

func (p PFPacket) Write(b []byte) (int, error) {
	return syscall.Write(p.fd, b)
}

func (p PFPacket) Close() error {
	return syscall.Close(p.fd)
}
