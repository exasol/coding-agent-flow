package domain

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
)

// LinkValidation is the result of the Link Checker. It is valid when all Links resolve. The zero value is invalid,
// so the result that accompanies an error is invalid.
type LinkValidation struct {
	valid       bool
	brokenLinks []BrokenLink
}

// NewLinkValidation creates a LinkValidation with the Broken Links sorted by source, line, and target.
func NewLinkValidation(brokenLinks []BrokenLink) LinkValidation {
	sorted := slices.Clone(brokenLinks)
	slices.SortFunc(sorted, compareBrokenLinks)
	return LinkValidation{valid: len(sorted) == 0, brokenLinks: sorted}
}

func compareBrokenLinks(first, second BrokenLink) int {
	return cmp.Or(
		cmp.Compare(first.link.source, second.link.source),
		cmp.Compare(first.link.line, second.link.line),
		cmp.Compare(first.link.target, second.link.target),
		cmp.Compare(first.status, second.status),
	)
}

func (validation LinkValidation) Valid() bool {
	return validation.valid
}

func (validation LinkValidation) BrokenLinks() []BrokenLink {
	return slices.Clone(validation.brokenLinks)
}

// ContainsBrokenLink reports whether the document at source has a Broken Link to target.
func (validation LinkValidation) ContainsBrokenLink(source, target string) bool {
	return slices.ContainsFunc(validation.brokenLinks, func(brokenLink BrokenLink) bool {
		return brokenLink.link.source == source && brokenLink.link.target == target
	})
}

// brokenLinkResponse is the JSON form of BrokenLink.
type brokenLinkResponse struct {
	Source string `json:"source"`
	Line   int    `json:"line"`
	Target string `json:"target"`
	Status string `json:"status"`
}

// linkValidationResponse is the JSON form of LinkValidation.
type linkValidationResponse struct {
	Valid       bool                 `json:"valid"`
	BrokenLinks []brokenLinkResponse `json:"brokenLinks"`
}

// MarshalJSON writes the fields valid and brokenLinks.
func (validation LinkValidation) MarshalJSON() ([]byte, error) {
	brokenLinks := make([]brokenLinkResponse, 0, len(validation.brokenLinks))
	for _, brokenLink := range validation.brokenLinks {
		brokenLinks = append(brokenLinks, brokenLinkResponse{
			Source: brokenLink.link.source,
			Line:   brokenLink.link.line,
			Target: brokenLink.link.target,
			Status: brokenLink.status,
		})
	}
	return json.Marshal(linkValidationResponse{Valid: validation.Valid(), BrokenLinks: brokenLinks})
}

// UnmarshalJSON reads the fields valid and brokenLinks. valid must be true exactly when brokenLinks is empty.
func (validation *LinkValidation) UnmarshalJSON(data []byte) error {
	var response linkValidationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	if response.Valid != (len(response.BrokenLinks) == 0) {
		return fmt.Errorf("LinkValidation reports valid %t with %d Broken Links", response.Valid, len(response.BrokenLinks))
	}
	brokenLinks := make([]BrokenLink, 0, len(response.BrokenLinks))
	for _, brokenLink := range response.BrokenLinks {
		brokenLinks = append(brokenLinks, NewBrokenLink(NewLink(brokenLink.Source, brokenLink.Line, brokenLink.Target), brokenLink.Status))
	}
	*validation = NewLinkValidation(brokenLinks)
	return nil
}
