package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/commandgroup/application/handlers"
	"gomander/internal/commandgroup/domain/test"
)

var cmdId = "cmd-123"

func TestDefaultCleanCommandGroupsOnCommandDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		mockRepo.On("DeleteEmpty").Return(nil).Once()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should return error if failing to remove command from command groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		expectedErr := errors.New("remove error")
		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(expectedErr).Once()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should return error if failing to remove empty groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		expectedErr := errors.New("delete empty error")
		mockRepo.On("DeleteEmpty").Return(expectedErr).Once()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)

		// Act
		err := handler.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDefaultCleanCommandGroupsOnCommandDeleted_GetEvent(t *testing.T) {
	// Arrange
	handler := handlers.NewCleanCommandGroupsOnCommandDeleted(nil)

	// Act
	event := handler.GetEvent()

	// Assert
	_, ok := event.(commanddomainevent.CommandDeletedEvent)
	assert.True(t, ok)
}
