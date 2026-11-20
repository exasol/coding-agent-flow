package lychee

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/global"
)

func TestArgumentsRequestAJSONReportWithFragmentsForTheSupportedExtensions(t *testing.T) {
	for _, scope := range []domain.LinkScope{domain.LocalLinks, domain.LocalAndExternalLinks} {
		arguments := strings.Join(arguments("/project", "/tmp/caf-lychee.toml", scope), " ")

		for _, expected := range []string{"--config /tmp/caf-lychee.toml", "--format json", "--include-fragments", "--extensions md,rst,txt,xml", "--root-dir /project"} {
			if !strings.Contains(arguments, expected) {
				t.Errorf("expected %q in %q", expected, arguments)
			}
		}
		if !strings.HasSuffix(arguments, " .") {
			t.Errorf("expected the Target Project as last argument in %q", arguments)
		}
	}
}

// [utest -> req~external-link-validation~1]
func TestArgumentsCheckLocalLinksOffline(t *testing.T) {
	arguments := strings.Join(arguments("/project", "/tmp/caf-lychee.toml", domain.LocalLinks), " ")

	if !strings.Contains(arguments, "--offline") || strings.Contains(arguments, "--accept") {
		t.Fatalf("expected an offline check in %q", arguments)
	}
}

func TestArgumentsCheckExternalLinksWithRateLimitedURLsAsValid(t *testing.T) {
	arguments := strings.Join(arguments("/project", "/tmp/caf-lychee.toml", domain.LocalAndExternalLinks), " ")

	if strings.Contains(arguments, "--offline") {
		t.Fatalf("expected no offline check in %q", arguments)
	}
	for _, expected := range []string{"--accept 200..=299,429", "--timeout 10", "--max-retries 2"} {
		if !strings.Contains(arguments, expected) {
			t.Errorf("expected %q in %q", expected, arguments)
		}
	}
}

func TestWriteConfigurationDisablesTheCacheOutsideTheTargetProject(t *testing.T) {
	configuration, err := writeConfiguration()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(configuration)

	content, err := os.ReadFile(configuration)
	if err != nil || !strings.Contains(string(content), "cache = false") {
		t.Fatalf("got %q, %v", content, err)
	}
}

func fakeLychee(t *testing.T, stdout, stderr string, exitCode int) string {
	t.Helper()
	directory := t.TempDir()
	for name, content := range map[string]string{"stdout": stdout, "stderr": stderr} {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	calls := filepath.Join(directory, "calls")
	script := fmt.Sprintf(`#!/bin/sh
if [ "$1" = "--version" ]; then echo "lychee %s"; exit 0; fi
{ pwd; echo "$@"; } > '%s'
cat '%s'
cat '%s' >&2
exit %d
`, lycheeVersion, calls, filepath.Join(directory, "stdout"), filepath.Join(directory, "stderr"), exitCode)
	binary := filepath.Join(directory, "lychee")
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(binaryEnvironmentVariable, binary)
	return calls
}

func findBrokenLinks(t *testing.T, root string) (domain.LinkValidation, error) {
	t.Helper()
	return NewAdapter().FindBrokenLinks(context.Background(), targetProject(root), domain.LocalLinks)
}

func TestFindBrokenLinksRunsLycheeInTheTargetProjectAndReadsItsReport(t *testing.T) {
	root := t.TempDir()
	report := `{"errors": 1, "error_map": {"./README.md": [
  {"url": "file://` + root + `/doc/missing.md", "status": {"text": "File not found"}, "span": {"line": 3}}
]}, "timeout_map": {}}`
	calls := fakeLychee(t, report, "", exitLinkFailures)

	validation, err := findBrokenLinks(t, root)

	if err != nil {
		t.Fatal(err)
	}
	if validation.Valid() || !validation.ContainsBrokenLink("README.md", "doc/missing.md") {
		t.Fatalf("expected the Broken Link doc/missing.md in README.md, got %+v", validation.BrokenLinks())
	}
	recorded, err := os.ReadFile(calls)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.SplitN(string(recorded), "\n", 2)
	if resolved, _ := filepath.EvalSymlinks(root); lines[0] != root && lines[0] != resolved {
		t.Errorf("expected lychee to run in %s, got %s", root, lines[0])
	}
	configuration := regexpConfiguration.FindStringSubmatch(lines[1])
	if configuration == nil || !strings.Contains(lines[1], "--offline") {
		t.Fatalf("expected --config and --offline in %q", lines[1])
	}
	if _, err := os.Stat(configuration[1]); !os.IsNotExist(err) {
		t.Errorf("expected the CAF configuration %s to be removed after the check", configuration[1])
	}
}

func TestFindBrokenLinksReportsLycheeFailures(t *testing.T) {
	for name, test := range map[string]struct {
		stdout, stderr string
		exitCode       int
		want           []string
	}{
		"unexpected runtime failure": {"", "Error: no inputs", 1, []string{"E-CAF-2", "Error: no inputs"}},
		"configuration error":        {"", "Error: invalid config", 3, []string{"E-CAF-2", "invalid config"}},
		"report is no JSON":          {"Error: crashed", "", exitLinkFailures, []string{"E-CAF-3", "no lychee JSON report"}},
		"success with failed links":  {report, "", exitSuccess, []string{"E-CAF-3", "exit code does not match"}},
		"failure without failed links": {`{"errors": 0, "error_map": {}, "timeout_map": {}}`, "", exitLinkFailures,
			[]string{"E-CAF-3", "exit code does not match"}},
	} {
		fakeLychee(t, test.stdout, test.stderr, test.exitCode)

		_, err := findBrokenLinks(t, t.TempDir())

		for _, want := range test.want {
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s: expected an error with %q, got %v", name, want, err)
			}
		}
	}
}

func TestFindBrokenLinksReportsAMissingTargetProject(t *testing.T) {
	fakeLychee(t, "", "", exitSuccess)

	_, err := findBrokenLinks(t, filepath.Join(t.TempDir(), "missing"))

	if err == nil || !strings.Contains(err.Error(), "E-CAF-2") {
		t.Fatalf("expected E-CAF-2, got %v", err)
	}
}

func TestFindBrokenLinksStopsWhenCancelled(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "lychee")
	script := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo \"lychee " + lycheeVersion + "\"; exit 0; fi\nexec sleep 30\n"
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(binaryEnvironmentVariable, binary)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	started := time.Now()

	_, err := NewAdapter().FindBrokenLinks(ctx, targetProject(t.TempDir()), domain.LocalLinks)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected the context error, got %v", err)
	}
	if elapsed := time.Since(started); elapsed > 5*time.Second {
		t.Fatalf("expected the check to stop at once, took %s", elapsed)
	}
}

var regexpConfiguration = regexp.MustCompile(`--config (\S+)`)

func targetProject(path string) global.TargetProject {
	return global.NewTargetProject(path)
}
