package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/application/usecases"
	"gomander/internal/config/domain"
)

func TestDefaultSaveUserConfig_Execute(t *testing.T) {
	t.Parallel()
	t.Run("Should save user configuration successfully", func(t *testing.T) {
		mockRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		sut := usecases.NewDefaultSaveUserConfig(mockRepository, mockLogger)

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

		err := sut.Execute(newUserConfig)
		assert.NoError(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
		mock.AssertExpectationsForObjects(t, mockRepository, mockLogger)
	})
	t.Run("Should fail to save user configuration", func(t *testing.T) {
		mockRepository := new(MockUserConfigRepository)
		mockLogger := new(MockLogger)

		sut := usecases.NewDefaultSaveUserConfig(mockRepository, mockLogger)

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

		err := sut.Execute(newUserConfig)
		assert.Error(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
		mock.AssertExpectationsForObjects(t, mockRepository, mockLogger)
	})
}
