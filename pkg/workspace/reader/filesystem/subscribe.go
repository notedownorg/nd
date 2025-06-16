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

package filesystem

import (
	"github.com/notedownorg/nd/pkg/workspace/reader"
)

// Internal event type to encapsulate the clock ordering functionality
type event struct {
	reader.Event
	clock uint64
}

type subscriber struct {
	// internal channels are buffered channels that allow us to drop events when a client is unable to keep up
	// they are used to prevent blocking to ensure per-client ordering WITHOUT blocking the main event dispatcher
	internal      chan event
	external      chan reader.Event
	loadCompleted bool
}

func (r *Reader) Subscribe(ch chan reader.Event, loadInitialDocuments bool) int {
	// Add to subscribe map first to ensure we don't miss any events
	sub := &subscriber{external: ch, loadCompleted: false, internal: make(chan event, 1000)}
	r.subscriberMutex.Lock()
	index := len(r.subscribers)
	r.subscribers[index] = sub
	r.subscriberMutex.Unlock()

	if loadInitialDocuments {
		go func(s *subscriber) {
			// Recover from a panic if the subscriber has been closed
			// Likely this will only happen in tests but its theoretically possible in regular usage
			defer func() {
				if recover() != nil {
					r.log.Warn("subscriber closed before initial documents could be loaded")
				}
			}()

			for key := range r.documents {
				_, content := r.loadDocument(key)
				if content == nil {
					continue
				}
				s.external <- reader.Event{Op: reader.Load, Id: key, Content: content}
			}
			s.external <- reader.Event{Op: reader.SubscriberLoadComplete}
			s.loadCompleted = true
		}(sub)
	} else {
		sub.loadCompleted = true
	}

	// Handle subscribers in a goroutine to avoid one subscriber blocking the main event dispatcher
	// but ensure that each subscriber receives events in the correct order.
	go func(s *subscriber) {
		loadBuffer := make([]event, 0, 1000)
		for ev := range s.internal {
			// Recover from a panic if the subscriber has been closed
			// Likely this will only happen in tests but its theoretically possible in regular usage
			defer func() {
				if recover() != nil {
					r.log.Warn("dropping event as subscriber has been closed", "path", ev.Id)
				}
			}()

			// If the subscriber has not completed loading, buffer everything we receive until it has
			if !s.loadCompleted {
				loadBuffer = append(loadBuffer, ev)
				continue
			}
			// If the subscriber has completed loading, send all buffered events
			if s.loadCompleted && loadBuffer != nil {
				for _, bufferedEvent := range loadBuffer {
					s.external <- bufferedEvent.Event
				}
				loadBuffer = nil // unallocate the buffer to free up memory
			}
			s.external <- ev.Event
		}
	}(sub)

	return index
}

func (r *Reader) Unsubscribe(index int) {
	r.subscriberMutex.Lock()
	defer r.subscriberMutex.Unlock()
	sub, ok := r.subscribers[index]
	if ok {
		close(sub.internal)
		delete(r.subscribers, index)
	}
}

func (r *Reader) eventDispatcher() {
	for ev := range r.events {
		r.subscriberMutex.RLock()
		for subId, subscriber := range r.subscribers {
			// Recover from a panic if the subscriber has been closed
			// Likely this will only happen in tests but its theoretically possible in regular usage
			defer func() {
				if recover() != nil {
					r.log.Warn("dropping event as subscriber has been closed", "path", ev.Id)
				}
			}()
			select {
			case subscriber.internal <- ev:
				// attempt to send the event to the subscriber but drop it if the channel is full
			default:
				r.log.Warn("subscriber channel is full, dropping event", "subscription", subId, "path", ev.Id)
			}
		}

		r.subscriberMutex.RUnlock()
	}
}
