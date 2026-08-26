package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/config/domain"
	configtest "gomander/internal/config/domain/test"
	loggertest "gomander/internal/logger/test"
	runnertest "gomander/internal/runner/test"
)

func TestApp_Startup(t *testing.T) {
	t.Run("Should successfully load configuration", func(t *testing.T) {
		// Arrange
		mockLogger := new(loggertest.MockLogger)
		mockUserConfigRepository := new(configtest.MockConfigRepository)

		sut := app.NewApp(mockLogger, nil, mockUserConfigRepository, nil, app.EventHandlers{})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: "123"}, nil)

		// Act & Assert
		assert.NotPanics(t, func() {
			sut.Startup(context.Background())
		})

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})

	t.Run("Should panic if configuration loading fails", func(t *testing.T) {
		// Arrange
		mockLogger := new(loggertest.MockLogger)
		mockUserConfigRepository := new(configtest.MockConfigRepository)

		sut := app.NewApp(mockLogger, nil, mockUserConfigRepository, nil, app.EventHandlers{})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(nil, assert.AnError)

		// Act & Assert
		assert.Panics(t, func() {
			sut.Startup(context.Background())
		})

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})
}

func TestApp_OnBeforeClose(t *testing.T) {
	t.Run("Should stop all running commands and stop successfully", func(t *testing.T) {
		// Arrange
		mockCommandRunner := new(runnertest.MockRunner)
		mockLogger := new(loggertest.MockLogger)

		sut := app.NewApp(mockLogger, mockCommandRunner, nil, nil, app.EventHandlers{})

		mockCommandRunner.On("StopAllRunningCommands").Return([]error{})

		// Act
		prevent := sut.OnBeforeClose(context.Background())

		// Assert
		assert.False(t, prevent)
		mock.AssertExpectationsForObjects(t, mockCommandRunner, mockLogger)
	})

	t.Run("Should prevent closing if there are errors stopping commands", func(t *testing.T) {
		// Arrange
		mockCommandRunner := new(runnertest.MockRunner)
		mockLogger := new(loggertest.MockLogger)

		sut := app.NewApp(mockLogger, mockCommandRunner, nil, nil, app.EventHandlers{})

		errs := []error{assert.AnError}
		mockCommandRunner.On("StopAllRunningCommands").Return(errs)

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		prevent := sut.OnBeforeClose(context.Background())

		// Assert
		assert.True(t, prevent)

		mock.AssertExpectationsForObjects(t, mockCommandRunner, mockLogger)
	})
}
