package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Masterminds/semver"
)

const CurrentRelease = "v1.0.0"

const RepoURL = "Lazylabz/gomander-app"

type ReleaseAssetJSON struct {
	DownloadURL string `json:"browser_download_url"`
	Name        string `json:"name"`
}

type ReleaseJSON struct {
	TagName string             `json:"tag_name"`
	Assets  []ReleaseAssetJSON `json:"assets"` // Will be used in the future to download the latest release automatically
}

func (a *App) GetCurrentRelease() string {
	return CurrentRelease
}

// IsThereANewRelease checks if there is a new release available.
// If a new release is available, it returns the version of the new release.
// If no new release is available, it returns an empty string.
func (a *App) IsThereANewRelease() (release string, err error) {
	release = ""

	a.logger.Info("Checking for new releases...")
	currentRelease := semver.MustParse(CurrentRelease)
	latestRelease, err := getLatestRelease()
	if err != nil {
		a.logger.Error("Failed to get latest release: " + err.Error())
		return "", fmt.Errorf("failed to get latest release: %w", err)
	}
	if latestRelease == nil {
		a.logger.Info("No new releases found.")
		return "", nil
	}

	if latestRelease.GreaterThan(currentRelease) {
		return latestRelease.String(), nil
	}

	return
}

func getLatestRelease() (release *semver.Version, err error) {
	latestReleaseUrl := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", RepoURL)

	res, err := http.Get(latestReleaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: " + err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: received status code %d", res.StatusCode)
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	defer func(Body io.ReadCloser) {
		bodyCloseError := Body.Close()
		if err == nil {
			err = bodyCloseError
		}
	}(res.Body)

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var releaseJSON ReleaseJSON
	err = json.Unmarshal(bodyBytes, &releaseJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	release, err = semver.NewVersion(releaseJSON.TagName)

	return
}
