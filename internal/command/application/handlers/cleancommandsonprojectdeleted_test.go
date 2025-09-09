package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/command/application/handlers"
	"gomander/internal/command/domain/test"
	projectdomainevent "gomander/internal/project/domain/event"
)

type FakeEvent struct{}

func (FakeEvent) GetName() string { return "fake" }

var pjId = "pj-123"

func TestDefaultCleanCommandsOnProjectDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandRepository)
		handler := handlers.NewCleanCommandOnProjectDeleted(mockRepo)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		mockRepo.On("DeleteAll", pjId).Return(nil).Once()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should return error if failing to remove command from command groups", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandRepository)
		handler := handlers.NewCleanCommandOnProjectDeleted(mockRepo)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		expectedErr := errors.New("remove error")
		mockRepo.On("DeleteAll", pjId).Return(expectedErr).Once()

		// Act
		err := handler.Execute(event)

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		// Arrange
		mockRepo := new(test.MockCommandRepository)
		handler := handlers.NewCleanCommandOnProjectDeleted(mockRepo)

		// Act
		err := handler.Execute(FakeEvent{})

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDefaultCleanCommandsOnProjectDeleted_GetEvent(t *testing.T) {
	// Arrange
	handler := handlers.NewCleanCommandOnProjectDeleted(nil)

	// Act
	event := handler.GetEvent()

	// Assert
	_, ok := event.(projectdomainevent.ProjectDeletedEvent)
	assert.True(t, ok)
}
