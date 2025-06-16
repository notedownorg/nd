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
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/notedownorg/nd/pkg/test/words"
	"github.com/notedownorg/nd/pkg/workspace/node"
	"github.com/notedownorg/nd/pkg/workspace/reader/mock"
	"github.com/stretchr/testify/assert"
)

func TestWorkspace_Events_SubscribeWithIntialDocuments_Sync(t *testing.T) {
	data := loadFilesToBytes(t, "testdata")

	reader := mock.NewReader(t)
	ws, err := NewWorkspace("test", reader)
	defer ws.Close()
	assert.NoError(t, err)

	for key, content := range data {
		reader.Add(key, content)
	}

	// Wait for workspace to process all documents
	time.Sleep(100 * time.Millisecond)

	// Create a subscriber and ensure it receives the load complete event
	sub := make(chan Event)
	done := false
	got := make(map[string]struct{})
	mut := sync.Mutex{}
	go func() {
		for ev := range sub {
			if ev.Op == SubscriberLoadComplete {
				mut.Lock()
				done = true
				mut.Unlock()
			}
			if ev.Op == Load {
				mut.Lock()
				got[strings.Split(ev.Id, node.KindDelimiter)[1]] = struct{}{}
				mut.Unlock()
			}
		}
	}()

	ws.Subscribe(sub, node.DocumentKind, true)

	// Ensure we eventually receive the load complete event
	waiter := func() bool { 
		mut.Lock()
		defer mut.Unlock()
		return done 
	}
	assert.Eventually(t, waiter, 3*time.Second, time.Millisecond*200, "wg didn't finish in time")

	// Finally check that we received all the documents
	want := make(map[string]struct{})
	for _, doc := range reader.ListFiles() {
		want[doc] = struct{}{}
	}
	mut.Lock()
	assert.Equal(t, want, got)
	mut.Unlock()
}

func TestWorkspace_Events_Fuzz(t *testing.T) {
	data := loadFilesToBytes(t, "testdata")

	reader := mock.NewReader(t)
	ws, err := NewWorkspace("test", reader)
	defer ws.Close()
	assert.NoError(t, err)

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}

	sub := make(chan Event)
	mut := sync.Mutex{}
	got := make(map[string]struct{})
	ws.Subscribe(sub, node.DocumentKind, false)
	go func() {
		for ev := range sub {
			id := strings.Split(ev.Id, node.KindDelimiter)[1]
			if ev.Op == Load || ev.Op == Change {
				t.Logf("received load/change %s (%d)", ev.Id, ev.Op)
				mut.Lock()
				got[id] = struct{}{}
				mut.Unlock()
			}
			if ev.Op == Delete {
				t.Logf("received delete %s (%d)", ev.Id, ev.Op)
				mut.Lock()
				delete(got, id)
				mut.Unlock()
			}
		}
	}()

	want := make(map[string]struct{})
	for range 30 {
		// Default to add only, open other operations if we have more than 1 files
		operations := 1
		if len(want) > 1 {
			operations = 4
		}
		content := data[keys[rand.Intn(len(keys))]]
		newFile, existingFile := words.Random()+".md", reader.RandomFile()

		switch rand.Intn(operations) {
		case 0:
			t.Logf("add %s", newFile)
			reader.Add(newFile, content)
			want[newFile] = struct{}{}
		case 1:
			t.Logf("update %s", existingFile)
			reader.Update(existingFile, content)
		case 2:
			t.Logf("remove %s", existingFile)
			reader.Remove(existingFile)
			delete(want, existingFile)
		case 3:
			t.Logf("rename %s to %s", existingFile, newFile)
			reader.Rename(existingFile, newFile)
			delete(want, existingFile)
			want[newFile] = struct{}{}
		}
	}

	time.Sleep(2 * time.Second)

	mut.Lock()
	assert.Equal(t, want, got)
	mut.Unlock()
}
