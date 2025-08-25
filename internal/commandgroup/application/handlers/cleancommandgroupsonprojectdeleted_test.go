package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/commandgroup/application/handlers"
	projectdomainevent "gomander/internal/project/domain/event"
)

var pjId = "pj-123"

func TestDefaultCleanCommandGroupsOnProjectDeleted(t *testing.T) {
	t.Run("Should remove command from command groups and delete empty groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultCleanCommandGroupsOnProjectDeleted(mockRepo)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		mockRepo.On("DeleteAll", pjId).Return(nil).Once()

		err := handler.Execute(event)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should return error if failing to remove command from command groups", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultCleanCommandGroupsOnProjectDeleted(mockRepo)
		event := projectdomainevent.ProjectDeletedEvent{ProjectId: pjId}

		expectedErr := errors.New("remove error")
		mockRepo.On("DeleteAll", pjId).Return(expectedErr).Once()

		err := handler.Execute(event)
		assert.ErrorIs(t, err, expectedErr)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Should do nothing if command is the wrong type", func(t *testing.T) {
		mockRepo := new(MockCommandGroupRepository)
		handler := handlers.NewDefaultCleanCommandGroupsOnProjectDeleted(mockRepo)

		err := handler.Execute(FakeEvent{})
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDefaultCleanCommandGroupsOnProjectDeleted_GetEvent(t *testing.T) {
	handler := handlers.NewDefaultCleanCommandGroupsOnProjectDeleted(nil)
	event := handler.GetEvent()
	_, ok := event.(projectdomainevent.ProjectDeletedEvent)
	assert.True(t, ok)
}
