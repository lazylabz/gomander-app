package releases

import (
	"encoding/xml"
	"errors"
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

type ReleaseHelper struct{}

func NewReleaseHelper() *ReleaseHelper {
	return &ReleaseHelper{}
}

func (rh *ReleaseHelper) GetCurrentRelease() string {
	return semver.MustParse(CurrentRelease).String()
}

// IsThereANewRelease checks if there is a new release available.
// If a new release is available, it returns the version of the new release.
// If no new release is available, it returns an empty string.
func (rh *ReleaseHelper) IsThereANewRelease() (release string, err error) {
	release = ""

	currentRelease := semver.MustParse(CurrentRelease)
	latestRelease, err := getLatestRelease()
	if err != nil {
		return "", errors.New("failed to get latest release: " + err.Error())
	}
	if latestRelease == nil {
		return "", nil
	}

	if latestRelease.GreaterThan(currentRelease) {
		return latestRelease.String(), nil
	}

	return
}

// LatestReleaseUrl is public so it can be overridden in tests
var LatestReleaseUrl = fmt.Sprintf("https://github.com/%s/releases.atom", RepoOwnerAndName)

func getLatestRelease() (release *semver.Version, err error) {

	res, err := http.Get(LatestReleaseUrl)
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
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	var releasesXml ReleasesFeedXML
	err = xml.Unmarshal(bodyBytes, &releasesXml)
	if err != nil {
		return nil, errors.New("failed to unmarshal response body: " + err.Error())
	}

	if len(releasesXml.Entry) == 0 {
		return nil, nil
	}

	release, err = semver.NewVersion(releasesXml.Entry[0].Title)

	return
}
