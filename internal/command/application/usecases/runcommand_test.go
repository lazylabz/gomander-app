package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/usecases"
	"gomander/internal/command/domain/test"
	configdomain "gomander/internal/config/domain"
	test3 "gomander/internal/config/domain/test"
	projectdomain "gomander/internal/project/domain"
	test2 "gomander/internal/project/domain/test"
)

func TestDefaultRunCommand_Execute(t *testing.T) {
	t.Run("Should run the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test2.MockProjectRepository)
		mockRunner := new(MockRunner)

		projectId := "project1"
		sut := usecases.NewRunCommand(mockUserConfigRepository, mockCommandRepository, mockProjectRepository, mockRunner)

		envPath := configdomain.EnvironmentPath{
			Id:   "1",
			Path: "/1",
		}

		envPaths := []configdomain.EnvironmentPath{envPath}

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{
			LastOpenedProjectId: projectId,
			EnvironmentPaths:    envPaths,
		}, nil)

		cmd := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(0).
			Build()

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		project := projectdomain.Project{
			Id:               projectId,
			Name:             "Test Project",
			WorkingDirectory: "/working/dir",
		}
		mockProjectRepository.On("Get", cmd.ProjectId).Return(&project, nil)

		mockRunner.On("RunCommand", &cmd, []string{"/1"}, project.WorkingDirectory).Return(nil)

		// Act
		err := sut.Execute(cmd.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockProjectRepository,
			mockRunner,
		)
	})

	t.Run("Should return an error if failing to retrieve the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test2.MockProjectRepository)
		mockRunner := new(MockRunner)

		sut := usecases.NewRunCommand(mockUserConfigRepository, mockCommandRepository, mockProjectRepository, mockRunner)

		cmdId := "command1"

		mockCommandRepository.On("Get", cmdId).Return(nil, errors.New("command not found"))

		// Act
		err := sut.Execute(cmdId)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test2.MockProjectRepository)
		mockRunner := new(MockRunner)

		sut := usecases.NewRunCommand(mockUserConfigRepository, mockCommandRepository, mockProjectRepository, mockRunner)

		cmd := test.NewCommandBuilder().Build()

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(nil, errors.New("failed to get user config"))

		// Act
		err := sut.Execute(cmd.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve the project", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test2.MockProjectRepository)
		mockRunner := new(MockRunner)

		sut := usecases.NewRunCommand(mockUserConfigRepository, mockCommandRepository, mockProjectRepository, mockRunner)

		projectId := "project1"
		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmd := test.NewCommandBuilder().WithProjectId(projectId).Build()

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("failed to get project"))

		// Act
		err := sut.Execute(cmd.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockProjectRepository,
		)
	})

	t.Run("Should return an error if failing to run the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test2.MockProjectRepository)
		mockRunner := new(MockRunner)

		projectId := "project1"
		sut := usecases.NewRunCommand(mockUserConfigRepository, mockCommandRepository, mockProjectRepository, mockRunner)

		envPath := configdomain.EnvironmentPath{
			Id:   "1",
			Path: "/1",
		}

		envPaths := []configdomain.EnvironmentPath{envPath}

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{
			LastOpenedProjectId: projectId,
			EnvironmentPaths:    envPaths,
		}, nil)

		cmd := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(0).
			Build()

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		project := projectdomain.Project{
			Id:               projectId,
			Name:             "Test Project",
			WorkingDirectory: "/working/dir",
		}
		mockProjectRepository.On("Get", cmd.ProjectId).Return(&project, nil)

		mockRunner.On("RunCommand", &cmd, []string{"/1"}, project.WorkingDirectory).Return(errors.New("failed to run command"))

		// Act
		err := sut.Execute(cmd.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockProjectRepository,
		)
	})
}
