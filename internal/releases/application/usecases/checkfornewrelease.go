package usecases

import (
	"errors"

	"github.com/Masterminds/semver"

	"gomander/internal/releases"
)

type CheckForNewRelease struct {
	releaseFeed releases.ReleaseFeed
}

func NewCheckForNewRelease(releaseFeed releases.ReleaseFeed) *CheckForNewRelease {
	return &CheckForNewRelease{
		releaseFeed: releaseFeed,
	}
}

// Execute answers with the version of the published Release when it is newer
// than the one running, and with an empty string when it is not.
func (uc *CheckForNewRelease) Execute() (string, error) {
	latest, err := uc.releaseFeed.GetLatestRelease()
	if err != nil {
		return "", errors.New("failed to get latest release: " + err.Error())
	}

	if latest == "" {
		return "", nil
	}

	latestRelease, err := semver.NewVersion(latest)
	if err != nil {
		return "", errors.New("failed to get latest release: " + err.Error())
	}

	if latestRelease.GreaterThan(semver.MustParse(releases.CurrentRelease)) {
		return latestRelease.String(), nil
	}

	return "", nil
}
