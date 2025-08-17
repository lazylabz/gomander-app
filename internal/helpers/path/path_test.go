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
		assert.Equal(t, base, path.GetComputedPath(base, ""))
	})

	t.Run("returns path if absolute", func(t *testing.T) {
		var abs string

		if runtime.GOOS == "windows" {
			abs = `C:\Program Files`
		} else {
			abs = "/etc/config"
		}
		assert.Equal(t, abs, path.GetComputedPath(base, abs))
	})

	t.Run("joins base and relative path", func(t *testing.T) {
		rel := "docs/readme.md"
		expected := filepath.Join(base, rel)
		assert.Equal(t, expected, path.GetComputedPath(base, rel))
	})
}
