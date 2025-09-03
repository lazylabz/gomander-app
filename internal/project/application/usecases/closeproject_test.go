package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
	"gomander/internal/project/application/usecases"
)

func TestDefaultCloseProject_Execute(t *testing.T) {
	t.Run("Should close the current project", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockConfigRepository)

		sut := usecases.NewCloseProject(mockConfigRepository)

		mockConfig := domain.Config{
			LastOpenedProjectId: "project1",
			EnvironmentPaths:    []domain.EnvironmentPath{{Id: "path1", Path: "TestPath"}},
		}
		mockUpdatedConfig := domain.Config{
			LastOpenedProjectId: "",
			EnvironmentPaths:    []domain.EnvironmentPath{{Id: "path1", Path: "TestPath"}},
		}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", &mockUpdatedConfig).Return(nil)

		// Act
		err := sut.Execute()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockConfigRepository)

		sut := usecases.NewCloseProject(mockConfigRepository)

		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))

		// Act
		err := sut.Execute()

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})

	t.Run("Should return an error if updating the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockConfigRepository)

		sut := usecases.NewCloseProject(mockConfigRepository)

		mockConfig := domain.Config{LastOpenedProjectId: "project1"}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", mock.Anything).Return(errors.New("update error"))

		// Act
		err := sut.Execute()

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})
}
