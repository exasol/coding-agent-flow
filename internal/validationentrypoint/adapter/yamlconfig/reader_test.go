package yamlconfig

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
)

func selection(t *testing.T, content string) ([]string, error) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, FileName), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return NewReader().ValidationSelection(context.Background(), global.NewTargetProject(root))
}

// [utest -> req~validation-selection~1]
func TestValidationSelectionReadsTheSelectedValidationsInOrder(t *testing.T) {
	got, err := selection(t, "version: 1\nvalidators:\n  - second\n  - ' first '\n")

	if err != nil || !reflect.DeepEqual(got, []string{"second", "first"}) {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestValidationSelectionReportsAMissingFile(t *testing.T) {
	_, err := NewReader().ValidationSelection(context.Background(), global.NewTargetProject(t.TempDir()))

	if err == nil || !strings.Contains(err.Error(), "E-CAF-5") {
		t.Fatalf("expected E-CAF-5, got %v", err)
	}
}

func TestValidationSelectionRefusesInvalidConfigurations(t *testing.T) {
	for name, content := range map[string]string{
		"unknown field":   "version: 1\nvalidators: [a]\nextra: true\n",
		"other version":   "version: 2\nvalidators: [a]\n",
		"no validator":    "version: 1\nvalidators: []\n",
		"empty name":      "version: 1\nvalidators: ['  ']\n",
		"duplicate name":  "version: 1\nvalidators: [a, a]\n",
		"two documents":   "version: 1\nvalidators: [a]\n---\nversion: 1\n",
		"no YAML mapping": "- a\n",
	} {
		if _, err := selection(t, content); err == nil || !strings.Contains(err.Error(), "E-CAF-6") {
			t.Errorf("%s: expected E-CAF-6, got %v", name, err)
		}
	}
}

// .caf.yaml is a regular file in the Target Project root; a symbolic link or a directory is a configuration error.
func TestValidationSelectionRequiresARegularFile(t *testing.T) {
	for name, create := range map[string]func(t *testing.T, path string){
		"symbolic link": func(t *testing.T, path string) {
			target := filepath.Join(t.TempDir(), "caf.yaml")
			if err := os.WriteFile(target, []byte("version: 1\nvalidators: [a]\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, path); err != nil {
				t.Fatal(err)
			}
		},
		"directory": func(t *testing.T, path string) {
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
		},
	} {
		root := t.TempDir()
		create(t, filepath.Join(root, FileName))

		_, err := NewReader().ValidationSelection(context.Background(), global.NewTargetProject(root))

		if err == nil || !strings.Contains(err.Error(), "E-CAF-10") {
			t.Errorf("%s: expected E-CAF-10, got %v", name, err)
		}
	}
}

// .caf.yaml holds at most maximumFileSize bytes.
func TestValidationSelectionLimitsTheFileSize(t *testing.T) {
	valid := "version: 1\nvalidators: [a]\n"
	padded := valid + strings.Repeat("#", maximumFileSize-len(valid))

	if got, err := selection(t, padded); err != nil || !reflect.DeepEqual(got, []string{"a"}) {
		t.Errorf("file of %d bytes: got %v, %v", len(padded), got, err)
	}
	if _, err := selection(t, padded+"#"); err == nil || !strings.Contains(err.Error(), "E-CAF-11") {
		t.Errorf("file of %d bytes: expected E-CAF-11, got %v", len(padded)+1, err)
	}
}
