package usecases_test

import (
	"embed"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"gomander/internal/localization/application/usecases"
)

//go:embed testdata
var fullTestFs embed.FS

func TestDefaultGetSupportedLanguages_Execute(t *testing.T) {
	t.Run("Should return list of supported languages from embedded filesystem", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFs, "testdata")
		sut := usecases.NewGetSupportedLanguages(testFs)

		// Act
		languages, err := sut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, languages, 3)
		assert.Contains(t, languages, "en-US")
		assert.Contains(t, languages, "fr-FR")
		assert.Contains(t, languages, "invalid")
	})

	t.Run("Should return error when locales directory does not exist", func(t *testing.T) {
		// Arrange
		emptyFs := make(fstest.MapFS)
		sut := usecases.NewGetSupportedLanguages(emptyFs)

		// Act
		languages, err := sut.Execute()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, languages)
		assert.Contains(t, err.Error(), "read locales directory")
	})
}
