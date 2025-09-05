package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/config/domain"
	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
)

func TestDefaultOpenProject_Execute(t *testing.T) {
	t.Run("Should open a project successfully", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		sut := usecases.NewOpenProject(mockConfigRepository, mockProjectRepository)

		projectId := "project1"
		projectId2 := "project2"

		mockConfig := domain.Config{
			LastOpenedProjectId: projectId,
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "path1",
					Path: "TestPath",
				},
			},
		}
		mockUpdatedConfig := domain.Config{
			LastOpenedProjectId: projectId2,
			EnvironmentPaths: []domain.EnvironmentPath{
				{
					Id:   "path1",
					Path: "TestPath",
				},
			},
		}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", &mockUpdatedConfig).Return(nil)
		mockProjectRepository.On("Get", projectId2).Return(&projectdomain.Project{Id: projectId2}, nil).Once()

		// Act
		err := sut.Execute(projectId2)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error if project does not exist", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(MockProjectRepository)

		sut := usecases.NewOpenProject(nil, mockProjectRepository)

		projectId := "nonexistent"
		mockProjectRepository.On("Get", projectId).Return(nil, errors.New("project not found"))

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})

	t.Run("Should return an error if getting the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		sut := usecases.NewOpenProject(mockConfigRepository, mockProjectRepository)

		projectId := "project1"
		mockProjectRepository.On("Get", projectId).Return(nil, nil)
		mockConfigRepository.On("GetOrCreate").Return(nil, errors.New("config error"))

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})

	t.Run("Should return an error if updating the config fails", func(t *testing.T) {
		// Arrange
		mockConfigRepository := new(MockConfigRepository)
		mockProjectRepository := new(MockProjectRepository)

		sut := usecases.NewOpenProject(mockConfigRepository, mockProjectRepository)

		projectId := "project1"
		mockConfig := domain.Config{LastOpenedProjectId: projectId}
		mockConfigRepository.On("GetOrCreate").Return(&mockConfig, nil)
		mockConfigRepository.On("Update", mock.Anything).Return(errors.New("update error"))
		mockProjectRepository.On("Get", projectId).Return(&projectdomain.Project{Id: projectId}, nil)

		// Act
		err := sut.Execute(projectId)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockConfigRepository, mockProjectRepository)
	})
}
