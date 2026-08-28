package infrastructure_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gomander/internal/facade"
	facadetest "gomander/internal/facade/test"
	"gomander/internal/releases/infrastructure"
)

func TestGithubReleaseFeed_GetLatestRelease(t *testing.T) {
	t.Run("Should return the version of the latest published release", func(t *testing.T) {
		// Arrange
		ts := feedServing(t, http.StatusOK, `<feed><entry><title>v9999.9.9</title></entry></feed>`)

		feed := infrastructure.NewGithubReleaseFeed(context.Background(), facade.DefaultIOFacade{}, ts.URL)

		// Act
		version, err := feed.GetLatestRelease()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "v9999.9.9", version)
	})

	t.Run("Should return an empty version when the feed lists no release", func(t *testing.T) {
		// Arrange
		ts := feedServing(t, http.StatusOK, "<feed></feed>")

		feed := infrastructure.NewGithubReleaseFeed(context.Background(), facade.DefaultIOFacade{}, ts.URL)

		// Act
		version, err := feed.GetLatestRelease()

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, version)
	})

	t.Run("Should return an error when the latest entry has no version", func(t *testing.T) {
		// Arrange
		ts := feedServing(t, http.StatusOK, "<feed><entry><title></title></entry></feed>")

		feed := infrastructure.NewGithubReleaseFeed(context.Background(), facade.DefaultIOFacade{}, ts.URL)

		// Act
		version, err := feed.GetLatestRelease()

		// Assert
		assert.ErrorContains(t, err, "no version")
		assert.Empty(t, version)
	})

	t.Run("Should return an error when the feed answers with a non-200 status code", func(t *testing.T) {
		// Arrange
		ts := feedServing(t, http.StatusInternalServerError, "")

		feed := infrastructure.NewGithubReleaseFeed(context.Background(), facade.DefaultIOFacade{}, ts.URL)

		// Act
		version, err := feed.GetLatestRelease()

		// Assert
		assert.ErrorContains(t, err, "failed to fetch latest release")
		assert.Empty(t, version)
	})

	t.Run("Should return an error when the request fails", func(t *testing.T) {
		// Arrange
		feed := infrastructure.NewGithubReleaseFeed(context.Background(), facade.DefaultIOFacade{}, "http://127.0.0.1:0")

		// Act
		version, err := feed.GetLatestRelease()

		// Assert
		assert.ErrorContains(t, err, "failed to fetch latest release")
		assert.Empty(t, version)
	})

	t.Run("Should return an error when reading the response body fails", func(t *testing.T) {
		// Arrange
		ts := feedServing(t, http.StatusOK, "<feed></feed>")

		mockIOFacade := new(facadetest.MockIOFacade)
		mockIOFacade.On("ReadAll", mock.Anything).Return([]byte(nil), errors.New("read failed"))

		feed := infrastructure.NewGithubReleaseFeed(context.Background(), mockIOFacade, ts.URL)

		// Act
		version, err := feed.GetLatestRelease()

		// Assert
		assert.ErrorContains(t, err, "failed to read response body")
		assert.Empty(t, version)
		mock.AssertExpectationsForObjects(t, mockIOFacade)
	})

	t.Run("Should return an error when the response body is not valid XML", func(t *testing.T) {
		// Arrange
		ts := feedServing(t, http.StatusOK, "not-xml-at-all <<>>")

		feed := infrastructure.NewGithubReleaseFeed(context.Background(), facade.DefaultIOFacade{}, ts.URL)

		// Act
		version, err := feed.GetLatestRelease()

		// Assert
		assert.ErrorContains(t, err, "failed to unmarshal response body")
		assert.Empty(t, version)
	})
}

func feedServing(t *testing.T, statusCode int, body string) *httptest.Server {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(ts.Close)

	return ts
}
