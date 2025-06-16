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

package node

import (
	"slices"
	"strings"
)

type branchNode struct {
	node
	children []Node

	// config
	nodeOpts []nodeOption
}

type branchNodeOption func(*branchNode)

func branchNodeWithId(id string) branchNodeOption {
	return func(n *branchNode) {
		n.nodeOpts = append(n.nodeOpts, nodeWithId(id))
	}
}

func newBranchNode(kind Kind, opts ...branchNodeOption) branchNode {
	n := branchNode{children: make([]Node, 0)}
	for _, opt := range opts {
		opt(&n)
	}

	n.node = newNode(kind, n.nodeOpts...)
	n.nodeOpts = nil // unallocate to save memory

	return n
}

func (n *branchNode) AddChild(child Node) {
	if child == nil {
		return
	}

	// Automatically merge placeholder nodes
	if len(n.children) > 0 {
		last := n.children[len(n.children)-1]
		if pLast, ok := last.(*Placeholder); ok {
			if pNew, ok := child.(*Placeholder); ok {
				pLast.content = append(pLast.content, pNew.content...)
				return
			}
		}
	}

	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n *branchNode) RemoveChild(child Node) {
	if child == nil {
		return
	}
	for i, c := range n.children {
		if c.ID() == child.ID() {
			child.SetParent(nil)
			n.children = slices.Delete(n.children, i, i+1)
		}
	}
}

// Although the graph may be cyclic, children/parent is not
// Note: This does not include the root node in its search
func (n *branchNode) DepthFirstSearch(fn func(Node)) {
	for _, child := range n.children {
		if branch, ok := child.(BranchNode); ok {
			branch.DepthFirstSearch(fn)
		} else {
			fn(child)
		}
	}
}

func (n *branchNode) Markdown() string {
	var builder strings.Builder
	for _, child := range n.children {
		builder.WriteString(child.Markdown())
	}
	return builder.String()
}
