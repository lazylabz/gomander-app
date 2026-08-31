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

// Execute is the last caller of the form that carries whole Commands. The
// controller narrows them to the ids the frontend receives, so what the wire
// speaks does not depend on it; the form goes when the embedded one does.
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
