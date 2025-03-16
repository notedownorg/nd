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
	DepthFirstSearch(func(Node))
}

type node struct {
	id     string
	kind   kind
	parent BranchNode
}

const kindDelimiter = "|"

func KindFromID(id string) kind {
	split := strings.Split(id, kindDelimiter)
	return kind(split[0])
}

type nodeOption func(*node)

// nodeWithId will still preserve the kind in the id
func nodeWithId(id string) nodeOption {
	return func(n *node) {
		if strings.HasPrefix(id, string(n.kind)) {
			n.id = id
			return
		}
		n.id = idFromKind(n.kind, id)
	}
}

func idFromKind(kind kind, id string) string {
	if strings.HasPrefix(id, string(kind)) {
		return id
	}
	return string(kind) + kindDelimiter + id
}

func newNode(kind kind, opts ...nodeOption) node {
	n := node{
		id:   uuid.New().String(),
		kind: kind,
	}
	for _, opt := range opts {
		opt(&n)
	}
	return n
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
