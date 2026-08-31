package handlers

import (
	commanddomain "gomander/internal/command/domain"
	commanddomainevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/eventbus"
	"gomander/internal/helpers/array"
)

type AddCommandToGroupOnCommandDuplicated struct {
	commandRepository      commanddomain.Repository
	commandGroupRepository commandgroupdomain.Repository
}

func (h *AddCommandToGroupOnCommandDuplicated) GetEvent() eventbus.Event {
	return commanddomainevent.CommandDuplicatedEvent{}
}

func NewAddCommandToGroupOnCommandDuplicated(
	commandRepository commanddomain.Repository,
	commandGroupRepository commandgroupdomain.Repository,
) *AddCommandToGroupOnCommandDuplicated {
	return &AddCommandToGroupOnCommandDuplicated{
		commandRepository:      commandRepository,
		commandGroupRepository: commandGroupRepository,
	}
}

func (h *AddCommandToGroupOnCommandDuplicated) Execute(e eventbus.Event) error {
	event, ok := e.(commanddomainevent.CommandDuplicatedEvent)
	if !ok {
		return nil
	}
	if event.InsideGroupId == "" {
		return nil
	}

	commandGroup, err := h.commandGroupRepository.GetWithCommandIds(event.InsideGroupId)
	if err != nil {
		return err
	}

	if array.Contains(commandGroup.CommandIds, event.CommandId) {
		return nil
	}

	// The Group names the Command rather than carrying it, but a Command that
	// is not there must not gain a membership row.
	if _, err := h.commandRepository.Get(event.CommandId); err != nil {
		return err
	}

	commandGroup.CommandIds = append(commandGroup.CommandIds, event.CommandId)

	return h.commandGroupRepository.UpdateWithCommandIds(&commandGroup)
}
