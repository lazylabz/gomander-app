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

func TestApp_GetCommands(t *testing.T) {
	t.Run("Should return the commands provided by the repository", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		command1Data := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		command2Data := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(1).
			Data()

		expectedCommandGroup := []commanddomain.Command{
			commandDataToDomain(command1Data),
			commandDataToDomain(command2Data),
		}

		mockCommandRepository.On("GetAll", projectId).Return(expectedCommandGroup, nil)

		got, err := a.GetCommands()
		assert.NoError(t, err)
		assert.Equal(t, got, expectedCommandGroup)

		mock.AssertExpectationsForObjects(t, mockCommandRepository, mockUserConfigRepository)
	})
}

func TestApp_AddCommand(t *testing.T) {
	t.Run("Should add the command", func(t *testing.T) {
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

		existingCommandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		newCommandData := testutils.
			NewCommand().
			WithProjectId(projectId)

		parameterCommand := commandDataToDomain(newCommandData.Data())

		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{
			commandDataToDomain(existingCommandData),
		}, nil)
		expectedCommandCall := commandDataToDomain(newCommandData.WithPosition(1).Data())
		mockCommandRepository.On("Create", &expectedCommandCall).Return(nil)

		mockLogger.On("Info", mock.Anything).Return(nil)

		err := a.AddCommand(parameterCommand)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to get all commands", func(t *testing.T) {
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

		newCommandData := testutils.
			NewCommand().
			WithProjectId(projectId)

		parameterCommand := commandDataToDomain(newCommandData.Data())

		mockCommandRepository.On("GetAll", projectId).Return(make([]commanddomain.Command, 0), errors.New("failed to get commands"))
		mockLogger.On("Error", mock.Anything).Return()

		err := a.AddCommand(parameterCommand)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to create commands", func(t *testing.T) {
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

		existingCommandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		newCommandData := testutils.
			NewCommand().
			WithProjectId(projectId)

		parameterCommand := commandDataToDomain(newCommandData.Data())

		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{
			commandDataToDomain(existingCommandData),
		}, nil)
		expectedCommandCall := commandDataToDomain(newCommandData.WithPosition(1).Data())
		mockCommandRepository.On("Create", &expectedCommandCall).Return(errors.New("failed to create command"))

		mockLogger.On("Error", mock.Anything).Return()

		err := a.AddCommand(parameterCommand)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
		)
	})
}

func TestApp_RemoveCommand(t *testing.T) {
	t.Run("Should remove the command", func(t *testing.T) {
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

		err := a.RemoveCommand(commandId)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockCommandRepository,
			mockEventBus,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to remove the command", func(t *testing.T) {
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

		err := a.RemoveCommand(commandId)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should return an error if side effect fail", func(t *testing.T) {
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

		err := a.RemoveCommand(commandId)
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

		err := a.EditCommand(commandToEdit)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to edit the command", func(t *testing.T) {
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

		err := a.EditCommand(commandToEdit)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
}

func TestApp_ReorderCommands(t *testing.T) {
	t.Run("Should reorder commands", func(t *testing.T) {
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

		err := a.ReorderCommands(orderedIds)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to retrieve commands", func(t *testing.T) {
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

		err := a.ReorderCommands([]string{})
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to update commands", func(t *testing.T) {
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

		err := a.ReorderCommands(orderedIds)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
}

func TestApp_RunCommand(t *testing.T) {
	t.Run("Should run the command", func(t *testing.T) {
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
			LastOpenedProjectId: "",
			EnvironmentPaths:    envPaths,
		}, nil)

		mockUserConfigRepository.On("Update", mock.Anything).Return(nil)

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

		err := a.OpenProject(projectId)
		assert.NoError(t, err)

		err = a.RunCommand(cmd.Id)
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

		err := a.RunCommand(cmdId)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if failing to retrieve the user config", func(t *testing.T) {
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

		err := a.RunCommand(cmd.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if failing to retrieve the project", func(t *testing.T) {
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

		err := a.RunCommand(cmd.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockProjectRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if failing to run the command", func(t *testing.T) {
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
			LastOpenedProjectId: "",
			EnvironmentPaths:    envPaths,
		}, nil)

		mockUserConfigRepository.On("Update", mock.Anything).Return(nil)

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

		err := a.OpenProject(projectId)
		assert.NoError(t, err)

		err = a.RunCommand(cmd.Id)
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

		a.StopCommand(cmd.Id)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
			mockEventEmitter,
			mockRunner,
		)
	})

	t.Run("Should return error if the command does not exist", func(t *testing.T) {
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

		a.StopCommand(commandId)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
		)
	})
	t.Run("Should return error if fails to stop the command", func(t *testing.T) {
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

		a.StopCommand(cmd.Id)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockLogger,
			mockRunner,
		)
	})
}

