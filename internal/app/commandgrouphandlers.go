package app

import (
	"sort"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/event"
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

	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}

func (a *App) RemoveCommandFromCommandGroups(id string) error {
	commandGroups, err := a.commandGroupRepository.GetAll(a.openedProjectId)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	for _, commandGroup := range commandGroups {
		commands := commandGroup.Commands
		commandIds := array.Map(commands, func(c commanddomain.Command) string { return c.Id })

		if array.Contains(commandIds, id) {
			commandGroup.Commands = array.Filter(commandGroup.Commands, func(c commanddomain.Command) bool {
				return c.Id != id
			})

			if len(commandGroup.Commands) == 0 {
				err := a.commandGroupRepository.Delete(commandGroup.Id)
				if err != nil {
					a.logger.Error(err.Error())
					return err
				}
			} else {
				err := a.commandGroupRepository.Update(&commandGroup)
				if err != nil {
					a.logger.Error(err.Error())
					return err
				}
			}
		}
	}

	return nil
}
