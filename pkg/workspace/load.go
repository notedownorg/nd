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
	"bytes"
	"fmt"
	"strings"

	utils "github.com/notedownorg/nd/pkg/goldmark"
	"github.com/notedownorg/nd/pkg/goldmark/extensions/frontmatter"
	. "github.com/notedownorg/nd/pkg/workspace/node"
	"github.com/notedownorg/nd/pkg/workspace/reader"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

var md = goldmark.New(goldmark.WithExtensions(frontmatter.Extension))

func (w *Workspace) processDocuments(subscription <-chan reader.Event) {
	for event := range subscription {
		// Traverse from the document in question to find the list of ids that currently exist in the graph
		// We take this approach because it's simpler to delete nodes than to try and work out the diff.
		// We may need to revisit this if we want to support undo/redo or hit performance issues in the future.
		deletionList, documentId := make(map[string]struct{}), DocumentId(event.Id)
		if doc, ok := w.documents.Get(documentId); ok {
			doc.DepthFirstSearch(func(n Node) { deletionList[n.ID()] = struct{}{} })
		}

		switch event.Op {
		case reader.Load:
		case reader.Change:
			doc := w.loadDocument(event.Id, event.Content)
			w.documents.Add(doc)
			w.log.Debug("updated document from event", "id", documentId, "event_id", event.Id)
		case reader.Delete:
			w.deleteNode(documentId)
			w.log.Debug("deleted document from event", "id", documentId, "event_id", event.Id)
		}

		// Delete the children of the original (or deleted) document that are still accessible from the workspace struct
		for id := range deletionList {
			w.deleteNode(id)
		}
	}
}

func (w *Workspace) loadDocument(id string, content []byte) *Document {
	doc := NewDocument(id)

	// Walk the ast building our graph
	// ast.Walk(tree, Debug())
	tree := md.Parser().Parse(text.NewReader([]byte(content)))

	// Keep track of the last position we have processed to ensure we don't lose any content between blocks
	// Keep track of the parents so we can add children correctly
	curr, parents := 0, []BranchNode{doc}
	ast.Walk(tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch node := node.(type) {
			case *frontmatter.Frontmatter:
				if len(node.Yaml) != 0 {
					var root yaml.Node
					if err := yaml.Unmarshal(node.Yaml, &root); err != nil {
						w.log.Error("failed to unmarshal yaml", "error", err)
					}
					doc.SetMetadata(&root)
				}
				curr = utils.End(0, node) + len(node.Closer)

			// If the node is a heading we need to create a new section
			// Setext headings are not currently supported and will be converted to ATX
			case *ast.Heading:
				start, end := node.Lines().At(0).Start, utils.End(curr, node)

				// Maintain the content (usually newlines) between the last block and the current block up to the start of the #
				prefix := content[curr:start]
				var buf bytes.Buffer
				for _, byte := range prefix {
					if byte == '#' {
						break
					}
					buf.WriteByte(byte)
				}
				if len(buf.Bytes()) > 0 {
					parents[len(parents)-1].AddChild(NewPlaceholder(buf.Bytes()))
				}

				// If the node level is smaller or equal to the latest section we need to pop the parents
				if section := RecurseToSection(parents[len(parents)-1]); section != nil && node.Level <= section.Level() {
					parents = parents[:len(parents)-1]
				}

				// Now we can create the section
				section := NewSection(node.Level, string(content[start:end])) // content[start:end] is the title in Goldmark
				parents[len(parents)-1].AddChild(section)
				parents = append(parents, section)
				curr = end

			default:
				// If the node is a block we're not currently interested in we need to persist so we can write back later
				if node.Type() == ast.TypeBlock {
					start := curr
					end := utils.End(start, node)
					parents[len(parents)-1].AddChild(NewPlaceholder(content[start:end]))
					curr = end
				}
			}
		}
		return ast.WalkContinue, nil
	})

	// Add the trailing content
	doc.AddChild(NewPlaceholder(content[curr:]))
	return doc
}

func Debug() func(node ast.Node, entering bool) (ast.WalkStatus, error) {
	depth := 0
	return func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			depth++
			if depth == 1 {
				fmt.Printf("%s\n", node.Kind().String())
			} else {
				fmt.Printf("%s%s\n", strings.Repeat("    ", depth-1), node.Kind().String())
			}
		} else {
			depth--
		}
		return ast.WalkContinue, nil
	}
}
