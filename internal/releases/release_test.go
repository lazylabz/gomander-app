package releases_test

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"

	"gomander/internal/releases"
)

func TestApp_GetCurrentRelease(t *testing.T) {
	// Arrange
	rh := releases.NewReleaseHelper()

	// Act
	currentRelease := rh.GetCurrentRelease()

	// Assert
	assert.Equal(t, "1.1.0", currentRelease)

	// Verify it's a valid semver
	_, err := semver.NewVersion(currentRelease)
	assert.NoError(t, err)
}

func TestApp_IsThereANewRelease(t *testing.T) {
	t.Run("Should return new release when available", func(t *testing.T) {
		// Arrange
		rh := releases.NewReleaseHelper()

		xmlResponse := releases.ReleasesFeedXML{
			XMLName: xml.Name{
				Space: "feed",
				Local: "feed",
			},
			Entry: []struct {
				Title string `xml:"title"`
			}{
				{
					Title: "v9999.9.9",
				},
			},
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/releases.atom", r.URL.Path)

			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)

			encoder := xml.NewEncoder(w)
			err := encoder.Encode(xmlResponse)
			assert.NoError(t, err)
		}))
		defer ts.Close()

		releases.LatestReleaseUrl = ts.URL + "/releases.atom"

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "9999.9.9", release)
	})

	t.Run("Should return empty string when no new release", func(t *testing.T) {
		// Arrange
		rh := releases.NewReleaseHelper()

		xmlResponse := releases.ReleasesFeedXML{
			XMLName: xml.Name{
				Space: "feed",
				Local: "feed",
			},
			Entry: []struct {
				Title string `xml:"title"`
			}{
				{
					Title: "v1.0.0",
				},
			},
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/releases.atom", r.URL.Path)

			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)

			encoder := xml.NewEncoder(w)
			err := encoder.Encode(xmlResponse)
			assert.NoError(t, err)
		}))
		defer ts.Close()

		releases.LatestReleaseUrl = ts.URL + "/releases.atom"

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", release)
	})

	t.Run("Should return nil when no releases found", func(t *testing.T) {
		// Arrange
		rh := releases.NewReleaseHelper()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/releases.atom", r.URL.Path)

			w.Header().Set("Content-Type", "application/atom+xml")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte("<feed></feed>"))
		}))
		defer ts.Close()

		releases.LatestReleaseUrl = ts.URL + "/releases.atom"

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", release)
	})

	t.Run("Should return error when failing to retrieve releases", func(t *testing.T) {
		// Arrange
		rh := releases.NewReleaseHelper()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/releases.atom", r.URL.Path)

			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		releases.LatestReleaseUrl = ts.URL + "/releases.atom"

		// Act
		release, err := rh.IsThereANewRelease()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "", release)
	})
}
