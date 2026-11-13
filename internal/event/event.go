// Package event defines provider-independent CAF lifecycle events.
package event

// Kind identifies a lifecycle event understood by CAF policy.
type Kind string

const (
	// TurnStop indicates that a coding agent is about to finish its turn.
	TurnStop Kind = "turn_stop"
	// CIValidation indicates an independent validation run in CI.
	CIValidation Kind = "ci_validation"
)

// Event contains only provider-independent lifecycle data.
type Event struct {
	Kind             Kind
	WorkingDirectory string
	SessionID        string
	AlreadyContinued bool
}
