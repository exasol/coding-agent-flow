//go:build oftvalidation

// OFT acceptance tests are disabled until CAF provides NewTargetProject and NewOFTTracingValidator.
package oftvalidation

import (
	"context"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
)

// When invokes the OFT tracing validator directly.
func When(ctx context.Context, directory string, _ acceptance.Values) (codingagentflow.OFTValidation, error) {
	targetProject := codingagentflow.NewTargetProject(directory)
	return codingagentflow.NewOFTTracingValidator().Validate(ctx, targetProject)
}
