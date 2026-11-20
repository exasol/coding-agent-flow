package openfasttrace

import (
	"encoding/xml"
	"fmt"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

// aspecReport is the XML report of the OpenFastTrace output format aspec. The XMLName field requires the root element
// specdocument, so other XML fails to read.
type aspecReport struct {
	XMLName     xml.Name `xml:"specdocument"`
	Collections []struct {
		Type  domain.OFTArtifactType `xml:"doctype,attr"`
		Items []struct {
			ID       string                   `xml:"id"`
			Revision int                      `xml:"version"`
			Needs    []domain.OFTArtifactType `xml:"coverage>needscoverage>needsobj"`
			// For a link to another revision, OFT lists the declared target and the existing revision.
			Covers []aspecLink `xml:"covering>coveredType"`
		} `xml:"specobject"`
	} `xml:"specobjects"`
}

// aspecLink is a coverage link to an OFT Specification Item.
type aspecLink struct {
	ID       string                 `xml:"id"`
	Revision int                    `xml:"version"`
	Type     domain.OFTArtifactType `xml:"doctype"`
}

// readASpecReport reads the OFT Specification Items with their OFT Coverage from an aspec report.
func readASpecReport(output []byte, result domain.OFTResult) (domain.OFTTrace, error) {
	var report aspecReport
	if err := xml.Unmarshal(output, &report); err != nil {
		return domain.OFTTrace{}, fmt.Errorf("read OFT trace: %w", err)
	}
	var items []domain.OFTSpecificationItem
	for _, collection := range report.Collections {
		for _, item := range collection.Items {
			id, err := domain.NewOFTSpecificationItemID(collection.Type, item.ID, item.Revision)
			if err != nil {
				return domain.OFTTrace{}, fmt.Errorf("read OFT trace: %w", err)
			}
			covers := make([]domain.OFTSpecificationItemID, 0, len(item.Covers))
			for _, link := range item.Covers {
				covered, err := domain.NewOFTSpecificationItemID(link.Type, link.ID, link.Revision)
				if err != nil {
					return domain.OFTTrace{}, fmt.Errorf("read OFT trace: coverage link of %s: %w", id, err)
				}
				covers = append(covers, covered)
			}
			items = append(items, domain.NewOFTSpecificationItem(id, domain.NewOFTCoverage(item.Needs, covers)))
		}
	}
	return domain.NewOFTTrace(items, result), nil
}
