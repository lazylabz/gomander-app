package apptest_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	configdomain "gomander/internal/config/domain"
)

func userConfigOf(t *testing.T, h *apptest.Harness) configdomain.Config {
	t.Helper()

	config, err := h.UseCases.GetUserConfig.Execute()
	assert.NoError(t, err)

	return *config
}

func TestAFreshUserConfig(t *testing.T) {
	t.Run("Should start at 100% zoom with the default logs font size", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		// Act
		config := userConfigOf(t, h)

		// Assert
		assert.Equal(t, configdomain.DefaultZoomLevel, config.ZoomLevel)
		assert.Equal(t, configdomain.DefaultLogsFontSize, config.LogsFontSize)
	})
}

func TestSavingTheUserConfig(t *testing.T) {
	t.Run("Should keep what the user saved", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		saved := configdomain.Config{
			LastOpenedProjectId: "project-1",
			EnvironmentPaths:    []configdomain.EnvironmentPath{{Id: "path-1", Path: "/usr/local/bin"}},
			Locale:              "es",
			LogsFontSize:        18,
			ZoomLevel:           configdomain.DefaultZoomLevel,
		}

		// Act
		err := h.UseCases.SaveUserConfig.Execute(saved)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, saved, userConfigOf(t, h))
	})

	t.Run("Should leave the zoom level where the shortcuts put it", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		assert.NoError(t, h.UseCases.ZoomIn.Execute())
		staleSettings := userConfigOf(t, h)
		staleSettings.ZoomLevel = configdomain.DefaultZoomLevel
		staleSettings.LogsFontSize = 20

		// Act
		err := h.UseCases.SaveUserConfig.Execute(staleSettings)

		// Assert
		assert.NoError(t, err)
		config := userConfigOf(t, h)
		assert.Equal(t, configdomain.ZoomLevel(1.1), config.ZoomLevel)
		assert.Equal(t, 20, config.LogsFontSize)
	})

	for _, fontSize := range []int{configdomain.MinLogsFontSize - 1, configdomain.MaxLogsFontSize + 1} {
		t.Run(fmt.Sprintf("Should refuse a logs font size of %d", fontSize), func(t *testing.T) {
			// Arrange
			h := apptest.New(t)
			config := userConfigOf(t, h)
			config.LogsFontSize = fontSize

			// Act
			err := h.UseCases.SaveUserConfig.Execute(config)

			// Assert
			assert.ErrorIs(t, err, configdomain.ErrLogsFontSizeOutOfRange)
			assert.Equal(t, configdomain.DefaultLogsFontSize, userConfigOf(t, h).LogsFontSize)
		})
	}
}

func TestZooming(t *testing.T) {
	t.Run("Should zoom the window in by one step and remember it", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		// Act
		err := h.UseCases.ZoomIn.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1.1, h.WindowZoomLevel())
		assert.Equal(t, configdomain.ZoomLevel(1.1), userConfigOf(t, h).ZoomLevel)
	})

	t.Run("Should zoom the window out below 100%", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		// Act
		assert.NoError(t, h.UseCases.ZoomOut.Execute())
		err := h.UseCases.ZoomOut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0.8, h.WindowZoomLevel())
		assert.Equal(t, configdomain.ZoomLevel(0.8), userConfigOf(t, h).ZoomLevel)
	})

	t.Run("Should land on whole steps however many it takes", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		// Act
		for range 3 {
			assert.NoError(t, h.UseCases.ZoomIn.Execute())
		}

		// Assert
		assert.Equal(t, 1.3, h.WindowZoomLevel())
	})

	t.Run("Should stop at the largest zoom level", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenZoomLevel(configdomain.MaxZoomLevel)

		// Act
		err := h.UseCases.ZoomIn.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, float64(configdomain.MaxZoomLevel), h.WindowZoomLevel())
		assert.Equal(t, configdomain.MaxZoomLevel, userConfigOf(t, h).ZoomLevel)
	})

	t.Run("Should stop at the smallest zoom level", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenZoomLevel(configdomain.MinZoomLevel)

		// Act
		err := h.UseCases.ZoomOut.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, float64(configdomain.MinZoomLevel), h.WindowZoomLevel())
		assert.Equal(t, configdomain.MinZoomLevel, userConfigOf(t, h).ZoomLevel)
	})

	t.Run("Should reset the window to 100% and remember it", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenZoomLevel(1.5)

		// Act
		err := h.UseCases.ResetZoom.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1.0, h.WindowZoomLevel())
		assert.Equal(t, configdomain.DefaultZoomLevel, userConfigOf(t, h).ZoomLevel)
	})

	t.Run("Should count every press of a held key, however they overlap", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		var presses sync.WaitGroup

		// Act
		for range 5 {
			presses.Go(func() {
				assert.NoError(t, h.UseCases.ZoomIn.Execute())
			})
		}
		presses.Wait()

		// Assert
		assert.Equal(t, 1.5, h.WindowZoomLevel())
		assert.Equal(t, configdomain.ZoomLevel(1.5), userConfigOf(t, h).ZoomLevel)
	})

	t.Run("Should leave the window where it was when the zoom level cannot be stored", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		refused := errors.New("disk full")
		h.GivenStorageThatRefusesToWriteTheZoomLevel(refused)

		// Act
		err := h.UseCases.ZoomIn.Execute()

		// Assert
		assert.ErrorIs(t, err, refused)
		assert.Equal(t, 1.0, h.WindowZoomLevel())
		assert.Equal(t, configdomain.DefaultZoomLevel, userConfigOf(t, h).ZoomLevel)
	})

	t.Run("Should keep the rest of the user config untouched", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenEnvironmentPaths("/usr/local/bin")
		before := userConfigOf(t, h)

		// Act
		err := h.UseCases.ZoomIn.Execute()

		// Assert
		assert.NoError(t, err)
		after := userConfigOf(t, h)
		after.ZoomLevel = before.ZoomLevel
		assert.Equal(t, before, after)
	})
}

func TestRestoringTheZoomLevel(t *testing.T) {
	t.Run("Should put the window back at the zoom level the user left it at", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenZoomLevel(0.7)

		// Act
		err := h.UseCases.RestoreZoom.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0.7, h.WindowZoomLevel())
	})

	t.Run("Should keep the window within the zoom range whatever is stored", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenZoomLevel(5)

		// Act
		err := h.UseCases.RestoreZoom.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, float64(configdomain.MaxZoomLevel), h.WindowZoomLevel())
	})
}
