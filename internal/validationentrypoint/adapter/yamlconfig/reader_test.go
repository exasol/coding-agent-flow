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

	if err == nil || !strings.Contains(err.Error(), "E-CAF-2") {
		t.Fatalf("expected E-CAF-2, got %v", err)
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
		if _, err := selection(t, content); err == nil || !strings.Contains(err.Error(), "E-CAF-3") {
			t.Errorf("%s: expected E-CAF-3, got %v", name, err)
		}
	}
}
