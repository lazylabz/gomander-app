package infrastructure_test

import (
	"embed"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"gomander/internal/localization/infrastructure"
)

//go:embed testdata
var fullTestFs embed.FS

func TestLocaleFiles_Translations(t *testing.T) {
	t.Run("Should return translation for valid locale", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFs, "testdata")
		sut := infrastructure.NewLocaleFiles(testFs)

		// Act
		translation, err := sut.Translations("en")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, translation)
		assert.Equal(t, "Commands", translation.SidebarCommandsTitle)
		assert.Equal(t, "Cancel", translation.CommonCancel)
	})

	t.Run("Should return error for non-existent locale", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFs, "testdata")
		sut := infrastructure.NewLocaleFiles(testFs)

		// Act
		translation, err := sut.Translations("non-existent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, translation)
		assert.Contains(t, err.Error(), "read locale json")
	})

	t.Run("Should return error for invalid JSON", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFs, "testdata")
		sut := infrastructure.NewLocaleFiles(testFs)

		// Act
		translation, err := sut.Translations("invalid")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, translation)
		assert.Contains(t, err.Error(), "unmarshal locale json")
	})
}

func TestLocaleFiles_Locales(t *testing.T) {
	t.Run("Should return list of supported languages from embedded filesystem", func(t *testing.T) {
		// Arrange
		testFs, _ := fs.Sub(fullTestFs, "testdata")
		sut := infrastructure.NewLocaleFiles(testFs)

		// Act
		languages, err := sut.Locales()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, languages, 3)
		assert.Contains(t, languages, "en")
		assert.Contains(t, languages, "es")
		assert.Contains(t, languages, "invalid")
	})

	t.Run("Should return error when locales directory does not exist", func(t *testing.T) {
		// Arrange
		emptyFs := make(fstest.MapFS)
		sut := infrastructure.NewLocaleFiles(emptyFs)

		// Act
		languages, err := sut.Locales()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, languages)
		assert.Contains(t, err.Error(), "read locales directory")
	})
}
