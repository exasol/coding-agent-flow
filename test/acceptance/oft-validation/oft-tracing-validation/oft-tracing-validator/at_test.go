//go:build oftvalidation

package ofttracingvalidator

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

var scenarios = acceptance.NewSuite(
	oftvalidation.When,
	[]acceptance.Scenario{
		{
			ID:          "scn~oft-tracing-validator~1",
			Description: "Return valid when a requirement covers a feature",
			Files: []acceptance.File{
				{
					Path: "features.md",
					Content: `
# Valid Fixture Feature

feat~valid-fixture-feature~1

Needs: req
`,
				},
				{
					Path: "requirements.md",
					Content: `
# Valid Fixture Requirement

req~valid-fixture-requirement~1

Covers:

- feat~valid-fixture-feature~1
`,
				},
			},
			Then: acceptance.Values{"valid": "true"},
		},
		{
			ID:          "scn~oft-tracing-validator~1",
			Description: "Return invalid when a feature stays uncovered",
			Files: []acceptance.File{{
				Path: "features.md",
				Content: `
# Invalid Fixture Feature

feat~invalid-fixture-feature~1

Needs: req
`,
			}},
			Then: acceptance.Values{"valid": "false"},
		},
	},
	checkValid,
)

func checkValid(response codingagentflow.OFTValidation, _ acceptance.Values, then acceptance.Values) error {
	return acceptance.Expect("valid", then.Bool("valid"), response.Valid())
}

// [at -> scn~oft-tracing-validator~1]
func TestOFTTracingValidator(t *testing.T) {
	scenarios.Run(t)
}
