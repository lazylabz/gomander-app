package usecases

import (
	"gomander/internal/commandgroup/domain"
	"gomander/internal/event"
)

type DeleteCommandGroup struct {
	commandGroupRepository domain.Repository
	eventEmitter           event.EventEmitter
}

func NewDeleteCommandGroup(
	commandGroupRepo domain.Repository,
	eventEmitter event.EventEmitter,
) *DeleteCommandGroup {
	return &DeleteCommandGroup{
		commandGroupRepository: commandGroupRepo,
		eventEmitter:           eventEmitter,
	}
}

func (uc *DeleteCommandGroup) Execute(commandGroupId string) error {
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

func (uc *DeleteCommandGroup) closeTheGapLeftIn(projectId string) error {
	remainingCommandGroups, err := uc.commandGroupRepository.GetAll(projectId)
	if err != nil {
		return err
	}

	return domain.Order.CloseGaps(remainingCommandGroups, uc.commandGroupRepository.Update)
}
