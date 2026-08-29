package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	eventbustest "gomander/internal/eventbus/test"
	"gomander/internal/project/application/usecases"
	"gomander/internal/project/domain/event"
	"gomander/internal/project/domain/test"
)

func TestDeleteProject_Execute(t *testing.T) {
	t.Run("Should delete a project and all its commands", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockEventBus := new(eventbustest.MockEventBus)

		sut := usecases.NewDeleteProject(
			mockProjectRepository,
			mockEventBus,
		)

		projectId := "1"

		mockEventBus.On("PublishSync", event.NewProjectDeletedEvent(projectId)).Return(make([]error, 0))
		mockProjectRepository.On("Delete", projectId).Return(nil)

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockEventBus)
	})

	t.Run("Should return an error if deleting the project fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockEventBus := new(eventbustest.MockEventBus)

		sut := usecases.NewDeleteProject(
			mockProjectRepository,
			mockEventBus,
		)

		projectId := "1"

		mockProjectRepository.On("Delete", projectId).Return(errors.New("some error occurred"))

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockEventBus)
	})

	t.Run("Should report a failing handler only through the returned error", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockEventBus := new(eventbustest.MockEventBus)

		sut := usecases.NewDeleteProject(
			mockProjectRepository,
			mockEventBus,
		)

		projectId := "1"
		handlerErr := errors.New("handler error")

		mockProjectRepository.On("Delete", projectId).Return(nil)
		mockEventBus.On("PublishSync", event.NewProjectDeletedEvent(projectId)).Return([]error{handlerErr})

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.ErrorIs(t, err, handlerErr)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockEventBus)
	})
}
