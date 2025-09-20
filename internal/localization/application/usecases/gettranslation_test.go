package usecases_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/localization/application/usecases"
)

//go:embed testdata
var fullTestFsTranslation embed.FS

func TestDefaultGetTranslation_Execute(t *testing.T) {
	t.Run("Should return translation for valid locale", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFsTranslation, "testdata")
		sut := usecases.NewGetTranslation(testFs)

		// Act
		translation, err := sut.Execute("en-US")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, translation)
		assert.Equal(t, "Commands", translation.SidebarTitle)
		assert.Equal(t, "Cancel", translation.ActionsCancel)
	})

	t.Run("Should return error for non-existent locale", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFsTranslation, "testdata")
		sut := usecases.NewGetTranslation(testFs)

		// Act
		translation, err := sut.Execute("non-existent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, translation)
		assert.Contains(t, err.Error(), "read locale json")
	})

	t.Run("Should return error for invalid JSON", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFsTranslation, "testdata")
		sut := usecases.NewGetTranslation(testFs)

		// Act
		translation, err := sut.Execute("invalid")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, translation)
		assert.Contains(t, err.Error(), "unmarshal locale json")
	})
}
