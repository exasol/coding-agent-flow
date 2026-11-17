//go:build oftvalidation

package oftrequirementsdetector

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
			ID:          "scn~oft-requirements-detection~1",
			Description: "Detect a requirement from doc/requirements/requirements.rst",
			Given:       acceptance.Values{"req_id": "req~valid-doc-requirements-rst~1"},
			Files: []acceptance.File{{
				Path: "doc/requirements/requirements.rst",
				Content: `
OFT Requirement document
========================

req~valid-doc-requirements-rst~1
`,
			}},
			Then: acceptance.Values{"found": "true"},
		},
		{
			ID:          "scn~oft-requirements-detection~1",
			Description: "Detect a requirement from requirements.rst",
			Given:       acceptance.Values{"req_id": "req~valid-root-rst~1"},
			Files: []acceptance.File{{
				Path: "requirements.rst",
				Content: `
OFT Requirement document
========================

req~valid-root-rst~1
`,
			}},
			Then: acceptance.Values{"found": "true"},
		},
		{
			ID:          "scn~oft-requirements-detection~1",
			Description: "Detect a requirement from requirements.md",
			Given:       acceptance.Values{"req_id": "req~valid-root-markdown~1"},
			Files: []acceptance.File{{
				Path: "requirements.md",
				Content: `
OFT Requirement document
========================

req~valid-root-markdown~1
`,
			}},
			Then: acceptance.Values{"found": "true"},
		},
		{
			ID:          "scn~oft-requirements-detection~1",
			Description: "Do not detect a requirement from doc/requirements/requirements.txt",
			Given:       acceptance.Values{"req_id": "req~invalid-doc-requirements-text~1"},
			Files: []acceptance.File{{
				Path: "doc/requirements/requirements.txt",
				Content: `
OFT Requirement document
========================

req~invalid-doc-requirements-text~1
`,
			}},
			Then: acceptance.Values{"found": "false"},
		},
		{
			ID:          "scn~oft-requirements-detection~1",
			Description: "Do not detect a requirement from requirements.txt",
			Given:       acceptance.Values{"req_id": "req~invalid-root-text~1"},
			Files: []acceptance.File{{
				Path: "requirements.txt",
				Content: `
OFT Requirement document
========================

req~invalid-root-text~1
`,
			}},
			Then: acceptance.Values{"found": "false"},
		},
	},
	checkRequirementFound,
)

func checkRequirementFound(response codingagentflow.OFTValidation, given acceptance.Values, then acceptance.Values) error {
	requirementID := given.String("req_id")
	id, err := codingagentflow.ParseOFTSpecificationItemID(requirementID)
	if err != nil {
		return err
	}
	return acceptance.Expect(requirementID+" found", then.Bool("found"), response.Contains(id))
}

// [at -> scn~oft-requirements-detection~1]
func TestOFTRequirementsDetector(t *testing.T) {
	scenarios.Run(t)
}
