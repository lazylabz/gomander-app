package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/testutils"
)

func TestApp_DeleteCommandGroup(t *testing.T) {
	t.Run("Should delete a command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("Delete", paramCommandGroup.Id).Return(nil)

		// Act
		err := a.DeleteCommandGroup(paramCommandGroup.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})

	t.Run("Should return an error if failing to delete the command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("Delete", paramCommandGroup.Id).Return(errors.New("failed to delete command group"))

		// Act
		err := a.DeleteCommandGroup(paramCommandGroup.Id)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})
}

func TestApp_ReorderCommandGroups(t *testing.T) {
	t.Run("Should reorder command groups based on the provided IDs", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData1 := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandGroupData2 := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandGroups := []domain.CommandGroup{
			commandGroupDataToDomain(commandGroupData1.Data()),
			commandGroupDataToDomain(commandGroupData2.Data()),
		}

		newOrder := []string{commandGroupData2.Data().Id, commandGroupData1.Data().Id}

		mockCommandGroupRepository.On("GetAll", projectId).Return(commandGroups, nil)

		expectedCommandGroup2Call := commandGroupDataToDomain(commandGroupData2.WithPosition(0).Data())
		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData1.WithPosition(1).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup2Call).Return(nil).Once()
		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(nil).Once()

		// Act
		err := a.ReorderCommandGroups(newOrder)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})

	t.Run("Should return an error if failing to retrieve user config", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})
		expectedError := errors.New("failed to get user config")
		mockUserConfigRepository.On("GetOrCreate").Return(nil, expectedError)

		newOrder := []string{"group1", "group2"}

		// Act
		err := a.ReorderCommandGroups(newOrder)

		// Assert
		assert.ErrorIs(t, err, expectedError)
		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})

	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		newOrder := []string{"group1", "group2"}

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), errors.New("failed to get command groups"))

		// Act
		err := a.ReorderCommandGroups(newOrder)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})

	t.Run("Should return an error if failing to update a command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData1 := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandGroups := []domain.CommandGroup{
			commandGroupDataToDomain(commandGroupData1.Data()),
		}

		newOrder := []string{commandGroupData1.Data().Id}

		mockCommandGroupRepository.On("GetAll", projectId).Return(commandGroups, nil)

		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData1.WithPosition(0).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(errors.New("failed to update command group")).Once()

		// Act
		err := a.ReorderCommandGroups(newOrder)

		// Assert
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})
}

func TestApp_RemoveCommandFromCommandGroup(t *testing.T) {
	t.Run("Should remove command from group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		cmdId := "cmd-1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

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
		err := a.RemoveCommandFromCommandGroup(cmdId, existingCommandGroup.Id)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})

	t.Run("Should return an error if failing to get command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		expectedError := errors.New("failed to get command group")
		mockCommandGroupRepository.On("Get", "group1").Return(nil, expectedError)

		// Act
		err := a.RemoveCommandFromCommandGroup("cmd-1", "group1")

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})

	t.Run("Should return an error if failing to update command group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		cmdId := "cmd-1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

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
		err := a.RemoveCommandFromCommandGroup("cmd-1", existingCommandGroup.Id)

		// Assert
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})

	t.Run("Should return an error when trying to remove the last command from the group", func(t *testing.T) {
		// Arrange
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		cmdId := "cmd-1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId)

		commandToBeDeletedData := testutils.NewCommand().WithId(cmdId).Data()
		existingCommandGroupData := commandGroupData.WithCommands(commandToBeDeletedData).Data()
		existingCommandGroup := commandGroupDataToDomain(existingCommandGroupData)

		mockCommandGroupRepository.On("Get", existingCommandGroup.Id).Return(&existingCommandGroup, nil)

		// Act
		err := a.RemoveCommandFromCommandGroup(cmdId, existingCommandGroup.Id)

		// Assert
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot remove the last command from the group")

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
}
