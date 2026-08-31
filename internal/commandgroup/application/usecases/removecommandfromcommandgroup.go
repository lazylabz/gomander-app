package usecases

import (
	"errors"

	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

type RemoveCommandFromCommandGroup struct {
	commandGroupRepository domain.Repository
}

func NewRemoveCommandFromCommandGroup(commandGroupRepo domain.Repository) *RemoveCommandFromCommandGroup {
	return &RemoveCommandFromCommandGroup{
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *RemoveCommandFromCommandGroup) Execute(commandId, commandGroupId string) error {
	commandGroup, err := uc.commandGroupRepository.GetWithCommandIds(commandGroupId)
	if err != nil {
		return err
	}
	if len(commandGroup.CommandIds) == 1 {
		return errors.New("cannot remove the last command from the group; delete the group instead")
	}

	commandGroup.CommandIds = array.Filter(commandGroup.CommandIds, func(heldCommandId string) bool {
		return heldCommandId != commandId
	})

	return uc.commandGroupRepository.UpdateWithCommandIds(&commandGroup)
}
