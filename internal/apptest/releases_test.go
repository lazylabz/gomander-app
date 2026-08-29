package apptest_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gomander/internal/apptest"
	"gomander/internal/releases"
)

func TestCheckingForANewRelease(t *testing.T) {
	t.Run("Should report the release it is running", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		// Act
		currentRelease := h.UseCases.GetCurrentRelease.Execute()

		// Assert
		assert.Equal(t, strings.TrimPrefix(releases.CurrentRelease, "v"), currentRelease)
	})

	t.Run("Should report the published release when it is newer than the one running", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenPublishedRelease("v9999.9.9")

		// Act
		newRelease, err := h.UseCases.CheckForNewRelease.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "9999.9.9", newRelease)
	})

	t.Run("Should report no release when the published one is not newer", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenPublishedRelease("v0.0.1")

		// Act
		newRelease, err := h.UseCases.CheckForNewRelease.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, newRelease)
	})

	t.Run("Should report no release when none has been published", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		// Act
		newRelease, err := h.UseCases.CheckForNewRelease.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, newRelease)
	})

	t.Run("Should report the failure when the feed cannot be read", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenAReleaseFeedThatCannotBeRead(assert.AnError)

		// Act
		newRelease, err := h.UseCases.CheckForNewRelease.Execute()

		// Assert
		assert.ErrorContains(t, err, assert.AnError.Error())
		assert.Empty(t, newRelease)
	})
}

func TestInstallingANewRelease(t *testing.T) {
	t.Run("Should download the release, hand its binary to the system and quit", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenPublishedRelease("v9999.9.9")

		newRelease, err := h.UseCases.CheckForNewRelease.Execute()
		require.NoError(t, err)

		// Act
		binaryPath, err := h.UseCases.DownloadRelease.Execute(newRelease)
		require.NoError(t, err)

		err = h.UseCases.InstallReleaseAndQuit.Execute(binaryPath)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{"9999.9.9"}, h.DownloadedReleases())
		assert.Equal(t, []string{binaryPath}, h.InstalledBinaries())
		assert.True(t, h.AppQuit())
	})

	t.Run("Should keep the app open when the install fails", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenAnInstallThatFails(assert.AnError)

		binaryPath, err := h.UseCases.DownloadRelease.Execute("9999.9.9")
		require.NoError(t, err)

		// Act
		err = h.UseCases.InstallReleaseAndQuit.Execute(binaryPath)

		// Assert
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, h.InstalledBinaries())
		assert.False(t, h.AppQuit())
	})
}
