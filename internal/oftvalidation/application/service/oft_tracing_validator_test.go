package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

func TestOFTTracingValidatorReturnsTheTraceOfTheTargetProject(t *testing.T) {
	targetProject := domain.NewTargetProject("target-project")
	expected := domain.NewOFTValidation(true, nil)
	tracer := oftTracerStub{validation: expected}

	validation, err := NewOFTTracingValidator(&tracer).Validate(context.Background(), targetProject)

	if err != nil {
		t.Fatalf("validate OFT tracing: %v", err)
	}
	if tracer.traced != targetProject {
		t.Fatalf("expected target project %q, got %q", targetProject.Path(), tracer.traced.Path())
	}
	if !reflect.DeepEqual(validation, expected) {
		t.Fatalf("expected validation %+v, got %+v", expected, validation)
	}
}

type oftTracerStub struct {
	traced     domain.TargetProject
	validation domain.OFTValidation
}

func (stub *oftTracerStub) Trace(_ context.Context, targetProject domain.TargetProject) (domain.OFTValidation, error) {
	stub.traced = targetProject
	return stub.validation, nil
}
