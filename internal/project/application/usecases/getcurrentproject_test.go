package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
	test2 "gomander/internal/config/domain/test"
	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/project/domain/test"
)

func TestDefaultGetCurrentProject_Execute(t *testing.T) {
	t.Run("Should return the current project", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockConfigRepository := new(test2.MockConfigRepository)

		projectId := "project1"
		project := projectdomain.Project{Id: projectId, Name: "Test", WorkingDirectory: "/tmp"}

		sut := usecases.NewGetCurrentProject(mockConfigRepository, mockProjectRepository)

		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: projectId}, nil)
		mockProjectRepository.On("Find", projectId).Return(project, true, nil)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &project, got)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockConfigRepository)
	})

	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(test2.MockConfigRepository)

		sut := usecases.NewGetCurrentProject(mockConfigRepository, nil)
		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))

		// Act
		_, err := sut.Execute()

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockConfigRepository)
	})

	t.Run("Should return nothing if no project is open", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockConfigRepository := new(test2.MockConfigRepository)

		sut := usecases.NewGetCurrentProject(mockConfigRepository, mockProjectRepository)

		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: ""}, nil)
		mockProjectRepository.On("Find", "").Return(nil, false, nil)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, got)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockConfigRepository)
	})

	t.Run("Should return an error if reading the project fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		mockConfigRepository := new(test2.MockConfigRepository)

		sut := usecases.NewGetCurrentProject(mockConfigRepository, mockProjectRepository)

		mockConfigRepository.On("GetOrCreate").Return(&domain.Config{LastOpenedProjectId: "project1"}, nil)
		mockProjectRepository.On("Find", "project1").Return(nil, false, errors.New("storage error"))

		// Act
		_, err := sut.Execute()

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository, mockConfigRepository)
	})
}
