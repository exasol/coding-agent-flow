// Package openfasttrace implements OFTTracer with the OpenFastTrace command line.
package openfasttrace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/application/port"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

// OpenFastTrace exit codes, see "oft help".
const (
	exitTracingDefects  = 1
	exitCommandLineFail = 2
)

type Adapter struct{}

var _ port.OFTTracer = Adapter{}

func NewAdapter() Adapter {
	return Adapter{}
}

// Trace runs OpenFastTrace on the Target Project and reads its ASpec report into an OFT Trace. Exit code 0 is the OFT Result success, exit code 1 is failure, and a cancellation or any other result is an error. It keeps the ASpec report and the diagnostics within the limits of runBounded.
func (Adapter) Trace(ctx context.Context, targetProject global.TargetProject) (domain.OFTTrace, error) {
	// [impl -> req~oft-requirements-detection~1]
	path, err := resolveDirectory(targetProject.Path())
	if err != nil {
		return domain.OFTTrace{}, err
	}
	jar, err := locateJar(ctx)
	if err != nil {
		return domain.OFTTrace{}, err
	}
	output, stderr, err := runBounded(ctx, maximumReportSize, "java", "-jar", jar, "trace", "--output-format", "aspec", path)
	if ctxErr := ctx.Err(); ctxErr != nil {
		return domain.OFTTrace{}, ctxErr
	}
	result, err := traceResult(err)
	if err != nil {
		return domain.OFTTrace{}, withDiagnostics(err, stderr)
	}
	trace, err := readASpecReport(output, result)
	if err != nil {
		return domain.OFTTrace{}, withDiagnostics(err, stderr)
	}
	return trace, nil
}

// resolveDirectory returns the Target Project path with symbolic links resolved, and fails when it is missing or is not
// a directory. OpenFastTrace reads a missing path or a symbolic link as an empty project, so CAF passes it the
// resolved directory.
func resolveDirectory(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("inspect Target Project path: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("inspect Target Project path: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("inspect Target Project path: %s is not a directory", path)
	}
	return resolved, nil
}

// traceResult maps the result of the OpenFastTrace run to the OFT Result: exit code 0 is success, exit code 1 is
// failure, and every other result is a tool failure. Exit code 1 also needs a readable ASpec report.
func traceResult(runError error) (domain.OFTResult, error) {
	if runError == nil {
		return domain.OFTSuccess, nil
	}
	var exitError *exec.ExitError
	if !errors.As(runError, &exitError) {
		return "", fmt.Errorf("run OpenFastTrace: %w", runError)
	}
	switch exitError.ExitCode() {
	case exitTracingDefects:
		return domain.OFTFailure, nil
	case exitCommandLineFail:
		return "", fmt.Errorf("run OpenFastTrace: command-line error (exit code %d)", exitCommandLineFail)
	default:
		return "", fmt.Errorf("run OpenFastTrace: unexpected exit code %d", exitError.ExitCode())
	}
}

func withDiagnostics(err error, stderr string) error {
	if message := strings.TrimSpace(stderr); message != "" {
		return fmt.Errorf("%w: %s", err, message)
	}
	return err
}
