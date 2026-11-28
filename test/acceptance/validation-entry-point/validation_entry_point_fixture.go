// Package validationentrypoint holds the shared fixture of the Validation Entry Point use cases.
package validationentrypoint

import (
	"context"
	"fmt"
	"strings"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
)

// Configuration creates the .caf.yaml file that selects the validations in order.
func Configuration(validations ...string) acceptance.File {
	return acceptance.File{Path: ".caf.yaml", Content: "version: 1\nvalidators:\n  - " + strings.Join(validations, "\n  - ") + "\n"}
}

// OFT Feature and OFT Requirement documents.
const (
	featureTemplate     = "\n# OFT Feature Document\n\n%s\n\nNeeds: req\n"
	requirementTemplate = "\n# OFT Requirement Document\n\n%s\n\nCovers:\n\n- %s\n"
)

// OFTTracing creates an OFT Feature that needs an OFT Requirement and, for valid tracing, the OFT Requirement that covers it.
func OFTTracing(valid bool) []acceptance.File {
	files := []acceptance.File{acceptance.NewFile("features.md", featureTemplate, "feat~fixture-feature~1")}
	if valid {
		files = append(files, acceptance.NewFile("requirements.md", requirementTemplate, "req~fixture-requirement~1", "feat~fixture-feature~1"))
	}
	return files
}

// Project creates .caf.yaml with the validations and OFT documents with valid or invalid tracing.
func Project(validTracing bool, validations ...string) []acceptance.File {
	return append([]acceptance.File{Configuration(validations...)}, OFTTracing(validTracing)...)
}

// Run runs the Validation Runner for a target-project fixture with one or more files.
func Run(t *testing.T, files []acceptance.File, then func(codingagentflow.ValidationReport) error) {
	t.Helper()
	runner, err := codingagentflow.NewValidationRunner()
	if err != nil {
		t.Fatal(err)
	}
	acceptance.Run(t, files, acceptance.TargetProject, runner.Run, then)
}

// Stop runs the Agent Turn Validator for a target-project fixture with a Stop Event that has the Continuation Active flag.
func Stop(t *testing.T, files []acceptance.File, continuationActive bool, then func(codingagentflow.TurnDecision) error) {
	t.Helper()
	validator, err := codingagentflow.NewAgentTurnValidator()
	if err != nil {
		t.Fatal(err)
	}
	stop := func(ctx context.Context, project codingagentflow.TargetProject) (codingagentflow.TurnDecision, error) {
		return validator.Validate(ctx, codingagentflow.NewStopEvent(project, continuationActive))
	}
	acceptance.Run(t, files, acceptance.TargetProject, stop, then)
}

// ExpectDecision checks the action of a TurnDecision and that its reason contains the text.
func ExpectDecision(action codingagentflow.TurnAction, reason string) func(codingagentflow.TurnDecision) error {
	return func(decision codingagentflow.TurnDecision) error {
		if decision.Action() != action || !strings.Contains(decision.Reason(), reason) {
			return fmt.Errorf("decision: got %s %q, want %s with %q", decision.Action(), decision.Reason(), action, reason)
		}
		return nil
	}
}

// Summary lists the results of a report as "validation: status".
func Summary(report codingagentflow.ValidationReport) []string {
	var lines []string
	for _, result := range report.Results() {
		lines = append(lines, fmt.Sprintf("%s: %s", result.Validation(), result.Status()))
	}
	return lines
}
