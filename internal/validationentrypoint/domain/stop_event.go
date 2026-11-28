package domain

import "github.com/exasol/coding-agent-flow/internal/global"

// StopEvent is what CAF receives from a Stop Hook: the Target Project and the Continuation Active flag.
type StopEvent struct {
	targetProject      global.TargetProject
	continuationActive bool
}

// NewStopEvent creates a StopEvent. continuationActive reports that the Agent Turn already continued because of a Stop
// Hook.
func NewStopEvent(targetProject global.TargetProject, continuationActive bool) StopEvent {
	return StopEvent{targetProject: targetProject, continuationActive: continuationActive}
}

func (event StopEvent) TargetProject() global.TargetProject {
	return event.targetProject
}

// ContinuationActive reports that the Agent Turn already continued because of a Stop Hook.
func (event StopEvent) ContinuationActive() bool {
	return event.continuationActive
}
