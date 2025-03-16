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

func (r *Reader) Subscribe(ch chan reader.Event, loadInitialDocuments bool) int {
	r.subscribers[len(r.subscribers)] = ch
	index := len(r.subscribers) - 1

	if loadInitialDocuments {
		go func(s chan reader.Event) {
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
				s <- reader.Event{Op: reader.Load, Id: key, Content: content}
			}
			s <- reader.Event{Op: reader.SubscriberLoadComplete}
		}(ch)
	}

	return index
}

func (r *Reader) Unsubscribe(index int) {
	delete(r.subscribers, index)
}

func (r *Reader) eventDispatcher() {
	for event := range r.events {
		for _, subscriber := range r.subscribers {
			go func(s chan reader.Event, e reader.Event) {
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
