package app

import (
	"sort"

	"gomander/internal/commandgroup/domain"
	"gomander/internal/event"
	"gomander/internal/helpers/array"
)

func (a *App) GetCommandGroups() ([]domain.CommandGroup, error) {
	return a.commandGroupRepository.GetCommandGroups(a.openedProjectId)
}

func (a *App) CreateCommandGroup(commandGroup *domain.CommandGroup) error {
	existingCommandGroups, err := a.commandGroupRepository.GetCommandGroups(a.openedProjectId)
	if err != nil {
		return err
	}

	newPosition := len(existingCommandGroups)

	commandGroup.Position = newPosition

	if err := a.commandGroupRepository.CreateCommandGroup(commandGroup); err != nil {
		return err
	}

	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}

func (a *App) UpdateCommandGroup(commandGroup *domain.CommandGroup) error {
	if err := a.commandGroupRepository.UpdateCommandGroup(commandGroup); err != nil {
		return err
	}

	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}

func (a *App) DeleteCommandGroup(commandGroupId string) error {
	if err := a.commandGroupRepository.DeleteCommandGroup(commandGroupId); err != nil {
		return err
	}

	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}

func (a *App) ReorderCommandGroups(newOrderedIds []string) error {
	existingCommandGroups, err := a.commandGroupRepository.GetCommandGroups(a.openedProjectId)
	if err != nil {
		return err
	}

	sort.Slice(existingCommandGroups, func(i, j int) bool {
		return array.IndexOf(newOrderedIds, existingCommandGroups[i].Id) < array.IndexOf(newOrderedIds, existingCommandGroups[j].Id)
	})

	for i := range existingCommandGroups {
		existingCommandGroups[i].Position = i

		err := a.commandGroupRepository.UpdateCommandGroup(&existingCommandGroups[i])
		if err != nil {
			return err
		}
	}

	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}
