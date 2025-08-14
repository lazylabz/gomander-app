package app

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/event"
)

func (a *App) GetCommandGroups() ([]domain.CommandGroup, error) {
	return a.commandGroupRepository.GetCommandGroups(a.openedProjectId)
}

func (a *App) CreateCommandGroup(commandGroup *domain.CommandGroup) error {
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

func (a *App) ReorderCommandGroups() error {
	// TODO: Implement reordering logic

	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}
