package usecases

import (
	"gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
)

type GetCommandGroups interface {
	Execute() ([]domain.CommandGroup, error)
}

type DefaultGetCommandGroups struct {
	configRepository       configdomain.Repository
	commandGroupRepository domain.Repository
}

func NewGetCommandGroups(configRepo configdomain.Repository, commandGroupRepo domain.Repository) *DefaultGetCommandGroups {
	return &DefaultGetCommandGroups{
		configRepository:       configRepo,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultGetCommandGroups) Execute() ([]domain.CommandGroup, error) {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return make([]domain.CommandGroup, 0), err
	}
	return uc.commandGroupRepository.GetAll(userConfig.LastOpenedProjectId)
}
