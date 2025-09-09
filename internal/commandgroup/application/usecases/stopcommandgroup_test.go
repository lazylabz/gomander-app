// filepath: /Users/moises/Code/personal/gomander/internal/commandgroup/application/usecases/stopcommandgroup_test.go
package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/domain/test"
	"gomander/internal/commandgroup/application/usecases"
	test2 "gomander/internal/commandgroup/domain/test"
	test3 "gomander/internal/runner/test"
)

func TestDefaultStopCommandGroup_Execute(t *testing.T) {
	t.Run("Should stop the command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockRunner := new(test3.MockRunner)

		sut := usecases.NewStopCommandGroup(mockCommandGroupRepository, mockRunner)

		cmd := test.NewCommandBuilder().
			WithId("command1").
			WithProjectId("project1").
			WithPosition(0).
			Build()

		cmdGroup := test2.NewCommandGroupBuilder().
			WithId("group1").
			WithProjectId("project1").
			WithPosition(0).
			WithCommands(cmd).
			Build()

		mockCommandGroupRepository.On("Get", cmdGroup.Id).Return(&cmdGroup, nil)
		mockRunner.On("StopRunningCommands", cmdGroup.Commands).Return(nil)

		// Act
		err := sut.Execute(cmdGroup.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockRunner,
		)
	})

	t.Run("Should return an error if failing to retrieve the command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockRunner := new(test3.MockRunner)

		sut := usecases.NewStopCommandGroup(mockCommandGroupRepository, mockRunner)

		cmdGroupId := "group1"
		expectedError := errors.New("command group not found")

		mockCommandGroupRepository.On("Get", cmdGroupId).Return(nil, expectedError)

		// Act
		err := sut.Execute(cmdGroupId)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})

	t.Run("Should return an error if failing to stop the commands", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(test2.MockCommandGroupRepository)
		mockRunner := new(test3.MockRunner)

		sut := usecases.NewStopCommandGroup(mockCommandGroupRepository, mockRunner)

		cmd := test.NewCommandBuilder().
			WithId("command1").
			WithProjectId("project1").
			WithPosition(0).
			Build()

		cmdGroup := test2.NewCommandGroupBuilder().
			WithId("group1").
			WithProjectId("project1").
			WithPosition(0).
			WithCommands(cmd).
			Build()

		expectedError := errors.New("failed to stop commands")

		mockCommandGroupRepository.On("Get", cmdGroup.Id).Return(&cmdGroup, nil)
		mockRunner.On("StopRunningCommands", cmdGroup.Commands).Return(expectedError)

		// Act
		err := sut.Execute(cmdGroup.Id)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockRunner,
		)
	})
}
