package usecases

import (
	"sort"

	"gomander/internal/command/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/helpers/array"
)

type ReorderCommands interface {
	Execute(orderedIds []string) error
}

type DefaultReorderCommands struct {
	configRepository  configdomain.Repository
	commandRepository domain.Repository
}

func NewReorderCommands(configRepo configdomain.Repository, commandRepo domain.Repository) *DefaultReorderCommands {
	return &DefaultReorderCommands{
		configRepository:  configRepo,
		commandRepository: commandRepo,
	}
}

func (uc *DefaultReorderCommands) Execute(orderedIds []string) error {
	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	existingCommands, err := uc.commandRepository.GetAll(userConfig.LastOpenedProjectId)
	if err != nil {
		return err
	}

	// Sort the existing commands based on the new order
	sort.Slice(existingCommands, func(i, j int) bool {
		return array.IndexOf(orderedIds, existingCommands[i].Id) < array.IndexOf(orderedIds, existingCommands[j].Id)
	})

	// Update the position of each command based on the new order
	for i := range existingCommands {
		existingCommands[i].Position = i
		err := uc.commandRepository.Update(&existingCommands[i])
		if err != nil {
			return err
		}
	}

	return nil
}
