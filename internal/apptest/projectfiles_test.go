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
	projectdomain "gomander/internal/project/domain"
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
			WithLink("http://localhost:3000").
			WithErrorPatterns([]string{"ERROR", "panic:"}).
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
			"version": 2,
			"name": "Gomander",
			"workingDirectory": "",
			"commands": [
				{
					"id": "`+command.Id+`",
					"name": "Dev server",
					"command": "pnpm dev",
					"workingDirectory": "`+command.WorkingDirectory+`",
					"link": "http://localhost:3000",
					"errorPatterns": ["ERROR", "panic:"]
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

	t.Run("Should give an exported project back with each command's link and error patterns", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().WithName("Gomander").Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		withLinkAndPatterns := commandtest.NewCommandBuilder().
			WithProjectId(project.Id).
			WithName("Dev server").
			WithLink("http://localhost:3000").
			WithErrorPatterns([]string{"ERROR", "panic:"}).
			WithPosition(0).
			Build()
		withNeither := commandtest.NewCommandBuilder().
			WithProjectId(project.Id).
			WithName("Tests").
			WithPosition(1).
			Build()
		h.GivenCommands(withLinkAndPatterns, withNeither)

		h.GivenExportDestination("/home/user/gomander.json")
		_, err := h.UseCases.ExportProject.Execute(project.Id)
		assert.NoError(t, err)

		exported, _ := h.ExportedFile("/home/user/gomander.json")
		h.GivenFileToImport("/home/user/gomander.json", exported)

		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)
		assert.NoError(t, err)

		// Act
		err = h.UseCases.ImportProject.Execute(*toImport, "Imported", "/work")

		// Assert
		assert.NoError(t, err)

		projects, err := h.UseCases.GetAvailableProjects.Execute()
		assert.NoError(t, err)
		imported := array.Filter(projects, func(p projectdomain.Project) bool { return p.Name == "Imported" })
		assert.Len(t, imported, 1)
		h.GivenOpenedProject(imported[0].Id)

		commands := commandsOf(t, h)
		assert.Equal(t, []string{"Dev server", "Tests"}, array.Map(commands, commandName))
		assert.Equal(t, "http://localhost:3000", commands[0].Link)
		assert.Equal(t, []string{"ERROR", "panic:"}, commands[0].ErrorPatterns)
		assert.Empty(t, commands[1].Link)
		assert.Empty(t, commands[1].ErrorPatterns)
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

	t.Run("Should import a file exported before links and error patterns were carried as commands without them", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/gomander.json", []byte(`{
			"version": 1,
			"name": "Exported",
			"commands": [
				{"id": "cmd-1", "name": "Dev server", "command": "pnpm dev", "workingDirectory": "web"}
			],
			"commandGroups": []
		}`))

		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)
		assert.NoError(t, err)

		// Act
		err = h.UseCases.ImportProject.Execute(*toImport, "Imported", "/work")

		// Assert
		assert.NoError(t, err)

		projects, err := h.UseCases.GetAvailableProjects.Execute()
		assert.NoError(t, err)
		h.GivenOpenedProject(projects[0].Id)

		commands := commandsOf(t, h)
		assert.Len(t, commands, 1)
		assert.Empty(t, commands[0].Link)
		assert.Empty(t, commands[0].ErrorPatterns)
	})

	t.Run("Should refuse a file written by a Gomander newer than this one", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/gomander.json", []byte(`{"version": 99, "name": "From the future"}`))

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)

		// Assert
		assert.ErrorContains(t, err, "version 99")
		assert.Nil(t, toImport)
	})

	t.Run("Should refuse a file that names no version", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/gomander.json", []byte(`{"name": "Unversioned", "commands": []}`))

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)

		// Assert
		assert.ErrorContains(t, err, "version 0")
		assert.Nil(t, toImport)
	})
}
