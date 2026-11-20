package lychee

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"

	exaerror "github.com/exasol/error-reporting-go"

	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
)

// jsonReport is the part of the lychee JSON report that lists the failed links per input file.
type jsonReport struct {
	ErrorMap   map[string][]jsonResponse `json:"error_map"`
	TimeoutMap map[string][]jsonResponse `json:"timeout_map"`
}

type jsonResponse struct {
	URL    string `json:"url"`
	Status struct {
		Text string `json:"text"`
	} `json:"status"`
	Span struct {
		Line int `json:"line"`
	} `json:"span"`
}

// readJSONReport reads the Broken Links from a lychee JSON report. linkFailures is true for lychee exit code 2. The
// report must list failed links exactly when linkFailures is true.
func readJSONReport(output []byte, root string, linkFailures bool) (domain.LinkValidation, error) {
	var report jsonReport
	if err := json.Unmarshal(output, &report); err != nil {
		return domain.LinkValidation{}, reportError("the report is no lychee JSON report: " + err.Error())
	}
	var brokenLinks []domain.BrokenLink
	for _, responsesByInput := range []map[string][]jsonResponse{report.ErrorMap, report.TimeoutMap} {
		for input, responses := range responsesByInput {
			source := projectPath(root, input)
			for _, response := range responses {
				link := domain.NewLink(source, response.Span.Line, linkTarget(root, response.URL))
				brokenLinks = append(brokenLinks, domain.NewBrokenLink(link, response.Status.Text))
			}
		}
	}
	if linkFailures != (len(brokenLinks) > 0) {
		return domain.LinkValidation{}, reportError("the exit code does not match the failed links in the report")
	}
	return domain.NewLinkValidation(brokenLinks), nil
}

// projectPath converts an input file of the report, such as ./doc/a.md, to a slash-separated, project-relative path.
func projectPath(root, input string) string {
	if filepath.IsAbs(input) {
		if relative, inRoot := relativeTo(root, input); inRoot {
			return relative
		}
		return filepath.ToSlash(input)
	}
	return filepath.ToSlash(filepath.Clean(input))
}

// linkTarget converts a file URL in the Target Project to a project-relative path with an optional #anchor.
// It keeps other URLs.
func linkTarget(root, target string) string {
	parsed, err := url.Parse(target)
	if err != nil || parsed.Scheme != "file" {
		return target
	}
	relative, inRoot := relativeTo(root, filepath.FromSlash(parsed.Path))
	if !inRoot {
		return target
	}
	if parsed.Fragment != "" {
		return relative + "#" + parsed.Fragment
	}
	return relative
}

func relativeTo(root, path string) (string, bool) {
	for _, base := range candidateRoots(root) {
		relative, err := filepath.Rel(base, path)
		if err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return filepath.ToSlash(relative), true
		}
	}
	return "", false
}

// candidateRoots returns the root as given and with symbolic links resolved, since lychee may report either.
func candidateRoots(root string) []string {
	if resolved, err := filepath.EvalSymlinks(root); err == nil && resolved != root {
		return []string{root, resolved}
	}
	return []string{root}
}

func reportError(reason string) error {
	return exaerror.New("E-CAF-3").
		Message("Failed to read the lychee report: {{reason|uq}}.").
		Parameter("reason", reason).
		Mitigation("Use lychee " + lycheeVersion + ", the version CAF supports.")
}
