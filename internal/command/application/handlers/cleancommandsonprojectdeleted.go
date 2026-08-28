package handlers

import (
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/eventbus"
	projectdomainevent "gomander/internal/project/domain/event"
)

type CleanCommandsOnProjectDeleted struct {
	commandRepository commanddomain.Repository
}

func (h *CleanCommandsOnProjectDeleted) GetEvent() eventbus.Event {
	return projectdomainevent.ProjectDeletedEvent{}
}

func NewCleanCommandOnProjectDeleted(commandRepository commanddomain.Repository) *CleanCommandsOnProjectDeleted {
	return &CleanCommandsOnProjectDeleted{
		commandRepository: commandRepository,
	}
}

func (h *CleanCommandsOnProjectDeleted) Execute(e eventbus.Event) error {
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
