package codingagentflow

import (
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/adapter/openfasttrace"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/application/service"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

type TargetProject = domain.TargetProject

// OFTElementType is the type of an OFT Element, such as req or feat.
type OFTElementType = domain.OFTElementType

// OFTElementID identifies an OFT Element formatted as type~name~revision.
type OFTElementID = domain.OFTElementID

type OFTValidation = domain.OFTValidation

type OFTTracingValidator = service.OFTTracingValidator

func NewTargetProject(path string) TargetProject {
	return domain.NewTargetProject(path)
}

// ParseOFTElementID reads an OFTElementID formatted as type~name~revision.
func ParseOFTElementID(text string) (OFTElementID, error) {
	return domain.ParseOFTElementID(text)
}

// NewOFTTracingValidator creates an OFTTracingValidator that uses OpenFastTrace.
func NewOFTTracingValidator() OFTTracingValidator {
	return service.NewOFTTracingValidator(openfasttrace.NewAdapter())
}
