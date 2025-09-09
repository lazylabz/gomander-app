package usecases_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/command/domain"
	"gomander/internal/command/domain/test"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	test2 "gomander/internal/commandgroup/domain/test"
	"gomander/internal/helpers/array"
	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/testutils/mocks"
)

func TestDefaultImportProject_Execute(t *testing.T) {
	t.Run("Should import the project", func(t *testing.T) {
		// Arrange
		projectId := "test-project-id"

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		cmd1 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithName("Name 1").
			WithCommand("echo 1").
			WithWorkingDirectory("/1").
			Build()
		cmd2 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithName("Name 2").
			WithCommand("echo 2").
			WithWorkingDirectory("/2").
			Build()
		cmd3 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithName("Name 3").
			WithCommand("echo 3").
			WithWorkingDirectory("/3").Build()

		cmdGroup1 := test2.NewCommandGroupBuilder().WithProjectId(projectId).WithCommands(cmd1, cmd2, cmd3).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithProjectId(projectId).WithCommands(cmd3, cmd1).Build()

		newName := "Imported Project"
		newWorkingDirectory := "/imported/project/dir"

		commands := []domain.Command{cmd1, cmd2, cmd3}
		commandGroups := []commandgroupdomain.CommandGroup{cmdGroup1, cmdGroup2}

		projectJSON := projectdomain.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: array.Map(commands, func(cmd domain.Command) projectdomain.CommandJSONv1 {
				return projectdomain.CommandJSONv1{
					Id:               cmd.Id,
					Name:             cmd.Name,
					Command:          cmd.Command,
					WorkingDirectory: cmd.WorkingDirectory,
				}
			}),
			CommandGroups: array.Map(commandGroups, func(group commandgroupdomain.CommandGroup) projectdomain.CommandGroupJSONv1 {
				return projectdomain.CommandGroupJSONv1{
					Name:       group.Name,
					CommandIds: array.Map(group.Commands, func(cmd domain.Command) string { return cmd.Id }),
				}
			}),
		}

		sut := usecases.NewImportProject(mockProjectRepository, mockCommandRepository, mockCommandGroupRepository)

		var newProjectId string

		mockProjectRepository.On("Create", mock.MatchedBy(func(p projectdomain.Project) bool {
			newProjectId = p.Id // Capture the new project ID
			return p.Name == newName && p.WorkingDirectory == newWorkingDirectory
		})).Return(nil).Once()

		// Testify doesn't allow for a custom matcher for pointers, so we capture all the calls and then check the values
		var capturedCommands []*domain.Command
		mockCommandRepository.On("Create", mock.Anything).Run(func(args mock.Arguments) {
			capturedCommands = append(capturedCommands, args.Get(0).(*domain.Command))
		}).Return(nil)

		var capturedCommandGroups []*commandgroupdomain.CommandGroup
		mockCommandGroupRepository.On("Create", mock.Anything).Run(func(args mock.Arguments) {
			capturedCommandGroups = append(capturedCommandGroups, args.Get(0).(*commandgroupdomain.CommandGroup))
		}).Return(nil)

		// Act
		err := sut.Execute(projectJSON, newName, newWorkingDirectory)

		// Assert
		assert.NoError(t, err)

		for i, expectedCmd := range commands {
			assert.Equal(t, expectedCmd.Name, capturedCommands[i].Name)
			assert.Equal(t, expectedCmd.Command, capturedCommands[i].Command)
			assert.Equal(t, expectedCmd.WorkingDirectory, capturedCommands[i].WorkingDirectory)
			assert.Equal(t, i, capturedCommands[i].Position)
			assert.Equal(t, newProjectId, capturedCommands[i].ProjectId) // Ensure the project ID is set correctly
			assert.NotEmpty(t, capturedCommands[i].Id)                   // Random ID exists
		}

		for i, expectedGroup := range commandGroups {
			assert.Equal(t, expectedGroup.Name, capturedCommandGroups[i].Name)
			assert.Equal(t, newProjectId, capturedCommandGroups[i].ProjectId) // Ensure the project ID is set correctly
			assert.Equal(t, i, capturedCommandGroups[i].Position)
			assert.NotEmpty(t, capturedCommandGroups[i].Id) // Random ID exists

			for j, command := range expectedGroup.Commands {
				assert.Equal(t, command.Name, capturedCommandGroups[i].Commands[j].Name)
				assert.Equal(t, command.Command, capturedCommandGroups[i].Commands[j].Command)
				assert.Equal(t, command.WorkingDirectory, capturedCommandGroups[i].Commands[j].WorkingDirectory)
			}
		}

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockFsFacade,
			mockRuntimeFacade,
		)
	})
	t.Run("Should return error if there is a problem saving the project", func(t *testing.T) {
		// Arrange

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		projectJSON := projectdomain.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: []projectdomain.CommandJSONv1{
				{Id: "cmd1", Name: "Command 1", Command: "echo 1", WorkingDirectory: "/1"},
				{Id: "cmd2", Name: "Command 2", Command: "echo 2", WorkingDirectory: "/2"},
			},
			CommandGroups: []projectdomain.CommandGroupJSONv1{
				{Name: "Group 1", CommandIds: []string{"cmd1", "cmd2"}},
			},
		}

		sut := usecases.NewImportProject(mockProjectRepository, mockCommandRepository, mockCommandGroupRepository)

		mockProjectRepository.On("Create", mock.Anything).Return(assert.AnError).Once()

		// Act
		err := sut.Execute(projectJSON, "Imported Project", "/imported/project/dir")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockFsFacade,
			mockRuntimeFacade,
		)
	})
	t.Run("Should return error if there is a problem saving the commands", func(t *testing.T) {
		// Arrange

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		projectJSON := projectdomain.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: []projectdomain.CommandJSONv1{
				{Id: "cmd1", Name: "Command 1", Command: "echo 1", WorkingDirectory: "/1"},
				{Id: "cmd2", Name: "Command 2", Command: "echo 2", WorkingDirectory: "/2"},
			},
			CommandGroups: []projectdomain.CommandGroupJSONv1{
				{Name: "Group 1", CommandIds: []string{"cmd1", "cmd2"}},
			},
		}

		sut := usecases.NewImportProject(mockProjectRepository, mockCommandRepository, mockCommandGroupRepository)

		mockProjectRepository.On("Create", mock.Anything).Return(nil)
		mockCommandRepository.On("Create", mock.Anything).Return(assert.AnError)

		// Act
		err := sut.Execute(projectJSON, "Imported Project", "/imported/project/dir")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockFsFacade,
			mockRuntimeFacade,
		)
	})
	t.Run("Should return error if there is a problem saving the command groups", func(t *testing.T) {
		// Arrange

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		projectJSON := projectdomain.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: []projectdomain.CommandJSONv1{
				{Id: "cmd1", Name: "Command 1", Command: "echo 1", WorkingDirectory: "/1"},
				{Id: "cmd2", Name: "Command 2", Command: "echo 2", WorkingDirectory: "/2"},
			},
			CommandGroups: []projectdomain.CommandGroupJSONv1{
				{Name: "Group 1", CommandIds: []string{"cmd1", "cmd2"}},
			},
		}

		sut := usecases.NewImportProject(mockProjectRepository, mockCommandRepository, mockCommandGroupRepository)

		mockProjectRepository.On("Create", mock.Anything).Return(nil)
		mockCommandRepository.On("Create", mock.Anything).Return(nil)
		mockCommandGroupRepository.On("Create", mock.Anything).Return(assert.AnError)

		// Act
		err := sut.Execute(projectJSON, "Imported Project", "/imported/project/dir")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockFsFacade,
			mockRuntimeFacade,
		)
	})
}
