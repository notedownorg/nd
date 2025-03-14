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
	"log/slog"
	"strings"
)

type section struct {
	branchNode

	title string
	level int
}

var sectionKind kind = "Section"

func NewSection(level int, title string) *section {
	return &section{
		// include the level when generating the id to avoid collisions
		branchNode: newBranchNode(sectionKind),
		title:      title,
		level:      level,
	}
}

func (n *section) Level() int {
	return n.level
}

// Move sets the parent AND updates the level to match the new parent
func (n *section) Move(parent BranchNode) {
	// Keep level in sync with the new parent
	// Travel up from the current parent to the root node whenever we find a section, increment the depth
	p, depth := parent, 1
	for p != nil {
		if p.Kind() == sectionKind {
			depth++
		}
		p = p.Parent()
	}
	if depth > 6 {
		slog.Warn("section depth is greater than 6, moving to correct parent", "depth", depth)
		for depth > 6 {
			parent = parent.Parent()
			depth--
		}
		depth = 6
	}
	n.branchNode.SetParent(parent)
}

func (n *section) Markdown() string {
	var builder strings.Builder
	builder.WriteString(strings.Repeat("#", n.level))
	builder.WriteString(" ")
	builder.WriteString(n.title)
	builder.WriteString(n.branchNode.Markdown())
	return builder.String()
}

func RecurseToSection(node BranchNode) *section {
	for node != nil {
		if node.Kind() == sectionKind {
			return node.(*section)
		}
		node = node.Parent()
	}
	return nil
}
