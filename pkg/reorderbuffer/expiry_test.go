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
