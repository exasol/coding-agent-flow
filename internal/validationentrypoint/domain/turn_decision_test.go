package domain

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
)

func stopEvent(continuationActive bool) StopEvent {
	return NewStopEvent(global.NewTargetProject("project"), continuationActive)
}

func report(statuses ...ValidationStatus) ValidationReport {
	var results []ValidationResult
	for index, status := range statuses {
		results = append(results, NewValidationResult(string(rune('a'+index)), status, "CODE-"+string(status), "message of "+string(status)))
	}
	return NewValidationReport(results)
}

// [utest -> req~agent-turn-validation~1]
// [utest -> req~agent-continuation-limit~1]
func TestNewTurnDecisionFollowsTheDecisionTable(t *testing.T) {
	for name, test := range map[string]struct {
		report             ValidationReport
		continuationActive bool
		want               TurnAction
	}{
		"pass":                                  {report(Pass), false, EndTurn},
		"not applicable":                        {report(NotApplicable), false, EndTurn},
		"no results":                            {report(), false, EndTurn},
		"finding":                               {report(Pass, Finding), false, ContinueTurn},
		"tool failure":                          {report(ToolFailure, Pass), false, ContinueTurn},
		"zero report":                           {ValidationReport{}, false, ContinueTurn},
		"pass with continuation active":         {report(Pass), true, EndTurn},
		"not applicable with continuation":      {report(NotApplicable), true, EndTurn},
		"finding with continuation active":      {report(Finding), true, EndTurn},
		"tool failure with continuation active": {report(ToolFailure), true, EndTurn},
	} {
		decision := NewTurnDecision(stopEvent(test.continuationActive), test.report)

		if decision.Action() != test.want {
			t.Errorf("%s: got %s, want %s", name, decision.Action(), test.want)
		}
	}
}

// [utest -> req~agent-turn-validation~1]
func TestNewTurnDecisionHasNoReasonWhenValidationPasses(t *testing.T) {
	for _, continuationActive := range []bool{false, true} {
		if reason := NewTurnDecision(stopEvent(continuationActive), report(Pass, NotApplicable)).Reason(); reason != "" {
			t.Errorf("continuation active %t: got reason %q", continuationActive, reason)
		}
	}
}

// [utest -> req~agent-turn-validation~1]
// [utest -> req~agent-continuation-limit~1]
func TestReasonNamesEachFindingAndToolFailure(t *testing.T) {
	for _, continuationActive := range []bool{false, true} {
		reason := NewTurnDecision(stopEvent(continuationActive), report(Finding, Pass, ToolFailure)).Reason()

		for _, want := range []string{"a: finding CODE-finding: message of finding", "c: tool_failure CODE-tool_failure: message of tool_failure"} {
			if !strings.Contains(reason, want) {
				t.Errorf("continuation active %t: reason %q lacks %q", continuationActive, reason, want)
			}
		}
		if strings.Contains(reason, "CODE-pass") {
			t.Errorf("continuation active %t: reason %q names a passing result", continuationActive, reason)
		}
	}
}

// [utest -> req~agent-turn-validation~1]
// [utest -> req~agent-continuation-limit~1]
func TestNewFailedTurnDecisionNamesTheErrorCode(t *testing.T) {
	for continuationActive, want := range map[bool]TurnAction{false: ContinueTurn, true: EndTurn} {
		decision := NewFailedTurnDecision(stopEvent(continuationActive), "E-CAF-5")

		if decision.Action() != want || !strings.Contains(decision.Reason(), "E-CAF-5") {
			t.Errorf("continuation active %t: got %s %q, want %s with E-CAF-5", continuationActive, decision.Action(), decision.Reason(), want)
		}
		if decision.Report().Status() != ToolFailure {
			t.Errorf("continuation active %t: got report status %s", continuationActive, decision.Report().Status())
		}
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
