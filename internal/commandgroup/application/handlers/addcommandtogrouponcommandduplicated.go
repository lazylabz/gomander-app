package handlers

import (
	commanddomain "gomander/internal/command/domain"
	commanddomainevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/eventbus"
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

	commandGroup, err := h.commandGroupRepository.Get(event.InsideGroupId)
	if err != nil {
		return err
	}

	// Check if the command is already in the group
	for _, cmd := range commandGroup.Commands {
		if cmd.Id == event.CommandId {
			return nil
		}
	}

	command, err := h.commandRepository.Get(event.CommandId)
	if err != nil {
		return err
	}

	commandGroup.Commands = append(commandGroup.Commands, command)

	if err := h.commandGroupRepository.Update(&commandGroup); err != nil {
		return err
	}

	return nil
}
