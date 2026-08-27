package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/commandgroup/application/handlers"
	"gomander/internal/commandgroup/domain/test"
	test2 "gomander/internal/event/test"
)

// What the handler does to the Command Groups is verified at the application
// seam, in internal/apptest: it takes several of them to say anything true
// about a rule that spans several of them.
func TestDefaultCleanCommandGroupsOnCommandDeleted(t *testing.T) {
	t.Run("Should do nothing if the event is of another type", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(test2.MockEventEmitter)
		sut := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo, mockEventEmitter)

		// Act
		err := sut.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockEventEmitter.AssertExpectations(t)
	})
}

func TestDefaultCleanCommandGroupsOnCommandDeleted_GetEvent(t *testing.T) {
	// Arrange
	sut := handlers.NewCleanCommandGroupsOnCommandDeleted(nil, nil)

	// Act
	event := sut.GetEvent()

	// Assert
	_, ok := event.(commanddomainevent.CommandDeletedEvent)
	assert.True(t, ok)
}
