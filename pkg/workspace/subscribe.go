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
	// internal channels are buffered channels that allow us to drop events when a client is unable to keep up
	// they are used to prevent blocking to ensure per-client ordering WITHOUT blocking the main event dispatcher
	internal      chan Event
	external      chan Event
	kind          Kind
	loadCompleted bool
}

func (w *Workspace) Subscribe(ch chan Event, kind Kind, loadInitialNodes bool) int {
	// Add to subscribe map first to ensure we don't miss any events
	// Multiple reporting of an event is better than missing one
	sub := &subscriber{external: ch, kind: kind, loadCompleted: false, internal: make(chan Event, 1000)}
	w.subscribersMu.Lock()
	index := len(w.subscribers)
	w.subscribers[index] = sub
	w.subscribersMu.Unlock()

	if loadInitialNodes {
		go func(s *subscriber) {
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
					s.external <- Event{Op: Load, Id: key, Node: node}
				}
			}
			s.external <- Event{Op: SubscriberLoadComplete}
			s.loadCompleted = true
		}(sub)
	} else {
		sub.loadCompleted = true
	}

	// Handle subscribers in a goroutine to prevent one subscriber from blocking the entire workspace
	// but ensure that each subscriber receives the events in the correct order.
	go func(s *subscriber) {
		// Recover from a panic if the subscriber has been closed
		// Likely this will only happen in tests but its theoretically possible in regular usage
		defer func() {
			if recover() != nil {
				w.log.Warn("subscriber goroutine panicked")
			}
		}()

		loadBuffer := make([]Event, 0, 1000)
		for ev := range s.internal {

			// If the subscriber has not completed loading, buffer everything we receive until it has
			if !s.loadCompleted {
				loadBuffer = append(loadBuffer, ev)
				continue
			}
			// If the subscriber has completed loading, send all buffered events
			if s.loadCompleted && loadBuffer != nil {
				for _, bufferedEvent := range loadBuffer {
					s.external <- bufferedEvent
				}
				loadBuffer = nil // unallocate the buffer to free up memory
			}
			s.external <- ev
		}

	}(sub)

	return index
}

func (w *Workspace) Unsubscribe(index int) {
	w.subscribersMu.Lock()
	defer w.subscribersMu.Unlock()
	sub, ok := w.subscribers[index]
	if ok {
		close(sub.external)
		delete(w.subscribers, index)
	}
}

func (w *Workspace) eventDispatcher() {
	for event := range w.events {
		w.subscribersMu.RLock()
		for subId, subscriber := range w.subscribers {
			if subscriber.kind != KindFromID(event.Id) {
				continue
			}
			// Recover from a panic if the subscriber has been closed
			// Likely this will only happen in tests but its theoretically possible in regular usage
			defer func() {
				if recover() != nil {
					w.log.Warn("dropping event as subscriber has been closed", "path", event.Id)

				}
			}()
			select {
			case subscriber.internal <- event:
				// attempt to send the event to the subscriber but drop it if the channel is full
			default:
				w.log.Warn("subscriber channel is full, dropping event", "subscription", subId, "path", event.Id)
			}
		}
		w.subscribersMu.RUnlock()
	}
}
