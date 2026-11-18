package lychee

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/exasol/coding-agent-flow/internal/codevalidation/domain"
	"github.com/exasol/coding-agent-flow/internal/global"
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

// fakeBinary writes a shell script that prints output and exits with exitCode, and configures it as LYCHEE_BINARY.
func fakeBinary(t *testing.T, output string, exitCode int) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "lychee")
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' '%s'\nexit %d\n", output, exitCode)
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(binaryEnvironmentVariable, binary)
	return binary
}

func TestLocateBinaryAcceptsTheSupportedVersion(t *testing.T) {
	expected := fakeBinary(t, "lychee "+lycheeVersion, 0)

	binary, err := locateBinary(context.Background())

	if err != nil || binary != expected {
		t.Fatalf("got %q, %v; want %q", binary, err, expected)
	}
}

func TestLocateBinaryRefusesOtherVersionsAndUnexpectedOutput(t *testing.T) {
	for name, test := range map[string]struct {
		output   string
		exitCode int
		want     string
	}{
		"other version":     {"lychee 0.18.1", 0, "reports version 0.18.1"},
		"no lychee version": {"garbage", 0, "reports no lychee version"},
		"failing binary":    {"lychee " + lycheeVersion, 1, "--version"},
	} {
		fakeBinary(t, test.output, test.exitCode)

		_, err := locateBinary(context.Background())

		if err == nil || !strings.Contains(err.Error(), "E-CAF-4") || !strings.Contains(err.Error(), test.want) {
			t.Errorf("%s: expected E-CAF-4 with %q, got %v", name, test.want, err)
		}
	}
}

func TestParseVersion(t *testing.T) {
	for output, want := range map[string]string{
		"lychee 0.24.2\n":      "0.24.2",
		"  lychee 1.0.0  ":     "1.0.0",
		"lychee 0.24.2-beta":   "",
		"lychee-lib 0.24.2":    "",
		"lychee 0.24.2\nextra": "",
		"":                     "",
	} {
		version, found := parseVersion([]byte(output))
		if version != want || found != (want != "") {
			t.Errorf("%q: got %q, %t; want %q", output, version, found, want)
		}
	}
}

func TestLocateBinaryReportsAMissingConfiguredBinary(t *testing.T) {
	t.Setenv(binaryEnvironmentVariable, filepath.Join(t.TempDir(), "lychee"))

	_, err := locateBinary(context.Background())

	if err == nil || !strings.Contains(err.Error(), binaryEnvironmentVariable) || !strings.Contains(err.Error(), "E-CAF-4") {
		t.Fatalf("expected E-CAF-4 naming %s, got %v", binaryEnvironmentVariable, err)
	}
}

func releaseArchiveWith(t *testing.T, name, content string) []byte {
	t.Helper()
	var archive bytes.Buffer
	compressed := gzip.NewWriter(&archive)
	files := tar.NewWriter(compressed)
	if err := files.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := files.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := files.Close(); err != nil {
		t.Fatal(err)
	}
	if err := compressed.Close(); err != nil {
		t.Fatal(err)
	}
	return archive.Bytes()
}

func checksumOf(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func TestInstallBinaryWritesTheExecutableFromTheArchive(t *testing.T) {
	archive := releaseArchiveWith(t, "lychee-x86_64-unknown-linux-musl/lychee", "binary")
	binary := filepath.Join(t.TempDir(), "cache", "lychee")

	if err := installBinary(archive, checksumOf(archive), binary); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(binary)
	if err != nil || string(content) != "binary" {
		t.Fatalf("got %q, %v", content, err)
	}
	if info, err := os.Stat(binary); err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Fatalf("expected an executable, got %v, %v", info, err)
	}
}

func TestInstallBinaryRefusesAnArchiveWithAnotherChecksum(t *testing.T) {
	archive := releaseArchiveWith(t, "lychee", "binary")
	binary := filepath.Join(t.TempDir(), "lychee")

	if err := installBinary(archive, checksumOf([]byte("other")), binary); err == nil || !strings.Contains(err.Error(), "checksum") {
		t.Fatalf("expected a checksum error, got %v", err)
	}
	if _, err := os.Stat(binary); !os.IsNotExist(err) {
		t.Fatal("expected no binary after a checksum mismatch")
	}
}

func TestInstallBinaryRefusesAnArchiveWithoutLychee(t *testing.T) {
	archive := releaseArchiveWith(t, "README.md", "text")

	if err := installBinary(archive, checksumOf(archive), filepath.Join(t.TempDir(), "lychee")); err == nil {
		t.Fatal("expected an error for an archive without lychee")
	}
}

// fakeLychee configures as LYCHEE_BINARY a shell script that reports the supported version for --version. For a
// check, it records its arguments and working directory, prints stdout and stderr, and exits with exitCode.
// It returns the file with the recorded arguments.
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
