// filepath: /Users/moises/Code/personal/gomander/internal/project/application/usecases/getavailableprojects.go
package usecases

import (
	"gomander/internal/project/domain"
)

type GetAvailableProjects interface {
	Execute() ([]domain.Project, error)
}

type DefaultGetAvailableProjects struct {
	projectRepository domain.Repository
}

func NewGetAvailableProjects(projectRepo domain.Repository) *DefaultGetAvailableProjects {
	return &DefaultGetAvailableProjects{
		projectRepository: projectRepo,
	}
}

func (uc *DefaultGetAvailableProjects) Execute() ([]domain.Project, error) {
	return uc.projectRepository.GetAll()
}
