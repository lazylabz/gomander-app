package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/command/domain/event"
)

func TestCommandDuplicatedEvent_GetName(t *testing.T) {
	e := event.CommandDuplicatedEvent{}
	assert.Equal(t, "domain_event.command.duplicate", e.GetName())
}

func TestNewCommandDuplicatedEvent(t *testing.T) {
	commandId := "1234"
	insideGroupId := "5678"
	e := event.NewCommandDuplicatedEvent(commandId, insideGroupId)
	assert.NotNil(t, e)
	assert.Equal(t, commandId, e.CommandId)
	assert.Equal(t, insideGroupId, e.InsideGroupId)
}
