// Package service holds one application service per use case.
package service

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/application/port"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// Codes of the results that the Validation Runner creates.
const (
	UnknownValidationCode = "CAF-UNKNOWN-VALIDATOR"
	FailedValidationCode  = "CAF-VALIDATION-FAILED"
)

type ValidationRunner struct {
	reader     port.ConfigurationReader
	validators map[string]port.Validator
}

// NewValidationRunner creates a ValidationRunner that reads the Validation Selection with reader and runs the validators.
func NewValidationRunner(reader port.ConfigurationReader, validators ...port.Validator) ValidationRunner {
	byName := make(map[string]port.Validator, len(validators))
	for _, validator := range validators {
		byName[validator.Name()] = validator
	}
	return ValidationRunner{reader: reader, validators: byName}
}

// Run runs each selected Validation in the selected order and returns the Validation Report.
func (runner ValidationRunner) Run(ctx context.Context, targetProject global.TargetProject) (domain.ValidationReport, error) {
	// [impl -> req~validation-selection~1]
	selection, err := runner.reader.ValidationSelection(ctx, targetProject)
	if err != nil {
		return domain.ValidationReport{}, err
	}
	results := make([]domain.ValidationResult, 0, len(selection))
	for _, name := range selection {
		if err := ctx.Err(); err != nil {
			return domain.ValidationReport{}, err
		}
		results = append(results, runner.run(ctx, targetProject, name))
	}
	return domain.NewValidationReport(results), nil
}

// run returns the result of one Validation. An Unknown Validation and a Validation error are tool failures.
func (runner ValidationRunner) run(ctx context.Context, targetProject global.TargetProject, name string) domain.ValidationResult {
	validator, known := runner.validators[name]
	if !known {
		return domain.NewValidationResult(name, domain.ToolFailure, UnknownValidationCode, "The Target Project selects an unknown validation.")
	}
	result, err := validator.Validate(ctx, targetProject)
	if err != nil {
		return domain.NewValidationResult(name, domain.ToolFailure, FailedValidationCode, err.Error())
	}
	return result
}
