package openfasttrace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocateJarReportsAMissingConfiguredJar(t *testing.T) {
	t.Setenv(jarEnvironmentVariable, filepath.Join(t.TempDir(), "openfasttrace.jar"))

	_, err := locateJar(context.Background())

	if err == nil || !strings.Contains(err.Error(), jarEnvironmentVariable) {
		t.Fatalf("expected an error naming %s for a missing JAR, got %v", jarEnvironmentVariable, err)
	}
}

// Regression: the download error keeps the first maximumDiagnosticsSize bytes of the Maven output.
func TestLocateJarKeepsTheFirstMavenOutput(t *testing.T) {
	directory := t.TempDir()
	script := "#!/bin/sh\nhead -c 100000 /dev/zero | tr '\\0' m\nexit 1\n"
	if err := os.WriteFile(filepath.Join(directory, "mvn"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("HOME", t.TempDir())
	t.Setenv(jarEnvironmentVariable, "")

	_, err := locateJar(context.Background())

	if err == nil || !strings.Contains(err.Error(), "download OpenFastTrace") {
		t.Fatalf("expected the download error, got %v", err)
	}
	if !strings.Contains(err.Error(), strings.Repeat("m", maximumDiagnosticsSize)) || strings.Contains(err.Error(), strings.Repeat("m", maximumDiagnosticsSize+1)) {
		t.Fatalf("expected exactly %d bytes of Maven output in the error, got %d bytes of error", maximumDiagnosticsSize, len(err.Error()))
	}
}

// configuredJar configures OPENFASTTRACE_JAR with an empty file and puts a java script first on the PATH that prints
// helpOutput and exits with exitCode for "java -jar <jar> help".
func configuredJar(t *testing.T, helpOutput string, exitCode int) string {
	t.Helper()
	directory := t.TempDir()
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' '%s'\nexit %d\n", helpOutput, exitCode)
	if err := os.WriteFile(filepath.Join(directory, "java"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	jar := filepath.Join(directory, "openfasttrace.jar")
	if err := os.WriteFile(jar, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv(jarEnvironmentVariable, jar)
	return jar
}

func TestLocateJarAcceptsAConfiguredJarOfTheSupportedVersion(t *testing.T) {
	expected := configuredJar(t, "OpenFastTrace "+openFastTraceVersion, 0)

	jar, err := locateJar(context.Background())

	if err != nil || jar != expected {
		t.Fatalf("got %q, %v; want %q", jar, err, expected)
	}
}

// Regression: a configured JAR of another OpenFastTrace version is refused.
func TestLocateJarRefusesAConfiguredJarOfAnotherVersion(t *testing.T) {
	configuredJar(t, "OpenFastTrace 4.10.0", 0)

	_, err := locateJar(context.Background())

	for _, want := range []string{jarEnvironmentVariable, "reports OpenFastTrace 4.10.0", "want " + openFastTraceVersion} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("expected an error with %q, got %v", want, err)
		}
	}
}

// Regression: a configured JAR that reports no OpenFastTrace version or fails is refused.
func TestLocateJarRefusesAnUnrelatedConfiguredJar(t *testing.T) {
	for name, test := range map[string]struct {
		output   string
		exitCode int
		want     string
	}{
		"unrelated JAR": {"Hello from another tool", 0, "reports no OpenFastTrace version"},
		"failing JAR":   {"Error: Invalid or corrupt jarfile", 1, "Invalid or corrupt jarfile"},
	} {
		configuredJar(t, test.output, test.exitCode)

		_, err := locateJar(context.Background())

		if err == nil || !strings.Contains(err.Error(), jarEnvironmentVariable) || !strings.Contains(err.Error(), test.want) {
			t.Errorf("%s: expected an error with %q, got %v", name, test.want, err)
		}
	}
}

func TestVersionPatternReadsTheVersion(t *testing.T) {
	for line, want := range map[string]string{
		"OpenFastTrace 4.11.0":       "4.11.0",
		"OpenFastTrace 4.11.0 extra": "",
		"OpenFastTraceX 4.11.0":      "",
		"Usage: oft":                 "",
	} {
		match := versionPattern.FindStringSubmatch(line)
		if got := ""; match != nil {
			got = match[1]
			if got != want {
				t.Errorf("%q: got %q, want %q", line, got, want)
			}
		} else if want != "" {
			t.Errorf("%q: got no match, want %q", line, want)
		}
	}
}

// fakeMaven puts an mvn script first on the PATH. It records its arguments and, if createJar is true, creates the JAR
// in the repository of -Dmaven.repo.local, like the Maven download. It returns the file with the recorded arguments.
func fakeMaven(t *testing.T, createJar bool, exitCode int) string {
	t.Helper()
	directory := t.TempDir()
	calls := filepath.Join(directory, "calls")
	create := ""
	if createJar {
		create = fmt.Sprintf("mkdir -p \"$repository/org/itsallcode/openfasttrace/openfasttrace/%[1]s\" && touch \"$repository/org/itsallcode/openfasttrace/openfasttrace/%[1]s/openfasttrace-%[1]s.jar\"\n", openFastTraceVersion)
	}
	script := fmt.Sprintf(`#!/bin/sh
echo "$@" >> '%s'
for argument in "$@"; do
  case "$argument" in -Dmaven.repo.local=*) repository="${argument#-Dmaven.repo.local=}";; esac
done
%sexit %d
`, calls, create, exitCode)
	if err := os.WriteFile(filepath.Join(directory, "mvn"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))
	return calls
}

func recordedCalls(t *testing.T, calls string) string {
	t.Helper()
	content, err := os.ReadFile(calls)
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

// Regression: when the JAR is missing, Maven downloads it into the repository that CAF passes explicitly.
func TestLocateJarInRepositoryDownloadsIntoTheGivenRepository(t *testing.T) {
	repository := t.TempDir()
	calls := fakeMaven(t, true, 0)

	jar, err := locateJarInRepository(context.Background(), repository)

	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(repository, "org", "itsallcode", "openfasttrace", "openfasttrace", openFastTraceVersion, "openfasttrace-"+openFastTraceVersion+".jar"); jar != want {
		t.Fatalf("got %q, want %q", jar, want)
	}
	recorded := recordedCalls(t, calls)
	for _, want := range []string{"--batch-mode", "-Dmaven.repo.local=" + repository, "-Dartifact=" + openFastTraceArtifact + openFastTraceVersion, "-Dtransitive=false"} {
		if !strings.Contains(recorded, want) {
			t.Errorf("expected %q in the Maven arguments %q", want, recorded)
		}
	}
}

func TestLocateJarInRepositoryUsesACachedJarWithoutMaven(t *testing.T) {
	repository := t.TempDir()
	calls := fakeMaven(t, true, 0)
	cached, err := locateJarInRepository(context.Background(), repository)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(calls); err != nil {
		t.Fatal(err)
	}

	jar, err := locateJarInRepository(context.Background(), repository)

	if err != nil || jar != cached || recordedCalls(t, calls) != "" {
		t.Fatalf("expected the cached JAR without a Maven run, got %q, %v, calls %q", jar, err, recordedCalls(t, calls))
	}
}

func TestLocateJarInRepositoryReportsADownloadWithoutJar(t *testing.T) {
	fakeMaven(t, false, 0)

	_, err := locateJarInRepository(context.Background(), t.TempDir())

	if err == nil || !strings.Contains(err.Error(), "find downloaded OpenFastTrace JAR") {
		t.Fatalf("expected the missing download error, got %v", err)
	}
}

func TestLocateJarUsesTheMavenRepositoryInTheHomeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(jarEnvironmentVariable, "")
	calls := fakeMaven(t, true, 0)

	jar, err := locateJar(context.Background())

	if err != nil || !strings.HasPrefix(jar, filepath.Join(home, ".m2", "repository")) {
		t.Fatalf("got %q, %v", jar, err)
	}
	if recorded := recordedCalls(t, calls); !strings.Contains(recorded, "-Dmaven.repo.local="+filepath.Join(home, ".m2", "repository")) {
		t.Fatalf("expected the home repository in the Maven arguments %q", recorded)
	}
}
