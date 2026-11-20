//go:build unix

package yamlconfig

import (
	"context"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/exasol/coding-agent-flow/internal/global"
)

// A FIFO as .caf.yaml is a configuration error, and reading stops at once instead of waiting for a writer.
func TestValidationSelectionRefusesAFIFO(t *testing.T) {
	root := t.TempDir()
	if err := syscall.Mkfifo(filepath.Join(root, FileName), 0o600); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)

	go func() {
		_, err := NewReader().ValidationSelection(context.Background(), global.NewTargetProject(root))
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "E-CAF-10") {
			t.Fatalf("expected E-CAF-10, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("reading the FIFO blocks")
	}
}
