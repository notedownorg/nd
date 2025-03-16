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

package mock

import (
	"sync"

	"github.com/notedownorg/nd/pkg/workspace/reader"
)

var _ reader.Reader = &Reader{}

type Reader struct {
	vfsMutex sync.RWMutex
	vfs      map[string][]byte

	subscribers map[int]chan reader.Event
	errors      chan error
}

func NewReader() *Reader {
	return &Reader{
		vfs:         make(map[string][]byte),
		subscribers: make(map[int]chan reader.Event),
		errors:      make(chan error),
	}
}

func (r *Reader) Subscribe(ch chan reader.Event, loadInitialDocuments bool) int {
	r.subscribers[len(r.subscribers)] = ch
	index := len(r.subscribers) - 1

	if loadInitialDocuments {
		go func(s chan reader.Event) {
			r.vfsMutex.RLock()
			for key, content := range r.vfs {
				ch <- reader.Event{Op: reader.Load, Id: key, Content: content}
			}
			r.vfsMutex.RUnlock()
			ch <- reader.Event{Op: reader.SubscriberLoadComplete}
		}(ch)
	}

	return index
}

func (r *Reader) Unsubscribe(index int) {
	delete(r.subscribers, index)
}

func (r *Reader) Errors() <-chan error {
	return r.errors
}

func (r *Reader) sendEvent(event reader.Event) {
	for _, ch := range r.subscribers {
		ch <- event
	}
}
