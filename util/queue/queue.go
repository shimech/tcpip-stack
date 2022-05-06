package queue

type Queue struct {
	head *QueueEntry
	tail *QueueEntry
	size uint
}

type QueueEntry struct {
	next *QueueEntry
	data any
}

func NewQueue() *Queue {
	return &Queue{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (q *Queue) Size() uint {
	return q.size
}

func (q *Queue) Push(data any) {
	entry := &QueueEntry{
		next: nil,
		data: data,
	}
	if q.tail != nil {
		q.tail.next = entry
	}
	q.tail = entry
	if q.head == nil {
		q.head = entry
	}
	q.size += 1
}

func (q *Queue) Pop() any {
	if q.head == nil {
		return nil
	}
	entry := q.head
	q.head = entry.next
	if q.head == nil {
		q.tail = nil
	}
	q.size -= 1
	d := entry.data
	return d
}

func (q *Queue) Peek() any {
	if q.head == nil {
		return nil
	}
	return q.head.data
}

func (q *Queue) ForEach(f func(data any, args ...any), args ...any) {
	if f == nil {
		return
	}
	for entry := q.head; entry != nil; entry = entry.next {
		f(entry.data, args...)
	}
}
