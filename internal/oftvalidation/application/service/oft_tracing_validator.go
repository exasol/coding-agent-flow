// Package service holds one application service per use case.
package service

import (
	"context"
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

func (validator OFTTracingValidator) Validate(ctx context.Context, targetProject global.TargetProject) (domain.OFTValidation, error) {
	// [impl -> req~oft-tracing-validation~1]
	return validator.tracer.Trace(ctx, targetProject)
}
