package usecases

import (
	"errors"

	"gomander/internal/command/domain"
	domainevent "gomander/internal/command/domain/event"
	"gomander/internal/eventbus"
)

type RemoveCommand struct {
	commandRepository domain.Repository
	eventBus          eventbus.EventBus
}

func NewRemoveCommand(commandRepo domain.Repository, eventBus eventbus.EventBus) *RemoveCommand {
	return &RemoveCommand{
		commandRepository: commandRepo,
		eventBus:          eventBus,
	}
}

func (uc *RemoveCommand) Execute(commandId string) error {
	removedCommand, err := uc.commandRepository.Get(commandId)
	if err != nil {
		return err
	}

	err = uc.commandRepository.Delete(commandId)
	if err != nil {
		return err
	}

	// The Command is gone, so the handlers that clean up after it run whatever
	// the renumbering behind them does: a Command Group still holding a deleted
	// Command is worse than a gap in the sequence.
	publishErr := eventbus.Combined(
		"Errors occurred while removing command:",
		uc.eventBus.PublishSync(domainevent.NewCommandDeletedEvent(commandId)),
	)

	if removedCommand != nil {
		err = uc.closeTheGapLeftIn(removedCommand.ProjectId)
		if err != nil {
			return errors.Join(publishErr, err)
		}
	}

	return publishErr
}

func (uc *RemoveCommand) closeTheGapLeftIn(projectId string) error {
	remainingCommands, err := uc.commandRepository.GetAll(projectId)
	if err != nil {
		return err
	}

	return domain.Order.CloseGaps(remainingCommands, uc.commandRepository.Update)
}
