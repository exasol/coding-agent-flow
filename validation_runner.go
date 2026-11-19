package codingagentflow

import (
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/adapter/ofttracingvalidator"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/adapter/yamlconfig"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/application/service"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// ValidationReport is the result of the Validation Runner: one Validation Result per selected Validation and the Overall Result.
type ValidationReport = domain.ValidationReport

// ValidationResult is the result of one Validation.
type ValidationResult = domain.ValidationResult

// ValidationStatus is the Validation Result Status: pass, finding, tool_failure, or not_applicable.
type ValidationStatus = domain.ValidationStatus

const (
	Pass          = domain.Pass
	Finding       = domain.Finding
	ToolFailure   = domain.ToolFailure
	NotApplicable = domain.NotApplicable
)

// ValidationRunner runs the Validations that a Target Project selects in .caf.yaml.
type ValidationRunner = service.ValidationRunner

// NewValidationRunner creates a ValidationRunner that reads .caf.yaml in the Target Project root and runs the oft-tracing-validator Validation with the OFT Tracing Validator.
func NewValidationRunner() ValidationRunner {
	return service.NewValidationRunner(yamlconfig.NewReader(), ofttracingvalidator.NewValidator(NewOFTTracingValidator()))
}
