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

func TestDefaultGetAvailableProjects_Execute(t *testing.T) {
	t.Run("Should return available projects", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)

		sut := usecases.NewGetAvailableProjects(mockProjectRepository)

		projects := []projectdomain.Project{{Id: "1", Name: "A", WorkingDirectory: "/a"}}
		mockProjectRepository.On("GetAll").Return(projects, nil)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, projects, got)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})

	t.Run("Should return an error if fetching projects fails", func(t *testing.T) {
		// Arrange
		mockProjectRepository := new(test.MockProjectRepository)

		sut := usecases.NewGetAvailableProjects(mockProjectRepository)

		mockProjectRepository.On("GetAll").Return(make([]projectdomain.Project, 0), errors.New("fail"))

		// Act
		_, err := sut.Execute()

		// Assert
		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockProjectRepository)
	})
}
