package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/openedproject"
)

type AddCommand interface {
	Execute(command domain.Command) error
}

type DefaultAddCommand struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
}

func NewAddCommand(openedProject openedproject.OpenedProject, commandRepo domain.Repository) *DefaultAddCommand {
	return &DefaultAddCommand{
		openedProject:     openedProject,
		commandRepository: commandRepo,
	}
}

func (uc *DefaultAddCommand) Execute(newCommand domain.Command) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	allCommands, err := uc.commandRepository.GetAll(project.Id)
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
