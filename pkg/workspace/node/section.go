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

type Section struct {
	branchNode

	title string
	level int
}

var sectionKind Kind = "Section"

func NewSection(level int, title string) *Section {
	return &Section{
		// include the level when generating the id to avoid collisions
		branchNode: newBranchNode(sectionKind),
		title:      title,
		level:      level,
	}
}

func (n *Section) Level() int {
	return n.level
}

// Move sets the parent AND updates the level to match the new parent
func (n *Section) Move(parent BranchNode) {
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

func (n *Section) Markdown() string {
	var builder strings.Builder
	builder.WriteString(strings.Repeat("#", n.level))
	builder.WriteString(" ")
	builder.WriteString(n.title)
	builder.WriteString(n.branchNode.Markdown())
	return builder.String()
}

// RecurseToSection finds the nearest section node in the parent hierarchy
func RecurseToSection(node Node) *Section {
	for node != nil {
		if node.Kind() == sectionKind {
			return node.(*Section)
		}
		node = node.Parent()
	}
	return nil
}
