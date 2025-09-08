package usecases_test

import (
	"context"
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
	"gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/testutils"
	"gomander/internal/testutils/mocks"
)

func TestDefaultExportProject_Execute(t *testing.T) {
	t.Run("Should export the project to the selected file", func(t *testing.T) {
		// Arrange
		projectId := "test-project-id"

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

		sut := usecases.NewExportProject(
			context.Background(),
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockRuntimeFacade,
			mockFsFacade,
		)

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

		// Act
		err = sut.Execute(projectId)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
		)
	})
	t.Run("Should return error if there is a problem opening the destination file", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		sut := usecases.NewExportProject(
			context.Background(),
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockRuntimeFacade,
			mockFsFacade,
		)

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)
		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("", assert.AnError)

		// Act
		err := sut.Execute("test-project-id")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return nil if the user cancels the save dialog", func(t *testing.T) {
		// Arrange

		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		sut := usecases.NewExportProject(
			context.Background(),
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockRuntimeFacade,
			mockFsFacade,
		)

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)
		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("", nil)

		// Act
		err := sut.Execute("test-project-id")

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade)
	})
	t.Run("Should return error if there is a problem reading the project data", func(t *testing.T) {
		// Arrange
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		sut := usecases.NewExportProject(
			context.Background(),
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockRuntimeFacade,
			mockFsFacade,
		)

		mockProjectRepository.On("Get", "test-project-id").Return(nil, assert.AnError)

		// Act
		err := sut.Execute("test-project-id")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade, mockProjectRepository)
	})
	t.Run("Should return error if there is a problem reading the commands ", func(t *testing.T) {
		// Arrange
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		sut := usecases.NewExportProject(
			context.Background(),
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockRuntimeFacade,
			mockFsFacade,
		)

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)
		mockCommandRepository.On("GetAll", "test-project-id").Return(make([]domain.Command, 0), assert.AnError)
		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("/somedir/file.json", nil)

		// Act
		err := sut.Execute("test-project-id")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade, mockProjectRepository, mockCommandRepository)
	})
	t.Run("Should return error if there is a problem reading the command groups", func(t *testing.T) {
		// Arrange
		mockRuntimeFacade := new(mocks.MockRuntimeFacade)
		mockFsFacade := new(MockFsFacade)
		mockProjectRepository := new(MockProjectRepository)
		mockCommandRepository := new(MockCommandRepository)
		mockCommandGroupRepository := new(MockCommandGroupRepository)

		sut := usecases.NewExportProject(
			context.Background(),
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockRuntimeFacade,
			mockFsFacade,
		)

		mockProjectRepository.On("Get", "test-project-id").Return(&projectdomain.Project{Id: "test-project-id", Name: "Test Project"}, nil)
		mockCommandRepository.On("GetAll", "test-project-id").Return([]domain.Command{}, nil)
		mockCommandGroupRepository.On("GetAll", "test-project-id").Return([]commandgroupdomain.CommandGroup{}, assert.AnError)
		mockRuntimeFacade.On("SaveFileDialog", mock.Anything, mock.Anything).Return("/somedir/file.json", nil)

		// Act
		err := sut.Execute("test-project-id")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockFsFacade, mockProjectRepository, mockCommandRepository, mockCommandGroupRepository)
	})
	t.Run("Should return error if there is a problem writing the file", func(t *testing.T) {
		projectId := "test-project-id"

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

		sut := usecases.NewExportProject(
			context.Background(),
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
			mockRuntimeFacade,
			mockFsFacade,
		)

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

		err = sut.Execute(projectId)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t,
			mockProjectRepository,
			mockCommandRepository,
			mockCommandGroupRepository,
		)
	})
}
