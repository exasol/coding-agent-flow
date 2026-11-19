// Package openfasttrace implements OFTTracer with the OpenFastTrace command line.
package openfasttrace

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

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

// Trace runs OpenFastTrace on the Target Project and reads its ASpec report. Exit code 0 is valid tracing, exit code 1 is tracing with defects, and a cancellation or any other result is an error. It keeps the ASpec report and the diagnostics within the limits of runBounded.
func (Adapter) Trace(ctx context.Context, targetProject domain.TargetProject) (domain.OFTValidation, error) {
	// [impl -> req~oft-requirements-detection~1]
	jar, err := locateJar(ctx)
	if err != nil {
		return domain.OFTValidation{}, err
	}
	output, stderr, err := runBounded(ctx, maximumReportSize, "java", "-jar", jar, "trace", "--output-format", "aspec", targetProject.Path())
	if ctxErr := ctx.Err(); ctxErr != nil {
		return domain.OFTValidation{}, ctxErr
	}
	traceSucceeded, err := traceResult(err)
	if err != nil {
		return domain.OFTValidation{}, withDiagnostics(err, stderr)
	}
	validation, err := readASpecReport(output, traceSucceeded)
	if err != nil {
		return domain.OFTValidation{}, withDiagnostics(err, stderr)
	}
	return validation, nil
}

// traceResult maps the result of the OpenFastTrace run: exit code 0 is valid tracing, exit code 1 is tracing with
// defects, and every other result is a tool failure. Exit code 1 also needs a readable ASpec report.
func traceResult(runError error) (bool, error) {
	if runError == nil {
		return true, nil
	}
	var exitError *exec.ExitError
	if !errors.As(runError, &exitError) {
		return false, fmt.Errorf("run OpenFastTrace: %w", runError)
	}
	switch exitError.ExitCode() {
	case exitTracingDefects:
		return false, nil
	case exitCommandLineFail:
		return false, fmt.Errorf("run OpenFastTrace: command-line error (exit code %d)", exitCommandLineFail)
	default:
		return false, fmt.Errorf("run OpenFastTrace: unexpected exit code %d", exitError.ExitCode())
	}
}

func withDiagnostics(err error, stderr string) error {
	if message := strings.TrimSpace(stderr); message != "" {
		return fmt.Errorf("%w: %s", err, message)
	}
	return err
}
