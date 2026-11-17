//go:build oftvalidation

// OFT acceptance tests are disabled until CAF provides NewTargetProject and NewOFTTracingValidator.
package oftvalidation

import (
	"fmt"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
)

type File struct {
	Path, Template string
	Parameters     []string
}

func Files(files ...File) []acceptance.File {
	result := make([]acceptance.File, len(files))
	for index, file := range files {
		result[index] = acceptance.File{Path: file.Path, Content: fmt.Sprintf(file.Template, file.Parameters)}
	}
	return result
}

func TargetProject(t *testing.T, files []acceptance.File) codingagentflow.TargetProject {
	return codingagentflow.NewTargetProject(acceptance.WriteFiles(t, files))
}

// Validate runs OFTTracingValidator.Validate for one file fixture.
func Validate(t *testing.T, files []acceptance.File, then func(codingagentflow.OFTValidation) error) {
	t.Helper()
	validator := codingagentflow.NewOFTTracingValidator()
	acceptance.Run(t, files, TargetProject, validator.Validate, then)
}
