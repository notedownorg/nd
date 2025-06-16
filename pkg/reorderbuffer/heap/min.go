// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
