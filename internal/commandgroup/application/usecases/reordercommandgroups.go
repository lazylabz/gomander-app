package usecases

import (
	"sort"

	"gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/helpers/array"
)

type ReorderCommandGroups interface {
	Execute(newOrderedIds []string) error
}

type DefaultReorderCommandGroups struct {
	configRepository       configdomain.Repository
	commandGroupRepository domain.Repository
}

func NewReorderCommandGroups(configRepo configdomain.Repository, commandGroupRepo domain.Repository) *DefaultReorderCommandGroups {
	return &DefaultReorderCommandGroups{
		configRepository:       configRepo,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultReorderCommandGroups) Execute(newOrderedIds []string) error {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	existingCommandGroups, err := uc.commandGroupRepository.GetAll(userConfig.LastOpenedProjectId)
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
