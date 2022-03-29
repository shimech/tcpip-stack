package net

import (
	"fmt"
	"log"
)

type ProtocolRxHandler func(dev *Device, data []byte, src, dst HardwareAddress) error

type Packet struct {
	dev  *Device
	data []byte
	src  HardwareAddress
	dst  HardwareAddress
}

type Entry struct {
	Type      EthernetType
	rxHandler ProtocolRxHandler
	rxQueue   chan *Packet
}

var protocols = map[EthernetType]*Entry{}

func RegisterProtocol(Type EthernetType, rxHandler ProtocolRxHandler) error {
	if protocols[Type] != nil {
		return fmt.Errorf("protocol `%d` is registered", Type)
	}
	entry := &Entry{
		Type:      Type,
		rxHandler: rxHandler,
		rxQueue:   make(chan *Packet),
	}
	protocols[Type] = entry
	go func() {
		for {
			select {
			case packet, _ := <-entry.rxQueue:
				entry.rxHandler(packet.dev, packet.data, packet.src, packet.dst)
			}
		}
	}()
	log.Printf("protocol registered: %x\n", entry.Type)
	return nil
}
