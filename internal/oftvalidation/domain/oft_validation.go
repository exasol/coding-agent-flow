package domain

import (
	"encoding/json"
	"fmt"
	"slices"
)

// OFTValidation is the result of validating a Target Project. It reports whether tracing is valid and lists the
// OFT Specification Item IDs by artifact type.
type OFTValidation struct {
	valid                bool
	specificationItemIDs map[OFTArtifactType][]OFTSpecificationItemID
}

func NewOFTValidation(valid bool, specificationItemIDs []OFTSpecificationItemID) OFTValidation {
	grouped := make(map[OFTArtifactType][]OFTSpecificationItemID)
	for _, id := range specificationItemIDs {
		grouped[id.Type()] = append(grouped[id.Type()], id)
	}
	return OFTValidation{valid: valid, specificationItemIDs: grouped}
}

func (validation OFTValidation) Valid() bool {
	return validation.valid
}

func (validation OFTValidation) Contains(id OFTSpecificationItemID) bool {
	return slices.Contains(validation.specificationItemIDs[id.Type()], id)
}

// oftValidationResponse is the JSON form of OFTValidation.
type oftValidationResponse struct {
	Valid                   bool                                         `json:"valid"`
	OFTSpecificationItemIDs map[OFTArtifactType][]OFTSpecificationItemID `json:"oftSpecificationItemIDs"`
}

// MarshalJSON writes the fields valid and oftSpecificationItemIDs.
func (validation OFTValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(oftValidationResponse{Valid: validation.valid, OFTSpecificationItemIDs: validation.specificationItemIDs})
}

// UnmarshalJSON reads the fields valid and oftSpecificationItemIDs. Each ID must have the artifact type it is listed under.
func (validation *OFTValidation) UnmarshalJSON(data []byte) error {
	var response oftValidationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	var ids []OFTSpecificationItemID
	for artifactType, typeIDs := range response.OFTSpecificationItemIDs {
		for _, id := range typeIDs {
			if id.Type() != artifactType {
				return fmt.Errorf("OFT Specification Item ID %s is listed under artifact type %s", id, artifactType)
			}
			ids = append(ids, id)
		}
	}
	*validation = NewOFTValidation(response.Valid, ids)
	return nil
}
