package reorderbuffer

import (
	"sync/atomic"
	"time"

	"github.com/notedownorg/nd/pkg/reorderbuffer/heap"
)

// A reorder buffer that allows for a timeout on operations.
// It receives out of order events based on a logical clock and emits them in order
func NewWithExpiry[T any](timeout time.Duration, input chan Event[T]) chan T {
	output := make(chan T)
	min := heap.NewMin[T]()
	var closed atomic.Bool

	// One goroutine to handle the input events
	go func() {
		for event := range input {
			min.Push(event)
		}
		closed.Store(true) // ensure the other goroutine doesn't leak
	}()

	// Another reading from the heap to process the events
	go func() {
		clock := uint64(0)
		for !closed.Load() {
			if min.Len() == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			ev, popped := min.PopIf(func(ev heap.Event[T]) bool {
				return ev.Clock() <= clock+1 || ev.Added().Add(timeout).Before(time.Now())
			})
			if popped {
				clock = ev.Clock()
				output <- ev.Payload()
				continue
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return output
}

var _ heap.Event[any] = &Event[any]{}

type Event[T any] struct {
	payload T
	clock   uint64
	added   time.Time
}

func NewEvent[T any](payload T, clock uint64) Event[T] {
	return Event[T]{payload: payload, clock: clock, added: time.Now()}
}

func (e Event[T]) Payload() T {
	return e.payload
}

func (e Event[T]) Clock() uint64 {
	return e.clock
}

func (e Event[T]) Added() time.Time {
	return e.added
}
