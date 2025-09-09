package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/config/domain"
	test2 "gomander/internal/config/domain/test"
	test4 "gomander/internal/logger/test"
	"gomander/internal/project/domain/test"
	test3 "gomander/internal/runner/test"
)

func TestApp_Startup(t *testing.T) {
	t.Run("Should successfully load configuration", func(t *testing.T) {
		// Arrange
		a := app.NewApp()
		ctx := context.Background()

		mockLogger := new(test4.MockLogger)
		mockUserConfigRepository := new(test2.MockConfigRepository)
		mockProjectRepository := new(test.MockProjectRepository)

		a.LoadDependencies(app.Dependencies{
			Logger:            mockLogger,
			ConfigRepository:  mockUserConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: "123"}, nil)

		// Act & Assert
		assert.NotPanics(t, func() {
			a.Startup(ctx)
		})

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})

	t.Run("Should panic if configuration loading fails", func(t *testing.T) {
		// Arrange
		a := app.NewApp()
		ctx := context.Background()

		mockLogger := new(test4.MockLogger)
		mockUserConfigRepository := new(test2.MockConfigRepository)

		a.LoadDependencies(app.Dependencies{
			Logger:           mockLogger,
			ConfigRepository: mockUserConfigRepository,
		})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(nil, assert.AnError)

		// Act & Assert
		assert.Panics(t, func() {
			a.Startup(ctx)
		})

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})
}

func TestApp_OnBeforeClose(t *testing.T) {
	t.Run("Should stop all running commands and stop successfully", func(t *testing.T) {
		// Arrange
		a := app.NewApp()

		mockCommandRunner := new(test3.MockRunner)
		mockLogger := new(test4.MockLogger)

		a.LoadDependencies(app.Dependencies{
			Runner: mockCommandRunner,
			Logger: mockLogger,
		})

		mockCommandRunner.On("StopAllRunningCommands").Return([]error{})

		// Act
		prevent := a.OnBeforeClose(context.Background())

		// Assert
		assert.False(t, prevent)
		mock.AssertExpectationsForObjects(t, mockCommandRunner, mockLogger)
	})

	t.Run("Should prevent closing if there are errors stopping commands", func(t *testing.T) {
		// Arrange
		a := app.NewApp()

		mockCommandRunner := new(test3.MockRunner)
		mockLogger := new(test4.MockLogger)

		a.LoadDependencies(app.Dependencies{
			Runner: mockCommandRunner,
			Logger: mockLogger,
		})

		errs := []error{assert.AnError}
		mockCommandRunner.On("StopAllRunningCommands").Return(errs)

		mockLogger.On("Error", mock.Anything).Return()

		// Act
		prevent := a.OnBeforeClose(context.Background())

		// Assert
		assert.True(t, prevent)

		mock.AssertExpectationsForObjects(t, mockCommandRunner, mockLogger)
	})
}
