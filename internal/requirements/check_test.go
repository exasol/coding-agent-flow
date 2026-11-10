package requirements

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasDocumentReturnsTrueForDirectMarkdownFile(t *testing.T) {
	projectRoot := t.TempDir()
	requirementsDirectory := filepath.Join(projectRoot, "doc", "requirements")
	err := os.MkdirAll(requirementsDirectory, 0o755)
	assert.NoError(t, err)

	requirementFile := filepath.Join(requirementsDirectory, "example.md")
	err = os.WriteFile(requirementFile, nil, 0o600)
	assert.NoError(t, err)

	found, err := HasDocument(projectRoot)

	assert.NoError(t, err)
	if !found {
		t.Fatal("expected to find a requirements document.")
	}
}

func TestHasDocumentReturnsFalseWhenDirectoryIsMissing(t *testing.T) {
	projectRoot := t.TempDir()

	found, err := HasDocument(projectRoot)
	assert.NoError(t, err)
	if found {
		t.Fatal("expected no requirements document")

	}
}
