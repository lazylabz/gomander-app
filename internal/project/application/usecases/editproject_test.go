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

func TestDefaultEditProject_Execute(t *testing.T) {
	t.Run("Should edit a project successfully", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)
		sut := usecases.NewEditProject(mockProjectRepository)

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Update", project).Return(nil)

		// Act
		err := sut.Execute(project)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})

	t.Run("Should return an error if project update fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)

		sut := usecases.NewEditProject(mockProjectRepository)

		project := projectdomain.Project{Id: "1", Name: "A", WorkingDirectory: "/a"}
		mockProjectRepository.On("Update", project).Return(errors.New("fail"))

		// Act
		err := sut.Execute(project)

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}
