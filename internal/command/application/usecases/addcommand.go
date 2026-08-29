package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/openedproject"
)

type AddCommand struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
}

func NewAddCommand(openedProject openedproject.OpenedProject, commandRepo domain.Repository) *AddCommand {
	return &AddCommand{
		openedProject:     openedProject,
		commandRepository: commandRepo,
	}
}

func (uc *AddCommand) Execute(newCommand domain.Command) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	allCommands, err := uc.commandRepository.GetAll(project.Id)
	if err != nil {
		return err
	}

	newCommand.Position = domain.Order.End(allCommands)

	err = uc.commandRepository.Create(&newCommand)
	if err != nil {
		return err
	}

	return nil
}
