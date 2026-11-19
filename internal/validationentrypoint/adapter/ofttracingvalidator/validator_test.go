package ofttracingvalidator

import (
	"context"
	"errors"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
	oftdomain "github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// fakeTracingValidator returns a fixed OFT Validation and records the Target Project Path.
type fakeTracingValidator struct {
	valid bool
	err   error
	path  *string
}

func (validator fakeTracingValidator) Validate(_ context.Context, targetProject global.TargetProject) (oftdomain.OFTValidation, error) {
	*validator.path = targetProject.Path()
	return oftdomain.NewOFTValidation(validator.valid, nil), validator.err
}

func TestValidatePassesForValidTracingAndReportsInvalidTracing(t *testing.T) {
	for valid, want := range map[bool]domain.ValidationStatus{true: domain.Pass, false: domain.Finding} {
		var path string

		result, err := NewValidator(fakeTracingValidator{valid: valid, path: &path}).Validate(context.Background(), global.NewTargetProject("project"))

		if err != nil || result.Status() != want || result.Validation() != Name || path != "project" {
			t.Errorf("valid %t: got %+v, %v, path %q; want %s", valid, result, err, path, want)
		}
	}
}

func TestValidateReturnsTracingErrors(t *testing.T) {
	var path string
	validator := NewValidator(fakeTracingValidator{err: errors.New("run OpenFastTrace: missing JAR"), path: &path})

	if _, err := validator.Validate(context.Background(), global.NewTargetProject("project")); err == nil {
		t.Fatal("expected the tracing error")
	}
}
