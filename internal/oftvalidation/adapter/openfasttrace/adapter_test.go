package openfasttrace

import (
	"context"
	"errors"
	"fmt"
	"github.com/exasol/coding-agent-flow/internal/global"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

const validReport = `<specdocument><specobjects doctype="feat"><specobject><id>f</id><version>1</version></specobject></specobjects></specdocument>`

// fakeJavaScript puts a java shell script with the body first on the PATH. The script answers "java -jar <jar> help"
// with the supported version. It also configures OPENFASTTRACE_JAR with an empty file, so the adapter finds a JAR.
func fakeJavaScript(t *testing.T, body string) {
	t.Helper()
	directory := t.TempDir()
	help := fmt.Sprintf("if [ \"$3\" = help ]; then echo 'OpenFastTrace %s'; exit 0; fi\n", openFastTraceVersion)
	if err := os.WriteFile(filepath.Join(directory, "java"), []byte("#!/bin/sh\n"+help+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	jar := filepath.Join(directory, "openfasttrace.jar")
	if err := os.WriteFile(jar, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv(jarEnvironmentVariable, jar)
}

// fakeJava puts a java script first on the PATH that prints the report, writes diagnostics, and exits with exitCode.
func fakeJava(t *testing.T, report string, exitCode int) {
	t.Helper()
	fakeJavaScript(t, fmt.Sprintf("cat <<'REPORT'\n%s\nREPORT\necho 'diagnostics' >&2\nexit %d", report, exitCode))
}

func trace(t *testing.T, ctx context.Context) (domain.OFTTrace, error) {
	t.Helper()
	return NewAdapter().Trace(ctx, global.NewTargetProject(t.TempDir()))
}

func TestTraceReportsACommandLineErrorAndUnexpectedExitCodes(t *testing.T) {
	for exitCode, want := range map[int]string{2: "command-line error", 3: "unexpected exit code 3", 127: "unexpected exit code 127"} {
		fakeJava(t, validReport, exitCode)

		_, err := trace(t, context.Background())

		if err == nil || !strings.Contains(err.Error(), want) || !strings.Contains(err.Error(), "diagnostics") {
			t.Errorf("exit code %d: expected an error with %q and the diagnostics, got %v", exitCode, want, err)
		}
	}
}

func TestTraceReturnsTheCancellation(t *testing.T) {
	fakeJavaScript(t, "exec sleep 30")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	started := time.Now()

	_, err := trace(t, ctx)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected the cancellation, got %v", err)
	}
	if elapsed := time.Since(started); elapsed > 5*time.Second {
		t.Fatalf("expected the trace to stop at once, took %s", elapsed)
	}
}

// Output with another root element and exit code 0 is a tool failure, not an empty valid result.
func TestTraceRefusesUnexpectedOutputWithExitCodeZero(t *testing.T) {
	fakeJava(t, "<html><body>OpenFastTrace</body></html>", 0)

	_, err := trace(t, context.Background())

	if err == nil || !strings.Contains(err.Error(), "read OFT trace") || !strings.Contains(err.Error(), "specdocument") {
		t.Fatalf("expected an unreadable ASpec report, got %v", err)
	}
}

// A report larger than maximumReportSize stops OpenFastTrace and is a tool failure.
func TestTraceStopsATooLargeReport(t *testing.T) {
	fakeJavaScript(t, fmt.Sprintf("head -c %d /dev/zero\nexit 0", maximumReportSize+1))

	_, err := trace(t, context.Background())

	if err == nil || !strings.Contains(err.Error(), fmt.Sprintf("exceeds %d bytes", maximumReportSize)) {
		t.Fatalf("expected the report limit error, got %v", err)
	}
}

// The error message keeps the first maximumDiagnosticsSize bytes of the diagnostics.
func TestTraceKeepsTheFirstDiagnostics(t *testing.T) {
	fakeJavaScript(t, "head -c 100000 /dev/zero | tr '\\0' x >&2\nexit 2")

	_, err := trace(t, context.Background())

	if err == nil || !strings.Contains(err.Error(), "command-line error") {
		t.Fatalf("expected the command-line error, got %v", err)
	}
	if !strings.Contains(err.Error(), strings.Repeat("x", maximumDiagnosticsSize)) || strings.Contains(err.Error(), strings.Repeat("x", maximumDiagnosticsSize+1)) {
		t.Fatalf("expected exactly %d bytes of diagnostics in the error, got %d bytes of error", maximumDiagnosticsSize, len(err.Error()))
	}
}
