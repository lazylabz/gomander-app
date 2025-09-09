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
	err := uc.commandRepository.Delete(commandId)
	if err != nil {
		return err
	}

	domainEvent := domainevent.NewCommandDeletedEvent(commandId)

	errs := uc.eventBus.PublishSync(domainEvent)

	if len(errs) > 0 {
		combinedErrMsg := "Errors occurred while removing command:"

		for _, pubErr := range errs {
			combinedErrMsg += "\n- " + pubErr.Error()
		}

		err = errors.New(combinedErrMsg)

		return err
	}

	return nil
}
