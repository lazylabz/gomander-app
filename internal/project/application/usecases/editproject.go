package usecases

import (
	"gomander/internal/project/domain"
)

type EditProject struct {
	projectRepository domain.Repository
}

func NewEditProject(projectRepo domain.Repository) *EditProject {
	return &EditProject{
		projectRepository: projectRepo,
	}
}

func (uc *EditProject) Execute(project domain.Project) error {
	return uc.projectRepository.Update(project)
}
