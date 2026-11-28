package agentcontinuationlimiter

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	validationentrypoint "github.com/exasol/coding-agent-flow/test/acceptance/validation-entry-point"
)

// [at -> scn~agent-continuation-limiter~1]
func TestAgentContinuationLimiter(t *testing.T) {
	validationentrypoint.Stop(t, validationentrypoint.Project(false, "oft-tracing-validator"), true,
		validationentrypoint.ExpectDecision(codingagentflow.EndTurn, "CAF-OFT-TRACING-INVALID"))
}
