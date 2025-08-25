package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/command/domain/event"
)

func TestCommandDeletedEvent_GetName(t *testing.T) {
	e := event.CommandDeletedEvent{}
	assert.Equal(t, "domain_event.command.delete", e.GetName())
}
