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

type Reader interface {
	Subscribe(ch chan Event, loadInitialDocuments bool) int
	Unsubscribe(int)
	Errors() <-chan error
}

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
