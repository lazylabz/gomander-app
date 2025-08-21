package app_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/app"
	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/testutils"
	"gomander/internal/testutils/mocks"
)

type MockFsFacade struct {
	mock.Mock
}

func (m *MockFsFacade) WriteFile(path string, data []byte, perm os.FileMode) error {
	args := m.Called(path, data, perm)
	return args.Error(0)
}

func (m *MockFsFacade) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

func TestApp_ExportProject(t *testing.T) {
	t.Run("Should export the project to the selected file", func(t *testing.T) {
		projectId := "test-project-id"

		a := app.NewApp()

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		project := projectdomain.Project{
			Id:   projectId,
			Name: "test",
		}

		cmd1Data := testutils.NewCommand().WithProjectId(projectId).Data()
		cmd2Data := testutils.NewCommand().WithProjectId(projectId).Data()
		cmd3Data := testutils.NewCommand().WithProjectId(projectId).Data()

		cmd1 := commandDataToDomain(cmd1Data)
		cmd2 := commandDataToDomain(cmd2Data)
		cmd3 := commandDataToDomain(cmd3Data)

		cmdGroup1Data := testutils.NewCommandGroup().WithProjectId(projectId).WithCommands(cmd1Data, cmd2Data, cmd3Data).Data()
		cmdGroup2Data := testutils.NewCommandGroup().WithProjectId(projectId).WithCommands(cmd3Data, cmd1Data).Data()

		cmdGroup1 := commandGroupDataToDomain(cmdGroup1Data)
		cmdGroup2 := commandGroupDataToDomain(cmdGroup2Data)

		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			FsFacade:               mockFsFacade,
			RuntimeFacade:          mockRuntimeFacade,
		})

		mockProjectRepository.On("Get", projectId).Return(&project, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]domain.Command{cmd1, cmd2, cmd3}, nil)
		mockCommandGroupRepository.On("GetAll", projectId).Return([]commandgroupdomain.CommandGroup{cmdGroup1, cmdGroup2}, nil)

		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("/somedir/file.json", nil)

		expectedExportJSON := app.ProjectExportJSONv1{
			Version: 1,
			Name:    project.Name,
			Commands: array.Map([]domain.Command{cmd1, cmd2, cmd3}, func(cmd domain.Command) app.CommandJSONv1 {
				return app.CommandJSONv1{
					Id:               cmd.Id,
					Name:             cmd.Name,
					Command:          cmd.Command,
					WorkingDirectory: cmd.WorkingDirectory,
				}
			}),
			CommandGroups: array.Map([]commandgroupdomain.CommandGroup{cmdGroup1, cmdGroup2}, func(group commandgroupdomain.CommandGroup) app.CommandGroupJSONv1 {
				return app.CommandGroupJSONv1{
					Name:       group.Name,
					CommandIds: array.Map(group.Commands, func(cmd domain.Command) string { return cmd.Id }),
				}
			}),
		}

		expectedBytes, err := json.MarshalIndent(expectedExportJSON, "", "  ")

		mockFsFacade.On("WriteFile", "/somedir/file.json", expectedBytes, os.FileMode(0644)).Return(nil)

		err = a.ExportProject(projectId)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should return error if there is a problem opening the destination file", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade:     mockRuntimeFacade,
			FsFacade:          mockFsFacade,
			ProjectRepository: mockProjectRepository,
		})

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)

		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("", assert.AnError)

		err := a.ExportProject("test-project-id")
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return nil if the user cancels the save dialog", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade:     mockRuntimeFacade,
			FsFacade:          mockFsFacade,
			ProjectRepository: mockProjectRepository,
		})

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)

		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("", nil)

		err := a.ExportProject("test-project-id")
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem reading the project data", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade:     mockRuntimeFacade,
			FsFacade:          mockFsFacade,
			ProjectRepository: mockProjectRepository,
		})

		mockProjectRepository.On("Get", "test-project-id").Return(nil, assert.AnError)

		err := a.ExportProject("test-project-id")
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade, mockProjectRepository)
	})
	t.Run("Should return error if there is a problem reading the commands ", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade:     mockRuntimeFacade,
			FsFacade:          mockFsFacade,
			ProjectRepository: mockProjectRepository,
			CommandRepository: mockCommandRepository,
		})

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)
		mockCommandRepository.On("GetAll", "test-project-id").Return(make([]domain.Command, 0), assert.AnError)
		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("/somedir/file.json", nil)

		err := a.ExportProject("test-project-id")
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade, mockProjectRepository, mockCommandRepository)
	})
	t.Run("Should return error if there is a problem reading the command groups", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade:          mockRuntimeFacade,
			FsFacade:               mockFsFacade,
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
		})

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)
		mockCommandRepository.On("GetAll", "test-project-id").Return([]domain.Command{}, nil)
		mockCommandGroupRepository.On("GetAll", "test-project-id").Return([]commandgroupdomain.CommandGroup{}, assert.AnError)
		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("/somedir/file.json", nil)

		err := a.ExportProject("test-project-id")
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade, mockProjectRepository, mockCommandRepository, mockCommandGroupRepository)
	})
	t.Run("Should return error if there is a problem writing the file", func(t *testing.T) {
		projectId := "test-project-id"

		a := app.NewApp()

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		project := projectdomain.Project{
			Id:   projectId,
			Name: "test",
		}

		cmd1Data := testutils.NewCommand().WithProjectId(projectId).Data()
		cmd2Data := testutils.NewCommand().WithProjectId(projectId).Data()
		cmd3Data := testutils.NewCommand().WithProjectId(projectId).Data()

		cmd1 := commandDataToDomain(cmd1Data)
		cmd2 := commandDataToDomain(cmd2Data)
		cmd3 := commandDataToDomain(cmd3Data)

		cmdGroup1Data := testutils.NewCommandGroup().WithProjectId(projectId).WithCommands(cmd1Data, cmd2Data, cmd3Data).Data()
		cmdGroup2Data := testutils.NewCommandGroup().WithProjectId(projectId).WithCommands(cmd3Data, cmd1Data).Data()

		cmdGroup1 := commandGroupDataToDomain(cmdGroup1Data)
		cmdGroup2 := commandGroupDataToDomain(cmdGroup2Data)

		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			FsFacade:               mockFsFacade,
			RuntimeFacade:          mockRuntimeFacade,
		})

		mockProjectRepository.On("Get", projectId).Return(&project, nil)
		mockCommandRepository.On("GetAll", projectId).Return([]domain.Command{cmd1, cmd2, cmd3}, nil)
		mockCommandGroupRepository.On("GetAll", projectId).Return([]commandgroupdomain.CommandGroup{cmdGroup1, cmdGroup2}, nil)

		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("/somedir/file.json", nil)

		expectedExportJSON := app.ProjectExportJSONv1{
			Version: 1,
			Name:    project.Name,
			Commands: array.Map([]domain.Command{cmd1, cmd2, cmd3}, func(cmd domain.Command) app.CommandJSONv1 {
				return app.CommandJSONv1{
					Id:               cmd.Id,
					Name:             cmd.Name,
					Command:          cmd.Command,
					WorkingDirectory: cmd.WorkingDirectory,
				}
			}),
			CommandGroups: array.Map([]commandgroupdomain.CommandGroup{cmdGroup1, cmdGroup2}, func(group commandgroupdomain.CommandGroup) app.CommandGroupJSONv1 {
				return app.CommandGroupJSONv1{
					Name:       group.Name,
					CommandIds: array.Map(group.Commands, func(cmd domain.Command) string { return cmd.Id }),
				}
			}),
		}

		expectedBytes, err := json.MarshalIndent(expectedExportJSON, "", "  ")

		mockFsFacade.On("WriteFile", "/somedir/file.json", expectedBytes, os.FileMode(0644)).Return(errors.New("problem writing file"))

		err = a.ExportProject(projectId)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
		)
	})
}

