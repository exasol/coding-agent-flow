package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/config"
	"github.com/exasol/coding-agent-flow/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// [itest~caf-machine-validation~1 -> req~caf-machine-validation~1]
// [itest~caf-agent-adapter~1 -> req~caf-agent-selection~1]
// [itest~caf-codex-stop-response~1 -> req~caf-codex-stop-response~1]
// [itest~caf-project-discovery~1 -> req~caf-project-discovery~1]
// [itest~caf-validation-routing~1 -> req~caf-validation-routing~1]
func TestRunCodexHookAllowsStopWhenValidationPasses(t *testing.T) {
	projectRoot := createProject(t, true)
	input := codexInput(projectRoot, false)

	exitCode, stdout, stderr := run(t, []string{"validate", "--agent", "codex"}, input, projectRoot)

	assert.Equal(t, 0, exitCode)
	assert.JSONEq(t, `{}`, stdout)
	assert.Empty(t, stderr)
}

func TestRunCodexHookRequestsContinuationForFinding(t *testing.T) {
	projectRoot := createProject(t, false)
	input := codexInput(projectRoot, false)

	exitCode, stdout, stderr := run(t, []string{"validate", "--agent", "codex"}, input, projectRoot)

	assert.Equal(t, 0, exitCode)
	assert.JSONEq(t, `{
		"decision":"block",
		"reason":"CAF-REQUIREMENTS-DOCUMENT-MISSING: Add a supported requirements document under doc/requirements."
	}`, stdout)
	assert.Empty(t, stderr)
}

func TestRunHookAllowsStopAfterPriorContinuation(t *testing.T) {
	projectRoot := createProject(t, false)

	exitCode, stdout, stderr := run(t, []string{"validate", "--agent", "codex"}, codexInput(projectRoot, true), projectRoot)

	assert.Equal(t, 0, exitCode)
	assert.JSONEq(t, `{}`, stdout)
	assert.Empty(t, stderr)
}

func TestRunHookConvertsMalformedInputToProviderToolFailure(t *testing.T) {
	projectRoot := createProject(t, true)

	exitCode, stdout, stderr := run(t, []string{"validate", "--agent", "codex"}, `{`, projectRoot)

	assert.Equal(t, 0, exitCode)
	assert.JSONEq(t, `{
		"decision":"block",
		"reason":"CAF-HOOK-INPUT-INVALID: CAF could not read the Stop-hook input."
	}`, stdout)
	assert.Empty(t, stderr)
}

func TestRunHookDoesNotLoopWhenInvalidInputSaysContinuationIsActive(t *testing.T) {
	projectRoot := createProject(t, true)
	input := `{"session_id":"session","cwd":"` + projectRoot + `","hook_event_name":"Wrong","turn_id":"turn","stop_hook_active":true}`

	exitCode, stdout, stderr := run(t, []string{"validate", "--agent", "codex"}, input, projectRoot)

	assert.Equal(t, 0, exitCode)
	assert.JSONEq(t, `{}`, stdout)
	assert.Empty(t, stderr)
}

func TestRunCIUsesProcessExitStatusForFinding(t *testing.T) {
	projectRoot := createProject(t, false)

	exitCode, stdout, stderr := run(t, []string{"validate"}, "", projectRoot)

	assert.Equal(t, 1, exitCode)
	assert.Empty(t, stderr)
	var report validation.Report
	require.NoError(t, json.Unmarshal([]byte(stdout), &report))
	assert.Equal(t, validation.StatusFinding, report.Status)
}

func TestRunCIDiscoversConfigurationFromNestedDirectory(t *testing.T) {
	projectRoot := createProject(t, true)
	nested := filepath.Join(projectRoot, "nested", "directory")
	require.NoError(t, os.MkdirAll(nested, 0o755))

	exitCode, stdout, stderr := run(t, []string{"validate"}, "", nested)

	assert.Equal(t, 0, exitCode)
	assert.Empty(t, stderr)
	var report validation.Report
	require.NoError(t, json.Unmarshal([]byte(stdout), &report))
	assert.Equal(t, validation.StatusPass, report.Status)
}

func TestRunRejectsUnknownAgent(t *testing.T) {
	projectRoot := createProject(t, true)

	exitCode, stdout, stderr := run(t, []string{"validate", "--agent", "unknown"}, "", projectRoot)

	assert.Equal(t, 2, exitCode)
	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "expected codex")
}

func TestRunCIReportsToolFailureWhenConfigurationIsMissing(t *testing.T) {
	projectRoot := t.TempDir()

	exitCode, stdout, stderr := run(t, []string{"validate"}, "", projectRoot)

	assert.Equal(t, 2, exitCode)
	assert.Empty(t, stderr)
	var report validation.Report
	require.NoError(t, json.Unmarshal([]byte(stdout), &report))
	assert.Equal(t, validation.StatusToolFailure, report.Status)
}

func run(t *testing.T, args []string, input string, workingDirectory string) (int, string, string) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := Run(
		context.Background(),
		args,
		bytes.NewBufferString(input),
		&stdout,
		&stderr,
		func() (string, error) { return workingDirectory, nil },
	)
	return exitCode, stdout.String(), stderr.String()
}

func createProject(t *testing.T, withRequirements bool) string {
	t.Helper()
	projectRoot := t.TempDir()
	configuration := "version: 1\nvalidators:\n  - requirements-document\n"
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, config.FileName), []byte(configuration), 0o600))
	if withRequirements {
		requirementsDirectory := filepath.Join(projectRoot, "doc", "requirements")
		require.NoError(t, os.MkdirAll(requirementsDirectory, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(requirementsDirectory, "system.md"), nil, 0o600))
	}
	return projectRoot
}

func codexInput(projectRoot string, alreadyContinued bool) string {
	payload := map[string]any{
		"session_id":             "session",
		"cwd":                    projectRoot,
		"hook_event_name":        "Stop",
		"turn_id":                "turn",
		"stop_hook_active":       alreadyContinued,
		"last_assistant_message": "done",
	}
	encoded, _ := json.Marshal(payload)
	return string(encoded)
}
