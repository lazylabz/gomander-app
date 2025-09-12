package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/commandgroup/application/handlers"
	"gomander/internal/commandgroup/domain/test"
	event2 "gomander/internal/event"
	test2 "gomander/internal/event/test"
)

var cmdId = "cmd-123"

func TestDefaultCleanCommandGroupsOnCommandDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(test2.MockEventEmitter)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo, mockEventEmitter)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		deletedCmdGroupId := "cmdGroupId"
		mockRepo.On("DeleteEmpty").Return([]string{deletedCmdGroupId}, nil).Once()
		mockEventEmitter.On("EmitEvent", event2.CommandGroupDeleted, deletedCmdGroupId).Return()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRepo, mockEventEmitter)
	})

	t.Run("Should return error if failing to remove command from command groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(test2.MockEventEmitter)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo, mockEventEmitter)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		expectedErr := errors.New("remove error")
		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(expectedErr).Once()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mock.AssertExpectationsForObjects(t, mockRepo, mockEventEmitter)
	})

	t.Run("Should return error if failing to remove empty groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(test2.MockEventEmitter)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo, mockEventEmitter)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		expectedErr := errors.New("delete empty error")
		mockRepo.On("DeleteEmpty").Return(make([]string, 0), expectedErr).Once()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mock.AssertExpectationsForObjects(t, mockRepo, mockEventEmitter)
	})

	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(test2.MockEventEmitter)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo, mockEventEmitter)

		// Act
		err := handler.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRepo, mockEventEmitter)
	})
}

func TestDefaultCleanCommandGroupsOnCommandDeleted_GetEvent(t *testing.T) {
	// Arrange
	handler := handlers.NewCleanCommandGroupsOnCommandDeleted(nil, nil)

	// Act
	event := handler.GetEvent()

	// Assert
	_, ok := event.(commanddomainevent.CommandDeletedEvent)
	assert.True(t, ok)
}
