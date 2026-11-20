// Package oftvalidation holds the shared fixture of the OFT Validation use cases.
package oftvalidation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
)

func File(path, template string, parameters ...any) acceptance.File {
	return acceptance.File{Path: path, Content: fmt.Sprintf(template, parameters...)}
}

// Given creates the Target Project of a scenario from the fixture files.
type Given func(*testing.T, []acceptance.File) codingagentflow.TargetProject

// MissingPath is a Target Project path that doesn't exist.
func MissingPath(t *testing.T, _ []acceptance.File) codingagentflow.TargetProject {
	return codingagentflow.NewTargetProject(filepath.Join(t.TempDir(), "missing"))
}

// RegularFile is a Target Project path that is a regular file.
func RegularFile(t *testing.T, _ []acceptance.File) codingagentflow.TargetProject {
	path := filepath.Join(t.TempDir(), "project.md")
	if err := os.WriteFile(path, []byte("# Project\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return codingagentflow.NewTargetProject(path)
}

// SymbolicLink is a Target Project path that is a symbolic link to a directory with the fixture files.
func SymbolicLink(t *testing.T, files []acceptance.File) codingagentflow.TargetProject {
	link := filepath.Join(t.TempDir(), "project")
	if err := os.Symlink(acceptance.WriteFiles(t, files), link); err != nil {
		t.Fatal(err)
	}
	return codingagentflow.NewTargetProject(link)
}

// Validate runs OFTTracingValidator.Validate for a target-project fixture with one or more files.
func Validate(t *testing.T, files []acceptance.File, then func(codingagentflow.OFTValidation) error) {
	t.Helper()
	ValidateIn(t, acceptance.TargetProject, files, then)
}

// ValidateIn runs OFTTracingValidator.Validate for the Target Project that given creates.
func ValidateIn(t *testing.T, given Given, files []acceptance.File, then func(codingagentflow.OFTValidation) error) {
	t.Helper()
	acceptance.Run(t, files, given, codingagentflow.NewOFTTracingValidator().Validate, then)
}

// ValidateFails runs OFTTracingValidator.Validate for the Target Project that given creates and checks the error.
func ValidateFails(t *testing.T, given Given, files []acceptance.File, then func(error) error) {
	t.Helper()
	acceptance.RunFails(t, files, given, codingagentflow.NewOFTTracingValidator().Validate, then)
}

// ValidateTargetProject runs OFTTracingValidator.Validate and passes the Target Project to then.
func ValidateTargetProject(t *testing.T, files []acceptance.File, then func(codingagentflow.OFTValidation, codingagentflow.TargetProject) error) {
	t.Helper()
	var targetProject codingagentflow.TargetProject
	given := func(t *testing.T, files []acceptance.File) codingagentflow.TargetProject {
		targetProject = acceptance.TargetProject(t, files)
		return targetProject
	}
	ValidateIn(t, given, files, func(response codingagentflow.OFTValidation) error {
		return then(response, targetProject)
	})
}
