package apptest_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	projectusecases "gomander/internal/project/application/usecases"
	projecttest "gomander/internal/project/domain/test"
)

func TestPickingAProjectFileToImport(t *testing.T) {
	t.Run("Should read a package.json as a project whose commands are its scripts", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/app/package.json", []byte(`{
			"name": "My NPM Project",
			"scripts": {"start": "node index.js"}
		}`))

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypePackageJSON)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "My NPM Project", toImport.Name)
		assert.Equal(t, "/home/user/app", toImport.WorkingDirectory)
		assert.Len(t, toImport.Commands, 1)
		assert.NotEmpty(t, toImport.Commands[0].Id)
		assert.Equal(t, "start", toImport.Commands[0].Name)
		assert.Equal(t, "node index.js", toImport.Commands[0].Command)
		assert.Empty(t, toImport.Commands[0].WorkingDirectory)
	})

	t.Run("Should read a package.json without scripts as a project without commands", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/app/package.json", []byte(`{"name": "Test Project", "version": "1.0.0"}`))

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypePackageJSON)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "Test Project", toImport.Name)
		assert.Empty(t, toImport.Commands)
	})

	t.Run("Should read a package.json without a name as an unnamed project", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/app/package.json", []byte(`{"scripts": {"test": "jest"}}`))

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypePackageJSON)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, toImport.Name)
		assert.Len(t, toImport.Commands, 1)
		assert.Equal(t, "test", toImport.Commands[0].Name)
	})

	t.Run("Should return nothing when the user cancels the dialog", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, toImport)
	})

	t.Run("Should report a file that is no longer there", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenMissingFileToImport("/home/user/gomander.json")

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)

		// Assert
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, toImport)
	})

	t.Run("Should report a malformed Gomander export", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/gomander.json", []byte(`{"version": 1, "name": "Test", "commands": [`))

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)
	})

	t.Run("Should report a malformed package.json", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenFileToImport("/home/user/app/package.json", []byte(`{"name": "Test", "scripts": {`))

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypePackageJSON)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, toImport)
	})

	t.Run("Should report a dialog the desktop cannot show", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		refused := errors.New("no dialog available")
		h.GivenDialogsThatFail(refused)

		// Act
		toImport, err := h.UseCases.GetProjectToImport.Execute(projectusecases.FileTypeGomander)

		// Assert
		assert.ErrorIs(t, err, refused)
		assert.Nil(t, toImport)
	})
}

func TestPickingAnExportDestination(t *testing.T) {
	t.Run("Should report a dialog the desktop cannot show", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)

		refused := errors.New("no dialog available")
		h.GivenDialogsThatFail(refused)

		// Act
		path, err := h.UseCases.ExportProject.Execute(project.Id)

		// Assert
		assert.ErrorIs(t, err, refused)
		assert.Empty(t, path)
	})
}
