package usecases

import (
	"gomander/internal/command/domain"
)

type EditCommand struct {
	commandRepository domain.Repository
}

func NewEditCommand(commandRepo domain.Repository) *EditCommand {
	return &EditCommand{
		commandRepository: commandRepo,
	}
}

func (uc *EditCommand) Execute(newCommand domain.Command) error {
	err := uc.commandRepository.Update(&newCommand)
	if err != nil {
		return err
	}

	return nil
}
