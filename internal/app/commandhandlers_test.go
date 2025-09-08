package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	commanddomain "gomander/internal/command/domain"
	commanddomainevent "gomander/internal/command/domain/event"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/event"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/testutils"
)

func TestApp_RemoveCommand(t *testing.T) {
	t.Run("Should remove the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockEventBus := new(MockEventBus)

		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			Logger:                 mockLogger,
			EventBus:               mockEventBus,
		})

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(nil)
		mockEventBus.On("PublishSync", commanddomainevent.NewCommandDeletedEvent(commandId)).Return(make([]error, 0))

		mockLogger.On("Info", "Command removed: "+commandId).Return(nil)

		// Act
		err := a.RemoveCommand(commandId)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockCommandRepository,
			mockEventBus,
			mockLogger,
		)
	})

	t.Run("Should return an error if fails to remove the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)

		mockLogger := new(MockLogger)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository:      mockCommandRepository,
			Logger:                 mockLogger,
			CommandGroupRepository: mockCommandGroupRepository,
		})

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(errors.New("failed to delete command"))

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.RemoveCommand(commandId)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
			mockCommandGroupRepository,
		)
	})

	t.Run("Should return an error if side effect fail", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockEventBus := new(MockEventBus)

		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			Logger:                 mockLogger,
			EventBus:               mockEventBus,
		})

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(nil)
		mockEventBus.On("PublishSync", commanddomainevent.NewCommandDeletedEvent(commandId)).Return([]error{errors.New("Something happened")})

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.RemoveCommand(commandId)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockCommandRepository,
			mockEventBus,
			mockLogger,
		)
	})
}

func TestApp_EditCommand(t *testing.T) {
	t.Run("Should edit the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandToEdit := commandDataToDomain(commandData)

		mockCommandRepository.On("Update", &commandToEdit).Return(nil)

		mockLogger.On("Info", "Command edited: "+commandToEdit.Id).Return(nil)

		// Act
		err := a.EditCommand(commandToEdit)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})

	t.Run("Should return an error if fails to edit the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandToEdit := commandDataToDomain(commandData)

		mockCommandRepository.On("Update", &commandToEdit).Return(errors.New("failed to update command"))

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.EditCommand(commandToEdit)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
}

func TestApp_ReorderCommands(t *testing.T) {
	t.Run("Should reorder commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmd1 := testutils.NewCommand().WithProjectId(projectId).WithPosition(0)
		cmd2 := testutils.NewCommand().WithProjectId(projectId).WithPosition(1)
		cmd3 := testutils.NewCommand().WithProjectId(projectId).WithPosition(2)

		orderedIds := []string{cmd3.Data().Id, cmd1.Data().Id, cmd2.Data().Id}

		cm1WithUpdatedPosition := commandDataToDomain(cmd1.WithPosition(1).Data())
		cm2WithUpdatedPosition := commandDataToDomain(cmd2.WithPosition(2).Data())
		cm3WithUpdatedPosition := commandDataToDomain(cmd3.WithPosition(0).Data())

		mockCommandRepository.On("GetAll", projectId).Return(
			[]commanddomain.Command{
				commandDataToDomain(cmd1.Data()),
				commandDataToDomain(cmd2.Data()),
				commandDataToDomain(cmd3.Data()),
			}, nil,
		)

		mockCommandRepository.On("Update", &cm1WithUpdatedPosition).Return(nil)
		mockCommandRepository.On("Update", &cm2WithUpdatedPosition).Return(nil)
		mockCommandRepository.On("Update", &cm3WithUpdatedPosition).Return(nil)

		mockLogger.On("Info", "Commands reordered").Return(nil)

		// Act
		err := a.ReorderCommands(orderedIds)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})

	t.Run("Should return an error if fails to get the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
		})
		expectedErr := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedErr)

		// Act
		err := a.ReorderCommands([]string{})

		// Assert
		assert.ErrorIs(t, err, expectedErr)
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if fails to retrieve commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		mockCommandRepository.On("GetAll", projectId).Return(
			make([]commanddomain.Command, 0),
			errors.New("failed to get commands"))

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.ReorderCommands([]string{})

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})

	t.Run("Should return an error if fails to update commands", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmd1 := testutils.NewCommand().WithProjectId(projectId).WithPosition(0)
		cmd2 := testutils.NewCommand().WithProjectId(projectId).WithPosition(1)
		cmd3 := testutils.NewCommand().WithProjectId(projectId).WithPosition(2)

		orderedIds := []string{cmd3.Data().Id, cmd1.Data().Id, cmd2.Data().Id}

		cm1WithUpdatedPosition := commandDataToDomain(cmd1.WithPosition(1).Data())
		cm2WithUpdatedPosition := commandDataToDomain(cmd2.WithPosition(2).Data())
		cm3WithUpdatedPosition := commandDataToDomain(cmd3.WithPosition(0).Data())

		mockCommandRepository.On("GetAll", projectId).Return(
			[]commanddomain.Command{
				commandDataToDomain(cmd1.Data()),
				commandDataToDomain(cmd2.Data()),
				commandDataToDomain(cmd3.Data()),
			}, nil,
		)

		mockCommandRepository.On("Update", &cm1WithUpdatedPosition).Return(nil)
		mockCommandRepository.On("Update", &cm2WithUpdatedPosition).Return(errors.New("failed to update command"))
		mockCommandRepository.On("Update", &cm3WithUpdatedPosition).Return(nil)

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.ReorderCommands(orderedIds)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
}

