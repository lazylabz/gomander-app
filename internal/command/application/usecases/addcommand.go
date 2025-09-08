package usecases

import (
	"gomander/internal/command/domain"
	configdomain "gomander/internal/config/domain"
)

type AddCommand interface {
	Execute(command domain.Command) error
}

type DefaultAddCommand struct {
	configRepository  configdomain.Repository
	commandRepository domain.Repository
}

func NewAddCommand(configRepo configdomain.Repository, commandRepo domain.Repository) *DefaultAddCommand {
	return &DefaultAddCommand{
		configRepository:  configRepo,
		commandRepository: commandRepo,
	}
}

func (uc *DefaultAddCommand) Execute(newCommand domain.Command) error {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	allCommands, err := uc.commandRepository.GetAll(userConfig.LastOpenedProjectId)
	if err != nil {
		return err
	}

	newPosition := len(allCommands)
	newCommand.Position = newPosition

	err = uc.commandRepository.Create(&newCommand)
	if err != nil {
		return err
	}

	return nil
}
