package path_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/uihelpers/path"
)

func TestGetComputedPath(t *testing.T) {
	base := "/home/user"
	helper := path.NewUiPathHelper()

	t.Run("returns base if path is empty", func(t *testing.T) {
		// Arrange
		emptyPath := ""

		// Act
		result := helper.GetComputedPath(base, emptyPath)

		// Assert
		assert.Equal(t, base, result)
	})

	t.Run("returns path if absolute", func(t *testing.T) {
		// Arrange
		var abs string

		if runtime.GOOS == "windows" {
			abs = `C:\Program Files`
		} else {
			abs = "/etc/config"
		}

		// Act
		result := helper.GetComputedPath(base, abs)

		// Assert
		assert.Equal(t, abs, result)
	})

	t.Run("joins base and relative path", func(t *testing.T) {
		// Arrange
		rel := "docs/readme.md"
		expected := filepath.Join(base, rel)

		// Act
		result := helper.GetComputedPath(base, rel)

		// Assert
		assert.Equal(t, expected, result)
	})
}
