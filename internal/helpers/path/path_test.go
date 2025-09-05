package path_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/helpers/path"
)

func TestGetComputedPath(t *testing.T) {
	base := "/home/user"

	t.Run("returns base if path is empty", func(t *testing.T) {
		// Arrange
		emptyPath := ""

		// Act
		result := path.GetComputedPath(base, emptyPath)

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
		result := path.GetComputedPath(base, abs)

		// Assert
		assert.Equal(t, abs, result)
	})

	t.Run("joins base and relative path", func(t *testing.T) {
		// Arrange
		rel := "docs/readme.md"
		expected := filepath.Join(base, rel)

		// Act
		result := path.GetComputedPath(base, rel)

		// Assert
		assert.Equal(t, expected, result)
	})
}
