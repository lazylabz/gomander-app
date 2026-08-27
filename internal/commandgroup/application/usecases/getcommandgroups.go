package usecases

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/openedproject"
)

type GetCommandGroups interface {
	Execute() ([]domain.CommandGroup, error)
}

type DefaultGetCommandGroups struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository domain.Repository
}

func NewGetCommandGroups(openedProject openedproject.OpenedProject, commandGroupRepo domain.Repository) *DefaultGetCommandGroups {
	return &DefaultGetCommandGroups{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultGetCommandGroups) Execute() ([]domain.CommandGroup, error) {
	project, open, err := uc.openedProject.Find()
	if err != nil {
		return make([]domain.CommandGroup, 0), err
	}
	if !open {
		return make([]domain.CommandGroup, 0), nil
	}

	return uc.commandGroupRepository.GetAll(project.Id)
}
