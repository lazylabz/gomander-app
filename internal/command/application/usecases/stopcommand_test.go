package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/application/usecases"
	"gomander/internal/command/domain/test"
)

func TestDefaultStopCommand_Execute(t *testing.T) {
	t.Run("Should stop the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockRunner := new(MockRunner)

		sut := usecases.NewStopCommand(mockCommandRepository, mockRunner)

		cmd := test.NewCommandBuilder().WithProjectId("project1").Build()

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)

		mockRunner.On("StopRunningCommand", cmd.Id).Return(nil)

		// Act
		err := sut.Execute(cmd.Id)
		assert.NoError(t, err)

		// Assert
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockRunner,
		)
	})

	t.Run("Should return error if the command does not exist", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockRunner := new(MockRunner)

		sut := usecases.NewStopCommand(mockCommandRepository, mockRunner)

		commandId := "non-existing-command"

		mockCommandRepository.On("Get", commandId).Return(nil, errors.New("command not found"))

		// Act
		err := sut.Execute(commandId)
		assert.Error(t, err, "command not found")

		// Assert
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
		)
	})

	t.Run("Should return error if fails to stop the command", func(t *testing.T) {
		// Arrange
		mockCommandRepository := new(MockCommandRepository)
		mockRunner := new(MockRunner)

		sut := usecases.NewStopCommand(mockCommandRepository, mockRunner)

		cmd := test.NewCommandBuilder().WithProjectId("project1").Build()

		mockCommandRepository.On("Get", cmd.Id).Return(&cmd, nil)

		mockRunner.On("StopRunningCommand", cmd.Id).Return(errors.New("failed to stop command"))

		// Act
		err := sut.Execute(cmd.Id)
		assert.Error(t, err, "failed to stop command")

		// Assert
		mock.AssertExpectationsForObjects(t,
			mockCommandRepository,
			mockRunner,
		)
	})
}
