package handlers

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/eventbus"
	projectdomainevent "gomander/internal/project/domain/event"
)

type CleanCommandGroupsOnProjectDeleted interface {
	Execute(e eventbus.Event) error
	GetEvent() eventbus.Event
}

type DefaultCleanCommandGroupsOnProjectDeleted struct {
	commandGroupRepository commandgroupdomain.Repository
}

func (h *DefaultCleanCommandGroupsOnProjectDeleted) GetEvent() eventbus.Event {
	return projectdomainevent.ProjectDeletedEvent{}
}

func NewDefaultCleanCommandGroupsOnProjectDeleted(commandGroupRepository commandgroupdomain.Repository) *DefaultCleanCommandGroupsOnProjectDeleted {
	return &DefaultCleanCommandGroupsOnProjectDeleted{
		commandGroupRepository: commandGroupRepository,
	}
}

func (h *DefaultCleanCommandGroupsOnProjectDeleted) Execute(e eventbus.Event) error {
	event, ok := e.(projectdomainevent.ProjectDeletedEvent)
	if !ok {
		return nil
	}

	err := h.commandGroupRepository.DeleteAll(event.ProjectId)
	if err != nil {
		return err
	}

	return nil
}
