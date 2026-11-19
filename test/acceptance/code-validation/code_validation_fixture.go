// Package codevalidation holds the shared fixture of the Code Validation use cases.
package codevalidation

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
)

// WebServer is a local web server that serves the page /ok, answers 429 Too Many Requests for /rate-limited, and
// 404 Not Found for other paths.
type WebServer struct {
	URL      string
	requests atomic.Int64
}

// Requests returns the number of requests the web server received.
func (server *WebServer) Requests() int64 {
	return server.requests.Load()
}

// NewWebServer starts a WebServer that stops at the end of the test.
func NewWebServer(t *testing.T) *WebServer {
	t.Helper()
	webServer := &WebServer{}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		webServer.requests.Add(1)
		switch request.URL.Path {
		case "/ok":
			fmt.Fprint(response, "<html><body>ok</body></html>")
		case "/rate-limited":
			http.Error(response, "rate limited", http.StatusTooManyRequests)
		default:
			http.NotFound(response, request)
		}
	}))
	t.Cleanup(server.Close)
	webServer.URL = server.URL
	return webServer
}

// Document templates with one Link each.
const (
	markdownTemplate         = "\n# Document\n\n## Setup\n\n[link](%s)\n"
	restructuredTextTemplate = "\nDocument\n========\n\nSee %s for details.\n"
	textTemplate             = "See %s for details.\n"
	sitemapTemplate          = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n  <url><loc>%s</loc></url>\n</urlset>\n"
)

// NoBrokenLink marks a valid expected result.
const NoBrokenLink = ""

// Document creates the file at path with the Link, using the template that matches its extension.
func Document(path, link string) acceptance.File {
	templates := map[string]string{
		".md":  markdownTemplate,
		".rst": restructuredTextTemplate,
		".txt": textTemplate,
		".xml": sitemapTemplate,
	}
	return acceptance.NewFile(path, templates[filepath.Ext(path)], link)
}

// ExpectResult checks valid and, if brokenLink is set, that the Broken Link is reported in source.
func ExpectResult(valid bool, source, brokenLink string) func(codingagentflow.LinkValidation) error {
	return func(response codingagentflow.LinkValidation) error {
		if err := acceptance.Expect("valid", valid, response.Valid()); err != nil {
			return fmt.Errorf("%w: Broken Links %+v", err, response.BrokenLinks())
		}
		if brokenLink != NoBrokenLink && !response.ContainsBrokenLink(source, brokenLink) {
			return fmt.Errorf("Broken Link %s in %s: not reported in %+v", brokenLink, source, response.BrokenLinks())
		}
		return nil
	}
}

// Check runs the LinkChecker for a target-project fixture with one or more files.
func Check(t *testing.T, checker codingagentflow.LinkChecker, files []acceptance.File,
	then func(codingagentflow.LinkValidation) error) {
	t.Helper()
	acceptance.Run(t, files, acceptance.TargetProject, checker.Check, then)
}

// CheckTargetProject works like Check and also passes the Target Project to then.
func CheckTargetProject(t *testing.T, checker codingagentflow.LinkChecker, files []acceptance.File,
	then func(codingagentflow.LinkValidation, codingagentflow.TargetProject) error) {
	t.Helper()
	var targetProject codingagentflow.TargetProject
	given := func(t *testing.T, files []acceptance.File) codingagentflow.TargetProject {
		targetProject = acceptance.TargetProject(t, files)
		return targetProject
	}
	acceptance.Run(t, files, given, checker.Check, func(response codingagentflow.LinkValidation) error {
		return then(response, targetProject)
	})
}
