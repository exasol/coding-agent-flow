package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// OFTElementType is the type of an OFT Element, such as req or feat.
type OFTElementType string

const oftElementIDSeparator = "~"

// OFTElementID identifies an OFT Element formatted as type~name~revision.
type OFTElementID struct {
	elementType OFTElementType
	name        string
	revision    int
}

func NewOFTElementID(elementType OFTElementType, name string, revision int) (OFTElementID, error) {
	switch {
	case elementType == "" || strings.Contains(string(elementType), oftElementIDSeparator):
		return OFTElementID{}, fmt.Errorf("invalid OFT Element Type %q", elementType)
	case name == "" || strings.Contains(name, oftElementIDSeparator):
		return OFTElementID{}, fmt.Errorf("invalid OFT Element name %q", name)
	case revision < 0:
		return OFTElementID{}, fmt.Errorf("invalid OFT Element revision %d", revision)
	}
	return OFTElementID{elementType: elementType, name: name, revision: revision}, nil
}

// ParseOFTElementID reads an OFTElementID formatted as type~name~revision.
func ParseOFTElementID(text string) (OFTElementID, error) {
	parts := strings.Split(text, oftElementIDSeparator)
	if len(parts) != 3 {
		return OFTElementID{}, fmt.Errorf("OFT Element ID %q must be formatted as type~name~revision", text)
	}
	revision, err := strconv.Atoi(parts[2])
	if err != nil {
		return OFTElementID{}, fmt.Errorf("OFT Element ID %q has no numeric revision", text)
	}
	return NewOFTElementID(OFTElementType(parts[0]), parts[1], revision)
}

func (id OFTElementID) Type() OFTElementType {
	return id.elementType
}

func (id OFTElementID) String() string {
	return string(id.elementType) + oftElementIDSeparator + id.name + oftElementIDSeparator + strconv.Itoa(id.revision)
}

func (id OFTElementID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *OFTElementID) UnmarshalText(text []byte) error {
	parsed, err := ParseOFTElementID(string(text))
	if err != nil {
		return err
	}
	*id = parsed
	return nil
}
