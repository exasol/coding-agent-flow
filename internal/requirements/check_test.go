package requirements

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHasDocumentReturnsTrueForDirectMarkdownFile(t *testing.T) {
	projectRoot := t.TempDir()
	requirementsDirectory := filepath.Join(projectRoot, "doc", "requirements")
	if err := os.MkdirAll(requirementsDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	requirementFile := filepath.Join(requirementsDirectory, "example.md")
	if err := os.WriteFile(requirementFile, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	found, err := HasDocument(projectRoot)

	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("expected to find a requirements document.")
	}
}

func TestHasDocumentReturnsFalseWhenDirectoryIsMissing(t *testing.T) {
	projectRoot := t.TempDir()

	found, err := HasDocument(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("expected no requirements document")

	}
}
