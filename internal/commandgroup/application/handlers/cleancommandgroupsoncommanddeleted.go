package handlers

import (
	commanddomainevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	internalEvent "gomander/internal/event"
	"gomander/internal/eventbus"
	"gomander/internal/helpers/array"
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

	emptiedCommandGroups, err := h.commandGroupRepository.DeleteEmpty()
	if err != nil {
		return err
	}

	for _, commandGroup := range emptiedCommandGroups {
		h.eventEmitter.EmitEvent(internalEvent.CommandGroupDeleted, commandGroup.Id)
	}

	return h.closeTheGapsLeftIn(projectIdsOf(emptiedCommandGroups))
}

func (h *DefaultCleanCommandGroupsOnCommandDeleted) closeTheGapsLeftIn(projectIds []string) error {
	for _, projectId := range projectIds {
		remainingCommandGroups, err := h.commandGroupRepository.GetAll(projectId)
		if err != nil {
			return err
		}

		err = commandgroupdomain.Order.CloseGaps(remainingCommandGroups, h.commandGroupRepository.Update)
		if err != nil {
			return err
		}
	}

	return nil
}

func projectIdsOf(commandGroups []commandgroupdomain.CommandGroup) []string {
	projectIds := make([]string, 0, len(commandGroups))

	for _, commandGroup := range commandGroups {
		if !array.Contains(projectIds, commandGroup.ProjectId) {
			projectIds = append(projectIds, commandGroup.ProjectId)
		}
	}

	return projectIds
}
