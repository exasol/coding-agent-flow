package ofttracingvalidator

import (
	"context"
	"errors"
	"strings"
	"testing"

	exaerror "github.com/exasol/error-reporting-go"

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

// A tracing error starts with the error code E-CAF-12 and keeps its cause in the error chain.
func TestValidateReturnsTracingErrorsWithTheErrorCode(t *testing.T) {
	var path string
	cause := errors.New("run OpenFastTrace: exec: \"java\": executable file not found in $PATH")
	validator := NewValidator(fakeTracingValidator{err: cause, path: &path})

	_, err := validator.Validate(context.Background(), global.NewTargetProject("project"))

	var exaError *exaerror.ExaError
	if !errors.As(err, &exaError) || !strings.HasPrefix(exaError.Error(), "E-CAF-12: ") {
		t.Fatalf("expected E-CAF-12 first in the error chain, got %v", err)
	}
	if strings.Contains(exaError.Error(), "executable file not found") {
		t.Errorf("the error code message holds details of the cause: %v", exaError)
	}
	if !errors.Is(err, cause) {
		t.Errorf("expected the cause in the error chain, got %v", err)
	}
}
