// Package validationentrypoint holds the shared fixture of the Validation Entry Point use cases.
package validationentrypoint

import (
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
	files := []acceptance.File{{Path: "features.md", Content: fmt.Sprintf(featureTemplate, "feat~fixture-feature~1")}}
	if valid {
		files = append(files, acceptance.File{Path: "requirements.md", Content: fmt.Sprintf(requirementTemplate, "req~fixture-requirement~1", "feat~fixture-feature~1")})
	}
	return files
}

func TargetProject(t *testing.T, files []acceptance.File) codingagentflow.TargetProject {
	return codingagentflow.NewTargetProject(acceptance.WriteFiles(t, files))
}

// Run runs the Validation Runner for a target-project fixture with one or more files.
func Run(t *testing.T, files []acceptance.File, then func(codingagentflow.ValidationReport) error) {
	t.Helper()
	runner := codingagentflow.NewValidationRunner()
	acceptance.Run(t, files, TargetProject, runner.Run, then)
}

// Summary lists the results of a report as "validation: status".
func Summary(report codingagentflow.ValidationReport) []string {
	var lines []string
	for _, result := range report.Results() {
		lines = append(lines, fmt.Sprintf("%s: %s", result.Validation(), result.Status()))
	}
	return lines
}
