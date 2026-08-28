package handlers

import (
	commandgroupdomain "gomander/internal/commandgroup/domain"
	internalEvent "gomander/internal/event"
	"gomander/internal/eventbus"
	projectdomainevent "gomander/internal/project/domain/event"
	"gomander/internal/unitofwork"
)

type CleanCommandGroupsOnProjectDeleted struct {
	unitOfWork   unitofwork.UnitOfWork
	eventEmitter EventEmitter
}

func (h *CleanCommandGroupsOnProjectDeleted) GetEvent() eventbus.Event {
	return projectdomainevent.ProjectDeletedEvent{}
}

func NewCleanCommandGroupsOnProjectDeleted(
	unitOfWork unitofwork.UnitOfWork,
	eventEmitter EventEmitter,
) *CleanCommandGroupsOnProjectDeleted {
	return &CleanCommandGroupsOnProjectDeleted{
		unitOfWork:   unitOfWork,
		eventEmitter: eventEmitter,
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

	err := h.unitOfWork.Do(func(repositories unitofwork.Repositories) error {
		commandGroups, err := repositories.CommandGroups.GetAll(projectId)
		if err != nil {
			return err
		}

		for _, commandGroup := range commandGroups {
			if err := repositories.CommandGroups.Delete(commandGroup.Id); err != nil {
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
