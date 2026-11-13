package oftvalidation

import "context"

// ValidationEngine validates OFT tracing for a Target Project.
type ValidationEngine interface {
	Validate(context.Context, TargetProject) (OFTValidation, error)
}

// OFTTracingValidator validates OFT tracing in a Target Project.
type OFTTracingValidator struct {
	engine ValidationEngine
}

// NewOFTTracingValidator creates an OFT tracing validator with its validation engine.
func NewOFTTracingValidator(engine ValidationEngine) OFTTracingValidator {
	return OFTTracingValidator{engine: engine}
}

// Validate validates OFT tracing for a Target Project.
func (validator OFTTracingValidator) Validate(ctx context.Context, targetProject TargetProject) (OFTValidation, error) {
	return validator.engine.Validate(ctx, targetProject)
}
