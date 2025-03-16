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
		name     string
		metadata map[string]any
		expected string
	}{
		{
			name:     "nil metadata",
			metadata: nil,
			expected: "",
		},
		{
			name:     "empty map metadata",
			metadata: map[string]any{},
			expected: "---\n---\n",
		},
		{
			name: "simple key-value pairs",
			metadata: map[string]any{
				"title":  "Test Document",
				"author": "John Doe",
			},
			expected: "---\nauthor: John Doe\ntitle: Test Document\n---\n",
		},
		{
			name: "array values",
			metadata: map[string]any{
				"tags": []string{"test", "document", "metadata"},
			},
			expected: "---\ntags:\n    - test\n    - document\n    - metadata\n---\n",
		},
		{
			name: "empty array",
			metadata: map[string]any{
				"tags": []string{},
			},
			expected: "---\ntags: []\n---\n",
		},
		{
			name: "nested maps",
			metadata: map[string]any{
				"metadata": map[string]any{
					"created": "2025-01-01",
					"updated": "2025-01-02",
				},
			},
			expected: "---\nmetadata:\n    created: \"2025-01-01\"\n    updated: \"2025-01-02\"\n---\n",
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
			expected: "---\nauthor:\n    email: john@example.com\n    name: John Doe\npublished: true\nrating: 4.5\ntags:\n    - test\ntitle: Mixed Types\nviews: 42\n---\n",
		},
		{
			name: "special characters",
			metadata: map[string]any{
				"title":       "Special: Characters!",
				"description": "Contains: colons, \"quotes\", 'apostrophes', #hashtags, @mentions",
			},
			expected: "---\ndescription: 'Contains: colons, \"quotes\", ''apostrophes'', #hashtags, @mentions'\ntitle: 'Special: Characters!'\n---\n",
		},
		{
			name: "empty strings",
			metadata: map[string]any{
				"title": "",
				"desc":  "   ",
			},
			expected: "---\ndesc: '   '\ntitle: \"\"\n---\n",
		},
		{
			name: "array with mixed types",
			metadata: map[string]any{
				"mixed": []any{42, "string", true, 3.14},
			},
			expected: "---\nmixed:\n    - 42\n    - string\n    - true\n    - 3.14\n---\n",
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
			expected: "---\nlevel1:\n    level2:\n        level3:\n            array:\n                - key: value\n---\n",
		},
		{
			name: "null values",
			metadata: map[string]any{
				"nullField": nil,
				"title":     "Document with null",
			},
			expected: "---\nnullField: null\ntitle: Document with null\n---\n",
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
			result := doc.Markdown()
			assert.Equal(t, tt.expected, result)
		})
	}
}
