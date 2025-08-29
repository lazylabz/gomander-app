package event

type CommandDuplicatedEvent struct {
	CommandId     string
	InsideGroupId string
}

func (CommandDuplicatedEvent) GetName() string {
	return "domain_event.command.duplicate"
}

func NewCommandDuplicatedEvent(commandId, insideGroupId string) CommandDuplicatedEvent {
	return CommandDuplicatedEvent{
		CommandId:     commandId,
		InsideGroupId: insideGroupId,
	}
}
