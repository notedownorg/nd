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
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/notedownorg/nd/pkg/test/words"
	. "github.com/notedownorg/nd/pkg/workspace/node"
	"github.com/notedownorg/nd/pkg/workspace/reader/mock"
)

// Test that we correctly handle graph updates by testing the invariants:
// - For any given random sequence of events, we should be able to recreate the workspace file system exactly from the graph
func TestWorkspace_Reader(t *testing.T) {
	data := loadFilesToBytes(t, "testdata")

	reader := mock.NewReader(t)
	ws, err := NewWorkspace("test", reader)
	defer ws.Close()
	assert.NoError(t, err)

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}

	for range 10000 {
		content := data[keys[rand.IntN(len(keys))]]
		switch rand.IntN(4) {
		case 0:
			reader.Add(words.Random()+".md", content)
		case 1:
			reader.Update(reader.RandomFile(), content)
		case 2:
			reader.Remove(reader.RandomFile())
		case 3:
			reader.Rename(reader.RandomFile(), words.Random()+".md")
		}
	}

	time.Sleep(1 * time.Second)
	want, got := make(map[string]string), make(map[string]string)
	for _, key := range reader.ListFiles() {
		want[DocumentId(key)] = string(reader.GetContent(key))
	}
	for _, doc := range ws.documents.Values() {
		got[doc.ID()] = doc.Markdown()
	}

	assert.Equal(t, want, got)
}

// Test we handle the different structures possible within markdown relative to our model
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
			ws, err := NewWorkspace("test", mock.NewReader(t))

			content, err := os.ReadFile(filepath.Join("testdata", tt.filename))
			require.NoError(t, err)

			ws.loadDocument("test.md", content)

			// Since we don't have access to the document ID directly,
			// we'll get the first (and only) document from the workspace
			var doc *Document
			for _, d := range ws.documents.Values() {
				doc = d
				break
			}

			// Its enough to check that the round trip works
			require.NotNil(t, doc)
			assert.Equal(t, string(content), doc.Markdown())
		})
	}
}
