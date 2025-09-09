package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/project/domain/test"
)

func TestDefaultCreateProject_Execute(t *testing.T) {
	t.Run("Should create a project successfully", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)

		sut := usecases.NewCreateProject(mockProjectRepository)

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Create", project).Return(nil)

		// Act
		err := sut.Execute(project)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})

	t.Run("Should return an error if project creation fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)

		sut := usecases.NewCreateProject(mockProjectRepository)

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Create", project).Return(errors.New("fail"))

		// Act
		err := sut.Execute(project)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}
