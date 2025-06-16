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
)

type fatalf interface {
	Fatalf(format string, args ...any)
}

func loadFilesToBytes(t fatalf, dir string) map[string][]byte {
	filesData := make(map[string][]byte)

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read directory: %s, error: %v\n", dir, err)
	}

	// Iterate over files
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		// Get full file path
		filePath := filepath.Join(dir, file.Name())

		// Read file content
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %s, error: %v\n", filePath, err)
		}

		filesData[file.Name()] = data
	}

	return filesData
}
