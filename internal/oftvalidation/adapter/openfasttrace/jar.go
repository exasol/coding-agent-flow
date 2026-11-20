package openfasttrace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const openFastTraceVersion = "4.11.0"

const jarEnvironmentVariable = "OPENFASTTRACE_JAR"

const (
	openFastTraceArtifact        = "org.itsallcode.openfasttrace:openfasttrace:"
	mavenDependencyPluginVersion = "3.11.0"
	mavenDependencyPluginGet     = "org.apache.maven.plugins:maven-dependency-plugin:" + mavenDependencyPluginVersion + ":get"
	mavenRepositoryDirectory     = ".m2/repository" // below the home directory
)

// versionPattern matches the first line of "oft help", such as "OpenFastTrace 4.11.0".
var versionPattern = regexp.MustCompile(`^OpenFastTrace (\S+)$`)

// versionTimeout limits the time "oft help" may take, including the Java start.
const versionTimeout = 30 * time.Second

// locateJar returns the JAR from OPENFASTTRACE_JAR if set and of the supported version, else from the local Maven
// repository, where Maven downloads it on first use.
func locateJar(ctx context.Context) (string, error) {
	if configured := os.Getenv(jarEnvironmentVariable); configured != "" {
		if _, err := os.Stat(configured); err != nil {
			return "", fmt.Errorf("inspect OpenFastTrace JAR from %s: %w", jarEnvironmentVariable, err)
		}
		if err := verifyVersion(ctx, configured); err != nil {
			return "", fmt.Errorf("verify OpenFastTrace JAR from %s: %w", jarEnvironmentVariable, err)
		}
		return configured, nil
	}
	repository, err := mavenRepository()
	if err != nil {
		return "", err
	}
	return locateJarInRepository(ctx, repository)
}

// mavenRepository returns the local Maven repository that CAF uses for the OpenFastTrace JAR: ~/.m2/repository.
func mavenRepository() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory for the Maven repository: %w", err)
	}
	return filepath.Join(homeDirectory, mavenRepositoryDirectory), nil
}

// locateJarInRepository returns the OpenFastTrace JAR from the local Maven repository. When the JAR is missing, Maven
// downloads it into this repository, which CAF passes explicitly, so the download lands in the same place for every
// Maven configuration. Maven runs in a new empty directory, because it reads .mvn/maven.config, .mvn/extensions.xml, and
// pom.xml from its working directory and its parents; a Target Project could otherwise change the command or load
// build extensions.
func locateJarInRepository(ctx context.Context, repository string) (string, error) {
	jar := filepath.Join(
		repository,
		"org",
		"itsallcode",
		"openfasttrace",
		"openfasttrace",
		openFastTraceVersion,
		"openfasttrace-"+openFastTraceVersion+".jar",
	)
	if _, err := os.Stat(jar); err == nil {
		return jar, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("inspect OpenFastTrace JAR: %w", err)
	}

	directory, err := os.MkdirTemp("", "caf-maven-")
	if err != nil {
		return "", fmt.Errorf("create working directory for Maven: %w", err)
	}
	defer os.RemoveAll(directory)
	output, err := runWithDiagnosticsIn(ctx, directory,
		"mvn",
		"--batch-mode",
		mavenDependencyPluginGet,
		"-Dmaven.repo.local="+repository,
		"-Dartifact="+openFastTraceArtifact+openFastTraceVersion,
		"-Dtransitive=false",
	)
	if err != nil {
		return "", fmt.Errorf("download OpenFastTrace %s: %w: %s", openFastTraceVersion, err, output)
	}
	if _, err := os.Stat(jar); err != nil {
		return "", fmt.Errorf("find downloaded OpenFastTrace JAR: %w", err)
	}
	return jar, nil
}

// verifyVersion requires "java -jar jar help" to report openFastTraceVersion in its first line.
func verifyVersion(ctx context.Context, jar string) error {
	ctx, cancel := context.WithTimeout(ctx, versionTimeout)
	defer cancel()
	output, err := runWithDiagnostics(ctx, "java", "-jar", jar, "help")
	if err != nil {
		return withDiagnostics(fmt.Errorf("run %s help: %w", jar, err), output)
	}
	firstLine, _, _ := strings.Cut(strings.TrimSpace(output), "\n")
	match := versionPattern.FindStringSubmatch(strings.TrimSpace(firstLine))
	if match == nil {
		return fmt.Errorf("%s reports no OpenFastTrace version: %q", jar, firstLine)
	}
	if match[1] != openFastTraceVersion {
		return fmt.Errorf("%s reports OpenFastTrace %s, want %s", jar, match[1], openFastTraceVersion)
	}
	return nil
}
