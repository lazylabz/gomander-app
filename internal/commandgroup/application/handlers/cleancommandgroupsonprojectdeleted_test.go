package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/commandgroup/application/handlers"
	handlerstest "gomander/internal/commandgroup/application/handlers/test"
	"gomander/internal/commandgroup/domain/test"
	projectdomainevent "gomander/internal/project/domain/event"
)

// What the handler does to the Project's Command Groups is verified at the
// application seam, in internal/apptest.
func TestCleanCommandGroupsOnProjectDeleted(t *testing.T) {
	t.Run("Should do nothing if the event is of another type", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(handlerstest.MockEventEmitter)
		sut := handlers.NewCleanCommandGroupsOnProjectDeleted(mockRepo, mockEventEmitter)

		// Act
		err := sut.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockEventEmitter.AssertExpectations(t)
	})
}

func TestCleanCommandGroupsOnProjectDeleted_GetEvent(t *testing.T) {
	// Arrange
	sut := handlers.NewCleanCommandGroupsOnProjectDeleted(nil, nil)

	// Act
	event := sut.GetEvent()

	// Assert
	_, ok := event.(projectdomainevent.ProjectDeletedEvent)
	assert.True(t, ok)
}
