package adapter

import (
	"fmt"
	"io"

	"github.com/exasol/coding-agent-flow/internal/event"
	"github.com/exasol/coding-agent-flow/internal/validation"
)

type codexStopInput struct {
	SessionID        string `json:"session_id"`
	WorkingDirectory string `json:"cwd"`
	HookEventName    string `json:"hook_event_name"`
	TurnID           string `json:"turn_id"`
	StopHookActive   bool   `json:"stop_hook_active"`
}

type codexStopOutput struct {
	Decision string `json:"decision,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func decodeCodexStop(input io.Reader) (event.Event, error) {
	var payload codexStopInput
	if err := decodeJSON(input, &payload); err != nil {
		return event.Event{}, err
	}
	result := event.Event{
		Kind:             event.TurnStop,
		WorkingDirectory: payload.WorkingDirectory,
		SessionID:        payload.SessionID,
		AlreadyContinued: payload.StopHookActive,
	}
	if payload.HookEventName != "Stop" {
		return result, fmt.Errorf("expected Codex Stop event")
	}
	if payload.SessionID == "" || payload.TurnID == "" || payload.WorkingDirectory == "" {
		return result, fmt.Errorf("Codex Stop input is missing required fields")
	}
	return result, nil
}

func encodeCodexStop(decision validation.Decision, output io.Writer) error {
	// [impl~caf-codex-stop-response~1 -> req~caf-codex-stop-response~1]
	response := codexStopOutput{}
	if decision.Action == validation.ContinueWork {
		response.Decision = "block"
		response.Reason = decision.Reason
	}
	return encodeJSON(output, response)
}
