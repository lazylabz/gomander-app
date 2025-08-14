package domain

type Repository interface {
	GetAllProjects() ([]*Project, error)
	GetProjectById(id string) (*Project, error)
	CreateProject(project Project) error
	UpdateProject(project Project) error
	DeleteProject(id string) error
}
