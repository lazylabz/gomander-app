package path_test

import (
	"path/filepath"
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
		abs := "/etc/config"
		assert.Equal(t, abs, path.GetComputedPath(base, abs))
	})

	t.Run("joins base and relative path", func(t *testing.T) {
		rel := "docs/readme.md"
		expected := filepath.Join(base, rel)
		assert.Equal(t, expected, path.GetComputedPath(base, rel))
	})
}
