package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/project/application/usecases"
	"gomander/internal/project/domain/event"
	"gomander/internal/project/domain/test"
)

func TestDefaultDeleteProject_Execute(t *testing.T) {
	t.Run("Should delete a project and all its commands", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockLogger := new(MockLogger)
		mockEventBus := new(MockEventBus)

		sut := usecases.NewDeleteProject(
			mockProjectRepository,
			mockEventBus,
			mockLogger,
		)

		projectId := "1"

		mockEventBus.On("PublishSync", event.NewProjectDeletedEvent(projectId)).Return(make([]error, 0))
		mockProjectRepository.On("Delete", projectId).Return(nil)

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger, mockEventBus)
	})

	t.Run("Should return an error if deleting the project fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockEventBus := new(MockEventBus)
		mockLogger := new(MockLogger)

		sut := usecases.NewDeleteProject(
			mockProjectRepository,
			mockEventBus,
			mockLogger,
		)

		projectId := "1"

		mockProjectRepository.On("Delete", projectId).Return(errors.New("some error occurred"))

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger)
	})

	t.Run("Should return an error if an async event handler fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockLogger := new(MockLogger)
		mockEventBus := new(MockEventBus)

		sut := usecases.NewDeleteProject(
			mockProjectRepository,
			mockEventBus,
			mockLogger,
		)

		projectId := "1"

		mockProjectRepository.On("Delete", projectId).Return(nil)
		mockEventBus.On("PublishSync", event.NewProjectDeletedEvent(projectId)).Return([]error{errors.New("handler error")})

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger, mockEventBus)
	})
}
