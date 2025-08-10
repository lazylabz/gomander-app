package app

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/Masterminds/semver"
)

const CurrentRelease = "v1.0.0"

const RepoOwnerAndName = "Lazylabz/gomander-app"

type ReleasesFeedXML struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []struct {
		Title string `xml:"title"`
	} `xml:"entry"`
}

func (a *App) GetCurrentRelease() string {
	return semver.MustParse(CurrentRelease).String()
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
	latestReleaseUrl := fmt.Sprintf("https://github.com/%s/releases.atom", RepoOwnerAndName)

	res, err := http.Get(latestReleaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: " + err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: received status code %d", res.StatusCode)
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

	var releasesXml ReleasesFeedXML
	err = xml.Unmarshal(bodyBytes, &releasesXml)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if len(releasesXml.Entry) == 0 {
		return nil, nil
	}

	release, err = semver.NewVersion(releasesXml.Entry[0].Title)

	return
}
