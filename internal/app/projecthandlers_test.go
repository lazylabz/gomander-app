package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/config/domain"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/project/domain/event"
)

func TestApp_OpenProject(t *testing.T) {
	t.Run("Should open a project successfully", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		projectId := "project1"
		projectId2 := "project2"

		mockConfig := domain.Config{
			LastOpenedProjectId: projectId,
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "path1",
					Path: "TestPath",
				},
			},
		}
		mockUpdatedConfig := domain.Config{
			LastOpenedProjectId: projectId2,
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "path1",
					Path: "TestPath",
				},
			},
		}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", &mockUpdatedConfig).Return(nil)
		mockProjectRepository.On("Get", projectId2).Return(&projectdomain.Project{Id: projectId2}, nil).Once()

		// Act
		err := a.OpenProject(projectId2)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error if project does not exist", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		projectId := "nonexistent"
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("project not found"))

		// Act
		err := a.OpenProject(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})

	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		projectId := "project1"
		mockProjectRepository.On("Get", projectId).Return(nil, nil)
		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))

		// Act
		err := a.OpenProject(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})

	t.Run("Should return an error if updating the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		projectId := "project1"
		mockConfig := domain.Config{LastOpenedProjectId: projectId}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", mock.Anything).Return(errors.New("update error"))
		mockProjectRepository.On("Get", projectId).Return(&projectdomain.Project{Id: projectId}, nil)

		// Act
		err := a.OpenProject(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})
}

func TestApp_CreateProject(t *testing.T) {
	t.Run("Should create a project successfully", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Create", project).Return(nil)

		// Act
		err := a.CreateProject(project)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})

	t.Run("Should return an error if project creation fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Create", project).Return(errors.New("fail"))

		// Act
		err := a.CreateProject(project)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}

func TestApp_EditProject(t *testing.T) {
	t.Run("Should edit a project successfully", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Update", project).Return(nil)

		// Act
		err := a.EditProject(project)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})

	t.Run("Should return an error if project update fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Update", project).Return(errors.New("fail"))

		// Act
		err := a.EditProject(project)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}

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
