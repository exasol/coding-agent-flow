package acceptance

import (
	"testing"
)

var oftTracingValidatorProgram = newTargetProjectProgram(
	"codingagentflow.NewOFTTracingValidator().Validate(context.Background(), codingagentflow.NewTargetProject(os.Args[1]))",
)

var oftTracingValidatorLocation = useCaseLocation{
	feature:     "oft-validation",
	requirement: "oft-tracing-validation",
	useCase:     "oft-tracing-validator",
}

var oftTracingValidationScenarios = newScenarioSuite(
	oftTracingValidatorLocation.scenarioCatalog(),
	oftTracingValidatorLocation.fixtureCatalog(),
	oftTracingValidatorProgram,
	decodeOFTValidation,
	matchOFTValidation,
)

func init() {
	oftTracingValidationScenarios.bind("scn~oft-requirement-covers-feature~1", oftValidationResponse{Valid: true})
	oftTracingValidationScenarios.bind("scn~oft-requirement-does-not-cover-feature~1", oftValidationResponse{Valid: false})
}

// [at -> scn~oft-requirement-covers-feature~1]
// [at -> scn~oft-requirement-does-not-cover-feature~1]
func TestOFTTracingValidationAcceptanceScenarios(t *testing.T) {
	oftTracingValidationScenarios.run(t)
}
