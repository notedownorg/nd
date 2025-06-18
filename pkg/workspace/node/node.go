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

type Kind string

type Node interface {
	ID() string
	Kind() Kind
	Parent() BranchNode
	SetParent(BranchNode)
	Markdown() string
}

type BranchNode interface {
	Node
	AddChild(Node)
	RemoveChild(Node)
	DepthFirstSearch(func(Node))
	DirectChildren() []Node
}

type node struct {
	id     string
	kind   Kind
	parent BranchNode
}

const KindDelimiter = "|"

func KindFromID(id string) Kind {
	split := strings.Split(id, KindDelimiter)
	return Kind(split[0])
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

func idFromKind(kind Kind, id string) string {
	if strings.HasPrefix(id, string(kind)) {
		return id
	}
	return string(kind) + KindDelimiter + id
}

func newNode(kind Kind, opts ...nodeOption) node {
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

func (n node) Kind() Kind {
	return n.kind
}

func (n node) Parent() BranchNode {
	return n.parent
}

func (n *node) SetParent(parent BranchNode) {
	n.parent = parent
}
