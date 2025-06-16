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
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/notedownorg/nd/pkg/workspace/reader"
	"github.com/stretchr/testify/assert"
)

func TestFilesystem_Client_Events_SubscribeWithInitialDocuments(t *testing.T) {
	// Do the setup and ensure its correct
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewReader("testclient", dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, client.documents, 1)
	go ensureNoErrors(t, client.Errors())

	// Create a subscriber and ensure it receives the load complete event
	sub := make(chan reader.Event)
	done := false
	loaded := 0
	go func() {
		for ev := range sub {
			if ev.Op == reader.SubscriberLoadComplete {
				done = true
			}
			if ev.Op == reader.Load {
				loaded++
			}
		}
	}()

	client.Subscribe(sub, true)

	// Ensure we eventually receive the load complete event and that an event was received for each document
	waiter := func(d bool) func() bool { return func() bool { return done } }(done)
	assert.Eventually(t, waiter, 3*time.Second, time.Millisecond*200, "wg didn't finish in time")
	assert.Len(t, client.documents, loaded)
}

func TestFilesystem_Client_Events_Fuzz(t *testing.T) {
	// Do the setup and ensure its correct
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewReader("testclient", dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, client.documents, 1)
	go ensureNoErrors(t, client.Errors())

	prexistingDocs := map[string]bool{}
	for k := range client.documents {
		prexistingDocs[k] = true
	}

	// Create two subscribers
	sub1 := make(chan reader.Event)
	sub2 := make(chan reader.Event)

	got1, got2 := map[string]bool{}, map[string]bool{}
	got1Mutex, got2Mutex := sync.Mutex{}, sync.Mutex{}
	go func() {
		for {
			select {
			case ev := <-sub1:
				got1Mutex.Lock()
				switch ev.Op {
				case reader.Change:
					got1[ev.Id] = true
				case reader.Delete:
					delete(got1, ev.Id)
				}
				got1Mutex.Unlock()
			case ev := <-sub2:
				got2Mutex.Lock()
				switch ev.Op {
				case reader.Change:
					got2[ev.Id] = true
				case reader.Delete:
					delete(got2, ev.Id)
				}
				got2Mutex.Unlock()
			}
		}
	}()

	// Hook them up to the client
	client.Subscribe(sub1, false)
	client.Subscribe(sub2, false)

	// Throw a bunch of events at the client and ensure the subscribers are notified correctly
	wantAbs := map[string]bool{}
	wantRel := map[string]bool{}

	actionCount := 1000
	for range actionCount {
		switch rand.Intn(4) {
		case 0:
			wantAbs[createFile(dir, "# Test Document")] = true
		case 1:
			wantAbs[createThenUpdateFile(dir, "# Test Document Updated")] = true
		case 2:
			createThenDeleteFile(dir)
		case 3:
			wantAbs[createThenRenameFile(dir, "# Test Document")] = true
		}
	}

	// We have to make the keys relative...
	for k := range wantAbs {
		rel, _ := client.relative(k)
		if err == nil {
			wantRel[rel] = true
		}
	}

	// To remove non-determinism we need to remove any pre-existing documents from the gots
	// Because of the way go schedules goroutines, we can't guarantee that the subscribers won't receive these events
	// but they dont actually matter for real use cases as the events are idempotent
	for k := range prexistingDocs {
		got1Mutex.Lock()
		delete(got1, k)
		got1Mutex.Unlock()
		got2Mutex.Lock()
		delete(got2, k)
		got2Mutex.Unlock()
	}

	// Wait until we have handled all the events
	assert.Eventually(t, func() bool { return len(wantRel) == len(got1) && len(wantRel) == len(got2) }, 5*time.Second, time.Millisecond*200, "expected %v documents, sub1 %v sub2 %v", len(wantAbs), len(got1), len(got2))
	// time.Sleep(10 * time.Millisecond * time.Duration(actionCount))

	// Check the subscribers got the expected events
	assert.Equal(t, wantRel, got1)
	assert.Equal(t, wantRel, got2)

}
