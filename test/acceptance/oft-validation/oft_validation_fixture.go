// Package oftvalidation provides file fixtures for OFT validation acceptance tests.
package oftvalidation

import (
	"os"
	"path/filepath"
	"testing"
)

// TargetProject writes files to a temporary target-project directory.
func TargetProject(t *testing.T, files map[string]string) string {
	t.Helper()
	directory := t.TempDir()
	for path, content := range files {
		file := filepath.Join(directory, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return directory
}
