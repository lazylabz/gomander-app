package usecases

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/openedproject"
)

type CreateCommandGroup struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository domain.Repository
}

func NewCreateCommandGroup(openedProject openedproject.OpenedProject, commandGroupRepo domain.Repository) *CreateCommandGroup {
	return &CreateCommandGroup{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *CreateCommandGroup) Execute(commandGroup *domain.CommandGroupWithCommandIds) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	existingCommandGroups, err := uc.commandGroupRepository.GetAllWithCommandIds(project.Id)
	if err != nil {
		return err
	}

	commandGroup.Position = domain.Order.End(existingCommandGroups)

	return uc.commandGroupRepository.CreateWithCommandIds(commandGroup)
}
