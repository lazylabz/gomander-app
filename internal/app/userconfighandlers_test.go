package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/config/domain"
	"gomander/internal/event"
)

func TestApp_GetUserConfig(t *testing.T) {
	mockRepository := new(MockUserConfigRepository)

	a := app.NewApp()

	a.LoadDependencies(app.Dependencies{
		ConfigRepository: mockRepository,
	})

	expectedResult := domain.Config{
		LastOpenedProjectId: "test-project-id",
		EnvironmentPaths: []domain.EnvironmentPath{
			{
				Id:   "test-env-path-id",
				Path: "test/path",
			},
		},
	}
	mockRepository.On("GetOrCreate").Return(&expectedResult, nil)

	config, err := a.GetUserConfig()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, *config)
}

func TestApp_SaveUserConfig(t *testing.T) {
	t.Parallel()
	t.Run("Should save user configuration successfully", func(t *testing.T) {
		mockRepository := new(MockUserConfigRepository)
		mockEventEmitter := new(MockEventEmitter)
		mockLogger := new(MockLogger)

		a := app.NewApp()

		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockRepository,
			EventEmitter:     mockEventEmitter,
			Logger:           mockLogger,
		})

		newUserConfig := domain.Config{
			LastOpenedProjectId: "new-project-id",
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "new-env-path-id",
					Path: "new/path",
				},
			},
		}

		mockRepository.On("Update", &newUserConfig).Return(nil)

		mockLogger.On("Info", mock.Anything).Return(nil)

		mockEventEmitter.On("EmitEvent", event.SuccessNotification, mock.Anything).Return(nil)
		mockEventEmitter.On("EmitEvent", event.GetUserConfig, nil).Return(nil)

		err := a.SaveUserConfig(newUserConfig)
		assert.NoError(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
	})
	t.Run("Should fail to save user configuration", func(t *testing.T) {
		mockRepository := new(MockUserConfigRepository)
		mockEventEmitter := new(MockEventEmitter)
		mockLogger := new(MockLogger)

		a := app.NewApp()

		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockRepository,
			EventEmitter:     mockEventEmitter,
			Logger:           mockLogger,
		})

		newUserConfig := domain.Config{
			LastOpenedProjectId: "new-project-id",
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "new-env-path-id",
					Path: "new/path",
				},
			},
		}

		mockRepository.On("Update", &newUserConfig).Return(errors.New("failed to save user configuration"))

		mockLogger.On("Error", mock.Anything).Return(nil)

		mockEventEmitter.On("EmitEvent", event.ErrorNotification, mock.Anything).Return(nil)

		err := a.SaveUserConfig(newUserConfig)
		assert.Error(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
	})
}
