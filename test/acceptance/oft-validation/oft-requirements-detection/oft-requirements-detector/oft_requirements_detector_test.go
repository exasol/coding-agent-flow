package oftrequirementsdetector

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

const requirementTemplate = `
OFT Requirement document
========================

%s
`

func requirementFiles(path, requirementID string) []acceptance.File {
	return []acceptance.File{oftvalidation.File(path, requirementTemplate, requirementID)}
}

func expectRequirementFound(requirementID string, found bool) func(codingagentflow.OFTValidation) error {
	return func(response codingagentflow.OFTValidation) error {
		id, err := codingagentflow.ParseOFTSpecificationItemID(requirementID)
		if err != nil {
			return err
		}
		return acceptance.Expect(requirementID+" found", found, response.Contains(id))
	}
}

// [at -> scn~oft-requirements-detector~1]
func TestOFTRequirementsDetector(t *testing.T) {
	oftvalidation.Validate(t,
		requirementFiles("doc/requirements/requirements.rst", "req~valid-doc-requirements-rst~1"),
		expectRequirementFound("req~valid-doc-requirements-rst~1", true),
	)

	oftvalidation.Validate(t,
		requirementFiles("requirements.rst", "req~valid-root-rst~1"),
		expectRequirementFound("req~valid-root-rst~1", true),
	)

	oftvalidation.Validate(t,
		requirementFiles("requirements.md", "req~valid-root-markdown~1"),
		expectRequirementFound("req~valid-root-markdown~1", true),
	)

	oftvalidation.Validate(t,
		requirementFiles("doc/requirements/requirements.txt", "req~invalid-doc-requirements-text~1"),
		expectRequirementFound("req~invalid-doc-requirements-text~1", false),
	)

	oftvalidation.Validate(t,
		requirementFiles("requirements.txt", "req~invalid-root-text~1"),
		expectRequirementFound("req~invalid-root-text~1", false),
	)
}
