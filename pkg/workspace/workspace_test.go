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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/notedownorg/nd/pkg/workspace/node"
)

func TestWorkspace_LoadDocument(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "document with frontmatter",
			filename: "frontmatter.md",
		},
		{
			name:     "document with code blocks",
			filename: "code_blocks.md",
		},
		{
			name:     "document with special characters",
			filename: "special_characters.md",
		},
		{
			name:     "document without frontmatter",
			filename: "no_frontmatter.md",
		},
		{
			name:     "empty document",
			filename: "empty.md",
		},
		{
			name:     "whitespace document",
			filename: "whitespace.md",
		},
		{
			name:     "document with single heading",
			filename: "single_heading.md",
		},
		{
			name:     "document with nested headings",
			filename: "nested_headings.md",
		},
		{
			name:     "document with multiple top-level headings",
			filename: "multiple_top_headings.md",
		},
		{
			name:     "document with mixed heading levels",
			filename: "mixed_headings.md",
		},
		{
			name:     "document with descending headers",
			filename: "descending_headers.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, err := NewWorkspace(".")
			require.NoError(t, err)

			content, err := os.ReadFile(filepath.Join("testdata", tt.filename))
			require.NoError(t, err)

			ws.loadDocument(content)

			// Since we don't have access to the document ID directly,
			// we'll get the first (and only) document from the workspace
			var doc *Document
			for _, d := range ws.documents {
				doc = d
				break
			}

			// Its enough to check that the round trip works
			assert.NotNil(t, doc)
			assert.Equal(t, string(content), doc.Markdown())
		})
	}
}
