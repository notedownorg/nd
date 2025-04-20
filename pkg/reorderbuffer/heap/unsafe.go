package heap

import (
	"container/heap"
)

var _ heap.Interface = &unsafeMinHeap[any]{}

type unsafeMinHeap[T any] []Event[T]

// heap.Interface
func (h unsafeMinHeap[T]) Len() int           { return len(h) }
func (h unsafeMinHeap[T]) Less(i, j int) bool { return h[i].Clock() < h[j].Clock() }
func (h unsafeMinHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *unsafeMinHeap[T]) Push(x any) {
	ev := x.(Event[T])
	*h = append(*h, ev)
}

func (h *unsafeMinHeap[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (h unsafeMinHeap[T]) peek() Event[T] {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}
