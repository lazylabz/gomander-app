package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/command/domain/event"
)

func TestCommandDuplicatedEvent_GetName(t *testing.T) {
	// Arrange
	e := event.CommandDuplicatedEvent{}

	// Act
	result := e.GetName()

	// Assert
	assert.Equal(t, "domain_event.command.duplicate", result)
}

func TestNewCommandDuplicatedEvent(t *testing.T) {
	// Arrange
	commandId := "1234"
	insideGroupId := "5678"

	// Act
	e := event.NewCommandDuplicatedEvent(commandId, insideGroupId)

	// Assert
	assert.NotNil(t, e)
	assert.Equal(t, commandId, e.CommandId)
	assert.Equal(t, insideGroupId, e.InsideGroupId)
}
