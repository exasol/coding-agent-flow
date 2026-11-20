package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

type oftTracerStub struct {
	traced global.TargetProject
	trace  domain.OFTTrace
	err    error
}

func (stub *oftTracerStub) Trace(_ context.Context, targetProject global.TargetProject) (domain.OFTTrace, error) {
	stub.traced = targetProject
	return stub.trace, stub.err
}

// uncoveredFeature is a trace with the OFT Feature feat~f~1 that needs an OFT Requirement and stays uncovered.
func uncoveredFeature(t *testing.T, result domain.OFTResult) domain.OFTTrace {
	t.Helper()
	feature, err := domain.ParseOFTSpecificationItemID("feat~f~1")
	if err != nil {
		t.Fatal(err)
	}
	return domain.NewOFTTrace([]domain.OFTSpecificationItem{
		domain.NewOFTSpecificationItem(feature, domain.NewOFTCoverage([]domain.OFTArtifactType{"req"}, nil)),
	}, result)
}

// [utest -> req~oft-tracing-validation~1]
// A difference between the domain result and the OFT Result is a tool failure with an error code.
func TestOFTTracingValidatorReportsADifferenceFromTheOFTResult(t *testing.T) {
	for name, trace := range map[string]domain.OFTTrace{
		"valid in the domain, failure in OFT":   domain.NewOFTTrace(nil, domain.OFTFailure),
		"invalid in the domain, success in OFT": uncoveredFeature(t, domain.OFTSuccess),
	} {
		_, err := NewOFTTracingValidator(&oftTracerStub{trace: trace}).Validate(context.Background(), global.NewTargetProject("target-project"))

		if err == nil || !strings.Contains(err.Error(), "E-CAF-9") {
			t.Errorf("%s: expected E-CAF-9, got %v", name, err)
		}
	}
}

func TestOFTTracingValidatorReturnsTheTracerError(t *testing.T) {
	cause := errors.New("OFT failed")

	_, err := NewOFTTracingValidator(&oftTracerStub{err: cause}).Validate(context.Background(), global.NewTargetProject("target-project"))

	if !errors.Is(err, cause) {
		t.Fatalf("expected %v, got %v", cause, err)
	}
}
