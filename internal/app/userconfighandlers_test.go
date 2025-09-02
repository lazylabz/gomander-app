package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/config/domain"
)

func TestApp_SaveUserConfig(t *testing.T) {
	t.Parallel()
	t.Run("Should save user configuration successfully", func(t *testing.T) {
		mockRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()

		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockRepository,
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

		err := a.SaveUserConfig(newUserConfig)
		assert.NoError(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
		mock.AssertExpectationsForObjects(t, mockRepository, mockLogger)
	})
	t.Run("Should fail to save user configuration", func(t *testing.T) {
		mockRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		a := app.NewApp()

		a.LoadDependencies(app.Dependencies{
			ConfigRepository: mockRepository,
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

		err := a.SaveUserConfig(newUserConfig)
		assert.Error(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
		mock.AssertExpectationsForObjects(t, mockRepository, mockLogger)
	})
}
