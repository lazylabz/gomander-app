package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	"gomander/internal/domainerrors"
	"gomander/internal/helpers/array"
	projectusecases "gomander/internal/project/application/usecases"
	projecttest "gomander/internal/project/domain/test"
)

func TestExportingAProject(t *testing.T) {
	t.Run("Should write the project, its commands and its groups to the file the user picked", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().WithName("Gomander").Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		command := commandtest.NewCommandBuilder().
			WithProjectId(project.Id).
			WithName("Dev server").
			WithCommand("pnpm dev").
			Build()
		h.GivenCommands(command)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithName("Everything").
			WithCommands(command).
			Build()
		h.GivenCommandGroups(group)

		h.GivenExportDestination("/home/user/gomander.json")

		// Act
		path, err := h.UseCases.ExportProject.Execute(project.Id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "/home/user/gomander.json", path)

		exported, written := h.ExportedFile("/home/user/gomander.json")
		assert.True(t, written)
		assert.JSONEq(t, `{
			"version": 1,
			"name": "Gomander",
			"workingDirectory": "",
			"commands": [
				{
					"id": "`+command.Id+`",
					"name": "Dev server",
					"command": "pnpm dev",
					"workingDirectory": "`+command.WorkingDirectory+`"
				}
			],
			"commandGroups": [
				{
					"id": "`+group.Id+`",
					"name": "Everything",
					"commandIds": ["`+command.Id+`"]
				}
			]
		}`, string(exported))
	})

	t.Run("Should write nothing when the user cancels the dialog", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)

		// Act
		path, err := h.UseCases.ExportProject.Execute(project.Id)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, path)
	})

	t.Run("Should report a project that no longer exists as missing", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenExportDestination("/home/user/gomander.json")

		// Act
		_, err := h.UseCases.ExportProject.Execute("deleted-project")

		// Assert
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)
		_, written := h.ExportedFile("/home/user/gomander.json")
		assert.False(t, written)
	})
}

func TestImportingAProject(t *testing.T) {
	t.Run("Should create the project the picked file describes, with its commands", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/gomander.json", []byte(`{
			"version": 1,
			"name": "Exported",
			"commands": [
				{"id": "cmd-1", "name": "Dev server", "command": "pnpm dev", "workingDirectory": "web"},
				{"id": "cmd-2", "name": "Tests", "command": "pnpm test", "workingDirectory": ""}
			],
			"commandGroups": [
				{"id": "group-1", "name": "Everything", "commandIds": ["cmd-1", "cmd-2"]}
			]
		}`))

		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)
		assert.NoError(t, err)
		assert.NotNil(t, toImport)

		// Act
		err = h.UseCases.ImportProject.Execute(*toImport, "Imported", "/work")

		// Assert
		assert.NoError(t, err)

		projects, err := h.UseCases.GetAvailableProjects.Execute()
		assert.NoError(t, err)
		assert.Len(t, projects, 1)
		assert.Equal(t, "Imported", projects[0].Name)
		assert.Equal(t, "/work", projects[0].WorkingDirectory)

		h.GivenOpenedProject(projects[0].Id)

		commands := commandsOf(t, h)
		assert.Equal(t, []string{"Dev server", "Tests"}, array.Map(commands, commandName))
		assert.Equal(t, []int{0, 1}, array.Map(commands, commandPosition))

		groups := commandGroupsOf(t, h)
		assert.Len(t, groups, 1)
		assert.Equal(t, "Everything", groups[0].Name)
		assert.Equal(t, []string{"Dev server", "Tests"}, commandNamesOf(commands, groups[0].CommandIds))
	})
}
