package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// OFTArtifactType is the type of an OFT Specification Item, such as req or feat.
type OFTArtifactType string

const oftSpecificationItemIDSeparator = "~"

// OFTSpecificationItemID identifies an OFT Specification Item formatted as type~name~revision.
type OFTSpecificationItemID struct {
	artifactType OFTArtifactType
	name         string
	revision     int
}

func NewOFTSpecificationItemID(artifactType OFTArtifactType, name string, revision int) (OFTSpecificationItemID, error) {
	switch {
	case artifactType == "" || strings.Contains(string(artifactType), oftSpecificationItemIDSeparator):
		return OFTSpecificationItemID{}, fmt.Errorf("invalid OFT Artifact Type %q", artifactType)
	case name == "" || strings.Contains(name, oftSpecificationItemIDSeparator):
		return OFTSpecificationItemID{}, fmt.Errorf("invalid OFT Specification Item name %q", name)
	case revision < 0:
		return OFTSpecificationItemID{}, fmt.Errorf("invalid OFT Specification Item revision %d", revision)
	}
	return OFTSpecificationItemID{artifactType: artifactType, name: name, revision: revision}, nil
}

// ParseOFTSpecificationItemID reads an OFTSpecificationItemID formatted as type~name~revision.
func ParseOFTSpecificationItemID(text string) (OFTSpecificationItemID, error) {
	parts := strings.Split(text, oftSpecificationItemIDSeparator)
	if len(parts) != 3 {
		return OFTSpecificationItemID{}, fmt.Errorf("OFT Specification Item ID %q must be formatted as type~name~revision", text)
	}
	revision, err := strconv.Atoi(parts[2])
	if err != nil {
		return OFTSpecificationItemID{}, fmt.Errorf("OFT Specification Item ID %q has no numeric revision", text)
	}
	return NewOFTSpecificationItemID(OFTArtifactType(parts[0]), parts[1], revision)
}

func (id OFTSpecificationItemID) Type() OFTArtifactType {
	return id.artifactType
}

func (id OFTSpecificationItemID) String() string {
	return string(id.artifactType) + oftSpecificationItemIDSeparator + id.name + oftSpecificationItemIDSeparator + strconv.Itoa(id.revision)
}

func (id OFTSpecificationItemID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *OFTSpecificationItemID) UnmarshalText(text []byte) error {
	parsed, err := ParseOFTSpecificationItemID(string(text))
	if err != nil {
		return err
	}
	*id = parsed
	return nil
}
