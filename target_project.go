package codingagentflow

import "github.com/exasol/coding-agent-flow/internal/global"

// TargetProject is the read-only input to a CAF operation.
type TargetProject = global.TargetProject

// NewTargetProject creates a TargetProject for the Target Project Path.
func NewTargetProject(path string) TargetProject {
	return global.NewTargetProject(path)
}
