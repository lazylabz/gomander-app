package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/config/domain"
	"gomander/internal/project/domain/event"
)

func TestApp_CloseProject(t *testing.T) {
	t.Run("Should close the current project", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		mockConfig := domain.Config{
			LastOpenedProjectId: "project1",
			EnvironmentPaths:    []domain.EnvironmentPath{{Id: "path1", Path: "TestPath"}},
		}
		mockUpdatedConfig := domain.Config{
			LastOpenedProjectId: "",
			EnvironmentPaths:    []domain.EnvironmentPath{{Id: "path1", Path: "TestPath"}},
		}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", &mockUpdatedConfig).Return(nil)

		// Act
		err := a.CloseProject()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockUserConfigRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockConfigRepository,
		})

		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))

		// Act
		err := a.CloseProject()

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})

	t.Run("Should return an error if updating the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockUserConfigRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockConfigRepository,
		})

		mockConfig := domain.Config{LastOpenedProjectId: "project1"}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", mock.Anything).Return(errors.New("update error"))

		// Act
		err := a.CloseProject()

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})
}

func TestApp_DeleteProject(t *testing.T) {
	t.Run("Should delete a project and all its commands", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)
		mockLogger := new(MockLogger)
		mockEventBus := new(MockEventBus)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			Logger:            mockLogger,
			EventBus:          mockEventBus,
		})

		projectId := "1"

		mockEventBus.On("PublishSync", event.NewProjectDeletedEvent(projectId)).Return(make([]error, 0))
		mockProjectRepository.On("Delete", projectId).Return(nil)

		// Act
		err := a.DeleteProject(projectId)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger, mockEventBus)
	})

	t.Run("Should return an error if deleting the project fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			Logger:            mockLogger,
		})

		projectId := "1"

		mockProjectRepository.On("Delete", projectId).Return(errors.New("some error occurred"))

		// Act
		err := a.DeleteProject(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger)
	})

	t.Run("Should return an error if an async event handler fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)
		mockLogger := new(MockLogger)
		mockEventBus := new(MockEventBus)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			Logger:            mockLogger,
			EventBus:          mockEventBus,
		})

		projectId := "1"

		mockProjectRepository.On("Delete", projectId).Return(nil)
		mockEventBus.On("PublishSync", event.NewProjectDeletedEvent(projectId)).Return([]error{errors.New("handler error")})

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.DeleteProject(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger, mockEventBus)
	})
}
