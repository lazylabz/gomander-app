package usecases

import (
	"gomander/internal/project/domain"
)

type EditProject interface {
	Execute(project domain.Project) error
}

type DefaultEditProject struct {
	projectRepository domain.Repository
}

func NewEditProject(projectRepo domain.Repository) *DefaultEditProject {
	return &DefaultEditProject{
		projectRepository: projectRepo,
	}
}

func (uc *DefaultEditProject) Execute(project domain.Project) error {
	return uc.projectRepository.Update(project)
}
