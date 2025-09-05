package usecases

import (
	configdomain "gomander/internal/config/domain"
	"gomander/internal/project/domain"
)

type OpenProject interface {
	Execute(projectId string) error
}

type DefaultOpenProject struct {
	configRepository  configdomain.Repository
	projectRepository domain.Repository
}

func NewOpenProject(configRepo configdomain.Repository, projectRepo domain.Repository) *DefaultOpenProject {
	return &DefaultOpenProject{
		configRepository:  configRepo,
		projectRepository: projectRepo,
	}
}

func (uc *DefaultOpenProject) Execute(projectId string) error {
	_, err := uc.projectRepository.Get(projectId)
	if err != nil {
		return err
	}

	config, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = projectId

	err = uc.configRepository.Update(config)
	if err != nil {
		return err
	}

	return nil
}
