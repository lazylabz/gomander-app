package usecases

import (
	"errors"

	"gomander/internal/command/domain"
	domainevent "gomander/internal/command/domain/event"
	"gomander/internal/eventbus"
)

type RemoveCommand interface {
	Execute(commandId string) error
}

type DefaultRemoveCommand struct {
	commandRepository domain.Repository
	eventBus          eventbus.EventBus
}

func NewRemoveCommand(commandRepo domain.Repository, eventBus eventbus.EventBus) *DefaultRemoveCommand {
	return &DefaultRemoveCommand{
		commandRepository: commandRepo,
		eventBus:          eventBus,
	}
}

func (uc *DefaultRemoveCommand) Execute(commandId string) error {
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
	publishErr := uc.publishDeleted(commandId)

	if removedCommand != nil {
		err = uc.closeTheGapLeftIn(removedCommand.ProjectId)
		if err != nil {
			return errors.Join(publishErr, err)
		}
	}

	return publishErr
}

func (uc *DefaultRemoveCommand) publishDeleted(commandId string) error {
	errs := uc.eventBus.PublishSync(domainevent.NewCommandDeletedEvent(commandId))

	if len(errs) == 0 {
		return nil
	}

	combinedErrMsg := "Errors occurred while removing command:"

	for _, pubErr := range errs {
		combinedErrMsg += "\n- " + pubErr.Error()
	}

	return errors.New(combinedErrMsg)
}

func (uc *DefaultRemoveCommand) closeTheGapLeftIn(projectId string) error {
	remainingCommands, err := uc.commandRepository.GetAll(projectId)
	if err != nil {
		return err
	}

	return domain.Order.CloseGaps(remainingCommands, uc.commandRepository.Update)
}
