package workspace

import (
	. "github.com/notedownorg/nd/pkg/workspace/node"
)

type Event struct {
	Op   Operation
	Id   string
	Node Node
}

type Operation uint32

const (
	// Signal that this node was present when the client was created or when the subscriber subscribed
	Load Operation = iota

	// Signal that this node has been updated or created
	Change

	// Signal that this node has been deleted
	Delete

	// Signal that the subscriber has received all existing nodes present at the time of subscription
	SubscriberLoadComplete
)

type subscriber struct {
	ch   chan Event
	kind Kind
}

func (w *Workspace) Subscribe(ch chan Event, kind Kind, loadInitialNodes bool) int {
	// Add to subscribe map first to ensure we don't miss any events
	// Multiple reporting of an event is better than missing one
	index := len(w.subscribers)
	w.subscribers[index] = &subscriber{ch: ch, kind: kind}

	if loadInitialNodes {
		go func(s chan Event) {
			// Recover from a panic if the subscriber has been closed
			// Likely this will only happen in tests but its theoretically possible in regular usage
			defer func() {
				if recover() != nil {
					w.log.Warn("subscriber closed before initial nodes could be loaded")
				}
			}()

			switch kind {
			case DocumentKind:
				for key, node := range w.documents.Entries() {
					s <- Event{Op: Load, Id: key, Node: node}
				}
			}
			s <- Event{Op: SubscriberLoadComplete}
		}(ch)
	}

	return index
}

func (w *Workspace) Unsubscribe(index int) {
	sub, ok := w.subscribers[index]
	if ok {
		delete(w.subscribers, index)
		close(sub.ch)
	}
}

func (w *Workspace) eventDispatcher() {
	for event := range w.events {
		for _, subscriber := range w.subscribers {
			if subscriber.kind != KindFromID(event.Id) {
				continue
			}
			go func(s chan Event, e Event) {
				// Recover from a panic if the subscriber has been closed
				// Likely this will only happen in tests but its theoretically possible in regular usage
				defer func() {
					if recover() != nil {
						w.log.Warn("dropping event as subscriber has been closed", "path", e.Id)
					}
				}()
				s <- e
			}(subscriber.ch, event)
		}
	}
}
