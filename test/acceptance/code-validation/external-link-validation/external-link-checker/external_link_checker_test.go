package externallinkchecker

import (
	"fmt"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	codevalidation "github.com/exasol/coding-agent-flow/test/acceptance/code-validation"
)

func expectNoRequest(webServer *codevalidation.WebServer) func(codingagentflow.LinkValidation) error {
	return func(response codingagentflow.LinkValidation) error {
		if err := codevalidation.ExpectResult(true, "", codevalidation.NoBrokenLink)(response); err != nil {
			return err
		}
		if requests := webServer.Requests(); requests != 0 {
			return fmt.Errorf("web server received %d requests, expected none", requests)
		}
		return nil
	}
}

// [at -> scn~external-link-checker_off-by-default~1]
func TestExternalLinkChecker_offByDefault(t *testing.T) {
	webServer := codevalidation.NewWebServer(t)
	checker := codingagentflow.NewLinkChecker()
	for _, source := range []string{"README.md", "doc/a.rst", "notes.txt", "sitemap.xml"} {
		codevalidation.Check(t, checker,
			[]acceptance.File{codevalidation.Document(source, webServer.URL+"/missing")},
			expectNoRequest(webServer),
		)
	}
}

// [at -> scn~external-link-checker_valid-links~1]
func TestExternalLinkChecker_validLinks(t *testing.T) {
	webServer := codevalidation.NewWebServer(t)
	checker := codingagentflow.NewLinkChecker(codingagentflow.WithExternalLinks())
	for _, row := range []struct{ source, path string }{
		{"README.md", "/ok"},
		{"README.md", "/rate-limited"},
		{"doc/a.rst", "/ok"},
		{"notes.txt", "/ok"},
		{"sitemap.xml", "/ok"},
	} {
		codevalidation.Check(t, checker,
			[]acceptance.File{codevalidation.Document(row.source, webServer.URL+row.path)},
			codevalidation.ExpectResult(true, row.source, codevalidation.NoBrokenLink),
		)
	}
}

// [at -> scn~external-link-checker_broken-link~1]
func TestExternalLinkChecker_brokenLink(t *testing.T) {
	webServer := codevalidation.NewWebServer(t)
	checker := codingagentflow.NewLinkChecker(codingagentflow.WithExternalLinks())
	for _, source := range []string{"README.md", "doc/a.rst", "notes.txt", "sitemap.xml"} {
		codevalidation.Check(t, checker,
			[]acceptance.File{codevalidation.Document(source, webServer.URL+"/missing")},
			codevalidation.ExpectResult(false, source, webServer.URL+"/missing"),
		)
	}
}
