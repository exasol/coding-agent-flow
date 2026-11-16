// Package openfasttrace implements OFTTracer with the OpenFastTrace command line.
package openfasttrace

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/application/port"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

type Adapter struct{}

var _ port.OFTTracer = Adapter{}

func NewAdapter() Adapter {
	return Adapter{}
}

// Trace treats a successful OpenFastTrace exit as valid tracing.
func (Adapter) Trace(ctx context.Context, targetProject domain.TargetProject) (domain.OFTValidation, error) {
	// [impl -> req~oft-requirements-detection~1]
	jar, err := locateJar(ctx)
	if err != nil {
		return domain.OFTValidation{}, err
	}
	command := exec.CommandContext(ctx, "java", "-jar", jar, "trace", "--output-format", "aspec", targetProject.Path())
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	var exitError *exec.ExitError
	if err != nil && !errors.As(err, &exitError) {
		return domain.OFTValidation{}, withDiagnostics(fmt.Errorf("run OpenFastTrace: %w", err), stderr.String())
	}
	traceSucceeded := exitError == nil
	validation, err := readASpecReport(output, traceSucceeded)
	if err != nil {
		return domain.OFTValidation{}, withDiagnostics(err, stderr.String())
	}
	return validation, nil
}

func withDiagnostics(err error, stderr string) error {
	if message := strings.TrimSpace(stderr); message != "" {
		return fmt.Errorf("%w: %s", err, message)
	}
	return err
}
