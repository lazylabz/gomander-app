package handlers

import (
	commanddomainevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	internalEvent "gomander/internal/event"
	"gomander/internal/eventbus"
)

type CleanCommandGroupsOnCommandDeleted interface {
	Execute(e eventbus.Event) error
	GetEvent() eventbus.Event
}

type DefaultCleanCommandGroupsOnCommandDeleted struct {
	commandGroupRepository commandgroupdomain.Repository
	eventEmitter           internalEvent.EventEmitter
}

func (h *DefaultCleanCommandGroupsOnCommandDeleted) GetEvent() eventbus.Event {
	return commanddomainevent.CommandDeletedEvent{}
}

func NewCleanCommandGroupsOnCommandDeleted(
	commandGroupRepository commandgroupdomain.Repository,
	eventEmitter internalEvent.EventEmitter,
) *DefaultCleanCommandGroupsOnCommandDeleted {
	return &DefaultCleanCommandGroupsOnCommandDeleted{
		commandGroupRepository: commandGroupRepository,
		eventEmitter:           eventEmitter,
	}
}

func (h *DefaultCleanCommandGroupsOnCommandDeleted) Execute(e eventbus.Event) error {
	event, ok := e.(commanddomainevent.CommandDeletedEvent)
	if !ok {
		return nil
	}

	err := h.commandGroupRepository.RemoveCommandFromCommandGroups(event.CommandId)
	if err != nil {
		return err
	}

	deletedIds, err := h.commandGroupRepository.DeleteEmpty()
	if err != nil {
		return err
	}

	for _, id := range deletedIds {
		h.eventEmitter.EmitEvent(internalEvent.CommandGroupDeleted, id)
	}

	return nil
}
