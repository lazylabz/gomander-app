package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/config/domain"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/testutils"
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

		mockProjectRepository.On("Get", "").Return(project, nil)

		got, err := a.GetCurrentProject()
		assert.NoError(t, err)
		assert.Equal(t, project, got)
	})
	t.Run("Should return an error if project does not exist", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockConfigRepository := new(MockUserConfigRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			ConfigRepository:  mockConfigRepository,
		})

		mockProjectRepository.On("Get", "").Return(nil, errors.New("project not found"))

		_, err := a.GetCurrentProject()
		assert.Error(t, err)
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
	})
	t.Run("Should return an error if project does not exist", func(t *testing.T) {
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		projectId := "nonexistent"
		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{}, nil)
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("project not found"))

		err := a.OpenProject(projectId)
		assert.Error(t, err)
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
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("project not found"))
		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))

		err := a.OpenProject(projectId)
		assert.Error(t, err)
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
	})
	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		mockConfigRepository := new(MockUserConfigRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockConfigRepository,
		})

		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))
		assert.Error(t, a.CloseProject())
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
	})
}

func TestApp_DeleteProject(t *testing.T) {
	t.Run("Should delete a project and all its commands", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
			Logger:                 mockLogger,
		})

		projectId := "1"
		mockProjectRepository.On("Get", projectId).Return(&projectdomain.Project{Id: projectId}, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: projectId}, nil)
		mockUserConfigRepository.On("Update", mock.Anything).Return(nil)

		err := a.OpenProject(projectId)
		assert.NoError(t, err)

		cmd1 := commandDataToDomain(testutils.NewCommand().WithProjectId(projectId).Data())
		cmd2 := commandDataToDomain(testutils.NewCommand().WithProjectId(projectId).Data())
		commands := []commanddomain.Command{
			cmd1,
			cmd2,
		}

		mockCommandRepository.On("GetAll", projectId).Return(commands, nil)
		mockCommandRepository.On("Delete", cmd1.Id).Return(nil)
		mockCommandRepository.On("Delete", cmd2.Id).Return(nil)
		mockProjectRepository.On("Delete", projectId).Return(nil)
		mockLogger.On("Info", mock.Anything).Return()

		mockCommandGroupRepository.On("GetAll", projectId).Return([]commandgroupdomain.CommandGroup{}, nil)

		assert.NoError(t, a.DeleteProject(projectId))
	})
	t.Run("Should return an error if getting project commands fails", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository: mockProjectRepository,
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
		})

		projectId := "1"
		mockProjectRepository.On("Get", projectId).Return(&projectdomain.Project{Id: projectId}, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: projectId}, nil)

		mockCommandRepository.On("GetAll", projectId).Return(make([]commanddomain.Command, 0), errors.New("fail"))

		err := a.DeleteProject(projectId)
		assert.Error(t, err)
	})
	t.Run("Should return an error if deleting a command fails", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			ConfigRepository:       mockUserConfigRepository,
			Logger:                 mockLogger,
			CommandGroupRepository: mockCommandGroupRepository,
		})

		projectId := "1"
		mockCommandGroupRepository.On("GetAll", mock.Anything).Return([]commandgroupdomain.CommandGroup{}, nil)

		cmd1 := commandDataToDomain(testutils.NewCommand().WithProjectId(projectId).Data())
		cmd2 := commandDataToDomain(testutils.NewCommand().WithProjectId(projectId).Data())
		commands := []commanddomain.Command{
			cmd1,
			cmd2,
		}

		mockCommandRepository.On("GetAll", projectId).Return(commands, nil)
		mockCommandRepository.On("Delete", cmd1.Id).Return(errors.New("fail"))

		mockLogger.On("Error", mock.Anything).Return()

		err := a.DeleteProject(projectId)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should return an error if deleting the project fails", func(t *testing.T) {
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
			Logger:                 mockLogger,
		})

		projectId := "1"
		mockProjectRepository.On("Get", projectId).Return(&projectdomain.Project{Id: projectId}, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: projectId}, nil)
		mockUserConfigRepository.On("Update", mock.Anything).Return(nil)

		err := a.OpenProject(projectId)
		assert.NoError(t, err)

		cmd1 := commandDataToDomain(testutils.NewCommand().WithProjectId(projectId).Data())
		cmd2 := commandDataToDomain(testutils.NewCommand().WithProjectId(projectId).Data())
		commands := []commanddomain.Command{
			cmd1,
			cmd2,
		}

		mockCommandRepository.On("GetAll", projectId).Return(commands, nil)
		mockCommandRepository.On("Delete", cmd1.Id).Return(nil)
		mockCommandRepository.On("Delete", cmd2.Id).Return(nil)
		mockProjectRepository.On("Delete", projectId).Return(errors.New("some error occurred"))
		mockLogger.On("Info", mock.Anything).Return()

		mockCommandGroupRepository.On("GetAll", projectId).Return([]commandgroupdomain.CommandGroup{}, nil)

		assert.Error(t, a.DeleteProject(projectId))
	})
}
