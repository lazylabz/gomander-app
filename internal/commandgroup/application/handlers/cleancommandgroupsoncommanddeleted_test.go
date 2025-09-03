package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	commanddomainevent "gomander/internal/command/domain/event"
	"gomander/internal/commandgroup/application/handlers"
)

var cmdId = "cmd-123"

func TestDefaultCleanCommandGroupsOnCommandDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		mockRepo.On("DeleteEmpty").Return(nil).Once()

		err := handler.Execute(event)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should return error if failing to remove command from command groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		expectedErr := errors.New("remove error")
		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(expectedErr).Once()

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should return error if failing to remove empty groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)
		event := commanddomainevent.CommandDeletedEvent{CommandId: cmdId}

		mockRepo.On("RemoveCommandFromCommandGroups", cmdId).Return(nil).Once()
		expectedErr := errors.New("delete empty error")
		mockRepo.On("DeleteEmpty").Return(expectedErr).Once()

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewCleanCommandGroupsOnCommandDeleted(mockRepo)

		err := handler.Execute(FakeEvent{})
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDefaultCleanCommandGroupsOnCommandDeleted_GetEvent(t *testing.T) {
	handler := handlers.NewCleanCommandGroupsOnCommandDeleted(nil)
	event := handler.GetEvent()
	_, ok := event.(commanddomainevent.CommandDeletedEvent)
	assert.True(t, ok)
}
