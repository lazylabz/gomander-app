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
	deletedCommandGroup, err := uc.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}

	err = uc.commandGroupRepository.Delete(commandGroupId)
	if err != nil {
		return err
	}

	// Delete has committed, so the UI is told before the renumbering that
	// follows gets a chance to fail - the way the cascade in
	// CleanCommandGroupsOnCommandDeleted already reports its own deletions.
	uc.eventEmitter.EmitEvent(event.CommandGroupDeleted, commandGroupId)

	if deletedCommandGroup != nil {
		return uc.closeTheGapLeftIn(deletedCommandGroup.ProjectId)
	}

	return nil
}

func (uc *DefaultDeleteCommandGroup) closeTheGapLeftIn(projectId string) error {
	remainingCommandGroups, err := uc.commandGroupRepository.GetAll(projectId)
	if err != nil {
		return err
	}

	return domain.Order.CloseGaps(remainingCommandGroups, uc.commandGroupRepository.Update)
}
