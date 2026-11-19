package linkchecker

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	codevalidation "github.com/exasol/coding-agent-flow/test/acceptance/code-validation"
)

const readmeTemplate = `
# Document

## Setup
`

const guideTemplate = `
# Guide

## Usage
`

// lycheeConfiguration would hide the Broken Links of README.md and write a lychee cache into the Target Project.
const lycheeConfiguration = `cache = true
exclude_path = ["README.md"]
`

// localProject creates README.md with the heading Setup, doc/guide.md with the heading Usage, and source with the Link.
func localProject(source, link string, extraFiles ...acceptance.File) []acceptance.File {
	files := []acceptance.File{acceptance.NewFile("doc/guide.md", guideTemplate), codevalidation.Document(source, link)}
	if source != "README.md" {
		files = append(files, acceptance.NewFile("README.md", readmeTemplate))
	}
	return append(files, extraFiles...)
}

func expectUnchangedResult(valid bool, source, brokenLink string, files []acceptance.File) func(codingagentflow.LinkValidation, codingagentflow.TargetProject) error {
	return func(response codingagentflow.LinkValidation, targetProject codingagentflow.TargetProject) error {
		if err := codevalidation.ExpectResult(valid, source, brokenLink)(response); err != nil {
			return err
		}
		return acceptance.ExpectUnchanged(targetProject, files)
	}
}

// [at -> scn~link-checker_local-links~1]
func TestLinkChecker_localLinks(t *testing.T) {
	checker := codingagentflow.NewLinkChecker()
	for _, row := range []struct {
		source, link  string
		configuration []acceptance.File
		valid         bool
		brokenLink    string
	}{
		{"README.md", "doc/guide.md", nil, true, codevalidation.NoBrokenLink},
		{"doc/a.md", "../README.md", nil, true, codevalidation.NoBrokenLink},
		{"README.md", "#setup", nil, true, codevalidation.NoBrokenLink},
		{"README.md", "doc/guide.md#usage", nil, true, codevalidation.NoBrokenLink},
		{"README.md", "doc/missing.md", nil, false, "doc/missing.md"},
		{"README.md", "#nope", nil, false, "README.md#nope"},
		{"README.md", "doc/guide.md#nope", nil, false, "doc/guide.md#nope"},
		{"README.md", "doc/missing.md", []acceptance.File{acceptance.NewFile("lychee.toml", lycheeConfiguration)}, false, "doc/missing.md"},
	} {
		files := localProject(row.source, row.link, row.configuration...)
		codevalidation.CheckTargetProject(t, checker, files, expectUnchangedResult(row.valid, row.source, row.brokenLink, files))
	}
}
