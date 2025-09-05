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

func TestApp_GetCurrentProject(t *testing.T) {
	t.Run("Should return the current project", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"
		project := &projectdomain.Project{Id: projectId, Name: "Test", WorkingDirectory: "/tmp"}

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			ConfigRepository:  mockConfigRepository,
		})

		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: projectId}, nil)
		mockProjectRepository.On("Get", projectId).Return(project, nil)

		got, err := a.GetCurrentProject()
		assert.NoError(t, err)
		assert.Equal(t, project, got)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockConfigRepository)
	})
	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		mockConfigRepository := new(MockUserConfigRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockConfigRepository,
		})
		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))
		_, err := a.GetCurrentProject()
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})
	t.Run("Should return an error if project does not exist", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockConfigRepository := new(MockUserConfigRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			ConfigRepository:  mockConfigRepository,
		})

		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: "nonexistent"}, nil)
		mockProjectRepository.On("Get", "nonexistent").Return(nil, errors.New("project not found"))

		_, err := a.GetCurrentProject()
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockConfigRepository)
	})
}

func TestApp_GetAvailableProjects(t *testing.T) {
	t.Run("Should return available projects", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		projects := []projectdomain.Project{{Id: "1", Name: "A", WorkingDirectory: "/a"}}
		mockProjectRepository.On("GetAll").Return(projects, nil)

		got, err := a.GetAvailableProjects()
		assert.NoError(t, err)
		assert.Equal(t, projects, got)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
	t.Run("Should return an error if fetching projects fails", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		mockProjectRepository.On("GetAll").Return(make([]projectdomain.Project, 0), errors.New("fail"))

		_, err := a.GetAvailableProjects()
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}

func TestApp_OpenProject(t *testing.T) {
	t.Run("Should open a project successfully", func(t *testing.T) {
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		projectId := "project1"
		mockConfig := domain.Config{
			LastOpenedProjectId: projectId,
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "path1",
					Path: "TestPath",
				},
			},
		}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", mock.Anything).Return(nil)
		mockProjectRepository.On("Get", projectId).Return(&projectdomain.Project{Id: projectId}, nil).Once()

		err := a.OpenProject(projectId)
		assert.NoError(t, err)

		// Assert that the project is saved in memory as the current project
		mockProjectRepository.On("Get", "project1").Return(&projectdomain.Project{}, nil).Once()
		_, err = a.GetCurrentProject()
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})
	t.Run("Should return an error if project does not exist", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		projectId := "nonexistent"
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("project not found"))

		err := a.OpenProject(projectId)
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
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

		err := a.OpenProject(projectId)
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})
	t.Run("Should return an error if updating the config fails", func(t *testing.T) {
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

		err := a.OpenProject(projectId)
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})
}

func TestApp_CreateProject(t *testing.T) {
	t.Run("Should create a project successfully", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Create", project).Return(nil)
		assert.NoError(t, a.CreateProject(project))
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
	t.Run("Should return an error if project creation fails", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Create", project).Return(errors.New("fail"))
		assert.Error(t, a.CreateProject(project))
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}

func TestApp_EditProject(t *testing.T) {
	t.Run("Should edit a project successfully", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Update", project).Return(nil)
		assert.NoError(t, a.EditProject(project))
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
	t.Run("Should return an error if project update fails", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
		})

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Update", project).Return(errors.New("fail"))
		assert.Error(t, a.EditProject(project))
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}

func TestApp_CloseProject(t *testing.T) {
	t.Run("Should close the current project", func(t *testing.T) {
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
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", mock.Anything).Return(nil)
		assert.NoError(t, a.CloseProject())

		// Assert that the project is closed in memory
		mockProjectRepository.On("Get", "").Return(&projectdomain.Project{}, nil).Once()
		_, err := a.GetCurrentProject()
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})
	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		mockConfigRepository := new(MockUserConfigRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockConfigRepository,
		})

		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))
		assert.Error(t, a.CloseProject())
		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})
	t.Run("Should return an error if updating the config fails", func(t *testing.T) {
		mockConfigRepository := new(MockUserConfigRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockConfigRepository,
		})

		mockConfig := domain.Config{LastOpenedProjectId: "project1"}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", mock.Anything).Return(errors.New("update error"))
		assert.Error(t, a.CloseProject())
		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})
}

func TestApp_DeleteProject(t *testing.T) {
	t.Run("Should delete a project and all its commands", func(t *testing.T) {
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

		assert.NoError(t, a.DeleteProject(projectId))
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger, mockEventBus)
	})
	t.Run("Should return an error if deleting the project fails", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			Logger:            mockLogger,
		})

		projectId := "1"

		mockProjectRepository.On("Delete", projectId).Return(errors.New("some error occurred"))

		assert.Error(t, a.DeleteProject(projectId))
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger)
	})
	t.Run("Should return an error if an async event handler fails", func(t *testing.T) {
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

		assert.Error(t, a.DeleteProject(projectId))
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockLogger, mockEventBus)
	})
}
