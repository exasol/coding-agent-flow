package openfasttrace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const openFastTraceVersion = "4.11.0"

const jarEnvironmentVariable = "OPENFASTTRACE_JAR"

const (
	openFastTraceArtifact        = "org.itsallcode.openfasttrace:openfasttrace:"
	mavenDependencyPluginVersion = "3.11.0"
	mavenDependencyPluginGet     = "org.apache.maven.plugins:maven-dependency-plugin:" + mavenDependencyPluginVersion + ":get"
	mavenRepositoryDirectory     = ".m2/repository"
)

// locateJar returns the JAR from OPENFASTTRACE_JAR if set, else from the local Maven repository, where Maven
// downloads it on first use.
func locateJar(ctx context.Context) (string, error) {
	if configured := os.Getenv(jarEnvironmentVariable); configured != "" {
		if _, err := os.Stat(configured); err != nil {
			return "", fmt.Errorf("inspect OpenFastTrace JAR from %s: %w", jarEnvironmentVariable, err)
		}
		return configured, nil
	}
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory for the Maven repository: %w", err)
	}
	jar := filepath.Join(
		homeDirectory,
		mavenRepositoryDirectory,
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

	command := exec.CommandContext(
		ctx,
		"mvn",
		"--batch-mode",
		mavenDependencyPluginGet,
		"-Dartifact="+openFastTraceArtifact+openFastTraceVersion,
		"-Dtransitive=false",
	)
	if output, err := command.CombinedOutput(); err != nil {
		return "", fmt.Errorf("download OpenFastTrace %s: %w: %s", openFastTraceVersion, err, output)
	}
	if _, err := os.Stat(jar); err != nil {
		return "", fmt.Errorf("find downloaded OpenFastTrace JAR: %w", err)
	}
	return jar, nil
}
