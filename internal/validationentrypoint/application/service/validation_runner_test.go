package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	exaerror "github.com/exasol/error-reporting-go"

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

// recordingSink keeps the diagnostics of each Validation.
type recordingSink map[string]string

func (sink recordingSink) Record(validation, diagnostics string) {
	sink[validation] = diagnostics
}

func run(t *testing.T, ctx context.Context, reader fakeReader) (domain.ValidationReport, error) {
	report, _, err := runRecording(t, ctx, reader)
	return report, err
}

// runRecording runs the fake validators and returns the diagnostics that the DiagnosticsSink receives.
func runRecording(t *testing.T, ctx context.Context, reader fakeReader) (domain.ValidationReport, recordingSink, error) {
	t.Helper()
	runner, err := NewValidationRunner(reader,
		fakeValidator{name: "passing", status: domain.Pass},
		fakeValidator{name: "finding", status: domain.Finding},
		fakeValidator{name: "failing", err: errors.New("broken at /home/user/project/secret.md")},
	)
	if err != nil {
		t.Fatal(err)
	}
	sink := recordingSink{}
	report, err := runner.WithDiagnosticsSink(sink).Run(ctx, global.NewTargetProject("project"))
	return report, sink, err
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

// A validator error can hold file paths, tool output, and credentials, so the DiagnosticsSink receives it, and no public
// form of the report contains it: not the JSON response, not the accessors, and not generic formatting.
func TestRunKeepsValidatorErrorsOutOfThePublicReport(t *testing.T) {
	const secret = "ghp_FakeSecret0123456789"
	cause := errors.New("broken at /home/user/project/.env with token " + secret)
	runner, err := NewValidationRunner(fakeReader{selection: []string{"failing"}}, fakeValidator{name: "failing", err: cause})
	if err != nil {
		t.Fatal(err)
	}
	sink := recordingSink{}

	report, err := runner.WithDiagnosticsSink(sink).Run(context.Background(), global.NewTargetProject("project"))

	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	failing, _ := report.Result("failing")
	for form, text := range map[string]string{
		"JSON":      string(encoded),
		"%v":        fmt.Sprintf("%v", report),
		"%+v":       fmt.Sprintf("%+v", report),
		"%#v":       fmt.Sprintf("%#v", report),
		"accessors": failing.Validation() + failing.Code() + failing.Message() + string(failing.Status()),
	} {
		if strings.Contains(text, secret[:4]) || strings.Contains(text, ".env") {
			t.Errorf("%s contains the validator error: %s", form, text)
		}
	}
	if sink["failing"] != cause.Error() {
		t.Errorf("expected the validator error in the DiagnosticsSink, got %q", sink["failing"])
	}
}

// The DiagnosticsSink receives at most MaximumDiagnosticsSize bytes of valid UTF-8 from the start of the error.
func TestRunRecordsTheFirstBytesOfTheDiagnostics(t *testing.T) {
	long := strings.Repeat("ä", MaximumDiagnosticsSize)
	runner, err := NewValidationRunner(fakeReader{selection: []string{"failing"}}, fakeValidator{name: "failing", err: errors.New(long)})
	if err != nil {
		t.Fatal(err)
	}
	sink := recordingSink{}

	if _, err := runner.WithDiagnosticsSink(sink).Run(context.Background(), global.NewTargetProject("project")); err != nil {
		t.Fatal(err)
	}

	diagnostics := sink["failing"]
	if len(diagnostics) > MaximumDiagnosticsSize || !utf8.ValidString(diagnostics) || !strings.HasPrefix(long, diagnostics) || diagnostics == "" {
		t.Fatalf("expected at most %d bytes of valid UTF-8 from the start, got %d bytes", MaximumDiagnosticsSize, len(diagnostics))
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

// Each validator has a unique, non-empty name without leading or trailing whitespace, so each selected name runs
// exactly the wired validator.
func TestNewValidationRunnerRejectsInvalidAndDuplicateNames(t *testing.T) {
	for name, test := range map[string]struct {
		validators []fakeValidator
		want       string
	}{
		"empty name":               {[]fakeValidator{{name: ""}}, "E-CAF-7"},
		"whitespace only":          {[]fakeValidator{{name: "  "}}, "E-CAF-7"},
		"leading whitespace":       {[]fakeValidator{{name: " check"}}, "E-CAF-7"},
		"trailing whitespace":      {[]fakeValidator{{name: "check\t"}}, "E-CAF-7"},
		"duplicate name":           {[]fakeValidator{{name: "check", status: domain.Pass}, {name: "check", status: domain.Finding}}, "E-CAF-8"},
		"duplicate after trimming": {[]fakeValidator{{name: "check", status: domain.Pass}, {name: " check ", status: domain.Finding}}, "E-CAF-7"},
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
	sink := recordingSink{}

	report, err := runner.WithDiagnosticsSink(sink).Run(context.Background(), global.NewTargetProject("project"))

	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"check: tool_failure " + MismatchingValidationCode}; !reflect.DeepEqual(summary(report), want) {
		t.Fatalf("got %v, want %v", summary(report), want)
	}
	if !strings.Contains(sink["check"], "other") {
		t.Errorf("expected the returned name in the diagnostics, got %q", sink["check"])
	}
}

// cancelingValidator cancels the context during Validate and returns the cancellation.
type cancelingValidator struct {
	cancel context.CancelFunc
}

func (cancelingValidator) Name() string {
	return "canceling"
}

func (validator cancelingValidator) Validate(ctx context.Context, _ global.TargetProject) (domain.ValidationResult, error) {
	validator.cancel()
	return domain.ValidationResult{}, ctx.Err()
}

// A cancellation during a validator stops the run and returns the cancellation, also for the last selected validator.
func TestRunReturnsACancellationDuringAValidator(t *testing.T) {
	for name, selection := range map[string][]string{
		"last validator":  {"canceling"},
		"inner validator": {"canceling", "next"},
	} {
		ctx, cancel := context.WithCancel(context.Background())
		next := &recordingValidator{name: "next"}
		runner, err := NewValidationRunner(fakeReader{selection: selection}, cancelingValidator{cancel: cancel}, next)
		if err != nil {
			t.Fatal(err)
		}

		_, err = runner.Run(ctx, global.NewTargetProject("project"))

		if !errors.Is(err, context.Canceled) {
			t.Errorf("%s: expected the cancellation, got %v", name, err)
		}
		if next.called {
			t.Errorf("%s: expected no run of the next validator", name)
		}
		cancel()
	}
}

// recordingValidator passes and records that it ran.
type recordingValidator struct {
	name   string
	called bool
}

func (validator *recordingValidator) Name() string {
	return validator.name
}

func (validator *recordingValidator) Validate(context.Context, global.TargetProject) (domain.ValidationResult, error) {
	validator.called = true
	return domain.NewValidationResult(validator.name, domain.Pass, "CODE", "message"), nil
}

// The message of a failed validator names the E-CAF code of its error, and the DiagnosticsSink receives the error text.
func TestRunNamesTheErrorCodeOfAFailedValidator(t *testing.T) {
	for name, test := range map[string]struct {
		err  error
		want string
	}{
		"error code":         {exaerror.New("E-CAF-5").Message("Failed to read /home/user/project/secret.md."), "The validator failed with E-CAF-5."},
		"wrapped error code": {fmt.Errorf("trace: %w", exaerror.New("E-CAF-9").Message("Paths differ.")), "The validator failed with E-CAF-9."},
		"no error code":      {errors.New("broken at /home/user/project/secret.md"), "The validator failed."},
	} {
		runner, err := NewValidationRunner(fakeReader{selection: []string{"failing"}}, fakeValidator{name: "failing", err: test.err})
		if err != nil {
			t.Fatal(err)
		}
		sink := recordingSink{}

		report, err := runner.WithDiagnosticsSink(sink).Run(context.Background(), global.NewTargetProject("project"))

		failing, _ := report.Result("failing")
		if err != nil || failing.Message() != test.want || failing.Code() != FailedValidationCode || sink["failing"] != test.err.Error() {
			t.Errorf("%s: got %q %q, diagnostics %q, %v; want %q", name, failing.Code(), failing.Message(), sink["failing"], err, test.want)
		}
	}
}
