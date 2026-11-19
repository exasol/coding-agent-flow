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

// [at -> scn~external-link-checker_opt-in~1]
func TestExternalLinkChecker_optIn(t *testing.T) {
	webServer := codevalidation.NewWebServer(t)
	checker := codingagentflow.NewLinkChecker(codingagentflow.WithExternalLinks())
	for _, row := range []struct {
		source, path string
		valid        bool
		brokenLink   string
	}{
		{"README.md", "/ok", true, codevalidation.NoBrokenLink},
		{"README.md", "/rate-limited", true, codevalidation.NoBrokenLink},
		{"README.md", "/missing", false, webServer.URL + "/missing"},
		{"doc/a.rst", "/ok", true, codevalidation.NoBrokenLink},
		{"doc/a.rst", "/missing", false, webServer.URL + "/missing"},
		{"notes.txt", "/ok", true, codevalidation.NoBrokenLink},
		{"notes.txt", "/missing", false, webServer.URL + "/missing"},
		{"sitemap.xml", "/ok", true, codevalidation.NoBrokenLink},
		{"sitemap.xml", "/missing", false, webServer.URL + "/missing"},
	} {
		codevalidation.Check(t, checker,
			[]acceptance.File{codevalidation.Document(row.source, webServer.URL+row.path)},
			codevalidation.ExpectResult(row.valid, row.source, row.brokenLink),
		)
	}
}
