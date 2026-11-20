// Package domain holds the Validation Entry Point domain classes and their rules.
package domain

// ValidationStatus is the Validation Result Status of one Validation.
type ValidationStatus string

const (
	// Pass reports that the Validation completed and accepts the Target Project.
	Pass ValidationStatus = "pass"
	// Finding reports that the Validation completed and found a problem.
	Finding ValidationStatus = "finding"
	// ToolFailure reports that the Validation could not complete reliably.
	ToolFailure ValidationStatus = "tool_failure"
	// NotApplicable reports that the Validation applies to other Target Projects.
	NotApplicable ValidationStatus = "not_applicable"
)

// severity orders the statuses from not_applicable to tool_failure. Any other status counts as tool_failure.
func (status ValidationStatus) severity() int {
	switch status {
	case NotApplicable:
		return 0
	case Pass:
		return 1
	case Finding:
		return 2
	default:
		return 3
	}
}

// supported returns the status, or ToolFailure for any other status.
func (status ValidationStatus) supported() ValidationStatus {
	switch status {
	case Pass, Finding, ToolFailure, NotApplicable:
		return status
	default:
		return ToolFailure
	}
}
