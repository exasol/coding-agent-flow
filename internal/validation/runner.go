package validation

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/requirements"
)

const RequirementsDocumentValidator = "requirements-document"

// Run executes the validators selected by project configuration.
func Run(ctx context.Context, projectRoot string, validators []string) Report {
	// [impl~caf-validation-routing~1 -> req~caf-validation-routing~1]
	results := make([]Result, 0, len(validators))
	for _, validator := range validators {
		if err := ctx.Err(); err != nil {
			results = append(results, Result{
				Validator: validator,
				Status:    StatusToolFailure,
				Code:      "CAF-VALIDATION-CANCELLED",
				Message:   "Validation was cancelled.",
			})
			continue
		}

		switch validator {
		case RequirementsDocumentValidator:
			results = append(results, runRequirementsDocument(projectRoot))
		default:
			results = append(results, Result{
				Validator: validator,
				Status:    StatusToolFailure,
				Code:      "CAF-UNKNOWN-VALIDATOR",
				Message:   "Project configuration selects an unsupported validator.",
			})
		}
	}
	return NewReport(results)
}

func runRequirementsDocument(projectRoot string) Result {
	found, err := requirements.HasDocument(projectRoot)
	if err != nil {
		return Result{
			Validator: RequirementsDocumentValidator,
			Status:    StatusToolFailure,
			Code:      "CAF-REQUIREMENTS-CHECK-FAILED",
			Message:   "The requirements-document check could not complete.",
		}
	}
	if !found {
		return Result{
			Validator: RequirementsDocumentValidator,
			Status:    StatusFinding,
			Code:      "CAF-REQUIREMENTS-DOCUMENT-MISSING",
			Message:   "Add a supported requirements document under doc/requirements.",
		}
	}
	return Result{
		Validator: RequirementsDocumentValidator,
		Status:    StatusPass,
		Code:      "CAF-REQUIREMENTS-DOCUMENT-PRESENT",
		Message:   "A supported requirements document is present.",
	}
}
