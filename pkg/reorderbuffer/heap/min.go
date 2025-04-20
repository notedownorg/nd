package heap

import (
	"container/heap"
	"sync"
	"time"
)

type Event[T any] interface {
	Payload() T
	Clock() uint64
	Added() time.Time
}

type min[T any] struct {
	mutex sync.RWMutex
	h     *unsafeMinHeap[T]
}

func NewMin[T any]() *min[T] {
	h := &unsafeMinHeap[T]{}
	heap.Init(h)
	return &min[T]{h: h}
}

func (h *min[T]) Len() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.h.Len()
}

func (m *min[T]) Push(x any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	heap.Push(m.h, x)
}

func (m *min[T]) Pop() any {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return heap.Pop(m.h)
}

func (m *min[T]) Peek() any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.h.Len() == 0 {
		return nil
	}
	return m.h.peek()
}

func (h *min[T]) PopIf(conditional func(Event[T]) bool) (Event[T], bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.h.Len() == 0 {
		return nil, false
	}
	if conditional(h.h.peek()) {
		ev := heap.Pop(h.h).(Event[T])
		return ev, true
	}
	return nil, false
}
