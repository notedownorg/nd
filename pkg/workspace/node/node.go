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

	"github.com/google/uuid"
)

type kind string

type Node interface {
	ID() string
	Kind() kind
	Parent() BranchNode
	SetParent(BranchNode)
	Markdown() string
}

type BranchNode interface {
	Node
	AddChild(Node)
	RemoveChild(Node)
}

type node struct {
	id     string
	kind   kind
	parent BranchNode
}

func kindFromID(id string) kind {
	split := strings.Split(id, "/")
	return kind(split[0])
}

func newNode(kind kind) node {
	return node{
		id:   uuid.New().String(),
		kind: kind,
	}
}

func (n node) ID() string {
	return n.id
}

func (n node) Kind() kind {
	return n.kind
}

func (n node) Parent() BranchNode {
	return n.parent
}

func (n *node) SetParent(parent BranchNode) {
	n.parent = parent
}

type branchNode struct {
	node
	children []Node
}

func newBranchNode(kind kind) branchNode {
	return branchNode{
		node:     newNode(kind),
		children: make([]Node, 0),
	}
}

func (n *branchNode) AddChild(child Node) {
	if child == nil {
		return
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

func (n *branchNode) Markdown() string {
	var builder strings.Builder
	for _, child := range n.children {
		builder.WriteString(child.Markdown())
	}
	return builder.String()
}
