package oftvalidation

import (
	"context"
	"reflect"
	"testing"
)

// [utest -> req~oft-requirements-detection~1]
func TestOFTTracingValidatorDelegatesTargetProjectToValidationEngine(t *testing.T) {
	expectedProject := NewTargetProject("target-project")
	expectedValidation := OFTValidation{Valid: true}
	engine := validationEngineStub{project: expectedProject, validation: expectedValidation}

	validation, err := NewOFTTracingValidator(&engine).Validate(context.Background(), expectedProject)

	if err != nil {
		t.Fatalf("validate OFT tracing: %v", err)
	}
	if engine.project.Path() != expectedProject.Path() {
		t.Fatalf("expected target project %q, got %q", expectedProject.Path(), engine.project.Path())
	}
	if !reflect.DeepEqual(validation, expectedValidation) {
		t.Fatalf("expected validation %#v, got %#v", expectedValidation, validation)
	}
}

type validationEngineStub struct {
	project    TargetProject
	validation OFTValidation
}

func (stub *validationEngineStub) Validate(_ context.Context, project TargetProject) (OFTValidation, error) {
	stub.project = project
	return stub.validation, nil
}