func TestApp_ImportProject(t *testing.T) {
	t.Run("Should import the project", func(t *testing.T) {
		projectId := "test-project-id"

		a := app.NewApp()

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		cmd1Data := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithName("Name 1").
			WithCommand("echo 1").
			WithWorkingDirectory("/1").
			Data()
		cmd2Data := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithName("Name 2").
			WithCommand("echo 2").
			WithWorkingDirectory("/2").
			Data()
		cmd3Data := testutils.
			NewCommand().
			WithProjectId(projectId).
			WithName("Name 3").
			WithCommand("echo 3").
			WithWorkingDirectory("/3").Data()

		cmd1 := commandDataToDomain(cmd1Data)
		cmd2 := commandDataToDomain(cmd2Data)
		cmd3 := commandDataToDomain(cmd3Data)

		cmdGroup1Data := testutils.NewCommandGroup().WithProjectId(projectId).WithCommands(cmd1Data, cmd2Data, cmd3Data).Data()
		cmdGroup2Data := testutils.NewCommandGroup().WithProjectId(projectId).WithCommands(cmd3Data, cmd1Data).Data()

		cmdGroup1 := commandGroupDataToDomain(cmdGroup1Data)
		cmdGroup2 := commandGroupDataToDomain(cmdGroup2Data)

		newName := "Imported Project"
		newWorkingDirectory := "/imported/project/dir"

		commands := []domain.Command{cmd1, cmd2, cmd3}
		commandGroups := []commandgroupdomain.CommandGroup{cmdGroup1, cmdGroup2}

		projectJSON := app.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: array.Map(commands, func(cmd domain.Command) app.CommandJSONv1 {
				return app.CommandJSONv1{
					Id:               cmd.Id,
					Name:             cmd.Name,
					Command:          cmd.Command,
					WorkingDirectory: cmd.WorkingDirectory,
				}
			}),
			CommandGroups: array.Map(commandGroups, func(group commandgroupdomain.CommandGroup) app.CommandGroupJSONv1 {
				return app.CommandGroupJSONv1{
					Name:       group.Name,
					CommandIds: array.Map(group.Commands, func(cmd domain.Command) string { return cmd.Id }),
				}
			}),
		}

		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			FsFacade:               mockFsFacade,
			RuntimeFacade:          mockRuntimeFacade,
		})

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

		err := a.ImportProject(projectJSON, newName, newWorkingDirectory)
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
		a := app.NewApp()

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		projectJSON := app.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: []app.CommandJSONv1{
				{Id: "cmd1", Name: "Command 1", Command: "echo 1", WorkingDirectory: "/1"},
				{Id: "cmd2", Name: "Command 2", Command: "echo 2", WorkingDirectory: "/2"},
			},
			CommandGroups: []app.CommandGroupJSONv1{
				{Name: "Group 1", CommandIds: []string{"cmd1", "cmd2"}},
			},
		}

		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			FsFacade:               mockFsFacade,
			RuntimeFacade:          mockRuntimeFacade,
		})

		mockProjectRepository.On("Create", mock.Anything).Return(assert.AnError).Once()

		err := a.ImportProject(projectJSON, "Imported Project", "/imported/project/dir")
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
		a := app.NewApp()

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		projectJSON := app.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: []app.CommandJSONv1{
				{Id: "cmd1", Name: "Command 1", Command: "echo 1", WorkingDirectory: "/1"},
				{Id: "cmd2", Name: "Command 2", Command: "echo 2", WorkingDirectory: "/2"},
			},
			CommandGroups: []app.CommandGroupJSONv1{
				{Name: "Group 1", CommandIds: []string{"cmd1", "cmd2"}},
			},
		}

		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			FsFacade:               mockFsFacade,
			RuntimeFacade:          mockRuntimeFacade,
		})

		mockProjectRepository.On("Create", mock.Anything).Return(nil)
		mockCommandRepository.On("Create", mock.Anything).Return(assert.AnError)

		err := a.ImportProject(projectJSON, "Imported Project", "/imported/project/dir")
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
		a := app.NewApp()

		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		mockFsFacade := new(MockFsFacade)
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)

		projectJSON := app.ProjectExportJSONv1{
			Version: 1,
			Name:    "test",
			Commands: []app.CommandJSONv1{
				{Id: "cmd1", Name: "Command 1", Command: "echo 1", WorkingDirectory: "/1"},
				{Id: "cmd2", Name: "Command 2", Command: "echo 2", WorkingDirectory: "/2"},
			},
			CommandGroups: []app.CommandGroupJSONv1{
				{Name: "Group 1", CommandIds: []string{"cmd1", "cmd2"}},
			},
		}

		a.LoadDependencies(app.Dependencies{
			ProjectRepository:      mockProjectRepository,
			CommandRepository:      mockCommandRepository,
			CommandGroupRepository: mockCommandGroupRepository,
			FsFacade:               mockFsFacade,
			RuntimeFacade:          mockRuntimeFacade,
		})

		mockProjectRepository.On("Create", mock.Anything).Return(nil)
		mockCommandRepository.On("Create", mock.Anything).Return(nil)
		mockCommandGroupRepository.On("Create", mock.Anything).Return(assert.AnError)

		err := a.ImportProject(projectJSON, "Imported Project", "/imported/project/dir")
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

