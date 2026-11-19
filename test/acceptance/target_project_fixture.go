package acceptance

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
)

// File is a file of a target-project fixture, with a slash-separated path relative to the Target Project.
type File struct{ Path, Content string }

// NewFile creates a target-project file whose content is template with the parameters filled in by fmt.Sprintf.
func NewFile(path, template string, parameters ...any) File {
	return File{Path: path, Content: fmt.Sprintf(template, parameters...)}
}

// TargetProject writes the files to a temporary directory and returns it as the Target Project.
func TargetProject(t *testing.T, files []File) codingagentflow.TargetProject {
	return codingagentflow.NewTargetProject(WriteFiles(t, files))
}

// WriteFiles writes the files to a temporary directory and returns its path.
func WriteFiles(t *testing.T, files []File) string {
	t.Helper()
	directory := t.TempDir()
	for _, file := range files {
		path := filepath.Join(directory, filepath.FromSlash(file.Path))
		if !strings.HasPrefix(path, filepath.Clean(directory)+string(filepath.Separator)) {
			t.Fatalf("file %q leaves target project", file.Path)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(file.Content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return directory
}

// ExpectUnchanged reports an error if the Target Project has added, modified, or deleted files compared with the
// fixture files.
func ExpectUnchanged(targetProject codingagentflow.TargetProject, files []File) error {
	expected := make(map[string]string, len(files))
	for _, file := range files {
		expected[file.Path] = file.Content
	}
	var unexpected []string
	found := make(map[string]bool, len(files))
	err := filepath.WalkDir(targetProject.Path(), func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		relative, err := filepath.Rel(targetProject.Path(), path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		path = filepath.ToSlash(relative)
		found[path] = true
		if want, isFixture := expected[path]; !isFixture {
			unexpected = append(unexpected, "added "+path)
		} else if string(content) != want {
			unexpected = append(unexpected, "modified "+path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for path := range expected {
		if !found[path] {
			unexpected = append(unexpected, "deleted "+path)
		}
	}
	if len(unexpected) > 0 {
		slices.Sort(unexpected)
		return fmt.Errorf("target project changed: %v", unexpected)
	}
	return nil
}
