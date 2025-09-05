package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/command/domain/event"
)

func TestCommandDeletedEvent_GetName(t *testing.T) {
	// Arrange
	e := event.CommandDeletedEvent{}

	// Act
	result := e.GetName()

	// Assert
	assert.Equal(t, "domain_event.command.delete", result)
}
