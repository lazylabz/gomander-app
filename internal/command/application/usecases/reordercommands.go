package usecases

import (
	"gomander/internal/command/domain"
	"gomander/internal/openedproject"
)

type ReorderCommands struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
}

func NewReorderCommands(openedProject openedproject.OpenedProject, commandRepo domain.Repository) *ReorderCommands {
	return &ReorderCommands{
		openedProject:     openedProject,
		commandRepository: commandRepo,
	}
}

func (uc *ReorderCommands) Execute(orderedIds []string) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	existingCommands, err := uc.commandRepository.GetAll(project.Id)
	if err != nil {
		return err
	}

	return domain.Order.Rearrange(existingCommands, orderedIds, uc.commandRepository.Update)
}
