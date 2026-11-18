package linkchecker

import (
	"fmt"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	codevalidation "github.com/exasol/coding-agent-flow/test/acceptance/code-validation"
)

const readmeTemplate = `
# Document

## Setup
`

const markdownTemplate = `
# Document

## Setup

[link](%s)
`

const guideTemplate = `
# Guide

## Usage
`

const restructuredTextTemplate = "\nDocument\n========\n\nSee %s for details.\n"

const textTemplate = "See %s for details.\n"

const sitemapTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>%s</loc></url>
</urlset>
`

// lycheeConfiguration would hide the Broken Links of README.md and write a lychee cache into the Target Project.
const lycheeConfiguration = `cache = true
exclude_path = ["README.md"]
`

// noBrokenLink marks a valid expected result.
const noBrokenLink = ""

// document creates the file at path with the Link, using the template that matches its extension.
func document(path, link string) acceptance.File {
	templates := map[string]string{
		"README.md":   markdownTemplate,
		"doc/a.md":    markdownTemplate,
		"doc/a.rst":   restructuredTextTemplate,
		"notes.txt":   textTemplate,
		"sitemap.xml": sitemapTemplate,
	}
	return acceptance.NewFile(path, templates[path], link)
}

// localProject creates README.md with the heading Setup, doc/guide.md with the heading Usage, and source with the Link.
func localProject(source, link string, extraFiles ...acceptance.File) []acceptance.File {
	files := []acceptance.File{acceptance.NewFile("doc/guide.md", guideTemplate), document(source, link)}
	if source != "README.md" {
		files = append(files, acceptance.NewFile("README.md", readmeTemplate))
	}
	return append(files, extraFiles...)
}

// expectResult checks valid and, if brokenLink is set, that the Broken Link is reported in source.
func expectResult(valid bool, source, brokenLink string) func(codingagentflow.LinkValidation) error {
	return func(response codingagentflow.LinkValidation) error {
		if err := acceptance.Expect("valid", valid, response.Valid()); err != nil {
			return fmt.Errorf("%w: Broken Links %+v", err, response.BrokenLinks())
		}
		if brokenLink != noBrokenLink && !response.ContainsBrokenLink(source, brokenLink) {
			return fmt.Errorf("Broken Link %s in %s: not reported in %+v", brokenLink, source, response.BrokenLinks())
		}
		return nil
	}
}

func expectUnchangedResult(valid bool, source, brokenLink string, files []acceptance.File) func(codingagentflow.LinkValidation, codingagentflow.TargetProject) error {
	return func(response codingagentflow.LinkValidation, targetProject codingagentflow.TargetProject) error {
		if err := expectResult(valid, source, brokenLink)(response); err != nil {
			return err
		}
		return acceptance.ExpectUnchanged(targetProject, files)
	}
}

func expectNoRequest(webServer *codevalidation.WebServer) func(codingagentflow.LinkValidation) error {
	return func(response codingagentflow.LinkValidation) error {
		if err := expectResult(true, "", noBrokenLink)(response); err != nil {
			return err
		}
		if requests := webServer.Requests(); requests != 0 {
			return fmt.Errorf("web server received %d requests, expected none", requests)
		}
		return nil
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
		{"README.md", "doc/guide.md", nil, true, noBrokenLink},
		{"doc/a.md", "../README.md", nil, true, noBrokenLink},
		{"README.md", "#setup", nil, true, noBrokenLink},
		{"README.md", "doc/guide.md#usage", nil, true, noBrokenLink},
		{"README.md", "doc/missing.md", nil, false, "doc/missing.md"},
		{"README.md", "#nope", nil, false, "README.md#nope"},
		{"README.md", "doc/guide.md#nope", nil, false, "doc/guide.md#nope"},
		{"README.md", "doc/missing.md", []acceptance.File{acceptance.NewFile("lychee.toml", lycheeConfiguration)}, false, "doc/missing.md"},
	} {
		files := localProject(row.source, row.link, row.configuration...)
		codevalidation.CheckTargetProject(t, checker, files, expectUnchangedResult(row.valid, row.source, row.brokenLink, files))
	}
}

// [at -> scn~link-checker_external-links-off-by-default~1]
func TestLinkChecker_externalLinksOffByDefault(t *testing.T) {
	webServer := codevalidation.NewWebServer(t)
	checker := codingagentflow.NewLinkChecker()
	for _, source := range []string{"README.md", "doc/a.rst", "notes.txt", "sitemap.xml"} {
		codevalidation.Check(t, checker,
			[]acceptance.File{document(source, webServer.URL+"/missing")},
			expectNoRequest(webServer),
		)
	}
}

// [at -> scn~link-checker_external-links-opt-in~1]
func TestLinkChecker_externalLinksOptIn(t *testing.T) {
	webServer := codevalidation.NewWebServer(t)
	checker := codingagentflow.NewLinkChecker(codingagentflow.WithExternalLinks())
	for _, row := range []struct {
		source, path string
		valid        bool
		brokenLink   string
	}{
		{"README.md", "/ok", true, noBrokenLink},
		{"README.md", "/rate-limited", true, noBrokenLink},
		{"README.md", "/missing", false, webServer.URL + "/missing"},
		{"doc/a.rst", "/ok", true, noBrokenLink},
		{"doc/a.rst", "/missing", false, webServer.URL + "/missing"},
		{"notes.txt", "/ok", true, noBrokenLink},
		{"notes.txt", "/missing", false, webServer.URL + "/missing"},
		{"sitemap.xml", "/ok", true, noBrokenLink},
		{"sitemap.xml", "/missing", false, webServer.URL + "/missing"},
	} {
		codevalidation.Check(t, checker,
			[]acceptance.File{document(row.source, webServer.URL+row.path)},
			expectResult(row.valid, row.source, row.brokenLink),
		)
	}
}
