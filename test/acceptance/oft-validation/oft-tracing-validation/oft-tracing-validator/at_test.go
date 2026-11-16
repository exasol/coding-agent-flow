package ofttracingvalidator

import (
	"testing"

	codingagentflow "github.com/exasol/coding-agent-flow"
	"github.com/exasol/coding-agent-flow/test/acceptance"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

var scenarios = acceptance.NewSuite(
	oftvalidation.When(
		"the caller validates OFT tracing for the target project path",          // the When step
		"codingagentflow.NewOFTTracingValidator().Validate(ctx, targetProject)", // the public method it calls
	),
	checkValid, // compares the response with the scenario values in `at_data.toml`
)

func checkValid(response codingagentflow.OFTValidation, _ acceptance.Values, then acceptance.Values) error {
	return acceptance.Expect("valid", then.Bool("valid"), response.Valid())
}

// [at -> scn~oft-requirement-covers-feature~1]
// [at -> scn~oft-requirement-does-not-cover-feature~1]
func TestOFTTracingValidator(t *testing.T) {
	scenarios.Run(t)
}
