// Package port holds the interfaces that the application services call.
package port

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// Validator runs one Validation on a Target Project.
type Validator interface {
	// Name is the Validation Name, such as oft-tracing-validator.
	Name() string
	Validate(ctx context.Context, targetProject global.TargetProject) (domain.ValidationResult, error)
}

// ConfigurationReader reads the Validation Selection of a Target Project.
type ConfigurationReader interface {
	ValidationSelection(ctx context.Context, targetProject global.TargetProject) ([]string, error)
}
