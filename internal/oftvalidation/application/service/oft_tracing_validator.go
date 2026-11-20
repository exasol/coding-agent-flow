// Package service holds one application service per use case.
package service

import (
	"context"

	exaerror "github.com/exasol/error-reporting-go"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/application/port"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

type OFTTracingValidator struct {
	tracer port.OFTTracer
}

func NewOFTTracingValidator(tracer port.OFTTracer) OFTTracingValidator {
	return OFTTracingValidator{tracer: tracer}
}

// Validate traces the Target Project with OFT and derives the OFT Validation in the domain. The domain result matches
// the OFT Result, so the domain rules and OFT agree on each Target Project.
func (validator OFTTracingValidator) Validate(ctx context.Context, targetProject global.TargetProject) (domain.OFTValidation, error) {
	// [impl -> req~oft-tracing-validation~1]
	trace, err := validator.tracer.Trace(ctx, targetProject)
	if err != nil {
		return domain.OFTValidation{}, err
	}
	validation := trace.Validation()
	if validation.Valid() != (trace.Result() == domain.OFTSuccess) {
		return domain.OFTValidation{}, exaerror.New("E-CAF-9").
			Message("The CAF tracing rules report valid {{valid}} for {{path}}, OFT reports {{result}}.").
			Parameter("valid", validation.Valid()).
			Parameter("path", targetProject.Path()).
			Parameter("result", string(trace.Result())).
			Mitigation("Report the Target Project and the OFT version to the CAF maintainers.")
	}
	return validation, nil
}
