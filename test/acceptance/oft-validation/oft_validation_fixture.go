// Package oftvalidation holds the shared When helper of the OFT validation use cases.
package oftvalidation

import (
	"slices"
	"strings"

	"github.com/exasol/coding-agent-flow/test/acceptance"
)

// OFTValidationResponse is the documented JSON response of OFT validation.
type OFTValidationResponse struct {
	Valid                bool                `json:"valid"`
	SpecificationItemIDs map[string][]string `json:"oftSpecificationItemIDs"`
}

// HasSpecificationItemID reports whether the response contains id.
func (response OFTValidationResponse) HasSpecificationItemID(id string) bool {
	artifactType, _, found := strings.Cut(id, "~")
	return found && slices.Contains(response.SpecificationItemIDs[artifactType], id)
}

// When creates the When binding of an OFT validation use case. The call can use targetProject, made of the
// run's input files.
func When(step, call string) acceptance.When[OFTValidationResponse] {
	return acceptance.When[OFTValidationResponse]{
		Step:    step,
		Call:    call,
		Imports: []string{`codingagentflow "github.com/exasol/coding-agent-flow"`},
		Inputs:  []acceptance.Input{{Name: "targetProject", Expression: "codingagentflow.NewTargetProject(files.Directory())"}},
		Decode:  acceptance.DecodeJSON[OFTValidationResponse],
	}
}
