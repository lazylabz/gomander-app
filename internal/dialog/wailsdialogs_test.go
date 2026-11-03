package dialog_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/pkg/application"

	"gomander/internal/dialog"
)

func TestOpenFileDialogOptions(t *testing.T) {
	t.Run("Should only let the user pick a file, narrowed by the request's filters", func(t *testing.T) {
		// Act
		options := dialog.OpenFileDialogOptions(dialog.OpenFileRequest{
			Title:   "Import project",
			Filters: []dialog.FileFilter{{DisplayName: "Gomander project", Pattern: "*.json"}},
		})

		// Assert
		assert.Equal(t, &application.OpenFileDialogOptions{
			Title:          "Import project",
			CanChooseFiles: true,
			Filters:        []application.FileFilter{{DisplayName: "Gomander project", Pattern: "*.json"}},
		}, options)
	})
}

func TestSaveFileDialogOptions(t *testing.T) {
	t.Run("Should suggest the request's filename and honour whether directories can be created", func(t *testing.T) {
		// Act
		options := dialog.SaveFileDialogOptions(dialog.SaveFileRequest{
			Title:                "Export project",
			DefaultFilename:      "my-project.json",
			CanCreateDirectories: true,
		})

		// Assert
		assert.Equal(t, &application.SaveFileDialogOptions{
			Title:                "Export project",
			Filename:             "my-project.json",
			CanCreateDirectories: true,
		}, options)
	})
}

func TestPickDirectoryDialogOptions(t *testing.T) {
	t.Run("Should only let the user pick a directory", func(t *testing.T) {
		// Act
		options := dialog.PickDirectoryDialogOptions(dialog.PickDirectoryRequest{Title: "Working directory"})

		// Assert
		assert.Equal(t, &application.OpenFileDialogOptions{
			Title:                "Working directory",
			CanChooseDirectories: true,
			CanCreateDirectories: true,
		}, options)
	})
}

func TestAnswer(t *testing.T) {
	t.Run("Should answer with an empty path, not an error, when the user cancels on Windows", func(t *testing.T) {
		// Act
		path, err := dialog.Answer("", errors.New("cancelled by user"))

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, path)
	})

	t.Run("Should report any other failure", func(t *testing.T) {
		// Arrange
		failure := errors.New("no window to attach the dialog to")

		// Act
		_, err := dialog.Answer("", failure)

		// Assert
		assert.Equal(t, failure, err)
	})

	t.Run("Should answer with the path the user picked", func(t *testing.T) {
		// Act
		path, err := dialog.Answer("/home/user/project.json", nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "/home/user/project.json", path)
	})
}
