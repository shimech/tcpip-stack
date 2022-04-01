package arp

import (
	"bytes"
	"sync"
	"time"

	"github.com/shimech/tcpip-stack/pkg/net"
)

type ARPEntry struct {
	protocolAddress []byte
	hardwareAddress []byte
	iface           net.ProtocolInterface
	timestamp       time.Time
}

type ARPTable struct {
	storage []*ARPEntry
	mutex   sync.RWMutex
}

var repo *ARPTable

func newARPTable() *ARPTable {
	return &ARPTable{
		storage: make([]*ARPEntry, 0, 1024),
	}
}

func (tbl *ARPTable) lookupUnlocked(protocolAddress []byte) *ARPEntry {
	for _, entry := range tbl.storage {
		if bytes.Compare(entry.protocolAddress, protocolAddress) == 0 {
			return entry
		}
	}
	return nil
}

func (tbl *ARPTable) lookup(protocolAddress []byte) *ARPEntry {
	tbl.mutex.RLock()
	defer tbl.mutex.RUnlock()
	return tbl.lookupUnlocked(protocolAddress)
}

func (tbl *ARPTable) update(protocolAddress []byte, hardwareAddress []byte) bool {
	tbl.mutex.Lock()
	defer tbl.mutex.Unlock()
	entry := tbl.lookupUnlocked(protocolAddress)
	if entry == nil {
		return false
	}
	entry.hardwareAddress = hardwareAddress
	entry.timestamp = time.Now()
	return true
}

func (tbl *ARPTable) insert(iface net.ProtocolInterface, protocolAddress []byte, hardwareAddress []byte) bool {
	tbl.mutex.Lock()
	defer tbl.mutex.Unlock()
	if tbl.lookupUnlocked(protocolAddress) != nil {
		return false
	}
	entry := &ARPEntry{
		protocolAddress: protocolAddress,
		hardwareAddress: hardwareAddress,
		iface:           iface,
		timestamp:       time.Now(),
	}
	tbl.storage = append(tbl.storage, entry)
	return true
}

func (tbl *ARPTable) length() int {
	tbl.mutex.RLock()
	defer tbl.mutex.RUnlock()
	return len(tbl.storage)
}
