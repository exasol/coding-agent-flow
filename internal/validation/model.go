// Package validation defines and executes provider-independent CAF validators.
package validation

// Status classifies a deterministic validation outcome.
type Status string

const (
	StatusPass          Status = "pass"
	StatusFinding       Status = "finding"
	StatusToolFailure   Status = "tool_failure"
	StatusNotApplicable Status = "not_applicable"
)

// Result is the outcome of one validator.
type Result struct {
	Validator string `json:"validator"`
	Status    Status `json:"status"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

// Report aggregates all validators selected by project policy.
type Report struct {
	Status  Status   `json:"status"`
	Results []Result `json:"results"`
}

// Action is a provider-independent decision about the current agent turn.
type Action string

const (
	AllowStop    Action = "allow_stop"
	ContinueWork Action = "continue_work"
)

// Decision is deliberately separate from validation results.
type Decision struct {
	Action Action
	Reason string
}

// NewReport derives the aggregate status using fail-safe precedence.
func NewReport(results []Result) Report {
	status := StatusNotApplicable
	for _, result := range results {
		switch result.Status {
		case StatusToolFailure:
			status = StatusToolFailure
		case StatusFinding:
			if status != StatusToolFailure {
				status = StatusFinding
			}
		case StatusPass:
			if status == StatusNotApplicable {
				status = StatusPass
			}
		case StatusNotApplicable:
			// Keep the current aggregate status.
		default:
			status = StatusToolFailure
		}
	}
	return Report{Status: status, Results: results}
}

// Decide maps a report to lifecycle intent while preventing continuation loops.
func Decide(report Report, alreadyContinued bool) Decision {
	if alreadyContinued {
		return Decision{Action: AllowStop}
	}
	if report.Status != StatusFinding && report.Status != StatusToolFailure {
		return Decision{Action: AllowStop}
	}
	for _, result := range report.Results {
		if result.Status == StatusToolFailure || result.Status == StatusFinding {
			return Decision{Action: ContinueWork, Reason: result.Code + ": " + result.Message}
		}
	}
	return Decision{Action: ContinueWork, Reason: "CAF validation requires more work."}
}
