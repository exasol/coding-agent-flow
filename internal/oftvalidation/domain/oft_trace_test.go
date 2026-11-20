package domain

import "testing"

// An OFT Specification Item without needed OFT Artifact Types is valid without coverage.
func TestOFTTraceValidationAcceptsAnItemWithoutNeeds(t *testing.T) {
	feature, err := ParseOFTSpecificationItemID("feat~f~1")
	if err != nil {
		t.Fatal(err)
	}

	validation := NewOFTTrace([]OFTSpecificationItem{NewOFTSpecificationItem(feature, NewOFTCoverage(nil, nil))}, OFTSuccess).Validation()

	if !validation.Valid() {
		t.Fatal("expected valid tracing")
	}
}
