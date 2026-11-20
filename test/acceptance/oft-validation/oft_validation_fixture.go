// Package oftvalidation holds the shared fixture of the OFT Validation use cases.
package oftvalidation

import (
	"fmt"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
)

func File(path, template string, parameters ...any) acceptance.File {
	return acceptance.File{Path: path, Content: fmt.Sprintf(template, parameters...)}
}

// Validate runs OFTTracingValidator.Validate for a target-project fixture with one or more files.
func Validate(t *testing.T, files []acceptance.File, then func(codingagentflow.OFTValidation) error) {
	t.Helper()
	validator := codingagentflow.NewOFTTracingValidator()
	acceptance.Run(t, files, acceptance.TargetProject, validator.Validate, then)
}
