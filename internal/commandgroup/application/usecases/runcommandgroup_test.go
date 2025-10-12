// filepath: /Users/moises/Code/personal/gomander/internal/commandgroup/application/usecases/runcommandgroup_test.go
package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/domain/test"
	"gomander/internal/commandgroup/application/usecases"
	test2 "gomander/internal/commandgroup/domain/test"
	configdomain "gomander/internal/config/domain"
	test3 "gomander/internal/config/domain/test"
	projectdomain "gomander/internal/project/domain"
	test4 "gomander/internal/project/domain/test"
	test5 "gomander/internal/runner/test"
)

func TestDefaultRunCommandGroup_Execute(t *testing.T) {
	t.Run("Should run the command group", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test4.MockProjectRepository)
		mockRunner := new(test5.MockRunner)

		projectId := "project1"
		sut := usecases.NewRunCommandGroup(
			mockUserConfigRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockProjectRepository,
			mockRunner,
		)

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

		cmdGroup := test2.NewCommandGroupBuilder().
			WithId("group1").
			WithProjectId(projectId).
			WithPosition(0).
			WithCommands(cmd).
			Build()

		mockCommandGroupRepository.On("Get", cmdGroup.Id).Return(&cmdGroup, nil)

		project := projectdomain.Project{
			Id:               projectId,
			Name:             "Test Project",
			WorkingDirectory: "/working/dir",
		}
		mockProjectRepository.On("Get", cmdGroup.ProjectId).Return(&project, nil)

		mockRunner.On("RunCommands", cmdGroup.Commands, []string{"/1"}, project.WorkingDirectory, mock.AnythingOfType("[]string")).Return(nil)

		// Act
		err := sut.Execute(cmdGroup.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockUserConfigRepository,
			mockProjectRepository,
			mockRunner,
		)
	})

	t.Run("Should return an error if failing to retrieve the command group", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test4.MockProjectRepository)
		mockRunner := new(test5.MockRunner)

		sut := usecases.NewRunCommandGroup(
			mockUserConfigRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockProjectRepository,
			mockRunner,
		)

		cmdGroupId := "group1"

		mockCommandGroupRepository.On("Get", cmdGroupId).Return(nil, errors.New("command group not found"))

		// Act
		err := sut.Execute(cmdGroupId)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test4.MockProjectRepository)
		mockRunner := new(test5.MockRunner)

		sut := usecases.NewRunCommandGroup(
			mockUserConfigRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockProjectRepository,
			mockRunner,
		)

		cmdGroup := test2.NewCommandGroupBuilder().WithId("group1").Build()

		mockCommandGroupRepository.On("Get", cmdGroup.Id).Return(&cmdGroup, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(nil, errors.New("failed to get user config"))

		// Act
		err := sut.Execute(cmdGroup.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve the project", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test4.MockProjectRepository)
		mockRunner := new(test5.MockRunner)

		sut := usecases.NewRunCommandGroup(
			mockUserConfigRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockProjectRepository,
			mockRunner,
		)

		projectId := "project1"
		cmdGroup := test2.NewCommandGroupBuilder().WithId("group1").WithProjectId(projectId).Build()

		mockCommandGroupRepository.On("Get", cmdGroup.Id).Return(&cmdGroup, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("failed to get project"))

		// Act
		err := sut.Execute(cmdGroup.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
			mockProjectRepository,
		)
	})

	t.Run("Should return an error if failing to run the commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(test.MockCommandRepository)
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)
		mockProjectRepository := new(test4.MockProjectRepository)
		mockRunner := new(test5.MockRunner)

		projectId := "project1"
		sut := usecases.NewRunCommandGroup(
			mockUserConfigRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockProjectRepository,
			mockRunner,
		)

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

		cmdGroup := test2.NewCommandGroupBuilder().
			WithId("group1").
			WithProjectId(projectId).
			WithPosition(0).
			WithCommands(cmd).
			Build()

		mockCommandGroupRepository.On("Get", cmdGroup.Id).Return(&cmdGroup, nil)

		project := projectdomain.Project{
			Id:               projectId,
			Name:             "Test Project",
			WorkingDirectory: "/working/dir",
		}
		mockProjectRepository.On("Get", cmdGroup.ProjectId).Return(&project, nil)

		mockRunner.On("RunCommands", cmdGroup.Commands, []string{"/1"}, project.WorkingDirectory, mock.AnythingOfType("[]string")).
			Return(errors.New("failed to run commands"))

		// Act
		err := sut.Execute(cmdGroup.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockUserConfigRepository,
			mockProjectRepository,
			mockRunner,
		)
	})
}
