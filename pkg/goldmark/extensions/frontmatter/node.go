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

package frontmatter

import (
	"github.com/yuin/goldmark/ast"
)

type Frontmatter struct {
	ast.BaseBlock

	// Allow consumer to choose how to process the yaml
	Yaml []byte

	// Include leading and trailing whitespace so we can reconstruct more accurately
	Opener, Closer string
}

// Dump implements Node.Dump .
func (n *Frontmatter) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindFrontmatter is a NodeKind of the Frontmatter node.
var KindFrontmatter = ast.NewNodeKind("Frontmatter")

// Kind implements Node.Kind.
func (n *Frontmatter) Kind() ast.NodeKind {
	return KindFrontmatter
}

// Text implements Node.Text.
//
// Deprecated: Use other properties of the node to get the text value(i.e. Frontmatter.Lines).
func (n *Frontmatter) Text(source []byte) []byte {
	return n.Lines().Value(source)
}

// NewFrontmatter returns a new Frontmatter node.
func NewFrontmatter(opener string) *Frontmatter {
	return &Frontmatter{
		BaseBlock: ast.BaseBlock{},
		Opener:    opener,
	}
}

func IsFrontmatter(node ast.Node) bool {
	_, ok := node.(*Frontmatter)
	return ok
}
