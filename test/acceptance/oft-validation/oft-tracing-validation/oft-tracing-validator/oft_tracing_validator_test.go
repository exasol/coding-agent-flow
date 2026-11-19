//go:build oftvalidation

package ofttracingvalidator

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

const featureTemplate = `
# OFT Feature Document

%s

Needs: req
`

const requirementTemplate = `
# OFT Requirement Document

%s

Covers:

- %s
`

func featureFile(featureID string) oftvalidation.File {
	return oftvalidation.File{Path: "features.md", Template: featureTemplate, Parameters: []string{featureID}}
}

func requirementFile(requirementID, featureID string) oftvalidation.File {
	return oftvalidation.File{Path: "requirements.md", Template: requirementTemplate, Parameters: []string{requirementID, featureID}}
}

func expectValid(valid bool) func(codingagentflow.OFTValidation) error {
	return func(response codingagentflow.OFTValidation) error {
		return acceptance.Expect("valid", valid, response.Valid())
	}
}

// [at -> scn~oft-tracing-validator_valid~1]
func TestOFTTracingValidator_valid(t *testing.T) {
	oftvalidation.Validate(t, oftvalidation.Files(
		featureFile("feat~valid-fixture-feature~1"),
		requirementFile("req~valid-fixture-requirement~1", "feat~valid-fixture-feature~1"),
	),
		expectValid(true),
	)
}

// [at -> scn~oft-tracing-validator_invalid~1]
func TestOFTTracingValidator_invalid(t *testing.T) {
	oftvalidation.Validate(t,
		oftvalidation.Files(featureFile("feat~invalid-fixture-feature~1")),
		expectValid(false),
	)
}

// [at -> scn~oft-tracing-validator_requirement-covers-other-feature~1]
func TestOFTTracingValidator_requirementCoversOtherFeature(t *testing.T) {
	oftvalidation.Validate(t, oftvalidation.Files(
		featureFile("feat~needed-fixture-feature~1"),
		requirementFile("req~other-feature-fixture-requirement~1", "feat~other-fixture-feature~1"),
	),
		expectValid(false),
	)
}
