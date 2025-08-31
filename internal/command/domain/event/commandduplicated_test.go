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
