package codingagentflow

import (
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/application/service"
	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// StopEvent is what CAF receives from a Stop Hook: the Target Project and the Continuation Active flag.
type StopEvent = domain.StopEvent

// TurnDecision is the result of the Agent Turn Validator: the action, the reason, and the Validation Report.
type TurnDecision = domain.TurnDecision

// TurnAction is end_turn or continue_turn.
type TurnAction = domain.TurnAction

const (
	EndTurn      = domain.EndTurn
	ContinueTurn = domain.ContinueTurn
)

// AgentTurnValidator validates the Target Project of a StopEvent and returns a TurnDecision.
type AgentTurnValidator = service.AgentTurnValidator

// NewStopEvent creates a StopEvent. continuationActive reports that the Agent Turn already continued because of a Stop Hook.
func NewStopEvent(targetProject TargetProject, continuationActive bool) StopEvent {
	return domain.NewStopEvent(targetProject, continuationActive)
}

// NewAgentTurnValidator creates an AgentTurnValidator that validates with the ValidationRunner of NewValidationRunner.
func NewAgentTurnValidator() (AgentTurnValidator, error) {
	runner, err := NewValidationRunner()
	if err != nil {
		return AgentTurnValidator{}, err
	}
	return service.NewAgentTurnValidator(runner), nil
}
