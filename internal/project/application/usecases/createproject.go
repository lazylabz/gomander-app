package usecases

import (
	"gomander/internal/project/domain"
)

type CreateProject interface {
	Execute(project domain.Project) error
}

type DefaultCreateProject struct {
	projectRepository domain.Repository
}

func NewCreateProject(projectRepo domain.Repository) *DefaultCreateProject {
	return &DefaultCreateProject{
		projectRepository: projectRepo,
	}
}

func (uc *DefaultCreateProject) Execute(project domain.Project) error {
	return uc.projectRepository.Create(project)
}
