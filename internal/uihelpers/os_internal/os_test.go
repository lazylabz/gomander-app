package os_internal_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/uihelpers/os_internal"
)

func TestGetOs(t *testing.T) {
	helper := os_internal.NewUIOsHelper()

	t.Run("returns the current runtime OS", func(t *testing.T) {
		// Act
		result := helper.GetOs()

		// Assert
		assert.Equal(t, runtime.GOOS, result)
	})
}
