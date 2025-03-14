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
	"bytes"

	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v2"
)

type parser struct{}

func (p *parser) Trigger() []byte {
	return []byte{'-'}
}

func (p *parser) Open(parent gast.Node, reader text.Reader, pc gparser.Context) (gast.Node, gparser.State) {
	linenum, _ := reader.Position()
	if linenum != 0 {
		return nil, gparser.NoChildren
	}
	line, _ := reader.PeekLine()
	if !isFence(line) {
		return nil, gparser.NoChildren
	}

	return NewFrontmatter(string(line)), gparser.NoChildren
}

func (p *parser) Continue(node gast.Node, reader text.Reader, pc gparser.Context) gparser.State {
	line, segment := reader.PeekLine()
	if isFence(line) {
		node.(*Frontmatter).Closer = string(line)
		reader.Advance(segment.Len())
		return gparser.Close
	}
	node.Lines().Append(segment)
	return gparser.Continue | gparser.NoChildren
}

func (p *parser) Close(node gast.Node, reader text.Reader, pc gparser.Context) {
	lines := node.Lines()
	var buf bytes.Buffer
	for i := range lines.Len() {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}

	// Verify its a valid yaml block or remove ourselves from the tree
	meta := map[string]any{}
	if err := yaml.Unmarshal(buf.Bytes(), &meta); err != nil {
		node.Parent().RemoveChild(node.Parent(), node)
		return
	}

	d := node.(*Frontmatter)
	d.Yaml = buf.Bytes()
}

func (p *parser) CanInterruptParagraph() bool {
	return false
}

func (p *parser) CanAcceptIndentedLine() bool {
	return false
}

func isFence(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	for i := range line {
		if line[i] != '-' {
			return false
		}
	}
	return len(line) >= 3
}
