package lychee

import (
	"reflect"
	"strings"
	"testing"
)

const report = `{
  "total": 4, "errors": 3, "timeouts": 0,
  "error_map": {
    "./README.md": [
      {"url": "file:///project/doc/missing.md", "status": {"text": "File not found"}, "span": {"line": 5, "column": 20}},
      {"url": "file:///project/README.md#nope", "status": {"text": "Cannot find fragment"}, "span": {"line": 3, "column": 1}}
    ],
    "./doc/notes.txt": [
      {"url": "http://127.0.0.1:8080/gone", "status": {"text": "Rejected status code: 404 Not Found"}, "span": {"line": 1, "column": 5}}
    ]
  },
  "timeout_map": {}
}`

type brokenLink struct {
	source string
	line   int
	target string
	status string
}

// [utest -> req~link-validation~1]
func TestReadJSONReportReturnsTheBrokenLinksWithProjectRelativeTargets(t *testing.T) {
	validation, err := readJSONReport([]byte(report), "/project", true)

	if err != nil {
		t.Fatal(err)
	}
	var got []brokenLink
	for _, link := range validation.BrokenLinks() {
		got = append(got, brokenLink{link.Link().Source(), link.Link().Line(), link.Link().Target(), link.Status()})
	}
	want := []brokenLink{
		{"README.md", 3, "README.md#nope", "Cannot find fragment"},
		{"README.md", 5, "doc/missing.md", "File not found"},
		{"doc/notes.txt", 1, "http://127.0.0.1:8080/gone", "Rejected status code: 404 Not Found"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestReadJSONReportKeepsFileURLsOutsideTheTargetProject(t *testing.T) {
	if got := linkTarget("/project", "file:///outside.md"); got != "file:///outside.md" {
		t.Fatalf("got %s", got)
	}
}

func TestReadJSONReportReturnsValidForAReportWithoutFailedLinks(t *testing.T) {
	validation, err := readJSONReport([]byte(`{"errors": 0, "error_map": {}, "timeout_map": {}}`), "/project", false)

	if err != nil || !validation.Valid() {
		t.Fatalf("expected valid, got %+v, %v", validation, err)
	}
}

func TestReadJSONReportRefusesUntrustworthyReports(t *testing.T) {
	for name, test := range map[string]struct {
		output       string
		linkFailures bool
	}{
		"no JSON":                       {"Error: invalid input", true},
		"failure exit without failures": {`{"error_map": {}}`, true},
		"success exit with failures":    {report, false},
	} {
		if _, err := readJSONReport([]byte(test.output), "/project", test.linkFailures); err == nil || !strings.Contains(err.Error(), "E-CAF-3") {
			t.Errorf("%s: expected E-CAF-3, got %v", name, err)
		}
	}
}
