package usecases

import (
	"gomander/internal/commandgroup/domain"
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

	return domain.Order.Rearrange(existingCommandGroups, newOrderedIds, uc.commandGroupRepository.Update)
}
