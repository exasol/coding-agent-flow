package domain

import "slices"

// OFTResult is the result that OFT reports for a trace.
type OFTResult string

const (
	// OFTSuccess reports that OFT finds no tracing defect.
	OFTSuccess OFTResult = "success"
	// OFTFailure reports that OFT finds at least one tracing defect.
	OFTFailure OFTResult = "failure"
)

// OFTCoverage lists the OFT Artifact Types an OFT Specification Item needs and the OFT Specification Items it covers.
type OFTCoverage struct {
	needs  []OFTArtifactType
	covers []OFTSpecificationItemID
}

func NewOFTCoverage(needs []OFTArtifactType, covers []OFTSpecificationItemID) OFTCoverage {
	return OFTCoverage{needs: slices.Clone(needs), covers: slices.Clone(covers)}
}

func (coverage OFTCoverage) Needs() []OFTArtifactType {
	return slices.Clone(coverage.needs)
}

func (coverage OFTCoverage) Covers() []OFTSpecificationItemID {
	return slices.Clone(coverage.covers)
}

// OFTSpecificationItem is one OFT Specification Item of a trace with its OFT Coverage.
type OFTSpecificationItem struct {
	id       OFTSpecificationItemID
	coverage OFTCoverage
}

func NewOFTSpecificationItem(id OFTSpecificationItemID, coverage OFTCoverage) OFTSpecificationItem {
	return OFTSpecificationItem{id: id, coverage: coverage}
}

func (item OFTSpecificationItem) ID() OFTSpecificationItemID {
	return item.id
}

func (item OFTSpecificationItem) Coverage() OFTCoverage {
	return item.coverage
}

// OFTTrace holds the OFT Specification Items of a Target Project with their OFT Coverage, and the OFT Result.
type OFTTrace struct {
	items  []OFTSpecificationItem
	result OFTResult
}

func NewOFTTrace(items []OFTSpecificationItem, result OFTResult) OFTTrace {
	return OFTTrace{items: slices.Clone(items), result: result}
}

func (trace OFTTrace) Items() []OFTSpecificationItem {
	return slices.Clone(trace.items)
}

func (trace OFTTrace) Result() OFTResult {
	return trace.result
}

// Validation derives the OFT Validation from the OFT Specification Items. Tracing is valid when:
//   - each OFT Specification Item ID is unique,
//   - each coverage link reaches exactly one OFT Specification Item, in the same revision,
//   - each covered OFT Specification Item needs the OFT Artifact Type of the covering one,
//   - each needed OFT Artifact Type has a covering OFT Specification Item.
func (trace OFTTrace) Validation() OFTValidation {
	// [impl -> req~oft-tracing-validation~1]
	ids := make([]OFTSpecificationItemID, 0, len(trace.items))
	byID := make(map[OFTSpecificationItemID][]OFTSpecificationItem, len(trace.items))
	for _, item := range trace.items {
		ids = append(ids, item.id)
		byID[item.id] = append(byID[item.id], item)
	}
	return NewOFTValidation(trace.uniqueIDs(byID) && trace.wantedLinks(byID) && trace.coveredNeeds(), ids)
}

func (trace OFTTrace) uniqueIDs(byID map[OFTSpecificationItemID][]OFTSpecificationItem) bool {
	return len(byID) == len(trace.items)
}

func (trace OFTTrace) wantedLinks(byID map[OFTSpecificationItemID][]OFTSpecificationItem) bool {
	for _, item := range trace.items {
		for _, covered := range item.coverage.covers {
			targets := byID[covered]
			if len(targets) != 1 || !slices.Contains(targets[0].coverage.needs, item.id.Type()) {
				return false
			}
		}
	}
	return true
}

func (trace OFTTrace) coveredNeeds() bool {
	for _, item := range trace.items {
		for _, needed := range item.coverage.needs {
			covered := slices.ContainsFunc(trace.items, func(covering OFTSpecificationItem) bool {
				return covering.id.Type() == needed && slices.Contains(covering.coverage.covers, item.id)
			})
			if !covered {
				return false
			}
		}
	}
	return true
}
