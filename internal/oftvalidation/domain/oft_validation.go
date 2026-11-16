package domain

import (
	"encoding/json"
	"fmt"
	"slices"
)

// OFTValidation is the result of validating a Target Project. It reports whether tracing is valid and lists the
// OFT Element IDs by type.
type OFTValidation struct {
	valid      bool
	elementIDs map[OFTElementType][]OFTElementID
}

func NewOFTValidation(valid bool, elementIDs []OFTElementID) OFTValidation {
	grouped := make(map[OFTElementType][]OFTElementID)
	for _, id := range elementIDs {
		grouped[id.Type()] = append(grouped[id.Type()], id)
	}
	return OFTValidation{valid: valid, elementIDs: grouped}
}

func (validation OFTValidation) Valid() bool {
	return validation.valid
}

func (validation OFTValidation) Contains(id OFTElementID) bool {
	return slices.Contains(validation.elementIDs[id.Type()], id)
}

// oftValidationResponse is the JSON form of OFTValidation.
type oftValidationResponse struct {
	Valid         bool                              `json:"valid"`
	OFTElementIDs map[OFTElementType][]OFTElementID `json:"oftElementIDs"`
}

// MarshalJSON writes the fields valid and oftElementIDs.
func (validation OFTValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(oftValidationResponse{Valid: validation.valid, OFTElementIDs: validation.elementIDs})
}

// UnmarshalJSON reads the fields valid and oftElementIDs. Each ID must have the type it is listed under.
func (validation *OFTValidation) UnmarshalJSON(data []byte) error {
	var response oftValidationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	var ids []OFTElementID
	for elementType, typeIDs := range response.OFTElementIDs {
		for _, id := range typeIDs {
			if id.Type() != elementType {
				return fmt.Errorf("OFT Element ID %s is listed under type %s", id, elementType)
			}
			ids = append(ids, id)
		}
	}
	*validation = NewOFTValidation(response.Valid, ids)
	return nil
}
