package handlers

import (
	commanddomain "gomander/internal/command/domain"
	commanddomainevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/eventbus"
)

type AddCommandToGroupOnCommandDuplicated interface {
	Execute(e eventbus.Event) error
	GetEvent() eventbus.Event
}

type DefaultAddCommandToGroupOnCommandDuplicated struct {
	commandRepository      commanddomain.Repository
	commandGroupRepository commandgroupdomain.Repository
}

func (h *DefaultAddCommandToGroupOnCommandDuplicated) GetEvent() eventbus.Event {
	return commanddomainevent.CommandDuplicatedEvent{}
}

func NewDefaultAddCommandToGroupOnCommandDuplicated(
	commandRepository commanddomain.Repository,
	commandGroupRepository commandgroupdomain.Repository,
) *DefaultAddCommandToGroupOnCommandDuplicated {
	return &DefaultAddCommandToGroupOnCommandDuplicated{
		commandRepository:      commandRepository,
		commandGroupRepository: commandGroupRepository,
	}
}

func (h *DefaultAddCommandToGroupOnCommandDuplicated) Execute(e eventbus.Event) error {
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

	commandGroup.Commands = append(commandGroup.Commands, *command)

	if err := h.commandGroupRepository.Update(commandGroup); err != nil {
		return err
	}

	return nil
}
