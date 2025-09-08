package usecases

import (
	"gomander/internal/command/domain"
)

type EditCommand interface {
	Execute(command domain.Command) error
}

type DefaultEditCommand struct {
	commandRepository domain.Repository
}

func NewEditCommand(commandRepo domain.Repository) *DefaultEditCommand {
	return &DefaultEditCommand{
		commandRepository: commandRepo,
	}
}

func (uc *DefaultEditCommand) Execute(newCommand domain.Command) error {
	err := uc.commandRepository.Update(&newCommand)
	if err != nil {
		return err
	}

	return nil
}
