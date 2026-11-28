package domain

import (
	"encoding/json"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
)

func stopEvent(continuationActive bool) StopEvent {
	return NewStopEvent(global.NewTargetProject("project"), continuationActive)
}

// The zero report accompanies a Validation Runner error and reports tool_failure.
func TestTurnDecisionOnAZeroReportContinuesTheTurn(t *testing.T) {
	if decision := NewTurnDecision(stopEvent(false), ValidationReport{}); decision.Action() != ContinueTurn || decision.Reason() == "" {
		t.Fatalf("got %s %q", decision.Action(), decision.Reason())
	}
}

func TestFailedTurnDecisionHasTheZeroReport(t *testing.T) {
	if status := NewFailedTurnDecision(stopEvent(false), "E-CAF-5").Report().Status(); status != ToolFailure {
		t.Fatalf("got %s", status)
	}
}

func TestStopEventKeepsTheTargetProjectAndTheContinuationActiveFlag(t *testing.T) {
	event := stopEvent(true)

	if event.TargetProject().Path() != "project" || !event.ContinuationActive() {
		t.Fatalf("got %+v", event)
	}
}

func TestZeroTurnDecisionEndsTheTurn(t *testing.T) {
	var decision TurnDecision

	if decision.Action() != EndTurn {
		t.Fatalf("got %s", decision.Action())
	}
}

func TestTurnDecisionUsesThePublishedResponseFormat(t *testing.T) {
	decision := NewTurnDecision(stopEvent(false), NewValidationReport([]ValidationResult{NewValidationResult("check", Finding, "CODE", "message")}))

	encoded, err := json.Marshal(decision)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"action":"continue_turn","reason":"CAF validation reports problems:\n- check: finding CODE: message","report":{"status":"finding","results":[{"validation":"check","status":"finding","code":"CODE","message":"message"}]}}`
	if string(encoded) != want {
		t.Fatalf("got %s, want %s", encoded, want)
	}
}
