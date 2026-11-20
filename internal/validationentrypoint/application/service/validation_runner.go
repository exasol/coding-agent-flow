// Package service holds one application service per use case.
package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	exaerror "github.com/exasol/error-reporting-go"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/application/port"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// Codes of the results that the Validation Runner creates.
const (
	UnknownValidationCode     = "CAF-UNKNOWN-VALIDATOR"
	FailedValidationCode      = "CAF-VALIDATION-FAILED"
	MismatchingValidationCode = "CAF-VALIDATION-NAME-MISMATCH"
)

// failedValidationMessage is the message of a failed Validation. Tool errors can hold file paths and command output, so
// the message names only the error code, and the error text goes to the diagnostics of the result.
const failedValidationMessage = "The validator failed."

// errorCodePattern matches the error code at the start of an ExaError text, such as E-CAF-5.
var errorCodePattern = regexp.MustCompile(`^[EWF]-[A-Z]+-[0-9]+`)

// failureMessage names the error code of the first ExaError in the error chain, or returns failedValidationMessage.
func failureMessage(err error) string {
	var exaError *exaerror.ExaError
	if errors.As(err, &exaError) {
		if code := errorCodePattern.FindString(exaError.Error()); code != "" {
			return "The validator failed with " + code + "."
		}
	}
	return failedValidationMessage
}

type ValidationRunner struct {
	reader     port.ConfigurationReader
	validators map[string]port.Validator
}

// NewValidationRunner creates a ValidationRunner that reads the Validation Selection with reader and runs the validators.
// Each validator needs a unique, non-empty name, so each selected name runs exactly the wired validator.
func NewValidationRunner(reader port.ConfigurationReader, validators ...port.Validator) (ValidationRunner, error) {
	byName := make(map[string]port.Validator, len(validators))
	for _, validator := range validators {
		name := validator.Name()
		if name == "" {
			return ValidationRunner{}, exaerror.New("E-CAF-7").
				Message("A validator of the Validation Runner has an empty name.").
				Mitigation("Give each validator a unique name.")
		}
		if _, duplicate := byName[name]; duplicate {
			return ValidationRunner{}, exaerror.New("E-CAF-8").
				Message("Two validators of the Validation Runner have the name {{name}}.").
				Parameter("name", name).
				Mitigation("Give each validator a unique name.")
		}
		byName[name] = validator
	}
	return ValidationRunner{reader: reader, validators: byName}, nil
}

// Run runs each selected Validation in the selected order and returns the Validation Report.
func (runner ValidationRunner) Run(ctx context.Context, targetProject global.TargetProject) (domain.ValidationReport, error) {
	// [impl -> req~validation-selection~1]
	selection, err := runner.reader.ValidationSelection(ctx, targetProject)
	if err != nil {
		return domain.ValidationReport{}, err
	}
	if err := ctx.Err(); err != nil {
		return domain.ValidationReport{}, err
	}
	results := make([]domain.ValidationResult, 0, len(selection))
	for _, name := range selection {
		results = append(results, runner.run(ctx, targetProject, name))
		// A cancellation during the validator is the result of the run, not a failure of the validator.
		if err := ctx.Err(); err != nil {
			return domain.ValidationReport{}, err
		}
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
		return domain.NewValidationResult(name, domain.ToolFailure, FailedValidationCode, failureMessage(err)).WithDiagnostics(err.Error())
	}
	if result.Validation() != name {
		return domain.NewValidationResult(name, domain.ToolFailure, MismatchingValidationCode, "The validator returned a result for another validation.").
			WithDiagnostics(fmt.Sprintf("selected %q, returned %q", name, result.Validation()))
	}
	return result
}
