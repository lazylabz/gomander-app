package app

import (
	"sort"

	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

func (a *App) GetCommandGroups() ([]domain.CommandGroup, error) {
	return a.commandGroupRepository.GetAll(a.openedProjectId)
}

func (a *App) CreateCommandGroup(commandGroup *domain.CommandGroup) error {
	existingCommandGroups, err := a.commandGroupRepository.GetAll(a.openedProjectId)
	if err != nil {
		return err
	}

	newPosition := len(existingCommandGroups)

	commandGroup.Position = newPosition

	if err := a.commandGroupRepository.Create(commandGroup); err != nil {
		return err
	}

	return nil
}

func (a *App) UpdateCommandGroup(commandGroup *domain.CommandGroup) error {
	if err := a.commandGroupRepository.Update(commandGroup); err != nil {
		return err
	}

	return nil
}

func (a *App) DeleteCommandGroup(commandGroupId string) error {
	if err := a.commandGroupRepository.Delete(commandGroupId); err != nil {
		return err
	}

	return nil
}

func (a *App) AddCommandToCommandGroup(commandId, commandGroupId string) error {
	commandGroup, err := a.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}

	for _, cmd := range commandGroup.Commands {
		if cmd.Id == commandId {
			return nil
		}
	}

	command, err := a.commandRepository.Get(commandId)
	if err != nil {
		return err
	}

	commandGroup.Commands = append(commandGroup.Commands, *command)

	if err := a.commandGroupRepository.Update(commandGroup); err != nil {
		return err
	}

	return nil
}

func (a *App) ReorderCommandGroups(newOrderedIds []string) error {
	existingCommandGroups, err := a.commandGroupRepository.GetAll(a.openedProjectId)
	if err != nil {
		return err
	}

	sort.Slice(existingCommandGroups, func(i, j int) bool {
		return array.IndexOf(newOrderedIds, existingCommandGroups[i].Id) < array.IndexOf(newOrderedIds, existingCommandGroups[j].Id)
	})

	for i := range existingCommandGroups {
		existingCommandGroups[i].Position = i

		err := a.commandGroupRepository.Update(&existingCommandGroups[i])
		if err != nil {
			return err
		}
	}

	return nil
}
