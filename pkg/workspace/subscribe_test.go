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

	reader := mock.NewReader(0)
	ws, err := NewWorkspace("test", reader)
	defer ws.Close()
	assert.NoError(t, err)

	clock := int64(0)
	for key, content := range data {
		clock += 1
		reader.Add(key, content, clock)
	}

	// Create a subscriber and ensure it receives the load complete event
	sub := make(chan Event)
	done := false
	got := make(map[string]struct{})
	go func() {
		for ev := range sub {
			if ev.Op == SubscriberLoadComplete {
				done = true
			}
			if ev.Op == Load {
				got[strings.Split(ev.Id, node.KindDelimiter)[1]] = struct{}{}
			}
		}
	}()

	ws.Subscribe(sub, node.DocumentKind, true)

	// Ensure we eventually receive the load complete event and that an event was received for each document node
	waiter := func(d bool) func() bool { return func() bool { return done } }(done)
	assert.Eventually(t, waiter, 3*time.Second, time.Millisecond*200, "wg didn't finish in time")

	// Finally check that we received all the documents
	want := make(map[string]struct{})
	for _, doc := range reader.ListFiles() {
		want[doc] = struct{}{}
	}
	assert.Equal(t, got, want)
}

func TestWorkspace_Events_Fuzz(t *testing.T) {
	data := loadFilesToBytes(t, "testdata")

	reader := mock.NewReader(0)
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
			if ev.Op == Delete {
				t.Logf("delete %s", ev.Id)
				mut.Lock()
				delete(got, id)
				mut.Unlock()
			}
			if ev.Op == Load || ev.Op == Change {
				t.Logf("load/change %s", ev.Id)
				mut.Lock()
				got[id] = struct{}{}
				mut.Unlock()
			}
		}
	}()

	want := make(map[string]struct{})
	clock := int64(0)
	for range 30 {
		// Default to add only, open other operations if we have more than 1 files
		operations := 1
		if len(want) > 1 {
			operations = 4
		}
		content := data[keys[rand.Intn(len(keys))]]
		newFile, existingFile := words.Random()+".md", reader.RandomFile()
		clock += 1

		switch rand.Intn(operations) {
		case 0:
			reader.Add(newFile, content, clock)
			want[newFile] = struct{}{}
		case 1:
			reader.Update(existingFile, content, clock)
		case 2:
			reader.Remove(existingFile, clock)
			delete(want, existingFile)
		case 3:
			reader.Rename(existingFile, newFile, clock, clock+1)
			clock += 1 // rename is two events (delete and create)
			delete(want, existingFile)
			want[newFile] = struct{}{}
		}
	}

	time.Sleep(5 * time.Second)

	assert.Equal(t, want, got)
}
