package oftrequirementsdetector

import (
	"testing"

	"github.com/exasol/coding-agent-flow/test/acceptance"
	oftvalidation "github.com/exasol/coding-agent-flow/test/acceptance/oft-validation"
)

var scenarios = acceptance.NewSuite(
	oftvalidation.When(
		"the caller validates OFT tracing for the target project path",          // the When step
		"codingagentflow.NewOFTTracingValidator().Validate(ctx, targetProject)", // the public method it calls
	),
	checkRequirementFound, // compares the response with the scenario values in `at_data.toml`
)

// checkRequirementFound checks that the response contains req_id exactly when found is true.
func checkRequirementFound(response oftvalidation.OFTValidationResponse, given acceptance.Values, then acceptance.Values) error {
	reqID := given.String("req_id")
	return acceptance.Expect(reqID+" found", then.Bool("found"), response.HasSpecificationItemID(reqID))
}

// [at -> scn~oft-requirements-detection~1]
func TestOFTRequirementsDetector(t *testing.T) {
	// Disabled until CAF provides NewTargetProject and NewOFTTracingValidator.
	t.Skip("OFT validation public API is not implemented")
	scenarios.Run(t)
}
