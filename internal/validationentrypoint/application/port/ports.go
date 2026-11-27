// Package port holds the interfaces that the application services call.
package port

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// Validator runs one Validation on a Target Project.
type Validator interface {
	// Name is the Validation Name, such as oft-tracing-validator, without leading or trailing whitespace.
	Name() string
	Validate(ctx context.Context, targetProject global.TargetProject) (domain.ValidationResult, error)
}

// DiagnosticsSink receives the raw cause of a tool failure, such as a validator error with file paths and tool output.
// The cause stays out of the ValidationReport, which callers print and log.
type DiagnosticsSink interface {
	Record(validation, diagnostics string)
}

// ConfigurationReader reads the Validation Selection of a Target Project.
type ConfigurationReader interface {
	// ValidationSelection returns the selected Validation Names in order, without leading or trailing whitespace.
	ValidationSelection(ctx context.Context, targetProject global.TargetProject) ([]string, error)
}
