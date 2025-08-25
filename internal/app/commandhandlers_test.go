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

type MockCommandRepository struct {
	mock.Mock
}

func (m *MockCommandRepository) Get(commandId string) (*commanddomain.Command, error) {
	args := m.Called(commandId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commanddomain.Command), args.Error(1)
}

func (m *MockCommandRepository) GetAll(projectId string) ([]commanddomain.Command, error) {
	args := m.Called(projectId)
	return args.Get(0).([]commanddomain.Command), args.Error(1)
}

func (m *MockCommandRepository) Create(command *commanddomain.Command) error {
	args := m.Called(command)
	return args.Error(0)
}

func (m *MockCommandRepository) Update(command *commanddomain.Command) error {
	args := m.Called(command)
	return args.Error(0)
}

func (m *MockCommandRepository) Delete(commandId string) error {
	args := m.Called(commandId)
	return args.Error(0)
}

func TestApp_GetCommands(t *testing.T) {
	t.Run("Should return the commands provided by the repository", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
		})

		a.SetOpenProjectId(projectId)

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

		mock.AssertExpectationsForObjects(t, mockCommandRepository)
	})
}

func TestApp_AddCommand(t *testing.T) {
	t.Run("Should add the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to get all commands", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to create commands", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			Logger:                 mockLogger,
			EventBus:               mockEventBus,
		})

		a.SetOpenProjectId(projectId)

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(nil)
		mockEventBus.On("PublishSync", commanddomainevent.NewCommandDeletedEvent(commandId)).Return(make([]error, 0))

		mockLogger.On("Info", "Command removed: "+commandId).Return(nil)

		err := a.RemoveCommand(commandId)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if fails to remove the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)

		mockLogger := new(MockLogger)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository:      mockCommandRepository,
			Logger:                 mockLogger,
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

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

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			Logger:                 mockLogger,
			EventBus:               mockEventBus,
		})

		a.SetOpenProjectId(projectId)

		commandId := "command1"

		mockCommandRepository.On("Delete", commandId).Return(nil)
		mockEventBus.On("PublishSync", commanddomainevent.NewCommandDeletedEvent(commandId)).Return([]error{errors.New("Something happened")})

		mockLogger.On("Error", mock.Anything).Return()

		err := a.RemoveCommand(commandId)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockCommandRepository,
			mockLogger,
		)
	})
}

func TestApp_EditCommand(t *testing.T) {
	t.Run("Should edit the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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

		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		a.SetOpenProjectId(projectId)

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
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)
		mockRunner := new(MockRunner)
		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
			Runner:            mockRunner,
			Logger:            mockLogger,
		})

		envPath := configdomain.EnvironmentPath{
			Id:   "1",
			Path: "/1",
		}

		envPaths := []configdomain.EnvironmentPath{envPath}

		mockConfigRepository.On("GetOrCreate").Return(&configdomain.Config{
			LastOpenedProjectId: "",
			EnvironmentPaths:    envPaths,
		}, nil)

		mockConfigRepository.On("Update", mock.Anything).Return(nil)

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
			mockConfigRepository,
			mockProjectRepository,
			mockRunner,
		)
	})
	t.Run("Should return an error if failing to retrieve the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
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
		mockConfigRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockConfigRepository,
			Logger:            mockLogger,
		})

		cmdData := testutils.NewCommand().Data()
		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("failed to get user config"))

		mockLogger.On("Error", mock.Anything).Return()

		err := a.RunCommand(cmd.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockConfigRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if failing to retrieve the project", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
			Logger:            mockLogger,
		})

		projectId := "project1"
		a.SetOpenProjectId(projectId)

		cmdData := testutils.NewCommand().WithProjectId(projectId).Data()
		cmd := commandDataToDomain(cmdData)

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)
		mockConfigRepository.On("GetOrCreate").Return(&configdomain.Config{}, nil)
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("failed to get project"))

		mockLogger.On("Error", mock.Anything).Return()

		err := a.RunCommand(cmd.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockConfigRepository,
			mockProjectRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if failing to run the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)
		mockRunner := new(MockRunner)
		mockLogger := new(MockLogger)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			ConfigRepository:  mockConfigRepository,
			ProjectRepository: mockProjectRepository,
			Runner:            mockRunner,
			Logger:            mockLogger,
		})

		envPath := configdomain.EnvironmentPath{
			Id:   "1",
			Path: "/1",
		}

		envPaths := []configdomain.EnvironmentPath{envPath}

		mockConfigRepository.On("GetOrCreate").Return(&configdomain.Config{
			LastOpenedProjectId: "",
			EnvironmentPaths:    envPaths,
		}, nil)

		mockConfigRepository.On("Update", mock.Anything).Return(nil)

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
			mockConfigRepository,
			mockProjectRepository,
			mockLogger,
		)
	})
}

func TestApp_StopCommand(t *testing.T) {
	t.Run("Should stop the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockLogger := new(MockLogger)
		mockEventEmitter := new(MockEventEmitter)
		mockRunner := new(MockRunner)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
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
			mockLogger,
			mockEventEmitter,
			mockRunner,
		)
	})

	t.Run("Should return error if the command does not exist", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
			Logger:            mockLogger,
		})

		commandId := "non-existing-command"

		mockCommandRepository.On("Get", commandId).Return(nil, errors.New("command not found"))

		mockLogger.On("Error", mock.Anything).Return()

		a.StopCommand(commandId)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockLogger,
		)
	})
	t.Run("Should return error if fails to stop the command", func(t *testing.T) {
		mockCommandRepository := new(MockCommandRepository)
		mockLogger := new(MockLogger)
		mockRunner := new(MockRunner)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandRepository: mockCommandRepository,
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
			mockLogger,
			mockRunner,
		)
	})
}
