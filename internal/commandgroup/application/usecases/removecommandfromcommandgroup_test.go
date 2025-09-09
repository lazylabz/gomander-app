package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/commandgroup/application/usecases"
	"gomander/internal/testutils"
)

func TestDefaultRemoveCommandFromCommandGroup_Execute(t *testing.T) {
	t.Run("Should remove command from group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		cmdId := "cmd-1"

		sut := usecases.NewRemoveCommandFromCommandGroup(mockCommandGroupRepository)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandToBeDeletedData := testutils.NewCommand().WithId(cmdId).Data()
		anotherCommandData := testutils.NewCommand().WithId("cmd-2").Data()
		existingCommandGroupData := commandGroupData.WithCommands(commandToBeDeletedData, anotherCommandData).Data()
		existingCommandGroup := commandGroupDataToDomain(existingCommandGroupData)

		expectedUpdatedGroup := existingCommandGroup
		expectedUpdatedGroup.Commands = []commanddomain.Command{
			commandDataToDomain(anotherCommandData),
		}

		mockCommandGroupRepository.On("Get", existingCommandGroup.Id).Return(&existingCommandGroup, nil)
		mockCommandGroupRepository.On("Update", &expectedUpdatedGroup).Return(nil)

		// Act
		err := sut.Execute(cmdId, existingCommandGroup.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})

	t.Run("Should return an error if failing to get command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		sut := usecases.NewRemoveCommandFromCommandGroup(mockCommandGroupRepository)

		expectedError := errors.New("failed to get command group")
		mockCommandGroupRepository.On("Get", "group1").Return(nil, expectedError)

		// Act
		err := sut.Execute("cmd-1", "group1")

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})

	t.Run("Should return an error if failing to update command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		cmdId := "cmd-1"

		sut := usecases.NewRemoveCommandFromCommandGroup(mockCommandGroupRepository)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId)
		commandToBeDeletedData := testutils.NewCommand().WithId(cmdId).Data()
		anotherCommandData := testutils.NewCommand().WithId("cmd-2").Data()
		existingCommandGroupData := commandGroupData.WithCommands(commandToBeDeletedData, anotherCommandData).Data()
		existingCommandGroup := commandGroupDataToDomain(existingCommandGroupData)

		expectedError := errors.New("failed to update command group")
		mockCommandGroupRepository.On("Get", existingCommandGroup.Id).Return(&existingCommandGroup, nil)
		mockCommandGroupRepository.On("Update", &existingCommandGroup).Return(expectedError)

		// Act
		err := sut.Execute("cmd-1", existingCommandGroup.Id)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})

	t.Run("Should return an error when trying to remove the last command from the group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		cmdId := "cmd-1"

		sut := usecases.NewRemoveCommandFromCommandGroup(mockCommandGroupRepository)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandToBeDeletedData := testutils.NewCommand().WithId(cmdId).Data()
		existingCommandGroupData := commandGroupData.WithCommands(commandToBeDeletedData).Data()
		existingCommandGroup := commandGroupDataToDomain(existingCommandGroupData)

		mockCommandGroupRepository.On("Get", existingCommandGroup.Id).Return(&existingCommandGroup, nil)

		// Act
		err := sut.Execute(cmdId, existingCommandGroup.Id)

		// Assert
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot remove the last command from the group")

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
}
