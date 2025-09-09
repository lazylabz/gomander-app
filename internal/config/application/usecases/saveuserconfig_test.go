package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/application/usecases"
	"gomander/internal/config/domain"
	"gomander/internal/config/domain/test"
)

func TestDefaultSaveUserConfig_Execute(t *testing.T) {
	t.Parallel()
	t.Run("Should save user configuration successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(test.MockConfigRepository)

		sut := usecases.NewSaveUserConfig(mockRepository)

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

		// Act
		err := sut.Execute(newUserConfig)

		// Assert
		assert.NoError(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
		mock.AssertExpectationsForObjects(t, mockRepository)
	})

	t.Run("Should fail to save user configuration", func(t *testing.T) {
		// Arrange
		mockRepository := new(test.MockConfigRepository)

		sut := usecases.NewSaveUserConfig(mockRepository)

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

		// Act
		err := sut.Execute(newUserConfig)

		// Assert
		assert.Error(t, err)

		mockRepository.AssertCalled(t, "Update", &newUserConfig)
		mock.AssertExpectationsForObjects(t, mockRepository)
	})
}
