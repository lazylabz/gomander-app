package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/event"
	"gomander/internal/testutils"
)

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
