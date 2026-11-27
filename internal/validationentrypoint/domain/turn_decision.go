package domain

import (
	"encoding/json"
	"strings"
)

// TurnAction tells the coding agent whether it may end its turn.
type TurnAction string

const (
	// EndTurn lets the coding agent end its turn.
	EndTurn TurnAction = "end_turn"
	// ContinueTurn asks the coding agent to continue its turn with the reason.
	ContinueTurn TurnAction = "continue_turn"
)

// TurnDecision is the result of the Agent Turn Validator: the action, the reason, and the Validation Report.
type TurnDecision struct {
	action TurnAction
	reason string
	report ValidationReport
}

// NewTurnDecision decides from the Validation Report: pass and not_applicable end the turn, finding and tool_failure
// continue it. The reason lists each finding and tool_failure result.
func NewTurnDecision(event StopEvent, report ValidationReport) TurnDecision {
	// [impl -> req~agent-turn-validation~1]
	var problems []string
	for _, result := range report.Results() {
		if status := result.Status(); status == Finding || status == ToolFailure {
			problems = append(problems, "- "+result.Validation()+": "+string(status)+" "+result.Code()+": "+result.Message())
		}
	}
	if report.Status() == Pass || report.Status() == NotApplicable {
		return TurnDecision{action: EndTurn, report: report}
	}
	reason := "CAF validation reports problems:\n" + strings.Join(problems, "\n")
	if len(problems) == 0 {
		reason = "CAF validation reports " + string(report.Status()) + "."
	}
	return TurnDecision{action: continueUnlessActive(event), reason: reason, report: report}
}

// NewFailedTurnDecision decides on a Validation Runner error with the given error code, such as E-CAF-5: it continues the
// turn, and the reason names the code.
func NewFailedTurnDecision(event StopEvent, errorCode string) TurnDecision {
	// [impl -> req~agent-turn-validation~1]
	reason := "CAF could not run the validations."
	if errorCode != "" {
		reason = "CAF could not run the validations: " + errorCode + "."
	}
	return TurnDecision{action: continueUnlessActive(event), reason: reason}
}

// continueUnlessActive returns ContinueTurn, or EndTurn when the turn already continued because of a Stop Hook, so CAF
// sends at most one Continuation Request per Agent Turn.
func continueUnlessActive(event StopEvent) TurnAction {
	// [impl -> req~agent-continuation-limit~1]
	if event.ContinuationActive() {
		return EndTurn
	}
	return ContinueTurn
}

// Action is end_turn or continue_turn; the zero value ends the turn.
func (decision TurnDecision) Action() TurnAction {
	if decision.action == "" {
		return EndTurn
	}
	return decision.action
}

// Reason names the Validation Name, status, code, and message of each finding and tool_failure result, or the error code
// of a Validation Runner error. It is empty when validation passes.
func (decision TurnDecision) Reason() string {
	return decision.reason
}

// Report is the Validation Report; after a Validation Runner error, it is the zero report with status tool_failure.
func (decision TurnDecision) Report() ValidationReport {
	return decision.report
}

// turnDecisionResponse is the JSON form of TurnDecision.
type turnDecisionResponse struct {
	Action TurnAction       `json:"action"`
	Reason string           `json:"reason"`
	Report ValidationReport `json:"report"`
}

// MarshalJSON writes the fields action, reason, and report.
func (decision TurnDecision) MarshalJSON() ([]byte, error) {
	return json.Marshal(turnDecisionResponse{Action: decision.Action(), Reason: decision.reason, Report: decision.report})
}
