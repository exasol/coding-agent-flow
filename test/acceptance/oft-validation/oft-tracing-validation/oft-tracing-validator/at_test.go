//go:build oftvalidation

package ofttracingvalidator

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

type tracingValidationScenario struct {
	name  string
	files map[string]string
	valid bool
}

// [at -> scn~oft-tracing-validator~1]
func TestOFTTracingValidator(t *testing.T) {
	// Disabled until CAF provides NewTargetProject and NewOFTTracingValidator.
	for _, scenario := range []tracingValidationScenario{
		{
			name: "returns valid when a requirement covers a feature",
			files: map[string]string{
				"features.md":     "# Valid Fixture Feature\n\nfeat~valid-fixture-feature~1\n\nNeeds: req\n",
				"requirements.md": "# Valid Fixture Requirement\n\nreq~valid-fixture-requirement~1\n\nCovers:\n\n- feat~valid-fixture-feature~1\n",
			},
			valid: true,
		},
		{
			name: "returns invalid when a feature stays uncovered",
			files: map[string]string{
				"features.md": "# Invalid Fixture Feature\n\nfeat~invalid-fixture-feature~1\n\nNeeds: req\n",
			},
			valid: false,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			project := codingagentflow.NewTargetProject(oftvalidation.TargetProject(t, scenario.files))
			response, err := codingagentflow.NewOFTTracingValidator().Validate(t.Context(), project)
			if err != nil {
				t.Fatal(err)
			}
			if got := response.Valid(); got != scenario.valid {
				t.Errorf("valid: got %t, want %t", got, scenario.valid)
			}
		})
	}
}
