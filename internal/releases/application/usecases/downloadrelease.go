package usecases

import (
	"gomander/internal/releases"
)

type DownloadRelease struct {
	releaseDownloader releases.ReleaseDownloader
}

func NewDownloadRelease(releaseDownloader releases.ReleaseDownloader) *DownloadRelease {
	return &DownloadRelease{
		releaseDownloader: releaseDownloader,
	}
}

func (uc *DownloadRelease) Execute(version string) (binaryPath string, err error) {
	return uc.releaseDownloader.Download(version)
}
