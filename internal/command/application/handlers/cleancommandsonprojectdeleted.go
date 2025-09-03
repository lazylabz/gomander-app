package handlers

import (
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/eventbus"
	projectdomainevent "gomander/internal/project/domain/event"
)

type CleanCommandsOnProjectDeleted interface {
	Execute(e eventbus.Event) error
	GetEvent() eventbus.Event
}

type DefaultCleanCommandsOnProjectDeleted struct {
	commandRepository commanddomain.Repository
}

func (h *DefaultCleanCommandsOnProjectDeleted) GetEvent() eventbus.Event {
	return projectdomainevent.ProjectDeletedEvent{}
}

func NewCleanCommandOnProjectDeleted(commandRepository commanddomain.Repository) *DefaultCleanCommandsOnProjectDeleted {
	return &DefaultCleanCommandsOnProjectDeleted{
		commandRepository: commandRepository,
	}
}

func (h *DefaultCleanCommandsOnProjectDeleted) Execute(e eventbus.Event) error {
	event, ok := e.(projectdomainevent.ProjectDeletedEvent)
	if !ok {
		return nil
	}

	err := h.commandRepository.DeleteAll(event.ProjectId)
	if err != nil {
		return err
	}

	return nil
}
