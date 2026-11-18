// Package codevalidation holds the shared fixture of the Code Validation use cases.
package codevalidation

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
