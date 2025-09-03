package handlers

import (
	commanddomainevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/eventbus"
)

type CleanCommandGroupsOnCommandDeleted interface {
	Execute(e eventbus.Event) error
	GetEvent() eventbus.Event
}

type DefaultCleanCommandGroupsOnCommandDeleted struct {
	commandGroupRepository commandgroupdomain.Repository
}

func (h *DefaultCleanCommandGroupsOnCommandDeleted) GetEvent() eventbus.Event {
	return commanddomainevent.CommandDeletedEvent{}
}

func NewCleanCommandGroupsOnCommandDeleted(commandGroupRepository commandgroupdomain.Repository) *DefaultCleanCommandGroupsOnCommandDeleted {
	return &DefaultCleanCommandGroupsOnCommandDeleted{
		commandGroupRepository: commandGroupRepository,
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

	err = h.commandGroupRepository.DeleteEmpty()
	if err != nil {
		return err
	}

	return nil
}
