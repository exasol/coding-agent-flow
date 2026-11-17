// Package oftvalidation holds the shared When helper of the OFT validation use cases.
package oftvalidation

import (
	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
)

// When creates the When binding of an OFT validation use case. The call can use targetProject, made of the
// run's input files.
func When(step, call string) acceptance.When[codingagentflow.OFTValidation] {
	return acceptance.When[codingagentflow.OFTValidation]{
		Step:    step,
		Call:    call,
		Imports: []string{`codingagentflow "github.com/exasol/coding-agent-flow"`},
		Inputs:  []acceptance.Input{{Name: "targetProject", Expression: "codingagentflow.NewTargetProject(files.Directory())"}},
		Decode:  acceptance.DecodeJSON[codingagentflow.OFTValidation],
	}
}
