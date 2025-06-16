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
	"sync"

	. "github.com/notedownorg/nd/pkg/workspace/node"
	"maps"
)

// Caches automatically handles locking and emitting change events in the correct order
type Cache[T Node] struct {
	ch    chan Event
	mu    sync.RWMutex
	nodes map[string]T
}

func newCache[T Node](ch chan Event) *Cache[T] {
	return &Cache[T]{ch: ch, nodes: make(map[string]T), mu: sync.RWMutex{}}
}

func (c *Cache[T]) Get(id string) (T, bool) {
	c.mu.RLock()
	node, ok := c.nodes[id]
	c.mu.RUnlock()
	return node, ok
}

func (c *Cache[T]) Keys() []string {
	c.mu.RLock()
	keys := make([]string, 0, len(c.nodes))
	for key := range c.nodes {
		keys = append(keys, key)
	}
	c.mu.RUnlock()
	return keys
}

func (c *Cache[T]) Values() []T {
	c.mu.RLock()
	nodes := make([]T, 0, len(c.nodes))
	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}
	c.mu.RUnlock()
	return nodes
}

func (c *Cache[T]) Entries() map[string]T {
	c.mu.RLock()
	nodes := make(map[string]T, len(c.nodes))
	maps.Copy(nodes, c.nodes)
	c.mu.RUnlock()
	return nodes
}

func (c *Cache[T]) Add(node T) {
	c.mu.Lock()
	c.nodes[node.ID()] = node
	c.mu.Unlock()
	c.ch <- Event{Op: Change, Id: node.ID(), Node: node}
}

func (c *Cache[T]) Remove(id string) {
	c.mu.Lock()
	delete(c.nodes, id)
	c.mu.Unlock()
	c.ch <- Event{Op: Delete, Id: id}
}
