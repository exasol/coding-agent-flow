// Package ofttracingvalidator implements the Validator oft-tracing-validator with the OFT Tracing Validator use case of the OFT Validation domain.
package ofttracingvalidator

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/global"
	oftdomain "github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/application/port"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// Name is the Validation Name of this Validator.
const Name = "oft-tracing-validator"

// Codes of the results of this Validator.
const (
	ValidCode   = "CAF-OFT-TRACING-VALID"
	InvalidCode = "CAF-OFT-TRACING-INVALID"
)

// tracingValidator is the OFT Tracing Validator use case of the OFT Validation domain.
type tracingValidator interface {
	Validate(ctx context.Context, targetProject global.TargetProject) (oftdomain.OFTValidation, error)
}

// Validator translates the OFT Validation of the OFT Tracing Validator into a Validation Result.
type Validator struct {
	validator tracingValidator
}

var _ port.Validator = Validator{}

// NewValidator creates a Validator that runs the given OFT Tracing Validator.
func NewValidator(validator tracingValidator) Validator {
	return Validator{validator: validator}
}

func (Validator) Name() string {
	return Name
}

// Validate passes when the OFT tracing of the Target Project is valid and reports a finding otherwise.
func (validator Validator) Validate(ctx context.Context, targetProject global.TargetProject) (domain.ValidationResult, error) {
	validation, err := validator.validator.Validate(ctx, targetProject)
	if err != nil {
		return domain.ValidationResult{}, err
	}
	if !validation.Valid() {
		return domain.NewValidationResult(Name, domain.Finding, InvalidCode, "OFT reports tracing defects. Run oft trace for the details."), nil
	}
	return domain.NewValidationResult(Name, domain.Pass, ValidCode, "The OFT tracing is valid."), nil
}
