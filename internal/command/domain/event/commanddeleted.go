package event

type CommandDeletedEvent struct {
	CommandId string
}

func (CommandDeletedEvent) GetName() string {
	return "domain_event.command.delete"
}

func NewCommandDeletedEvent(commandId string) CommandDeletedEvent {
	return CommandDeletedEvent{
		CommandId: commandId,
	}
}
