package codingagentflow

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation"
	"github.com/exasol/coding-agent-flow/internal/openfasttrace"
)

// TargetProjectPath identifies a Target Project for one validation operation.
type TargetProjectPath = oftvalidation.TargetProjectPath

// TargetProject is the project selected for OFT tracing validation.
type TargetProject = oftvalidation.TargetProject

// OFTValidation is the immutable result of validating a Target Project with OFT.
type OFTValidation = oftvalidation.OFTValidation

// OFTElementType identifies an OFT element type.
type OFTElementType = oftvalidation.OFTElementType

// OFTElementID identifies one OFT Element.
type OFTElementID = oftvalidation.OFTElementID

// OFTElementTypeRequirement identifies OFT Requirement elements.
const OFTElementTypeRequirement = oftvalidation.OFTElementTypeRequirement

// OFTTracingValidator validates OFT tracing in a Target Project.
type OFTTracingValidator struct {
	validator oftvalidation.OFTTracingValidator
}

// NewTargetProject creates a Target Project from its caller-supplied path.
func NewTargetProject(path string) TargetProject {
	return oftvalidation.NewTargetProject(path)
}

// NewOFTTracingValidator creates an OFT tracing validator.
func NewOFTTracingValidator() OFTTracingValidator {
	return OFTTracingValidator{validator: oftvalidation.NewOFTTracingValidator(openfasttrace.NewAdapter())}
}

// Validate validates OFT tracing for a Target Project.
func (validator OFTTracingValidator) Validate(ctx context.Context, targetProject TargetProject) (OFTValidation, error) {
	// [impl -> req~oft-tracing-validation~1]
	// [impl -> req~oft-requirements-detection~1]
	return validator.validator.Validate(ctx, targetProject)
}

// ValidateOFTTracing validates the OFT traceability of a Target Project.
func ValidateOFTTracing(ctx context.Context, targetProjectPath string) (OFTValidation, error) {
	return NewOFTTracingValidator().Validate(ctx, NewTargetProject(targetProjectPath))
}
