package usecases

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/event"
)

type DeleteCommandGroup interface {
	Execute(commandGroupId string) error
}

type DefaultDeleteCommandGroup struct {
	commandGroupRepository domain.Repository
	eventEmitter           event.EventEmitter
}

func NewDeleteCommandGroup(
	commandGroupRepo domain.Repository,
	eventEmitter event.EventEmitter,
) *DefaultDeleteCommandGroup {
	return &DefaultDeleteCommandGroup{
		commandGroupRepository: commandGroupRepo,
		eventEmitter:           eventEmitter,
	}
}

func (uc *DefaultDeleteCommandGroup) Execute(commandGroupId string) error {
	err := uc.commandGroupRepository.Delete(commandGroupId)
	if err != nil {
		return err
	}

	uc.eventEmitter.EmitEvent(event.CommandGroupDeleted, commandGroupId)

	return nil
}
