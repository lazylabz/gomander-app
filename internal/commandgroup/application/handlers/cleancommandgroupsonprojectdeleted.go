package handlers

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	internalEvent "gomander/internal/event"
	"gomander/internal/eventbus"
	projectdomainevent "gomander/internal/project/domain/event"
)

type CleanCommandGroupsOnProjectDeleted struct {
	commandGroupRepository commandgroupdomain.Repository
	eventEmitter           EventEmitter
}

func (h *CleanCommandGroupsOnProjectDeleted) GetEvent() eventbus.Event {
	return projectdomainevent.ProjectDeletedEvent{}
}

func NewCleanCommandGroupsOnProjectDeleted(
	commandGroupRepository commandgroupdomain.Repository,
	eventEmitter EventEmitter,
) *CleanCommandGroupsOnProjectDeleted {
	return &CleanCommandGroupsOnProjectDeleted{
		commandGroupRepository: commandGroupRepository,
		eventEmitter:           eventEmitter,
	}
}

func (h *CleanCommandGroupsOnProjectDeleted) Execute(e eventbus.Event) error {
	event, ok := e.(projectdomainevent.ProjectDeletedEvent)
	if !ok {
		return nil
	}

	deleted, err := h.deleteTheCommandGroupsOf(event.ProjectId)
	if err != nil {
		return err
	}

	for _, commandGroup := range deleted {
		h.eventEmitter.EmitEvent(internalEvent.CommandGroupDeleted, commandGroup.Id)
	}

	return nil
}

// deleteTheCommandGroupsOf works in one transaction, so a Project never
// outlives only some of its Command Groups.
func (h *CleanCommandGroupsOnProjectDeleted) deleteTheCommandGroupsOf(projectId string) ([]commandgroupdomain.CommandGroup, error) {
	var deleted []commandgroupdomain.CommandGroup

	err := h.commandGroupRepository.Atomically(func(commandGroupRepository commandgroupdomain.Repository) error {
		commandGroups, err := commandGroupRepository.GetAll(projectId)
		if err != nil {
			return err
		}

		for _, commandGroup := range commandGroups {
			if err := commandGroupRepository.Delete(commandGroup.Id); err != nil {
				return err
			}
		}

		deleted = commandGroups

		return nil
	})

	if err != nil {
		return nil, err
	}

	return deleted, nil
}
