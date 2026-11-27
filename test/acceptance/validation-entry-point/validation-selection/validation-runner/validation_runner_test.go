package validationrunner

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	validationentrypoint "github.com/exasol/coding-agent-flow/test/acceptance/validation-entry-point"
)

// project creates .caf.yaml with the validations and OFT documents with valid or invalid tracing.
func project(validTracing bool, validations ...string) []acceptance.File {
	return append([]acceptance.File{validationentrypoint.Configuration(validations...)}, validationentrypoint.OFTTracing(validTracing)...)
}

func expectResults(want ...string) func(codingagentflow.ValidationReport) error {
	return func(report codingagentflow.ValidationReport) error {
		if got := validationentrypoint.Summary(report); !slices.Equal(got, want) {
			return fmt.Errorf("results: got %v, want %v", got, want)
		}
		return nil
	}
}

func expectStatus(want codingagentflow.ValidationStatus) func(codingagentflow.ValidationReport) error {
	return func(report codingagentflow.ValidationReport) error {
		if report.Status() != want {
			return fmt.Errorf("status: got %s, want %s with results %v", report.Status(), want, validationentrypoint.Summary(report))
		}
		return nil
	}
}

// [at -> scn~validation-runner_selected-validations~1]
func TestValidationRunner_selectedValidations(t *testing.T) {
	validationentrypoint.Run(t, project(true, "oft-tracing-validator"), expectResults("oft-tracing-validator: pass"))
	validationentrypoint.Run(t, project(false, "oft-tracing-validator"), expectResults("oft-tracing-validator: finding"))
}

// [at -> scn~validation-runner_unknown-validation~1]
func TestValidationRunner_unknownValidation(t *testing.T) {
	validationentrypoint.Run(t, project(true, "oft-tracing-validator", "unknown-check"),
		func(report codingagentflow.ValidationReport) error {
			if err := expectResults("oft-tracing-validator: pass", "unknown-check: tool_failure")(report); err != nil {
				return err
			}
			if unknown, _ := report.Result("unknown-check"); unknown.Code() != "CAF-UNKNOWN-VALIDATOR" {
				return fmt.Errorf("code of unknown-check: got %q, want CAF-UNKNOWN-VALIDATOR", unknown.Code())
			}
			return nil
		},
	)
}

// [at -> scn~validation-runner_failed-validation~1]
func TestValidationRunner_failedValidation(t *testing.T) {
	// OpenFastTrace cannot run with a configured JAR that does not exist.
	t.Setenv("OPENFASTTRACE_JAR", filepath.Join(t.TempDir(), "missing.jar"))

	validationentrypoint.Run(t, project(true, "oft-tracing-validator"),
		func(report codingagentflow.ValidationReport) error {
			if err := expectResults("oft-tracing-validator: tool_failure")(report); err != nil {
				return err
			}
			failed, _ := report.Result("oft-tracing-validator")
			if failed.Code() != "CAF-VALIDATION-FAILED" || !strings.Contains(failed.Message(), "E-CAF-12") {
				return fmt.Errorf("result of oft-tracing-validator: got %q %q, want CAF-VALIDATION-FAILED with E-CAF-12", failed.Code(), failed.Message())
			}
			// The report holds no details of the cause, such as the path of the JAR, so a caller can print and log it.
			if printed := fmt.Sprintf("%+v", report); strings.Contains(printed, "missing.jar") {
				return fmt.Errorf("the report contains the cause: %s", printed)
			}
			return nil
		},
	)
}

// [at -> scn~validation-runner_overall-result~1]
func TestValidationRunner_overallResult(t *testing.T) {
	validationentrypoint.Run(t, project(true, "oft-tracing-validator"), expectStatus(codingagentflow.Pass))
	validationentrypoint.Run(t, project(false, "oft-tracing-validator"), expectStatus(codingagentflow.Finding))
	validationentrypoint.Run(t, project(false, "oft-tracing-validator", "unknown-check"), expectStatus(codingagentflow.ToolFailure))
}
