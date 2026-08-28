package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	projectusecases "gomander/internal/project/application/usecases"
)

func TestAFailedImport(t *testing.T) {
	t.Run("Should leave no project, commands or groups behind", func(t *testing.T) {
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

		h.GivenStorageThatRefusesToWriteCommandGroups(assert.AnError)

		// Act
		err = h.UseCases.ImportProject.Execute(*toImport, "Imported", "/work")

		// Assert
		assert.ErrorIs(t, err, assert.AnError)

		projects, err := h.UseCases.GetAvailableProjects.Execute()
		assert.NoError(t, err)
		assert.Empty(t, projects)

		assert.Equal(t, apptest.StoredRows{}, h.StoredRows())
	})
}
