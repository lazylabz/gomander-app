package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	appmocks "gomander/internal/app/test"
	"gomander/internal/config/domain"
	configtest "gomander/internal/config/domain/test"
	executiontest "gomander/internal/execution/test"
)

func TestApp_Startup(t *testing.T) {
	t.Run("Should successfully load configuration", func(t *testing.T) {
		// Arrange
		mockLogger := new(appmocks.MockLogger)
		mockUserConfigRepository := new(configtest.MockConfigRepository)

		sut := app.NewApp(mockLogger, nil, mockUserConfigRepository, nil, app.EventHandlers{})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: "123"}, nil)

		// Act
		err := sut.Startup(context.Background())

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})

	t.Run("Should report an error if configuration loading fails", func(t *testing.T) {
		// Arrange
		mockLogger := new(appmocks.MockLogger)
		mockUserConfigRepository := new(configtest.MockConfigRepository)

		sut := app.NewApp(mockLogger, nil, mockUserConfigRepository, nil, app.EventHandlers{})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(nil, assert.AnError)

		// Act
		err := sut.Startup(context.Background())

		// Assert
		assert.ErrorIs(t, err, assert.AnError)

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})
}

func TestApp_OnBeforeClose(t *testing.T) {
	t.Run("Should stop all running commands and stop successfully", func(t *testing.T) {
		// Arrange
		mockCommandRunner := new(executiontest.MockRunner)
		mockLogger := new(appmocks.MockLogger)

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
		mockCommandRunner := new(executiontest.MockRunner)
		mockLogger := new(appmocks.MockLogger)

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
