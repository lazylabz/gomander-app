package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/usecases"
	"gomander/internal/command/domain/test"
)

func TestApp_EditCommand(t *testing.T) {
	t.Run("Should edit the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)

		projectId := "project1"
		sut := usecases.NewEditCommand(mockCommandRepository)

		commandToEdit := test.NewCommandBuilder().
			WithProjectId(projectId).
			Build()

		mockCommandRepository.On("Update", &commandToEdit).Return(nil)

		// Act
		err := sut.Execute(commandToEdit)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return an error if fails to edit the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)

		projectId := "project1"
		sut := usecases.NewEditCommand(mockCommandRepository)

		commandToEdit := test.NewCommandBuilder().
			WithProjectId(projectId).
			Build()

		mockCommandRepository.On("Update", &commandToEdit).Return(errors.New("failed to update command"))

		// Act
		err := sut.Execute(commandToEdit)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})
}