func TestApp_RunCommand(t *testing.T) {
	t.Run("Should run the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)
		mockRunner := new(MockRunner)
		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			ProjectRepository: mockProjectRepository,
			Runner:            mockRunner,
			Logger:            mockLogger,
		})

		envPath := configdomain.EnvironmentPath{
			Id:   "1",
			Path: "/1",
		}

		envPaths := []configdomain.EnvironmentPath{envPath}

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{
			LastOpenedProjectId: projectId,
			EnvironmentPaths:    envPaths,
		}, nil)

		cmdData := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		project := projectdomain.Project{
			Id:               projectId,
			Name:             "Test Project",
			WorkingDirectory: "/working/dir",
		}
		mockProjectRepository.On("Get", cmd.ProjectId).Return(&project, nil)

		mockLogger.On("Info", mock.Anything).Return()

		mockRunner.On("RunCommand", &cmd, []string{"/1"}, project.WorkingDirectory).Return(nil)

		// Act
		err := a.RunCommand(cmd.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockProjectRepository,
			mockRunner,
			mockLogger,
		)
	})

	t.Run("Should return an error if failing to retrieve the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})
		cmdId := "command1"

		mockCommandRepository.On("Get", cmdId).Return(nil, errors.New("command not found"))

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.RunCommand(cmdId)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})

	t.Run("Should return an error if failing to retrieve the user config", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})

		cmdData := testutils.NewCommand().Data()
		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(nil, errors.New("failed to get user config"))

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.RunCommand(cmd.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
		)
	})

	t.Run("Should return an error if failing to retrieve the project", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			ProjectRepository: mockProjectRepository,
			Logger:            mockLogger,
		})

		projectId := "project1"
		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		cmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{}, nil)
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("failed to get project"))

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		err := a.RunCommand(cmd.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockProjectRepository,
			mockLogger,
		)
	})

	t.Run("Should return an error if failing to run the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)
		mockRunner := new(MockRunner)
		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			ProjectRepository: mockProjectRepository,
			Runner:            mockRunner,
			Logger:            mockLogger,
		})

		envPath := configdomain.EnvironmentPath{
			Id:   "1",
			Path: "/1",
		}

		envPaths := []configdomain.EnvironmentPath{envPath}

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{
			LastOpenedProjectId: projectId,
			EnvironmentPaths:    envPaths,
		}, nil)

		cmdData := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		project := projectdomain.Project{
			Id:               projectId,
			Name:             "Test Project",
			WorkingDirectory: "/working/dir",
		}
		mockProjectRepository.On("Get", cmd.ProjectId).Return(&project, nil)

		mockLogger.On("Error", mock.Anything).Return()

		mockRunner.On("RunCommand", &cmd, []string{"/1"}, project.WorkingDirectory).Return(errors.New("failed to run command"))

		// Act
		err := a.RunCommand(cmd.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockProjectRepository,
			mockLogger,
		)
	})
}

func TestApp_StopCommand(t *testing.T) {
	t.Run("Should stop the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)
		mockEventEmitter := new(MockEventEmitter)
		mockRunner := new(MockRunner)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
			EventEmitter:      mockEventEmitter,
			Runner:            mockRunner,
		})

		cmdData := testutils.NewCommand().WithProjectId("project1").Data()
		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)

		mockLogger.On("Info", mock.Anything).Return(nil)
		mockEventEmitter.On("EmitEvent", event.ProcessFinished, cmd.Id).Return(nil)

		mockRunner.On("StopRunningCommand", cmd.Id).Return(nil)

		// Act
		a.StopCommand(cmd.Id)

		// Assert
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
			mockEventEmitter,
			mockRunner,
		)
	})

	t.Run("Should return error if the command does not exist", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
		})

		commandId := "non-existing-command"

		mockCommandRepository.On("Get", commandId).Return(nil, errors.New("command not found"))

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		a.StopCommand(commandId)

		// Assert
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
		)
	})

	t.Run("Should return error if fails to stop the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)
		mockRunner := new(MockRunner)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			Logger:            mockLogger,
			Runner:            mockRunner,
		})

		cmdData := testutils.NewCommand().WithProjectId("project1").Data()
		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)

		mockLogger.On("Error", mock.Anything).Return()

		mockRunner.On("StopRunningCommand", cmd.Id).Return(errors.New("failed to stop command"))

		// Act
		a.StopCommand(cmd.Id)

		// Assert
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
			mockRunner,
		)
	})
}
