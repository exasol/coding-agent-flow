package openfasttrace

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

func TestWithDiagnosticsAddsStderr(t *testing.T) {
	cause := errors.New("read OFT trace: EOF")

	err := withDiagnostics(cause, "  Error: Invalid or corrupt jarfile\n")

	if !errors.Is(err, cause) || err.Error() != "read OFT trace: EOF: Error: Invalid or corrupt jarfile" {
		t.Fatalf("got %v", err)
	}
}

func TestWithDiagnosticsKeepsErrorWithoutStderr(t *testing.T) {
	cause := errors.New("read OFT trace: EOF")

	if err := withDiagnostics(cause, " \n"); err != cause {
		t.Fatalf("got %v, want %v", err, cause)
	}
}

func TestTraceReportsAnUnreadableJar(t *testing.T) {
	jar := filepath.Join(t.TempDir(), "openfasttrace.jar")
	if err := os.WriteFile(jar, []byte("no JAR content"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(jarEnvironmentVariable, jar)

	_, err := NewAdapter().Trace(context.Background(), domain.NewTargetProject(t.TempDir()))

	if err == nil || !strings.Contains(err.Error(), "read OFT trace") || !strings.Contains(err.Error(), "Invalid or corrupt jarfile") {
		t.Fatalf("expected an unreadable OFT trace with the Java diagnostics, got %v", err)
	}
}
