package adapter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/event"
	"github.com/exasol/coding-agent-flow/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeCodexStopNormalizesProviderInput(t *testing.T) {
	input := `{"session_id":"session","cwd":"/project","hook_event_name":"Stop","turn_id":"turn","stop_hook_active":true,"model":"example"}`

	result, err := DecodeStop(Codex, strings.NewReader(input))

	require.NoError(t, err)
	assert.Equal(t, event.TurnStop, result.Kind)
	assert.Equal(t, "/project", result.WorkingDirectory)
	assert.True(t, result.AlreadyContinued)
}

func TestDecodeCodexStopRejectsOversizedInput(t *testing.T) {
	input := strings.Repeat(" ", maxHookInputBytes+1)

	_, err := DecodeStop(Codex, strings.NewReader(input))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds")
}

func TestEncodeStopWritesProviderProtocolOnly(t *testing.T) {
	var output bytes.Buffer
	err := EncodeStop(Codex, validation.Decision{Action: validation.ContinueWork, Reason: "Run tests."}, &output)

	require.NoError(t, err)
	assert.JSONEq(t, `{"decision":"block","reason":"Run tests."}`, output.String())
}

func TestEncodeStopAllowsProviderToStopWithEmptyObject(t *testing.T) {
	var output bytes.Buffer
	err := EncodeStop(Codex, validation.Decision{Action: validation.AllowStop}, &output)

	require.NoError(t, err)
	assert.JSONEq(t, `{}`, output.String())
}
