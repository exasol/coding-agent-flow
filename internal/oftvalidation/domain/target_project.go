package domain

type TargetProject struct {
	path string
}

func NewTargetProject(path string) TargetProject {
	return TargetProject{path: path}
}

func (project TargetProject) Path() string {
	return project.path
}
