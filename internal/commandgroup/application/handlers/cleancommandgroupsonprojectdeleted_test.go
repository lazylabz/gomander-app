package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/commandgroup/application/handlers"
	"gomander/internal/commandgroup/domain/test"
	event2 "gomander/internal/event"
	test2 "gomander/internal/event/test"
	projectdomainevent "gomander/internal/project/domain/event"
)

var pjId = "pj-123"

func TestDefaultCleanCommandGroupsOnProjectDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandGroupRepository)
		mockEventEmitter := new(test2.MockEventEmitter)
		handler := handlers.NewCleanCommandGroupsOnProjectDeleted(mockRepo, mockEventEmitter)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		deletedCmdGroupId := "cmdGroupId"
		mockRepo.On("DeleteAll", pjId).Return([]string{deletedCmdGroupId}, nil).Once()
		mockEventEmitter.On("EmitEvent", event2.CommandGroupDeleted, deletedCmdGroupId).Return().Once()

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
		handler := handlers.NewCleanCommandGroupsOnProjectDeleted(mockRepo, mockEventEmitter)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		expectedErr := errors.New("remove error")
		mockRepo.On("DeleteAll", pjId).Return(make([]string, 0), expectedErr).Once()

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
		handler := handlers.NewCleanCommandGroupsOnProjectDeleted(mockRepo, mockEventEmitter)

		// Act
		err := handler.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRepo, mockEventEmitter)
	})
}

func TestDefaultCleanCommandGroupsOnProjectDeleted_GetEvent(t *testing.T) {
	// Arrange
	handler := handlers.NewCleanCommandGroupsOnProjectDeleted(nil, nil)

	// Act
	event := handler.GetEvent()

	// Assert
	_, ok := event.(projectdomainevent.ProjectDeletedEvent)
	assert.True(t, ok)
}
