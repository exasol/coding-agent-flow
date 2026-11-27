package service

import (
	"context"

	"github.com/exasol/coding-agent-flow/internal/validationentrypoint/domain"
)

// AgentTurnValidator validates the Target Project of a Stop Event and returns a Turn Decision. It keeps no state between
// calls.
type AgentTurnValidator struct {
	runner ValidationRunner
}

// NewAgentTurnValidator creates an AgentTurnValidator that validates with the runner.
func NewAgentTurnValidator(runner ValidationRunner) AgentTurnValidator {
	return AgentTurnValidator{runner: runner}
}

// Validate runs the Validation Runner on the Target Project of the event and decides from the Validation Report, or from
// the error code of a Validation Runner error. It returns an error only when ctx ends.
func (validator AgentTurnValidator) Validate(ctx context.Context, event domain.StopEvent) (domain.TurnDecision, error) {
	// [impl -> req~agent-turn-validation~1]
	report, err := validator.runner.Run(ctx, event.TargetProject())
	if ctxErr := ctx.Err(); ctxErr != nil {
		return domain.TurnDecision{}, ctxErr
	}
	if err != nil {
		return domain.NewFailedTurnDecision(event, errorCode(err)), nil
	}
	return domain.NewTurnDecision(event, report), nil
}
