package handlers

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	internalEvent "gomander/internal/event"
	"gomander/internal/eventbus"
	projectdomainevent "gomander/internal/project/domain/event"
)

type CleanCommandGroupsOnProjectDeleted interface {
	Execute(e eventbus.Event) error
	GetEvent() eventbus.Event
}

type DefaultCleanCommandGroupsOnProjectDeleted struct {
	commandGroupRepository commandgroupdomain.Repository
	eventEmitter           internalEvent.EventEmitter
}

func (h *DefaultCleanCommandGroupsOnProjectDeleted) GetEvent() eventbus.Event {
	return projectdomainevent.ProjectDeletedEvent{}
}

func NewCleanCommandGroupsOnProjectDeleted(
	commandGroupRepository commandgroupdomain.Repository,
	eventEmitter internalEvent.EventEmitter,
) *DefaultCleanCommandGroupsOnProjectDeleted {
	return &DefaultCleanCommandGroupsOnProjectDeleted{
		commandGroupRepository: commandGroupRepository,
		eventEmitter:           eventEmitter,
	}
}

func (h *DefaultCleanCommandGroupsOnProjectDeleted) Execute(e eventbus.Event) error {
	event, ok := e.(projectdomainevent.ProjectDeletedEvent)
	if !ok {
		return nil
	}

	deletedIds, err := h.commandGroupRepository.DeleteAll(event.ProjectId)
	if err != nil {
		return err
	}

	for _, id := range deletedIds {
		h.eventEmitter.EmitEvent(internalEvent.CommandGroupDeleted, id)
	}

	return nil
}