func TestApp_DuplicateCommand(t *testing.T) {
	t.Run("Should duplicate the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockEventBus := new(MockEventBus)
		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			EventBus:          mockEventBus,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		originalCmd := commandDataToDomain(originalCmdData)

		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		mockCommandRepository.On("Create", mock.MatchedBy(func(cmd *commanddomain.Command) bool {
			return cmd.Id != originalCmd.Id &&
				cmd.ProjectId == originalCmd.ProjectId &&
				cmd.Name == originalCmd.Name+" (copy)" &&
				cmd.Command == originalCmd.Command &&
				cmd.WorkingDirectory == originalCmd.WorkingDirectory &&
				cmd.Position == 1
		})).Return(nil)

		mockEventBus.On("PublishSync", mock.MatchedBy(func(event interface{}) bool {
			e, ok := event.(commanddomainevent.CommandDuplicatedEvent)
			return ok && e.InsideGroupId == ""
		})).Return(make([]error, 0))

		mockLogger.On("Info", mock.Anything).Return(nil)

		err := a.DuplicateCommand(originalCmd.Id, "")
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockEventBus,
			mockLogger,
		)
	})
	t.Run("Should duplicate the command with a target group", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockEventBus := new(MockEventBus)
		mockLogger := new(MockLogger)

		projectId := "project1"
		targetGroupId := "group1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			EventBus:          mockEventBus,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		originalCmd := commandDataToDomain(originalCmdData)

		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		mockCommandRepository.On("Create", mock.MatchedBy(func(cmd *commanddomain.Command) bool {
			return cmd.Id != originalCmd.Id &&
				cmd.ProjectId == originalCmd.ProjectId &&
				cmd.Name == originalCmd.Name+" (copy)" &&
				cmd.Command == originalCmd.Command &&
				cmd.WorkingDirectory == originalCmd.WorkingDirectory &&
				cmd.Position == 1
		})).Return(nil)

		mockEventBus.On("PublishSync", mock.MatchedBy(func(event interface{}) bool {
			e, ok := event.(commanddomainevent.CommandDuplicatedEvent)
			return ok && e.InsideGroupId == targetGroupId
		})).Return(make([]error, 0))

		mockLogger.On("Info", mock.Anything).Return(nil)

		err := a.DuplicateCommand(originalCmd.Id, targetGroupId)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockEventBus,
			mockLogger,
		)
	})
	t.Run("Should return an error if the command does not exist", func(t *testing.T) {
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

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: "project1"}, nil)
		expectedErr := errors.New("command not found")
		mockCommandRepository.On("Get", commandId).Return(nil, expectedErr)
		mockLogger.On("Error", "command not found").Return()

		err := a.DuplicateCommand(commandId, "")
		assert.ErrorIs(t, err, expectedErr)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to get all commands", func(t *testing.T) {
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

		originalCmd := commandDataToDomain(testutils.NewCommand().Data())
		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)

		expectedErr := errors.New("failed to get commands")
		mockCommandRepository.On("GetAll", projectId).Return(nil, expectedErr)

		mockLogger.On("Error", "failed to get commands").Return()

		err := a.DuplicateCommand(originalCmd.Id, "")
		assert.ErrorIs(t, err, expectedErr)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to create the duplicated command", func(t *testing.T) {
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

		originalCmd := commandDataToDomain(testutils.NewCommand().Data())
		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		expectedErr := errors.New("failed to create command")
		mockCommandRepository.On("Create", mock.Anything).Return(expectedErr)
		mockLogger.On("Error", "failed to create command").Return()

		err := a.DuplicateCommand(originalCmd.Id, "")
		assert.ErrorIs(t, err, expectedErr)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if side effects fail", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockEventBus := new(MockEventBus)
		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockUserConfigRepository,
			EventBus:          mockEventBus,
			Logger:            mockLogger,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		originalCmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		originalCmd := commandDataToDomain(originalCmdData)

		mockCommandRepository.On("Get", originalCmd.Id).Return(&originalCmd, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]commanddomain.Command{originalCmd}, nil)

		mockCommandRepository.On("Create", mock.Anything).Return(nil)

		sideEffectErrorMsg := "Something happened"
		mockEventBus.On("PublishSync", mock.MatchedBy(func(event interface{}) bool {
			e, ok := event.(commanddomainevent.CommandDuplicatedEvent)
			return ok && e.InsideGroupId == ""
		})).Return([]error{errors.New(sideEffectErrorMsg)})

		mockLogger.On("Error", sideEffectErrorMsg).Return()

		err := a.DuplicateCommand(originalCmd.Id, "")
		assert.ErrorContains(t, err, sideEffectErrorMsg)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockUserConfigRepository,
			mockEventBus,
			mockLogger,
		)
	})
}
