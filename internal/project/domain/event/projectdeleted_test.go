package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/project/domain/event"
)

func TestProjectDeletedEvent_GetName(t *testing.T) {
	// Arrange
	e := event.ProjectDeletedEvent{}

	// Act
	result := e.GetName()

	// Assert
	assert.Equal(t, "domain_event.project.delete", result)
}
