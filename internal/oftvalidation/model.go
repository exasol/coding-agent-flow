// Package oftvalidation contains the OFT tracing validation model.
package oftvalidation

// TargetProjectPath identifies a Target Project for one validation operation.
type TargetProjectPath string

// TargetProject is the project selected for OFT tracing validation.
type TargetProject struct {
	path TargetProjectPath
}

// NewTargetProject creates a Target Project from its caller-supplied path.
func NewTargetProject(path string) TargetProject {
	return TargetProject{path: TargetProjectPath(path)}
}

// Path returns the Target Project path.
func (project TargetProject) Path() string {
	return string(project.path)
}

// OFTElementType identifies an OFT element type.
type OFTElementType string

// OFTElementID identifies one OFT Element.
type OFTElementID string

// OFTElementTypeRequirement identifies OFT Requirement elements.
const OFTElementTypeRequirement OFTElementType = "req"

// OFTValidation is the immutable result of validating a Target Project with OFT.
type OFTValidation struct {
	Valid         bool                              `json:"valid"`
	OFTElementIDs map[OFTElementType][]OFTElementID `json:"oftElementIDs"`
}
