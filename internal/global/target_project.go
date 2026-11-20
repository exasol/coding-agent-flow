// Package global holds the domain classes shared by all CAF domains, such as the Target Project.
package global

// TargetProject is the read-only input to a CAF operation.
type TargetProject struct {
	path string
}

func NewTargetProject(path string) TargetProject {
	return TargetProject{path: path}
}

func (project TargetProject) Path() string {
	return project.path
}
