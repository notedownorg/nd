package workspace

import (
	"sync"

	. "github.com/notedownorg/nd/pkg/workspace/node"
	"maps"
)

// Caches automatically handles locking and emitting change events
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
	c.ch <- Event{Op: Load, Id: node.ID(), Node: node}
}

func (c *Cache[T]) Remove(id string) {
	c.mu.Lock()
	delete(c.nodes, id)
	c.mu.Unlock()
	c.ch <- Event{Op: Delete, Id: id}
}
