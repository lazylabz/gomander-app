package usecases

import (
	"gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
)

type CreateCommandGroup interface {
	Execute(commandGroup *domain.CommandGroup) error
}

type DefaultCreateCommandGroup struct {
	configRepository       configdomain.Repository
	commandGroupRepository domain.Repository
}

func NewCreateCommandGroup(configRepo configdomain.Repository, commandGroupRepo domain.Repository) *DefaultCreateCommandGroup {
	return &DefaultCreateCommandGroup{
		configRepository:       configRepo,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultCreateCommandGroup) Execute(commandGroup *domain.CommandGroup) error {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	existingCommandGroups, err := uc.commandGroupRepository.GetAll(userConfig.LastOpenedProjectId)
	if err != nil {
		return err
	}

	newPosition := len(existingCommandGroups)
	commandGroup.Position = newPosition

	return uc.commandGroupRepository.Create(commandGroup)
}
