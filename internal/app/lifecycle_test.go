package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/config/domain"
	domain2 "gomander/internal/project/domain"
)

func TestApp_Startup(t *testing.T) {
	t.Run("Should successfully load configuration", func(t *testing.T) {
		a := app.NewApp()
		ctx := context.Background()

		mockLogger := new(MockLogger)
		mockUserConfigRepository := new(MockUserConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		a.LoadDependencies(app.Dependencies{
			Logger:            mockLogger,
			ConfigRepository:  mockUserConfigRepository,
			ProjectRepository: mockProjectRepository,
		})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: "123"}, nil)

		assert.NotPanics(t, func() {
			a.Startup(ctx)
		})

		// Verify that the openedProjectId is set correctly
		mockProjectRepository.On("Get", "123").Return(&domain2.Project{}, nil)

		_, err := a.GetCurrentProject()
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})
	t.Run("Should panic if configuration loading fails", func(t *testing.T) {
		a := app.NewApp()
		ctx := context.Background()

		mockLogger := new(MockLogger)
		mockUserConfigRepository := new(MockUserConfigRepository)

		a.LoadDependencies(app.Dependencies{
			Logger:           mockLogger,
			ConfigRepository: mockUserConfigRepository,
		})

		mockLogger.On("Info", mock.Anything).Return()
		mockUserConfigRepository.On("GetOrCreate").Return(nil, assert.AnError)

		assert.Panics(t, func() {
			a.Startup(ctx)
		})

		mock.AssertExpectationsForObjects(t, mockUserConfigRepository, mockLogger)
	})
}

func TestApp_OnBeforeClose(t *testing.T) {
	t.Run("Should stop all running commands and stop successfully", func(t *testing.T) {
		a := app.NewApp()

		mockCommandRunner := new(MockRunner)
		mockLogger := new(MockLogger)
		mockEventEmitter := new(MockEventEmitter)

		a.LoadDependencies(app.Dependencies{
			Runner:       mockCommandRunner,
			Logger:       mockLogger,
			EventEmitter: mockEventEmitter,
		})

		mockCommandRunner.On("StopAllRunningCommands").Return([]error{})

		prevent := a.OnBeforeClose(context.Background())

		assert.False(t, prevent)
		mock.AssertExpectationsForObjects(t, mockCommandRunner, mockLogger, mockEventEmitter)
	})
	t.Run("Should prevent closing if there are errors stopping commands", func(t *testing.T) {
		a := app.NewApp()

		mockCommandRunner := new(MockRunner)
		mockLogger := new(MockLogger)
		mockEventEmitter := new(MockEventEmitter)

		a.LoadDependencies(app.Dependencies{
			Runner:       mockCommandRunner,
			Logger:       mockLogger,
			EventEmitter: mockEventEmitter,
		})

		errs := []error{assert.AnError}
		mockCommandRunner.On("StopAllRunningCommands").Return(errs)

		mockLogger.On("Error", mock.Anything).Return()

		prevent := a.OnBeforeClose(context.Background())

		assert.True(t, prevent)

		mock.AssertExpectationsForObjects(t, mockCommandRunner, mockLogger, mockEventEmitter)
	})
}
