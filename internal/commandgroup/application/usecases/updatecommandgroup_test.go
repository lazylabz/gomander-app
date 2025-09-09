package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/commandgroup/application/usecases"
	"gomander/internal/testutils"
)

func TestDefaultUpdateCommandGroup_Execute(t *testing.T) {
	t.Run("Should update a command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		sut := usecases.NewUpdateCommandGroup(mockCommandGroupRepository)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("Update", &paramCommandGroup).Return(nil)

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})

	t.Run("Should return an error if failing to update the command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		sut := usecases.NewUpdateCommandGroup(mockCommandGroupRepository)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("Update", &paramCommandGroup).Return(errors.New("failed to update command group"))

		// Act
		err := sut.Execute(&paramCommandGroup)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
}
