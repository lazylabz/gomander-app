package infrastructure_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gomander/internal/facade"
	facadetest "gomander/internal/facade/test"
	"gomander/internal/releases/infrastructure"
)

func TestGithubReleaseDownloader_Download(t *testing.T) {
	t.Run("Should download the binary to the temp directory", func(t *testing.T) {
		// Arrange
		payload := []byte("fake-binary-content")
		version := "9.9.9"

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v"+version+"/"+expectedBinaryName(), r.URL.Path)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(payload)
		}))
		t.Cleanup(ts.Close)

		downloader := infrastructure.NewGithubReleaseDownloader(context.Background(), facade.DefaultOSFacade{}, facade.DefaultIOFacade{}, ts.URL)

		// Act
		binaryPath, err := downloader.Download(version)
		t.Cleanup(func() { _ = os.Remove(binaryPath) })

		// Assert
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(os.TempDir(), expectedBinaryName()), binaryPath)

		written, readErr := os.ReadFile(binaryPath)
		require.NoError(t, readErr)
		assert.Equal(t, payload, written)
	})

	t.Run("Should return an error when the server returns a non-200 status code", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		t.Cleanup(ts.Close)

		downloader := infrastructure.NewGithubReleaseDownloader(context.Background(), facade.DefaultOSFacade{}, facade.DefaultIOFacade{}, ts.URL)

		// Act
		binaryPath, err := downloader.Download("9.9.9")

		// Assert
		assert.ErrorContains(t, err, "404")
		assert.Empty(t, binaryPath)
	})

	t.Run("Should return an error when the request fails", func(t *testing.T) {
		// Arrange
		downloader := infrastructure.NewGithubReleaseDownloader(context.Background(), facade.DefaultOSFacade{}, facade.DefaultIOFacade{}, "http://127.0.0.1:0")

		// Act
		binaryPath, err := downloader.Download("9.9.9")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, binaryPath)
	})

	t.Run("Should return an error when creating the destination file fails", func(t *testing.T) {
		// Arrange
		ts := serverServingAPayload(t)

		mockOSFacade := new(facadetest.MockOSFacade)
		mockOSFacade.On("TempDir").Return(os.TempDir())
		mockOSFacade.On("Create", mock.Anything).Return(nil, errors.New("create failed"))

		downloader := infrastructure.NewGithubReleaseDownloader(context.Background(), mockOSFacade, facade.DefaultIOFacade{}, ts.URL)

		// Act
		binaryPath, err := downloader.Download("9.9.9")

		// Assert
		assert.EqualError(t, err, "create failed")
		assert.Empty(t, binaryPath)
		mock.AssertExpectationsForObjects(t, mockOSFacade)
	})

	t.Run("Should return an error when copying the response body fails", func(t *testing.T) {
		// Arrange
		ts := serverServingAPayload(t)

		mockIOFacade := new(facadetest.MockIOFacade)
		mockIOFacade.On("Copy", mock.Anything, mock.Anything).Return(int64(0), errors.New("copy failed"))

		downloader := infrastructure.NewGithubReleaseDownloader(context.Background(), facade.DefaultOSFacade{}, mockIOFacade, ts.URL)

		// Act
		binaryPath, err := downloader.Download("9.9.9")
		t.Cleanup(func() {
			_ = os.Remove(filepath.Join(os.TempDir(), expectedBinaryName()))
		})

		// Assert
		assert.EqualError(t, err, "copy failed")
		assert.Empty(t, binaryPath)
		mock.AssertExpectationsForObjects(t, mockIOFacade)
	})
}

func serverServingAPayload(t *testing.T) *httptest.Server {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("payload"))
	}))
	t.Cleanup(ts.Close)

	return ts
}

// expectedBinaryName spells out, per platform, the name the release assets are
// published under, so the test fails if the built-in one drifts from it.
func expectedBinaryName() string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("gomander-linux-%s", runtime.GOARCH)
	case "darwin":
		return fmt.Sprintf("gomander-darwin-%s.dmg", runtime.GOARCH)
	case "windows":
		return fmt.Sprintf("gomander-windows-%s-installer.exe", runtime.GOARCH)
	default:
		return ""
	}
}
