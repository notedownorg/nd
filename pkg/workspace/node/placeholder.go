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

// Placeholder nodes are used to represent areas of the document that are not parsed.
// They are primarily used to maintain the original markdown structure when we write back to disk
type placeholder struct {
	node
	content []byte
}

var placeholderKind kind = "Placeholder"

func NewPlaceholder(data []byte) *placeholder {
	return &placeholder{
		node:    newNode(placeholderKind),
		content: data,
	}
}

func (p placeholder) Markdown() string {
	return string(p.content)
}
