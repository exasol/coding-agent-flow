// Package turndecision checks the scenarios of turn_decision.feature on TurnDecision, without the service.
package turndecision

import (
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/global"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// report creates a Validation Report with one result per status, named a, b, c, and so on.
func report(statuses ...domain.ValidationStatus) domain.ValidationReport {
	var results []domain.ValidationResult
	for index, status := range statuses {
		results = append(results, domain.NewValidationResult(string(rune('a'+index)), status, "CODE-"+string(status), "message of "+string(status)))
	}
	return domain.NewValidationReport(results)
}

// turnDecisionCase decides from the report, or from the error code when it is set.
type turnDecisionCase struct {
	name               string
	report             domain.ValidationReport
	errorCode          string
	continuationActive bool
	want               domain.TurnAction
	reasonContains     []string
	reasonLacks        []string
}

func (test turnDecisionCase) decide() domain.TurnDecision {
	event := domain.NewStopEvent(global.NewTargetProject("project"), test.continuationActive)
	if test.errorCode != "" {
		return domain.NewFailedTurnDecision(event, test.errorCode)
	}
	return domain.NewTurnDecision(event, test.report)
}

// Each row checks one example of a scenario.
var turnDecisionCases = []turnDecisionCase{
	// [at -> scn~domain-turn-decision_end-turn-when-validation-passes~1]
	{name: "pass", report: report(domain.Pass), want: domain.EndTurn},
	// [at -> scn~domain-turn-decision_end-turn-when-validation-passes~1]
	{name: "not applicable", report: report(domain.NotApplicable), want: domain.EndTurn},
	// [at -> scn~domain-turn-decision_continue-turn-on-problems~1]
	{name: "finding and tool failure", report: report(domain.Finding, domain.Pass, domain.ToolFailure), want: domain.ContinueTurn,
		reasonContains: []string{"a: finding CODE-finding: message of finding", "c: tool_failure CODE-tool_failure: message of tool_failure"},
		reasonLacks:    []string{"CODE-pass"}},
	// [at -> scn~domain-turn-decision_continue-turn-on-validation-error~1]
	{name: "validation runner error", errorCode: "E-CAF-5", want: domain.ContinueTurn, reasonContains: []string{"E-CAF-5"}},
	// [at -> scn~domain-turn-decision_end-turn-when-continuation-active~1]
	{name: "pass with continuation active", report: report(domain.Pass), continuationActive: true, want: domain.EndTurn},
	// [at -> scn~domain-turn-decision_end-turn-when-continuation-active~1]
	{name: "finding with continuation active", report: report(domain.Finding), continuationActive: true, want: domain.EndTurn,
		reasonContains: []string{"a: finding"}},
	// [at -> scn~domain-turn-decision_end-turn-when-continuation-active~1]
	{name: "validation runner error with continuation active", errorCode: "E-CAF-5", continuationActive: true, want: domain.EndTurn,
		reasonContains: []string{"E-CAF-5"}},
}

func TestTurnDecision(t *testing.T) {
	for _, test := range turnDecisionCases {
		t.Run(test.name, func(t *testing.T) {
			decision := test.decide()

			if decision.Action() != test.want {
				t.Errorf("action: got %s, want %s", decision.Action(), test.want)
			}
			if len(test.reasonContains) == 0 && decision.Reason() != "" {
				t.Errorf("reason: got %q, want none", decision.Reason())
			}
			for _, text := range test.reasonContains {
				if !strings.Contains(decision.Reason(), text) {
					t.Errorf("reason %q lacks %q", decision.Reason(), text)
				}
			}
			for _, text := range test.reasonLacks {
				if strings.Contains(decision.Reason(), text) {
					t.Errorf("reason %q contains %q", decision.Reason(), text)
				}
			}
		})
	}
}
