//go:build oftvalidation

package oftrequirementsdetector

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

type requirementsDetectionScenario struct {
	name          string
	files         map[string]string
	requirementID string
	found         bool
}

// [at -> scn~oft-requirements-detection~1]
func TestOFTRequirementsDetector(t *testing.T) {
	// Disabled until CAF provides NewTargetProject and NewOFTTracingValidator.
	for _, scenario := range []requirementsDetectionScenario{
		{
			name: "detects a requirement from doc/requirements/requirements.rst",
			files: map[string]string{
				"doc/requirements/requirements.rst": "OFT Requirement document\n========================\n\nreq~valid-doc-requirements-rst~1\n",
			},
			requirementID: "req~valid-doc-requirements-rst~1",
			found:         true,
		},
		{
			name: "detects a requirement from requirements.rst",
			files: map[string]string{
				"requirements.rst": "OFT Requirement document\n========================\n\nreq~valid-root-rst~1\n",
			},
			requirementID: "req~valid-root-rst~1",
			found:         true,
		},
		{
			name: "detects a requirement from requirements.md",
			files: map[string]string{
				"requirements.md": "OFT Requirement document\n========================\n\nreq~valid-root-markdown~1\n",
			},
			requirementID: "req~valid-root-markdown~1",
			found:         true,
		},
		{
			name: "does not detect a requirement from doc/requirements/requirements.txt",
			files: map[string]string{
				"doc/requirements/requirements.txt": "OFT Requirement document\n========================\n\nreq~invalid-doc-requirements-text~1\n",
			},
			requirementID: "req~invalid-doc-requirements-text~1",
			found:         false,
		},
		{
			name: "does not detect a requirement from requirements.txt",
			files: map[string]string{
				"requirements.txt": "OFT Requirement document\n========================\n\nreq~invalid-root-text~1\n",
			},
			requirementID: "req~invalid-root-text~1",
			found:         false,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			project := codingagentflow.NewTargetProject(oftvalidation.TargetProject(t, scenario.files))
			response, err := codingagentflow.NewOFTTracingValidator().Validate(t.Context(), project)
			if err != nil {
				t.Fatal(err)
			}
			id, err := codingagentflow.ParseOFTSpecificationItemID(scenario.requirementID)
			if err != nil {
				t.Fatal(err)
			}
			if got := response.Contains(id); got != scenario.found {
				t.Errorf("%s found: got %t, want %t", scenario.requirementID, got, scenario.found)
			}
		})
	}
}
