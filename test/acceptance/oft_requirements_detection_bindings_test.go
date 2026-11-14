package acceptance

import (
	"testing"
)

var oftRequirementsDetectorProgram = newTargetProjectProgram(
	"codingagentflow.NewOFTTracingValidator().Validate(context.Background(), codingagentflow.NewTargetProject(os.Args[1]))",
)

var oftRequirementsDetectorLocation = useCaseLocation{
	feature:     "oft-validation",
	requirement: "oft-requirements-detection",
	useCase:     "oft-requirements-detector",
}

var oftRequirementsDetectionScenarios = newScenarioSuite(
	oftRequirementsDetectorLocation.scenarioCatalog(),
	oftRequirementsDetectorLocation.fixtureCatalog(),
	oftRequirementsDetectorProgram,
	decodeOFTValidation,
	matchOFTValidation,
)

func init() {
	oftRequirementsDetectionScenarios.bind("scn~oft-requirements-rst-in-doc-requirements~1", oftValidationResponse{Valid: true, fixtureElementIDTypes: []oftElementType{oftRequirementElementType}})
	oftRequirementsDetectionScenarios.bind("scn~oft-requirements-rst-in-root~1", oftValidationResponse{Valid: true, fixtureElementIDTypes: []oftElementType{oftRequirementElementType}})
	oftRequirementsDetectionScenarios.bind("scn~oft-requirements-markdown-in-root~1", oftValidationResponse{Valid: true, fixtureElementIDTypes: []oftElementType{oftRequirementElementType}})
	oftRequirementsDetectionScenarios.bind("scn~oft-requirements-text-in-doc-requirements~1", oftValidationResponse{Valid: true, OFTElementIDs: map[oftElementType][]oftElementID{oftRequirementElementType: {}}})
	oftRequirementsDetectionScenarios.bind("scn~oft-requirements-text-in-root~1", oftValidationResponse{Valid: true, OFTElementIDs: map[oftElementType][]oftElementID{oftRequirementElementType: {}}})
}

// [at -> scn~oft-requirements-rst-in-doc-requirements~1]
// [at -> scn~oft-requirements-rst-in-root~1]
// [at -> scn~oft-requirements-markdown-in-root~1]
// [at -> scn~oft-requirements-text-in-doc-requirements~1]
// [at -> scn~oft-requirements-text-in-root~1]
func TestOFTRequirementsDetectionAcceptanceScenarios(t *testing.T) {
	oftRequirementsDetectionScenarios.run(t)
}
