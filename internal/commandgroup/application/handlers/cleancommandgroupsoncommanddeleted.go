package handlers

import (
	commanddomainevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	internalEvent "gomander/internal/event"
	"gomander/internal/eventbus"
	"gomander/internal/helpers/array"
	"gomander/internal/unitofwork"
)

// EventEmitter tells the UI what a handler has already committed.
type EventEmitter interface {
	EmitEvent(event internalEvent.Event, payload interface{})
}

type CleanCommandGroupsOnCommandDeleted struct {
	unitOfWork             unitofwork.UnitOfWork
	commandGroupRepository commandgroupdomain.Repository
	eventEmitter           EventEmitter
}

func (h *CleanCommandGroupsOnCommandDeleted) GetEvent() eventbus.Event {
	return commanddomainevent.CommandDeletedEvent{}
}

func NewCleanCommandGroupsOnCommandDeleted(
	unitOfWork unitofwork.UnitOfWork,
	commandGroupRepository commandgroupdomain.Repository,
	eventEmitter EventEmitter,
) *CleanCommandGroupsOnCommandDeleted {
	return &CleanCommandGroupsOnCommandDeleted{
		unitOfWork:             unitOfWork,
		commandGroupRepository: commandGroupRepository,
		eventEmitter:           eventEmitter,
	}
}

func (h *CleanCommandGroupsOnCommandDeleted) Execute(e eventbus.Event) error {
	event, ok := e.(commanddomainevent.CommandDeletedEvent)
	if !ok {
		return nil
	}

	cascade, err := h.applyTheCascadeFor(event.CommandId)
	if err != nil {
		return err
	}

	// The cascade has committed, so the UI is told before the renumbering that
	// follows gets a chance to fail.
	for _, commandGroup := range cascade.Deleted {
		h.eventEmitter.EmitEvent(internalEvent.CommandGroupDeleted, commandGroup.Id)
	}

	return h.closeTheGapsLeftIn(projectIdsOf(cascade.Deleted))
}

// applyTheCascadeFor lets the domain decide which Command Groups survive losing
// the Command, and writes that answer back in one Unit of Work, so a Group is
// never left holding a Command that is gone.
func (h *CleanCommandGroupsOnCommandDeleted) applyTheCascadeFor(commandId string) (commandgroupdomain.Cascade, error) {
	var cascade commandgroupdomain.Cascade

	err := h.unitOfWork.Do(func(repositories unitofwork.Repositories) error {
		commandGroups, err := repositories.CommandGroups.GetAllContaining(commandId)
		if err != nil {
			return err
		}

		cascade = commandgroupdomain.RemoveCommandFrom(commandGroups, commandId)

		for i := range cascade.Survived {
			if err := repositories.CommandGroups.Update(&cascade.Survived[i]); err != nil {
				return err
			}
		}

		for _, commandGroup := range cascade.Deleted {
			if err := repositories.CommandGroups.Delete(commandGroup.Id); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return commandgroupdomain.Cascade{}, err
	}

	return cascade, nil
}

func (h *CleanCommandGroupsOnCommandDeleted) closeTheGapsLeftIn(projectIds []string) error {
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
