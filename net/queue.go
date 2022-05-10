package net

type QueueEntry struct {
	Device Device
	Data   []byte
}

func NewQueueEntry(d Device, data []byte) *QueueEntry {
	return &QueueEntry{
		Device: d,
		Data:   data,
	}
}
