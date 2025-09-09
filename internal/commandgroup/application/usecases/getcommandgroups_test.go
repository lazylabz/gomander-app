package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/domain/test"
	"gomander/internal/commandgroup/application/usecases"
	"gomander/internal/commandgroup/domain"
	test2 "gomander/internal/commandgroup/domain/test"
	configdomain "gomander/internal/config/domain"
	test3 "gomander/internal/config/domain/test"
)

func TestDefaultGetCommandGroups_Execute(t *testing.T) {
	t.Run("Should return the command groups provided by the command group repository", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)

		projectId := "project1"
		sut := usecases.NewGetCommandGroups(mockUserConfigRepository, mockCommandGroupRepository)

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		expectedCommandGroup := test2.NewCommandGroupBuilder().
			WithProjectId(projectId).
			WithCommands(test.NewCommandBuilder().
				WithProjectId(projectId).
				Build(),
			).Build()

		mockCommandGroupRepository.On("GetAll", projectId).Return([]domain.CommandGroup{expectedCommandGroup}, nil)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, expectedCommandGroup, got[0])
		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})

	t.Run("Should return an error if failing to retrieve user config", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockUserConfigRepository := new(test3.MockConfigRepository)

		sut := usecases.NewGetCommandGroups(mockUserConfigRepository, mockCommandGroupRepository)

		expectedError := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedError)

		// Act
		got, err := sut.Execute()

		// Assert
		assert.ErrorIs(t, err, expectedError)
		assert.Len(t, got, 0)
		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})
}
