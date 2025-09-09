package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/project/domain/test"
)

func TestDefaultGetCurrentProject_Execute(t *testing.T) {
	t.Run("Should return the current project", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockConfigRepository := new(MockConfigRepository)

		projectId := "project1"
		project := &projectdomain.Project{Id: projectId, Name: "Test", WorkingDirectory: "/tmp"}

		sut := usecases.NewGetCurrentProject(mockConfigRepository, mockProjectRepository)

		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: projectId}, nil)
		mockProjectRepository.On("Get", projectId).Return(project, nil)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, project, got)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockConfigRepository)
	})

	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockConfigRepository)

		sut := usecases.NewGetCurrentProject(mockConfigRepository, nil)
		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))

		// Act
		_, err := sut.Execute()

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})

	t.Run("Should return an error if project does not exist", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockConfigRepository := new(MockConfigRepository)

		sut := usecases.NewGetCurrentProject(mockConfigRepository, mockProjectRepository)

		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: "nonexistent"}, nil)
		mockProjectRepository.On("Get", "nonexistent").Return(nil, errors.New("project not found"))

		// Act
		_, err := sut.Execute()

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockConfigRepository)
	})
}
