package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/project/domain/event"
)

func TestProjectDeletedEvent_GetName(t *testing.T) {
	e := event.ProjectDeletedEvent{}
	assert.Equal(t, "domain_event.project.delete", e.GetName())
}
