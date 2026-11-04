package apptest

import configdomain "gomander/internal/config/domain"

// configRepositoryThatCanFail refuses to store the Zoom Level when a test has
// arranged it to, and answers like the repository it wraps otherwise.
type configRepositoryThatCanFail struct {
	configdomain.Repository
	failures *storageFailures
}

func (r configRepositoryThatCanFail) UpdateZoomLevel(level configdomain.ZoomLevel) error {
	if r.failures.zoomLevelWrite != nil {
		return r.failures.zoomLevelWrite
	}

	return r.Repository.UpdateZoomLevel(level)
}
