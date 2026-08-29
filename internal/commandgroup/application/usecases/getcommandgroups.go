package usecases

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/openedproject"
)

type GetCommandGroups struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository domain.Repository
}

func NewGetCommandGroups(openedProject openedproject.OpenedProject, commandGroupRepo domain.Repository) *GetCommandGroups {
	return &GetCommandGroups{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *GetCommandGroups) Execute() ([]domain.CommandGroup, error) {
	project, open, err := uc.openedProject.Find()
	if err != nil {
		return make([]domain.CommandGroup, 0), err
	}
	if !open {
		return make([]domain.CommandGroup, 0), nil
	}

	return uc.commandGroupRepository.GetAll(project.Id)
}
