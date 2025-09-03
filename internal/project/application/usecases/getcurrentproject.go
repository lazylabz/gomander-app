package usecases

import (
	configdomain "gomander/internal/config/domain"
	"gomander/internal/project/domain"
)

type GetCurrentProject interface {
	Execute() (*domain.Project, error)
}

type DefaultGetCurrentProject struct {
	configRepository  configdomain.Repository
	projectRepository domain.Repository
}

func NewGetCurrentProject(configRepo configdomain.Repository, projectRepo domain.Repository) *DefaultGetCurrentProject {
	return &DefaultGetCurrentProject{
		configRepository:  configRepo,
		projectRepository: projectRepo,
	}
}

func (uc *DefaultGetCurrentProject) Execute() (*domain.Project, error) {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return nil, err
	}

	return uc.projectRepository.Get(userConfig.LastOpenedProjectId)
}
