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
// the Command, and writes that answer back in one transaction, so a Group is
// never left holding a Command that is gone.
func (h *DefaultCleanCommandGroupsOnCommandDeleted) applyTheCascadeFor(commandId string) (commandgroupdomain.Cascade, error) {
	var cascade commandgroupdomain.Cascade

	err := h.commandGroupRepository.Atomically(func(commandGroupRepository commandgroupdomain.Repository) error {
		commandGroups, err := commandGroupRepository.GetAllContaining(commandId)
		if err != nil {
			return err
		}

		cascade = commandgroupdomain.RemoveCommandFrom(commandGroups, commandId)

		for i := range cascade.Survived {
			if err := commandGroupRepository.Update(&cascade.Survived[i]); err != nil {
				return err
			}
		}

		for _, commandGroup := range cascade.Deleted {
			if err := commandGroupRepository.Delete(commandGroup.Id); err != nil {
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
