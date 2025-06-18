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
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDocument_Markdown(t *testing.T) {
	tests := []struct {
		name             string
		metadata         map[string]any
		expectedMarkdown string
		expectedMetadata map[string]interface{}
	}{
		{
			name:             "nil metadata",
			metadata:         nil,
			expectedMarkdown: "",
			expectedMetadata: nil,
		},
		{
			name:             "empty map metadata",
			metadata:         map[string]any{},
			expectedMarkdown: "---\n---\n",
			expectedMetadata: map[string]interface{}{},
		},
		{
			name: "simple key-value pairs",
			metadata: map[string]any{
				"title":  "Test Document",
				"author": "John Doe",
			},
			expectedMarkdown: "---\nauthor: John Doe\ntitle: Test Document\n---\n",
			expectedMetadata: map[string]interface{}{
				"title":  "Test Document",
				"author": "John Doe",
			},
		},
		{
			name: "array values",
			metadata: map[string]any{
				"tags": []string{"test", "document", "metadata"},
			},
			expectedMarkdown: "---\ntags:\n    - test\n    - document\n    - metadata\n---\n",
			expectedMetadata: map[string]interface{}{
				"tags": []interface{}{"test", "document", "metadata"},
			},
		},
		{
			name: "empty array",
			metadata: map[string]any{
				"tags": []string{},
			},
			expectedMarkdown: "---\ntags: []\n---\n",
			expectedMetadata: map[string]interface{}{
				"tags": []interface{}{},
			},
		},
		{
			name: "nested maps",
			metadata: map[string]any{
				"metadata": map[string]any{
					"created": "2025-01-01",
					"updated": "2025-01-02",
				},
			},
			expectedMarkdown: "---\nmetadata:\n    created: \"2025-01-01\"\n    updated: \"2025-01-02\"\n---\n",
			expectedMetadata: map[string]interface{}{
				"metadata": map[string]interface{}{
					"created": "2025-01-01",
					"updated": "2025-01-02",
				},
			},
		},
		{
			name: "mixed types",
			metadata: map[string]any{
				"title":     "Mixed Types",
				"published": true,
				"views":     42,
				"rating":    4.5,
				"tags":      []string{"test"},
				"author": map[string]any{
					"name":  "John Doe",
					"email": "john@example.com",
				},
			},
			expectedMarkdown: "---\nauthor:\n    email: john@example.com\n    name: John Doe\npublished: true\nrating: 4.5\ntags:\n    - test\ntitle: Mixed Types\nviews: 42\n---\n",
			expectedMetadata: map[string]interface{}{
				"title":     "Mixed Types",
				"published": true,
				"views":     42,
				"rating":    4.5,
				"tags":      []interface{}{"test"},
				"author": map[string]interface{}{
					"name":  "John Doe",
					"email": "john@example.com",
				},
			},
		},
		{
			name: "special characters",
			metadata: map[string]any{
				"title":       "Special: Characters!",
				"description": "Contains: colons, \"quotes\", 'apostrophes', #hashtags, @mentions",
			},
			expectedMarkdown: "---\ndescription: 'Contains: colons, \"quotes\", ''apostrophes'', #hashtags, @mentions'\ntitle: 'Special: Characters!'\n---\n",
			expectedMetadata: map[string]interface{}{
				"title":       "Special: Characters!",
				"description": "Contains: colons, \"quotes\", 'apostrophes', #hashtags, @mentions",
			},
		},
		{
			name: "empty strings",
			metadata: map[string]any{
				"title": "",
				"desc":  "   ",
			},
			expectedMarkdown: "---\ndesc: '   '\ntitle: \"\"\n---\n",
			expectedMetadata: map[string]interface{}{
				"title": "",
				"desc":  "   ",
			},
		},
		{
			name: "array with mixed types",
			metadata: map[string]any{
				"mixed": []any{42, "string", true, 3.14},
			},
			expectedMarkdown: "---\nmixed:\n    - 42\n    - string\n    - true\n    - 3.14\n---\n",
			expectedMetadata: map[string]interface{}{
				"mixed": []interface{}{42, "string", true, 3.14},
			},
		},
		{
			name: "deeply nested structure",
			metadata: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"array": []any{
								map[string]any{
									"key": "value",
								},
							},
						},
					},
				},
			},
			expectedMarkdown: "---\nlevel1:\n    level2:\n        level3:\n            array:\n                - key: value\n---\n",
			expectedMetadata: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": map[string]interface{}{
							"array": []interface{}{
								map[string]interface{}{
									"key": "value",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "null values",
			metadata: map[string]any{
				"nullField": nil,
				"title":     "Document with null",
			},
			expectedMarkdown: "---\nnullField: null\ntitle: Document with null\n---\n",
			expectedMetadata: map[string]interface{}{
				"nullField": nil,
				"title":     "Document with null",
			},
		},
		{
			name: "round-trip consistency",
			metadata: map[string]any{
				"title":     "Test Document",
				"author":    "Test Author",
				"tags":      []string{"test", "metadata"},
				"priority":  1,
				"published": true,
				"config": map[string]any{
					"enabled": true,
					"value":   42,
				},
			},
			expectedMarkdown: "---\nauthor: Test Author\nconfig:\n    enabled: true\n    value: 42\npriority: 1\npublished: true\ntags:\n    - test\n    - metadata\ntitle: Test Document\n---\n",
			expectedMetadata: map[string]interface{}{
				"title":     "Test Document",
				"author":    "Test Author",
				"tags":      []interface{}{"test", "metadata"},
				"priority":  1,
				"published": true,
				"config": map[string]interface{}{
					"enabled": true,
					"value":   42,
				},
			},
		},
		// {
		// 	name: "unicode + emoji characters",
		// 	metadata: map[string]any{
		// 		"title": "Unicode ♥ Test 🚀",
		// 		"tags":  []string{"emoji 👍", "unicode ★"},
		// 	},
		// 	expected: "---\ntags:\n    - \"emoji 👍\"\n    - \"unicode ★\"\ntitle: \"Unicode ♥ Test 🚀\"\n---\n",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create document
			doc := NewDocument("test.md")

			// Set metadata if provided
			if tt.metadata != nil {
				var node yaml.Node
				err := node.Encode(tt.metadata)
				assert.NoError(t, err)
				doc.SetMetadata(&node)
			}

			// Test the Markdown method
			markdownResult := doc.Markdown()
			assert.Equal(t, tt.expectedMarkdown, markdownResult)

			// Test the GetMetadata method
			metadataResult := doc.GetMetadata()
			assert.Equal(t, tt.expectedMetadata, metadataResult)
		})
	}
}
