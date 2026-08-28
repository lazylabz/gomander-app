package usecases

import (
	"gomander/internal/project/domain"
)

type CreateProject struct {
	projectRepository domain.Repository
}

func NewCreateProject(projectRepo domain.Repository) *CreateProject {
	return &CreateProject{
		projectRepository: projectRepo,
	}
}

func (uc *CreateProject) Execute(project domain.Project) error {
	return uc.projectRepository.Create(project)
}
