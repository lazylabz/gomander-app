package usecases

import (
	"sort"

	"gomander/internal/command/domain"
	"gomander/internal/helpers/array"
	"gomander/internal/openedproject"
)

type ReorderCommands interface {
	Execute(orderedIds []string) error
}

type DefaultReorderCommands struct {
	openedProject     openedproject.OpenedProject
	commandRepository domain.Repository
}

func NewReorderCommands(openedProject openedproject.OpenedProject, commandRepo domain.Repository) *DefaultReorderCommands {
	return &DefaultReorderCommands{
		openedProject:     openedProject,
		commandRepository: commandRepo,
	}
}

func (uc *DefaultReorderCommands) Execute(orderedIds []string) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	existingCommands, err := uc.commandRepository.GetAll(project.Id)
	if err != nil {
		return err
	}

	// Sort the existing commands based on the new order
	sort.Slice(existingCommands, func(i, j int) bool {
		return array.IndexOf(orderedIds, existingCommands[i].Id) < array.IndexOf(orderedIds, existingCommands[j].Id)
	})

	// Update the position of each command based on the new order
	for i := range existingCommands {
		existingCommands[i].Position = i
		err := uc.commandRepository.Update(&existingCommands[i])
		if err != nil {
			return err
		}
	}

	return nil
}
