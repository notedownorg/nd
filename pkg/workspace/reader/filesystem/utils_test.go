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

package filesystem

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/notedownorg/nd/pkg/test/words"
	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

func setupTestDir(name string) (string, error) {
	// If we're running in a CI environment, we dont want to create temp directories
	// This ensures we can store the artifacts for debugging
	dir := os.Getenv("GITHUB_WORKSPACE")
	if dir == "" {
		var err error
		dir, err = os.MkdirTemp("", fmt.Sprintf("nl-%v-", name))
		if err != nil {
			return "", err
		}
	} else {
		dir = fmt.Sprintf("%v/testdata/%v", dir, name)
		if err := os.MkdirAll(dir, 0777); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func copyTestData(name string) (string, error) {
	dir, err := setupTestDir(name)
	if err != nil {
		return "", err
	}
	if err := cp.Copy("testdata/workspace", dir); err != nil {
		return "", err
	}
	return dir, nil
}

func generateTestData(name string, fileCount int) (string, error) {
	dir, err := setupTestDir(name)
	if err != nil {
		return "", err
	}
	for i := range fileCount {
		content := fmt.Sprintf("# Test Document %v", i) // maybe put more meaningful content here
		if err := writeFile(dir, fmt.Sprintf("%v.md", i), content); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func writeFile(dir string, name string, content string) error {
	return os.WriteFile(path.Join(dir, name), []byte(content), 0644)
}

func ensureNoErrors(t *testing.T, ch <-chan error) {
	for err := range ch {
		assert.NoError(t, err)
	}
}

func createFile(dir string, content string) string {
	filename := fmt.Sprintf("%v.md", words.Random())
	extraDir := words.Random()
	fullDir := fmt.Sprintf("%v/%v", dir, extraDir)

	// If the directory doesn't exist, create it
	// If it does exist we need to override fullDir to the existing directory
	// This is to avoid issues were certain OSes handle case sensitivity differently
	_, err := os.Stat(fullDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(fullDir, 0777); err != nil {
			slog.Error("failed to create directory", slog.String("dir", fullDir), slog.Any("error", err))
		}
	} else if err != nil {
		slog.Error("failed to stat directory", slog.String("dir", fullDir), slog.Any("error", err))
	} else {
		// Keep the original fullDir path - directory already exists
	}
	path := fmt.Sprintf("%v/%v", fullDir, filename)

	slog.Debug("creating file", slog.String("file", path))
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		slog.Error("failed to create file", slog.String("file", path), slog.Any("error", err))
	}
	return path
}

func createThenDeleteFile(dir string) string {
	content := "some random text"
	path := createFile(dir, content)
	slog.Debug("deleting file", slog.String("file", path))
	if err := os.Remove(path); err != nil {
		slog.Error("failed to delete file", slog.String("file", path), slog.Any("error", err))
	}
	return path
}

func createThenUpdateFile(dir string, content string) string {
	path := createFile(dir, content)
	slog.Debug("updating file", slog.String("file", path))
	if err := os.WriteFile(path, []byte("some random updated text"), 0644); err != nil {
		slog.Error("failed to update file", slog.String("file", path), slog.Any("error", err))
	}
	return path
}

func createThenRenameFile(dir string, content string) string {
	path := createFile(dir, content)
	newPath := fmt.Sprintf("%v/%v.md", dir, words.Random())
	slog.Debug("renaming file", slog.String("file", path), slog.String("new", newPath))
	if err := os.Rename(path, newPath); err != nil {
		slog.Error("failed to rename file", slog.String("file", path), slog.Any("error", err))
	}
	return newPath
}
