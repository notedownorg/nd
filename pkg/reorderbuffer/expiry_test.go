package reorderbuffer_test

import (
	"testing"
	"time"

	"github.com/notedownorg/nd/pkg/reorderbuffer"
)

// Verify the invarients, take x events and randomise the order in which they are sent, expect all events in the correct order.
// We intentionally drop a percentage of events to ensure that the algorithm is robust to missing events (i.e. uses the timout).
// As we use uint64 as the clock its unlikely that we will overflow the clock.
func TestWithExpiryInvariants(t *testing.T) {
	input := make(chan reorderbuffer.Event[int])
	output := reorderbuffer.NewWithExpiry(time.Second, input)

	// Send 10000 events in random order
	eventCount, dropMod := 10000, 5
	for i := range eventCount {
		go func(i int) {
			if i%dropMod == 0 {
				return
			}
			input <- reorderbuffer.NewEvent(i, uint64(i))
		}(i)
	}

	// Read the events from the output channel
	for i := range eventCount {
		if i%dropMod == 0 {
			continue
		}
		emitted := <-output
		if emitted != i {
			t.Errorf("expected %d, got %d", i, emitted)
		}
	}
	close(input)
}
