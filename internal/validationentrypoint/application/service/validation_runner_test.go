package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/application/port"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

type fakeReader struct {
	selection []string
	err       error
}

func (reader fakeReader) ValidationSelection(context.Context, global.TargetProject) ([]string, error) {
	return reader.selection, reader.err
}

type fakeValidator struct {
	name       string
	resultName string // the Validation Name of the result; name if empty
	status     domain.ValidationStatus
	err        error
}

func (validator fakeValidator) Name() string {
	return validator.name
}

func (validator fakeValidator) Validate(context.Context, global.TargetProject) (domain.ValidationResult, error) {
	resultName := validator.resultName
	if resultName == "" {
		resultName = validator.name
	}
	return domain.NewValidationResult(resultName, validator.status, "CODE", "message"), validator.err
}

func run(t *testing.T, ctx context.Context, reader fakeReader) (domain.ValidationReport, error) {
	t.Helper()
	runner, err := NewValidationRunner(reader,
		fakeValidator{name: "passing", status: domain.Pass},
		fakeValidator{name: "finding", status: domain.Finding},
		fakeValidator{name: "failing", err: errors.New("broken at /home/user/project/secret.md")},
	)
	if err != nil {
		t.Fatal(err)
	}
	return runner.Run(ctx, global.NewTargetProject("project"))
}

func summary(report domain.ValidationReport) []string {
	var lines []string
	for _, result := range report.Results() {
		lines = append(lines, result.Validation()+": "+string(result.Status())+" "+result.Code())
	}
	return lines
}

// [utest -> req~validation-selection~1]
func TestRunRunsEachSelectedValidationInOrder(t *testing.T) {
	report, err := run(t, context.Background(), fakeReader{selection: []string{"finding", "passing"}})

	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"finding: finding CODE", "passing: pass CODE"}; !reflect.DeepEqual(summary(report), want) {
		t.Fatalf("got %v, want %v", summary(report), want)
	}
	if report.Status() != domain.Finding {
		t.Fatalf("got overall %s", report.Status())
	}
}

func TestRunReportsUnknownValidationsAndValidationErrorsAsToolFailures(t *testing.T) {
	report, err := run(t, context.Background(), fakeReader{selection: []string{"unknown", "failing", "passing"}})

	if err != nil {
		t.Fatal(err)
	}
	want := []string{"unknown: tool_failure " + UnknownValidationCode, "failing: tool_failure " + FailedValidationCode, "passing: pass CODE"}
	if !reflect.DeepEqual(summary(report), want) {
		t.Fatalf("got %v, want %v", summary(report), want)
	}
}

// A validator error can hold file paths and tool output, so the result has a fixed message and keeps the error in the
// diagnostics.
func TestRunKeepsValidatorErrorsOutOfTheMessage(t *testing.T) {
	report, err := run(t, context.Background(), fakeReader{selection: []string{"failing"}})

	if err != nil {
		t.Fatal(err)
	}
	failing, _ := report.Result("failing")
	if strings.Contains(failing.Message(), "secret") || failing.Message() == "" {
		t.Errorf("expected a fixed message without the validator error, got %q", failing.Message())
	}
	if failing.Diagnostics() != "broken at /home/user/project/secret.md" {
		t.Errorf("expected the validator error in the diagnostics, got %q", failing.Diagnostics())
	}
}

func TestRunReturnsConfigurationErrors(t *testing.T) {
	if _, err := run(t, context.Background(), fakeReader{err: errors.New("unreadable")}); err == nil {
		t.Fatal("expected the configuration error")
	}
}

func TestRunStopsWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := run(t, ctx, fakeReader{selection: []string{"passing"}}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected the cancellation, got %v", err)
	}
}

// Each validator has a unique, non-empty name, so each selected name runs exactly the wired validator.
func TestNewValidationRunnerRejectsEmptyAndDuplicateNames(t *testing.T) {
	for name, test := range map[string]struct {
		validators []fakeValidator
		want       string
	}{
		"empty name":     {[]fakeValidator{{name: ""}}, "E-CAF-7"},
		"duplicate name": {[]fakeValidator{{name: "check", status: domain.Pass}, {name: "check", status: domain.Finding}}, "E-CAF-8"},
	} {
		validators := make([]port.Validator, 0, len(test.validators))
		for _, validator := range test.validators {
			validators = append(validators, validator)
		}

		_, err := NewValidationRunner(fakeReader{}, validators...)

		if err == nil || !strings.Contains(err.Error(), test.want) {
			t.Errorf("%s: expected %s, got %v", name, test.want, err)
		}
	}
}

// Each selected Validation has one result under its name; a result under another name is a tool failure.
func TestRunReportsAResultUnderAnotherNameAsToolFailure(t *testing.T) {
	runner, err := NewValidationRunner(fakeReader{selection: []string{"check"}},
		fakeValidator{name: "check", resultName: "other", status: domain.Pass})
	if err != nil {
		t.Fatal(err)
	}

	report, err := runner.Run(context.Background(), global.NewTargetProject("project"))

	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"check: tool_failure " + MismatchingValidationCode}; !reflect.DeepEqual(summary(report), want) {
		t.Fatalf("got %v, want %v", summary(report), want)
	}
	if check, _ := report.Result("check"); !strings.Contains(check.Diagnostics(), "other") {
		t.Errorf("expected the returned name in the diagnostics, got %q", check.Diagnostics())
	}
}
