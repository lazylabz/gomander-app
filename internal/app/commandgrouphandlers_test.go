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

func TestApp_GetCommandGroups(t *testing.T) {
	t.Run("Should return the command groups provided by the command group repository", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(testutils.
				NewCommand().
				WithProjectId(projectId).
				Data(),
			).Data()

		expectedCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("GetAll", projectId).Return([]domain.CommandGroup{expectedCommandGroup}, nil)

		got, err := a.GetCommandGroups()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, expectedCommandGroup, got[0])
		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})
}

func TestApp_CreateCommandGroup(t *testing.T) {
	t.Run("Should create a command group", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"

		a := app.NewApp()

		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		otherCommandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			WithPosition(0).
			WithCommands(testutils.
				NewCommand().
				WithProjectId(projectId).
				Data(),
			)

		otherCommandGroup := commandGroupDataToDomain(otherCommandGroupData)
		paramCommandGroup := commandGroupDataToDomain(commandGroupData.Data())
		expectedCommandGroupCall := commandGroupDataToDomain(commandGroupData.WithPosition(1).Data())

		mockCommandGroupRepository.On("GetAll", projectId).Return([]domain.CommandGroup{otherCommandGroup}, nil)
		mockCommandGroupRepository.On("Create", &expectedCommandGroupCall).Return(nil)

		err := a.CreateCommandGroup(&paramCommandGroup)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})
	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), errors.New("failed to get command groups"))

		err := a.CreateCommandGroup(&paramCommandGroup)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})
	t.Run("Should return an error if failing to save the command group", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockUserConfigRepository := new(MockUserConfigRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			ConfigRepository:       mockUserConfigRepository,
		})

		mockUserConfigRepository.On("GetOrCreate").Return(&configdomain.Config{LastOpenedProjectId: projectId}, nil)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), nil)
		mockCommandGroupRepository.On("Create", &paramCommandGroup).Return(errors.New("failed to create command group"))

		err := a.CreateCommandGroup(&paramCommandGroup)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})
}

func TestApp_UpdateCommandGroup(t *testing.T) {
	t.Run("Should update a command group", func(t *testing.T) {
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

		mockCommandGroupRepository.On("Update", &paramCommandGroup).Return(nil)

		err := a.UpdateCommandGroup(&paramCommandGroup)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should return an error if failing to update the command group", func(t *testing.T) {
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

		mockCommandGroupRepository.On("Update", &paramCommandGroup).Return(errors.New("failed to update command group"))

		err := a.UpdateCommandGroup(&paramCommandGroup)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
}

func TestApp_DeleteCommandGroup(t *testing.T) {
	t.Run("Should delete a command group", func(t *testing.T) {
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

		err := a.DeleteCommandGroup(paramCommandGroup.Id)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should return an error if failing to delete the command group", func(t *testing.T) {
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

		err := a.DeleteCommandGroup(paramCommandGroup.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})
}

func TestApp_ReorderCommandGroups(t *testing.T) {
	t.Run("Should reorder command groups based on the provided IDs", func(t *testing.T) {
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

		err := a.ReorderCommandGroups(newOrder)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockUserConfigRepository,
		)
	})
	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
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
		err := a.ReorderCommandGroups(newOrder)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})
	t.Run("Should return an error if failing to update a command group", func(t *testing.T) {
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

		err := a.ReorderCommandGroups(newOrder)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository, mockUserConfigRepository)
	})
}

func TestApp_RemoveCommandFromCommandGroup(t *testing.T) {
	t.Run("Should remove command from group", func(t *testing.T) {
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

		err := a.RemoveCommandFromCommandGroup(cmdId, existingCommandGroup.Id)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
	t.Run("Should return an error if failing to get command group", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		expectedError := errors.New("failed to get command group")
		mockCommandGroupRepository.On("Get", "group1").Return(nil, expectedError)

		err := a.RemoveCommandFromCommandGroup("cmd-1", "group1")
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
	t.Run("Should return an error if failing to update command group", func(t *testing.T) {
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

		err := a.RemoveCommandFromCommandGroup("cmd-1", existingCommandGroup.Id)
		assert.ErrorIs(t, err, expectedError)

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
	t.Run("Should return an error when trying to remove the last command from the group", func(t *testing.T) {
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

		err := a.RemoveCommandFromCommandGroup(cmdId, existingCommandGroup.Id)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot remove the last command from the group")

		mock.AssertExpectationsForObjects(t, mockCommandGroupRepository)
	})
}
