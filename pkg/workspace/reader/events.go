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

package reader

import "slices"

type Event struct {
	Op      Operation
	Id      string
	Content []byte
}

type Operation uint32

const (
	// Signal that this document was present when the client was created or when the subscriber subscribed
	Load Operation = iota

	// Signal that this document has been updated or created
	Change

	// Signal that this document has been deleted
	Delete

	// Signal that the subscriber has received all existing documents present at the time of subscription
	SubscriberLoadComplete
)

type subscribeOptions func(*Reader, chan Event)

// Load all existing documents as events to the new subscriber.
// Once all events have been sent, the a LoadComplete event is sent.
func WithInitialDocuments() subscribeOptions {
	return func(client *Reader, sub chan Event) {
		go func(s chan Event) {
			for key := range client.documents {
				_, content := client.loadDocument(key)
				if content == nil {
					continue
				}
				s <- Event{Op: Load, Id: key, Content: content}
			}
			s <- Event{Op: SubscriberLoadComplete}
		}(sub)
	}
}

func (r *Reader) Subscribe(ch chan Event, opts ...subscribeOptions) int {
	r.subscribers = append(r.subscribers, ch)
	index := len(r.subscribers) - 1

	// Apply any subscribeOptions
	for _, opt := range opts {
		opt(r, ch)
	}
	return index
}

func (r *Reader) Unsubscribe(index int) {
	r.subscribers = slices.Delete(r.subscribers, index, index+1)
}

func (r *Reader) eventDispatcher() {
	for event := range r.events {
		for _, subscriber := range r.subscribers {
			go func(s chan Event, e Event) {
				// Recover from a panic if the subscriber has been closed
				// Likely this will only happen in tests but its theoretically possible in regular usage
				defer func() {
					if recover() != nil {
						r.log.Warn("dropping event as subscriber has been closed", "path", e.Id)
					}
				}()
				s <- e
			}(subscriber, event)
		}
	}
}
