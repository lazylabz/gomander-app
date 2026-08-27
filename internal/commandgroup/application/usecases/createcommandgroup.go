package usecases

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/openedproject"
)

type CreateCommandGroup interface {
	Execute(commandGroup *domain.CommandGroup) error
}

type DefaultCreateCommandGroup struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository domain.Repository
}

func NewCreateCommandGroup(openedProject openedproject.OpenedProject, commandGroupRepo domain.Repository) *DefaultCreateCommandGroup {
	return &DefaultCreateCommandGroup{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultCreateCommandGroup) Execute(commandGroup *domain.CommandGroup) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	existingCommandGroups, err := uc.commandGroupRepository.GetAll(project.Id)
	if err != nil {
		return err
	}

	newPosition := len(existingCommandGroups)
	commandGroup.Position = newPosition

	return uc.commandGroupRepository.Create(commandGroup)
}