func TestApp_GetProjectToImport(t *testing.T) {
	t.Run("Should return project import", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		basicProjectJson := app.ProjectExportJSONv1{
			Version:       1,
			Name:          "Name",
			Commands:      make([]app.CommandJSONv1, 0),
			CommandGroups: make([]app.CommandGroupJSONv1, 0),
		}

		basicProjectJsonBytes, err := json.Marshal(basicProjectJson)
		assert.NoError(t, err)

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("/path/to/project.json", nil)
		mockFsFacade.On("ReadFile", "/path/to/project.json").Return(basicProjectJsonBytes, nil)

		toImport, err := a.GetProjectToImport()
		assert.NoError(t, err)
		assert.Equal(t, &basicProjectJson, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem opening the file dialog", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", assert.AnError)

		toImport, err := a.GetProjectToImport()
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return nil if the user cancels the file dialog", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("", nil)

		toImport, err := a.GetProjectToImport()
		assert.NoError(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem reading the file", func(t *testing.T) {
		a := app.NewApp()

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)

		a.LoadDependencies(app.Dependencies{
			RuntimeFacade: mockRuntimeFacade,
			FsFacade:      mockFsFacade,
		})

		mockRuntimeFacade.On("OpenFileDialog", mock.Anything, mock.Anything).Return("/path/to/project.json", nil)
		mockFsFacade.On("ReadFile", "/path/to/project.json").Return([]byte{}, assert.AnError)

		toImport, err := a.GetProjectToImport()
		assert.Error(t, err)
		assert.Nil(t, toImport)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
}
