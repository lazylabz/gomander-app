package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/testutils"
)

func TestApp_GetCommandGroups(t *testing.T) {
	t.Run("Should return the command groups provided by the command group repository", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"
		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

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
	})
}

func TestApp_CreateCommandGroup(t *testing.T) {
	t.Run("Should create a command group", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()

		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

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
		)
	})
	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

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
		)
	})
	t.Run("Should return an error if failing to save the command group", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

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

		a.SetOpenProjectId(projectId)

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

		a.SetOpenProjectId(projectId)

		commandGroupData := testutils.
			NewCommandGroup().
			WithProjectId(projectId).
			Data()

		paramCommandGroup := commandGroupDataToDomain(commandGroupData)

		mockCommandGroupRepository.On("Update", &paramCommandGroup).Return(errors.New("failed to update command group"))

		err := a.UpdateCommandGroup(&paramCommandGroup)
		assert.Error(t, err)

		mockCommandGroupRepository.AssertExpectations(t)
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

		a.SetOpenProjectId(projectId)

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

		a.SetOpenProjectId(projectId)

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

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

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
		)
	})
	t.Run("Should return an error if failing to retrieve existing command groups", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

		newOrder := []string{"group1", "group2"}

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), errors.New("failed to get command groups"))
		err := a.ReorderCommandGroups(newOrder)
		assert.Error(t, err)
	})
	t.Run("Should return an error if failing to update a command group", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

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

		mockCommandGroupRepository.AssertExpectations(t)
	})
}

func TestApp_RemoveCommandFromCommandGroups(t *testing.T) {
	t.Run("Should remove command from command groups", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

		commandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandData2 := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandGroupData := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData, commandData2)
		commandGroup2Data := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData, commandData2)

		mockCommandGroupRepository.On("GetAll", projectId).Return(
			[]domain.CommandGroup{
				commandGroupDataToDomain(commandGroupData.Data()),
				commandGroupDataToDomain(commandGroup2Data.Data()),
			}, nil)

		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData.WithCommands(commandData2).Data())
		expectedCommandGroup2Call := commandGroupDataToDomain(commandGroup2Data.WithCommands(commandData2).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(nil).Once()
		mockCommandGroupRepository.On("Update", &expectedCommandGroup2Call).Return(nil).Once()

		err := a.RemoveCommandFromCommandGroups(commandData.Id)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should remove remove a command group if it becomes empty", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
		})

		a.SetOpenProjectId(projectId)

		commandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandData2 := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandGroupData := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData, commandData2)
		commandGroup2Data := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData)

		mockCommandGroupRepository.On("GetAll", projectId).Return(
			[]domain.CommandGroup{
				commandGroupDataToDomain(commandGroupData.Data()),
				commandGroupDataToDomain(commandGroup2Data.Data()),
			}, nil)

		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData.WithCommands(commandData2).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(nil).Once()
		mockCommandGroupRepository.On("Delete", commandGroup2Data.Data().Id).Return(nil).Once()

		err := a.RemoveCommandFromCommandGroups(commandData.Id)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should return an error if failing to retrieve command groups", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockLogger := new(MockLogger)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			Logger:                 mockLogger,
		})

		a.SetOpenProjectId(projectId)

		commandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		mockCommandGroupRepository.On("GetAll", projectId).Return(make([]domain.CommandGroup, 0), errors.New("failed to get command groups"))
		mockLogger.On("Error", mock.Anything)

		err := a.RemoveCommandFromCommandGroups(commandData.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if failing to update command groups", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockLogger := new(MockLogger)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			Logger:                 mockLogger,
		})

		a.SetOpenProjectId(projectId)

		commandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandData2 := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandGroupData := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData, commandData2)
		commandGroup2Data := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData, commandData2)

		mockCommandGroupRepository.On("GetAll", projectId).Return(
			[]domain.CommandGroup{
				commandGroupDataToDomain(commandGroupData.Data()),
				commandGroupDataToDomain(commandGroup2Data.Data()),
			}, nil)

		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData.WithCommands(commandData2).Data())
		expectedCommandGroup2Call := commandGroupDataToDomain(commandGroup2Data.WithCommands(commandData2).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(nil).Once()
		mockCommandGroupRepository.On("Update", &expectedCommandGroup2Call).Return(errors.New("failed to update command group")).Once()
		mockLogger.On("Error", mock.Anything)

		err := a.RemoveCommandFromCommandGroups(commandData.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockLogger,
		)
	})
	t.Run("Should return an error if failing to delete a command group", func(t *testing.T) {
		mockCommandGroupRepository := new(MockCommandGroupRepository)
		mockLogger := new(MockLogger)

		projectId := "project1"

		a := app.NewApp()
		a.LoadDependencies(app.Dependencies{
			CommandGroupRepository: mockCommandGroupRepository,
			Logger:                 mockLogger,
		})

		a.SetOpenProjectId(projectId)

		commandData := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandData2 := testutils.
			NewCommand().
			WithProjectId(projectId).
			Data()

		commandGroupData := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData, commandData2)
		commandGroup2Data := testutils.NewCommandGroup().
			WithProjectId(projectId).
			WithCommands(commandData)

		mockCommandGroupRepository.On("GetAll", projectId).Return(
			[]domain.CommandGroup{
				commandGroupDataToDomain(commandGroupData.Data()),
				commandGroupDataToDomain(commandGroup2Data.Data()),
			}, nil)

		expectedCommandGroup1Call := commandGroupDataToDomain(commandGroupData.WithCommands(commandData2).Data())

		mockCommandGroupRepository.On("Update", &expectedCommandGroup1Call).Return(nil).Once()
		mockCommandGroupRepository.On("Delete", commandGroup2Data.Data().Id).Return(errors.New("failed to delete command group")).Once()
		mockLogger.On("Error", mock.Anything)

		err := a.RemoveCommandFromCommandGroups(commandData.Id)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockCommandGroupRepository,
			mockLogger,
		)
	})
}
