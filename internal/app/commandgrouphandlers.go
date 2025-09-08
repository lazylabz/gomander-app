package app

import (
	"errors"
	"sort"

	commandDomain "gomander/internal/command/domain"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

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

func (a *App) RemoveCommandFromCommandGroup(commandId, commandGroupId string) error {
	commandGroup, err := a.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}
	if len(commandGroup.Commands) == 1 {
		return errors.New("cannot remove the last command from the group; delete the group instead")
	}

	commandGroup.Commands = array.Filter(commandGroup.Commands, func(cmd commandDomain.Command) bool {
		return cmd.Id != commandId
	})

	if err := a.commandGroupRepository.Update(commandGroup); err != nil {
		return err
	}

	return nil
}

func (a *App) ReorderCommandGroups(newOrderedIds []string) error {
	userConfig, err := a.userConfigRepository.GetOrCreate()
	if err != nil {
		return err
	}

	existingCommandGroups, err := a.commandGroupRepository.GetAll(userConfig.LastOpenedProjectId)
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
