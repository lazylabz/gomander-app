package domain_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	commandtest "gomander/internal/command/domain/test"
)

func TestCommand_MatchesErrorPattern(t *testing.T) {
	t.Run("Should not match a line when the command has no error patterns", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithErrorPatterns([]string{}).Build()

		// Act
		matches := command.MatchesErrorPattern("ERROR: something went wrong")

		// Assert
		assert.False(t, matches)
	})

	t.Run("Should match a line that contains an error pattern as a substring", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithErrorPatterns([]string{"ERROR"}).Build()

		// Act
		matches := command.MatchesErrorPattern("2 problems found, ERROR in module")

		// Assert
		assert.True(t, matches)
	})

	t.Run("Should match a line containing any of the error patterns", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithErrorPatterns([]string{"ERROR", "FAILED"}).Build()

		// Act
		matches := command.MatchesErrorPattern("build FAILED")

		// Assert
		assert.True(t, matches)
	})

	t.Run("Should not match a line that contains none of the error patterns", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithErrorPatterns([]string{"ERROR", "FAILED"}).Build()

		// Act
		matches := command.MatchesErrorPattern("build succeeded")

		// Assert
		assert.False(t, matches)
	})

	t.Run("Should treat an error pattern as a plain substring and not as a regex", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithErrorPatterns([]string{"e.ror"}).Build()

		// Act
		matches := command.MatchesErrorPattern("error")

		// Assert
		assert.False(t, matches)
	})

	t.Run("Should match case sensitively", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithErrorPatterns([]string{"ERROR"}).Build()

		// Act
		matches := command.MatchesErrorPattern("error")

		// Assert
		assert.False(t, matches)
	})
}

func TestCommand_ResolveWorkingDirectory(t *testing.T) {
	base := "/home/user"

	t.Run("Should run in the base working directory when the command has none of its own", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithWorkingDirectory("").Build()

		// Act
		workingDirectory := command.ResolveWorkingDirectory(base)

		// Assert
		assert.Equal(t, base, workingDirectory)
	})

	t.Run("Should run in its own working directory when that one is absolute", func(t *testing.T) {
		// Arrange
		absolute := "/etc/config"
		if runtime.GOOS == "windows" {
			absolute = `C:\Program Files`
		}
		command := commandtest.NewCommandBuilder().WithWorkingDirectory(absolute).Build()

		// Act
		workingDirectory := command.ResolveWorkingDirectory(base)

		// Assert
		assert.Equal(t, absolute, workingDirectory)
	})

	t.Run("Should hang a relative working directory off the base one", func(t *testing.T) {
		// Arrange
		command := commandtest.NewCommandBuilder().WithWorkingDirectory("packages/api").Build()

		// Act
		workingDirectory := command.ResolveWorkingDirectory(base)

		// Assert
		assert.Equal(t, filepath.Join(base, "packages/api"), workingDirectory)
	})
}
