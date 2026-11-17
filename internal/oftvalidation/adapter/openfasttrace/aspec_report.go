package openfasttrace

import (
	"encoding/xml"
	"fmt"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

// aspecReport is the XML report of the OpenFastTrace output format aspec.
type aspecReport struct {
	Collections []struct {
		Type  domain.OFTArtifactType `xml:"doctype,attr"`
		Items []struct {
			ID       string `xml:"id"`
			Revision int    `xml:"version"`
		} `xml:"specobject"`
	} `xml:"specobjects"`
}

// readASpecReport reads the specification item IDs from an aspec report.
func readASpecReport(output []byte, valid bool) (domain.OFTValidation, error) {
	var report aspecReport
	if err := xml.Unmarshal(output, &report); err != nil {
		return domain.OFTValidation{}, fmt.Errorf("read OFT trace: %w", err)
	}
	var ids []domain.OFTSpecificationItemID
	for _, collection := range report.Collections {
		for _, item := range collection.Items {
			id, err := domain.NewOFTSpecificationItemID(collection.Type, item.ID, item.Revision)
			if err != nil {
				return domain.OFTValidation{}, fmt.Errorf("read OFT trace: %w", err)
			}
			ids = append(ids, id)
		}
	}
	return domain.NewOFTValidation(valid, ids), nil
}
