package usecases

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/openedproject"
)

type ReorderCommandGroups struct {
	openedProject          openedproject.OpenedProject
	commandGroupRepository domain.Repository
}

func NewReorderCommandGroups(openedProject openedproject.OpenedProject, commandGroupRepo domain.Repository) *ReorderCommandGroups {
	return &ReorderCommandGroups{
		openedProject:          openedProject,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *ReorderCommandGroups) Execute(newOrderedIds []string) error {
	project, err := uc.openedProject.Get()
	if err != nil {
		return err
	}

	existingCommandGroups, err := uc.commandGroupRepository.GetAllWithCommandIds(project.Id)
	if err != nil {
		return err
	}

	return domain.Order.Rearrange(existingCommandGroups, newOrderedIds, uc.commandGroupRepository.UpdateWithCommandIds)
}
