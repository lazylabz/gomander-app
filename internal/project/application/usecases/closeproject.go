package usecases

import (
	configdomain "gomander/internal/config/domain"
)

type CloseProject interface {
	Execute() error
}

type DefaultCloseProject struct {
	configRepository configdomain.Repository
}

func NewCloseProject(configRepo configdomain.Repository) *DefaultCloseProject {
	return &DefaultCloseProject{
		configRepository: configRepo,
	}
}

func (uc *DefaultCloseProject) Execute() error {
	config, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = ""

	err = uc.configRepository.Update(config)
	if err != nil {
		return err
	}

	return nil
}
