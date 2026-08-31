package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/commandgroup/application/handlers"
	handlerstest "gomander/internal/commandgroup/application/handlers/test"
	"gomander/internal/commandgroup/domain/test"
	unitofworktest "gomander/internal/unitofwork/test"
)

// What the handler does to the Command Groups is verified at the application
// seam, in internal/apptest: it takes several of them to say anything true
// about a rule that spans several of them.
func TestCleanCommandGroupsOnCommandDeleted(t *testing.T) {
	t.Run("Should do nothing if the event is of another type", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(handlerstest.MockEventEmitter)
		sut := handlers.NewCleanCommandGroupsOnCommandDeleted(unitofworktest.NewMockUnitOfWork(nil, nil, mockRepo), mockRepo, mockEventEmitter)

		// Act
		err := sut.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockEventEmitter.AssertExpectations(t)
	})
}

func TestCleanCommandGroupsOnCommandDeleted_GetEvent(t *testing.T) {
	// Arrange
	sut := handlers.NewCleanCommandGroupsOnCommandDeleted(nil, nil, nil)

	// Act
	event := sut.GetEvent()

	// Assert
	_, ok := event.(commanddomainevent.CommandDeletedEvent)
	assert.True(t, ok)
}
