package app

import (
	"sort"

	"gomander/internal/helpers/array"
)

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
