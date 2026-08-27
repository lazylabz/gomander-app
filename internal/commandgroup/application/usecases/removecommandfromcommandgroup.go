package usecases

import (
	"errors"

	commanddomain "gomander/internal/command/domain"
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
	commandGroup, err := uc.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}
	if len(commandGroup.Commands) == 1 {
		return errors.New("cannot remove the last command from the group; delete the group instead")
	}

	commandGroup.Commands = array.Filter(commandGroup.Commands, func(cmd commanddomain.Command) bool {
		return cmd.Id != commandId
	})

	return uc.commandGroupRepository.Update(commandGroup)
}
