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
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

var _ BranchNode = &Document{}

type Document struct {
	branchNode
	metadata
}

var DocumentKind Kind = "Document"

func NewDocument(key string) *Document {
	return &Document{
		branchNode: newBranchNode(DocumentKind, branchNodeWithId(key)),
	}
}

func DocumentId(key string) string {
	return idFromKind(DocumentKind, key)
}

// Override markdown to include metadata
func (d *Document) Markdown() string {
	if d.metadata.root == nil {
		return d.branchNode.Markdown()
	}

	var builder strings.Builder
	builder.WriteString("---\n")

	// Maintains empty metadata vs no metadata
	if d.metadata.root.Kind == yaml.MappingNode && len(d.metadata.root.Content) == 0 {
		builder.WriteString("---\n")
	} else {
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(4)
		encoder.Encode(d.metadata.root)
		builder.Write(buf.Bytes())
		builder.WriteString("---\n")
	}

	builder.WriteString(d.branchNode.Markdown())
	return builder.String()
}

func (d *Document) SetMetadata(root *yaml.Node) {
	d.metadata.root = root
}

type metadata struct {
	root *yaml.Node
}
