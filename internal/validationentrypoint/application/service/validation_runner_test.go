package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
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
	name   string
	status domain.ValidationStatus
	err    error
}

func (validator fakeValidator) Name() string {
	return validator.name
}

func (validator fakeValidator) Validate(context.Context, global.TargetProject) (domain.ValidationResult, error) {
	return domain.NewValidationResult(validator.name, validator.status, "CODE", "message"), validator.err
}

func run(t *testing.T, ctx context.Context, reader fakeReader) (domain.ValidationReport, error) {
	t.Helper()
	runner := NewValidationRunner(reader,
		fakeValidator{name: "passing", status: domain.Pass},
		fakeValidator{name: "finding", status: domain.Finding},
		fakeValidator{name: "failing", err: errors.New("broken")},
	)
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
	if failing, _ := report.Result("failing"); failing.Message() != "broken" {
		t.Fatalf("got message %q", failing.Message())
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
