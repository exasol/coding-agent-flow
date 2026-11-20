package codingagentflow

import (
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/adapter/openfasttrace"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/application/service"
	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

// OFTArtifactType is the type of an OFT Specification Item, such as req or feat.
type OFTArtifactType = domain.OFTArtifactType

// OFTSpecificationItemID identifies an OFT Specification Item formatted as type~name~revision.
type OFTSpecificationItemID = domain.OFTSpecificationItemID

type OFTValidation = domain.OFTValidation

type OFTTracingValidator = service.OFTTracingValidator

// ParseOFTSpecificationItemID reads an OFTSpecificationItemID formatted as type~name~revision.
func ParseOFTSpecificationItemID(text string) (OFTSpecificationItemID, error) {
	return domain.ParseOFTSpecificationItemID(text)
}

// NewOFTTracingValidator creates an OFTTracingValidator that uses OpenFastTrace.
func NewOFTTracingValidator() OFTTracingValidator {
	return service.NewOFTTracingValidator(openfasttrace.NewAdapter())
}
