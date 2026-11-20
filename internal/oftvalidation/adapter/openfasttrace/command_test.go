package openfasttrace

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestLimitedBufferKeepsTheFirstBytesAndReportsTheExcessOnce(t *testing.T) {
	calls := 0
	limited := &limitedBuffer{limit: 5, onExceed: func() { calls++ }}

	for _, part := range []string{"abc", "defg", "hij"} {
		if written, err := limited.Write([]byte(part)); err != nil || written != len(part) {
			t.Fatalf("got %d, %v for %q", written, err, part)
		}
	}

	if limited.buffer.String() != "abcde" || !limited.exceeded || calls != 1 {
		t.Fatalf("got %q, exceeded %t, %d calls", limited.buffer.String(), limited.exceeded, calls)
	}
}

func TestRunBoundedReturnsSmallOutput(t *testing.T) {
	output, stderr, err := runBounded(context.Background(), 1000, "sh", "-c", "echo report; echo diagnostics >&2")

	if err != nil || string(output) != "report\n" || stderr != "diagnostics\n" {
		t.Fatalf("got %q, %q, %v", output, stderr, err)
	}
}

func TestRunBoundedStopsTheCommandForTooMuchOutput(t *testing.T) {
	started := time.Now()

	output, _, err := runBounded(context.Background(), 1000, "sh", "-c", "yes report")

	if err == nil || !strings.Contains(err.Error(), "exceeds 1000 bytes") || output != nil {
		t.Fatalf("expected the output limit error, got %q, %v", output, err)
	}
	if elapsed := time.Since(started); elapsed > 5*time.Second {
		t.Fatalf("expected the command to stop at once, took %s", elapsed)
	}
}

func TestRunBoundedKeepsTheFirstDiagnostics(t *testing.T) {
	_, stderr, err := runBounded(context.Background(), 1000, "sh", "-c", "head -c 200000 /dev/zero | tr '\\\\0' x >&2")

	if err != nil || len(stderr) != maximumDiagnosticsSize {
		t.Fatalf("got %d bytes of diagnostics, %v", len(stderr), err)
	}
}

func TestRunBoundedReturnsTheExitError(t *testing.T) {
	_, _, err := runBounded(context.Background(), 1000, "sh", "-c", "exit 1")

	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
		t.Fatalf("expected exit code 1, got %v", err)
	}
}

func TestRunWithDiagnosticsKeepsTheFirstOutputAndRunsToTheEnd(t *testing.T) {
	output, err := runWithDiagnostics(context.Background(), "sh", "-c", "head -c 20000 /dev/zero | tr '\\0' x; echo done >&2; exit 3")

	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 3 {
		t.Fatalf("expected exit code 3 after the full run, got %v", err)
	}
	if len(output) != maximumDiagnosticsSize {
		t.Fatalf("got %d bytes, want %d", len(output), maximumDiagnosticsSize)
	}
}
