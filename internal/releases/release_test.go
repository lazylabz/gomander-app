package releases_test

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

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gomander/internal/facade"
	"gomander/internal/facade/test"
	"gomander/internal/releases"
)

func TestApp_GetCurrentRelease(t *testing.T) {
	// Arrange
	rh := newDefaultReleaseHelper("", "")

	// Act
	currentRelease := rh.GetCurrentRelease()

	// Assert
	expected, err := semver.NewVersion(releases.CurrentRelease)
	require.NoError(t, err)
	assert.Equal(t, expected.String(), currentRelease)
}

func TestApp_IsThereANewRelease(t *testing.T) {
	t.Run("Should return new release when available", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<feed><entry><title>v9999.9.9</title></entry></feed>`))
		}))
		defer ts.Close()

		rh := newDefaultReleaseHelper(ts.URL, "")

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "9999.9.9", release)
	})

	t.Run("Should return empty string when no new release", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<feed><entry><title>v0.0.1</title></entry></feed>`))
		}))
		defer ts.Close()

		rh := newDefaultReleaseHelper(ts.URL, "")

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", release)
	})

	t.Run("Should return empty when no releases found", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("<feed></feed>"))
		}))
		defer ts.Close()

		rh := newDefaultReleaseHelper(ts.URL, "")

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", release)
	})

	t.Run("Should return error when failing to retrieve releases", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		rh := newDefaultReleaseHelper(ts.URL, "")

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "", release)
	})

	t.Run("Should return error when the HTTP request fails", func(t *testing.T) {
		// Arrange
		rh := newDefaultReleaseHelper("http://127.0.0.1:0", "")

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch latest release")
		assert.Equal(t, "", release)
	})

	t.Run("Should return error when reading the response body fails", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<feed></feed>`))
		}))
		defer ts.Close()

		mockIOFacade := new(test.MockIOFacade)
		mockIOFacade.On("ReadAll", mock.Anything).Return([]byte(nil), errors.New("read failed"))

		rh := releases.NewReleaseHelper(
			&test.MockRuntimeFacade{},
			&test.MockOpenFacade{},
			facade.DefaultOSFacade{},
			mockIOFacade,
			ts.URL,
			"",
		)

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read response body")
		assert.Equal(t, "", release)
		mock.AssertExpectationsForObjects(t, mockIOFacade)
	})

	t.Run("Should return error when the response body is not valid XML", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`not-xml-at-all <<>>`))
		}))
		defer ts.Close()

		rh := newDefaultReleaseHelper(ts.URL, "")

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal response body")
		assert.Equal(t, "", release)
	})
}

func TestReleaseHelper_DownloadLatestRelease(t *testing.T) {
	t.Run("Should download the binary to the temp directory", func(t *testing.T) {
		// Arrange
		payload := []byte("fake-binary-content")
		version := "9.9.9"

		binaryName := getBinaryName()

		var rh *releases.ReleaseHelper
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedPath := "/v" + version + "/" + binaryName
			assert.Equal(t, expectedPath, r.URL.Path)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(payload)
		}))
		defer ts.Close()

		rh = newDefaultReleaseHelper("", ts.URL)

		// Act
		binaryPath, err := rh.DownloadLatestRelease(version)
		t.Cleanup(func() { _ = os.Remove(binaryPath) })

		// Assert
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(os.TempDir(), binaryName), binaryPath)

		written, readErr := os.ReadFile(binaryPath)
		require.NoError(t, readErr)
		assert.Equal(t, payload, written)
	})

	t.Run("Should return an error when the server returns a non-200 status code", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		rh := newDefaultReleaseHelper("", ts.URL)

		// Act
		binaryPath, err := rh.DownloadLatestRelease("9.9.9")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "", binaryPath)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("Should return an error when the HTTP request fails", func(t *testing.T) {
		// Arrange
		rh := newDefaultReleaseHelper("", "http://127.0.0.1:0")

		// Act
		binaryPath, err := rh.DownloadLatestRelease("9.9.9")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "", binaryPath)
	})

	t.Run("Should return an error when creating the destination file fails", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("payload"))
		}))
		defer ts.Close()

		mockOSFacade := new(test.MockOSFacade)
		mockOSFacade.On("TempDir").Return(os.TempDir())
		mockOSFacade.On("Create", mock.Anything).Return(nil, errors.New("create failed"))

		rh := releases.NewReleaseHelper(
			&test.MockRuntimeFacade{},
			&test.MockOpenFacade{},
			mockOSFacade,
			facade.DefaultIOFacade{},
			"",
			ts.URL,
		)

		// Act
		binaryPath, err := rh.DownloadLatestRelease("9.9.9")

		// Assert
		assert.EqualError(t, err, "create failed")
		assert.Equal(t, "", binaryPath)
		mock.AssertExpectationsForObjects(t, mockOSFacade)
	})

	t.Run("Should return an error when copying the response body fails", func(t *testing.T) {
		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("payload"))
		}))
		defer ts.Close()

		mockIOFacade := new(test.MockIOFacade)
		mockIOFacade.On("Copy", mock.Anything, mock.Anything).Return(int64(0), errors.New("copy failed"))

		rh := releases.NewReleaseHelper(
			&test.MockRuntimeFacade{},
			&test.MockOpenFacade{},
			facade.DefaultOSFacade{},
			mockIOFacade,
			"",
			ts.URL,
		)

		// Act
		binaryPath, err := rh.DownloadLatestRelease("9.9.9")
		t.Cleanup(func() {
			_ = os.Remove(filepath.Join(os.TempDir(), getBinaryName()))
		})

		// Assert
		assert.EqualError(t, err, "copy failed")
		assert.Equal(t, "", binaryPath)
		mock.AssertExpectationsForObjects(t, mockIOFacade)
	})
}

func TestReleaseHelper_InstallLatestReleaseAndQuit(t *testing.T) {
	binaryPath := "/some/path/to/binary"

	t.Run("Should return an error when the binary does not exist", func(t *testing.T) {
		// Arrange
		mockRuntimeFacade := new(test.MockRuntimeFacade)
		mockOSFacade := new(test.MockOSFacade)
		mockOpenFacade := new(test.MockOpenFacade)

		mockOSFacade.On("Stat", binaryPath).Return(nil, os.ErrNotExist)

		rh := releases.NewReleaseHelper(
			mockRuntimeFacade,
			mockOpenFacade,
			mockOSFacade,
			facade.DefaultIOFacade{},
			"",
			"",
		)

		// Act
		err := rh.InstallLatestReleaseAndQuit(binaryPath)

		// Assert
		assert.ErrorIs(t, err, os.ErrNotExist)
		mockRuntimeFacade.AssertNotCalled(t, "CloseApp", mock.Anything)
		mockOpenFacade.AssertNotCalled(t, "Run", mock.Anything)
		mock.AssertExpectationsForObjects(t, mockOSFacade)
	})

	t.Run("Should run the binary and close the app when the binary exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRuntimeFacade := new(test.MockRuntimeFacade)
		mockOSFacade := new(test.MockOSFacade)
		mockOpenFacade := new(test.MockOpenFacade)

		mockOSFacade.On("Stat", binaryPath).Return(nil, nil)
		mockOpenFacade.On("Run", expectedOpenArg(binaryPath)).Return(nil)
		mockRuntimeFacade.On("CloseApp", ctx).Return()

		rh := releases.NewReleaseHelper(
			mockRuntimeFacade,
			mockOpenFacade,
			mockOSFacade,
			facade.DefaultIOFacade{},
			"",
			"",
		)
		rh.SetContext(ctx)

		// Act
		err := rh.InstallLatestReleaseAndQuit(binaryPath)

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRuntimeFacade, mockOSFacade, mockOpenFacade)
	})

	t.Run("Should return the error and not close the app when running the binary fails", func(t *testing.T) {
		// Arrange
		mockRuntimeFacade := new(test.MockRuntimeFacade)
		mockOSFacade := new(test.MockOSFacade)
		mockOpenFacade := new(test.MockOpenFacade)

		mockOSFacade.On("Stat", binaryPath).Return(nil, nil)
		mockOpenFacade.On("Run", expectedOpenArg(binaryPath)).Return(assert.AnError)

		rh := releases.NewReleaseHelper(
			mockRuntimeFacade,
			mockOpenFacade,
			mockOSFacade,
			facade.DefaultIOFacade{},
			"",
			"",
		)
		rh.SetContext(context.Background())

		// Act
		err := rh.InstallLatestReleaseAndQuit(binaryPath)

		// Assert
		assert.ErrorIs(t, err, assert.AnError)
		mockRuntimeFacade.AssertNotCalled(t, "CloseApp", mock.Anything)
		mock.AssertExpectationsForObjects(t, mockOSFacade, mockOpenFacade)
	})
}

func newDefaultReleaseHelper(latestReleaseUrl, binaryDownloadBaseUrl string) *releases.ReleaseHelper {
	return releases.NewReleaseHelper(
		&test.MockRuntimeFacade{},
		facade.DefaultOpenFacade{},
		facade.DefaultOSFacade{},
		facade.DefaultIOFacade{},
		latestReleaseUrl,
		binaryDownloadBaseUrl,
	)
}

func expectedOpenArg(binaryPath string) string {
	if runtime.GOOS == "linux" {
		return filepath.Dir(binaryPath)
	}
	return binaryPath
}

func getBinaryName() string {
	binaryName := ""

	if runtime.GOOS == "linux" {
		binaryName = fmt.Sprintf("gomander-linux-%s", runtime.GOARCH)
	}
	if runtime.GOOS == "darwin" {
		binaryName = fmt.Sprintf("gomander-darwin-%s.dmg", runtime.GOARCH)
	}
	if runtime.GOOS == "windows" {
		binaryName = fmt.Sprintf("gomander-windows-%s-installer.exe", runtime.GOARCH)
	}
	return binaryName
}
