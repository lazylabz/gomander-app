package usecases

import (
	"sort"

	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
	"gomander/internal/openedproject"
)

type ReorderCommandGroups interface {
	Execute(newOrderedIds []string) error
}

type DefaultReorderCommandGroups struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository domain.Repository
}

func NewReorderCommandGroups(openedProject openedproject.OpenedProject, commandGroupRepo domain.Repository) *DefaultReorderCommandGroups {
	return &DefaultReorderCommandGroups{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultReorderCommandGroups) Execute(newOrderedIds []string) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	existingCommandGroups, err := uc.commandGroupRepository.GetAll(project.Id)
	if err != nil {
		return err
	}

	sort.Slice(existingCommandGroups, func(i, j int) bool {
		return array.IndexOf(newOrderedIds, existingCommandGroups[i].Id) < array.IndexOf(newOrderedIds, existingCommandGroups[j].Id)
	})

	for i := range existingCommandGroups {
		existingCommandGroups[i].Position = i

		err := uc.commandGroupRepository.Update(&existingCommandGroups[i])
		if err != nil {
			return err
		}
	}

	return nil
}
