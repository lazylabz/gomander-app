package usecases

import (
	"gomander/internal/project/domain"
)

type GetAvailableProjects struct {
	projectRepository domain.Repository
}

func NewGetAvailableProjects(projectRepo domain.Repository) *GetAvailableProjects {
	return &GetAvailableProjects{
		projectRepository: projectRepo,
	}
}

func (uc *GetAvailableProjects) Execute() ([]domain.Project, error) {
	return uc.projectRepository.GetAll()
}
