package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// [utest~caf-project-discovery~1 -> req~caf-project-discovery~1]
func TestDiscoverFindsNearestConfigurationAndDerivesRoot(t *testing.T) {
	outer := t.TempDir()
	writeConfig(t, outer, "requirements-document")
	projectRoot := filepath.Join(outer, "project")
	workingDirectory := filepath.Join(projectRoot, "nested", "directory")
	require.NoError(t, os.MkdirAll(workingDirectory, 0o755))
	writeConfig(t, projectRoot, "requirements-document")

	project, err := Discover(workingDirectory)

	require.NoError(t, err)
	assert.Equal(t, projectRoot, project.Root)
	assert.Equal(t, filepath.Join(projectRoot, FileName), project.ConfigPath)
	assert.Equal(t, Config{Version: 1, Validators: []string{"requirements-document"}}, project.Config)
}

func TestDiscoverRejectsMissingConfiguration(t *testing.T) {
	_, err := Discover(t.TempDir())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "E-CAF-2")
}

func TestDiscoverRejectsInvalidConfiguration(t *testing.T) {
	projectRoot := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(projectRoot, FileName),
		[]byte("version: 1\nunknown: true\nvalidators:\n  - requirements-document\n"),
		0o600,
	))

	_, err := Discover(projectRoot)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "E-CAF-3")
}

func writeConfig(t *testing.T, projectRoot string, validators ...string) {
	t.Helper()
	content := "version: 1\nvalidators:\n"
	for _, validator := range validators {
		content += "  - " + validator + "\n"
	}
	require.NoError(t, os.MkdirAll(projectRoot, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, FileName), []byte(content), 0o600))
}
