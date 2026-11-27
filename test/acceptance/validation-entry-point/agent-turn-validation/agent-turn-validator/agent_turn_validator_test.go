package agentturnvalidator

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	validationentrypoint "github.com/exasol/coding-agent-flow/test/acceptance/validation-entry-point"
)

// [at -> scn~agent-turn-validator_end-turn~1]
func TestAgentTurnValidator_endTurn(t *testing.T) {
	validationentrypoint.Stop(t, validationentrypoint.Project(true, "oft-tracing-validator"), false,
		validationentrypoint.ExpectDecision(codingagentflow.EndTurn, ""))
}

// [at -> scn~agent-turn-validator_continue-turn~1]
func TestAgentTurnValidator_continueTurn(t *testing.T) {
	validationentrypoint.Stop(t, validationentrypoint.Project(false, "oft-tracing-validator"), false,
		validationentrypoint.ExpectDecision(codingagentflow.ContinueTurn, "CAF-OFT-TRACING-INVALID"))
	validationentrypoint.Stop(t, validationentrypoint.Project(true, "unknown-check"), false,
		validationentrypoint.ExpectDecision(codingagentflow.ContinueTurn, "CAF-UNKNOWN-VALIDATOR"))
}
