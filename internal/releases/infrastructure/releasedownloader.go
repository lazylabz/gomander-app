package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"gomander/internal/facade"
)

const downloadRequestTimeout = 10 * time.Minute

// GithubReleaseDownloader writes the binary GitHub publishes for a Release into
// the temporary directory, under the name the platform this app was built for
// expects.
type GithubReleaseDownloader struct {
	ctx      context.Context
	osFacade facade.OSFacade
	ioFacade facade.IOFacade
	baseUrl  string
}

func NewGithubReleaseDownloader(ctx context.Context, osFacade facade.OSFacade, ioFacade facade.IOFacade, baseUrl string) *GithubReleaseDownloader {
	return &GithubReleaseDownloader{
		ctx:      ctx,
		osFacade: osFacade,
		ioFacade: ioFacade,
		baseUrl:  baseUrl,
	}
}

func (d *GithubReleaseDownloader) Download(version string) (binaryPath string, err error) {
	binaryPath = filepath.Join(d.osFacade.TempDir(), binaryName())

	resp, err := httpGet(d.ctx, d.binaryUrl(version), downloadRequestTimeout)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download release: received status code %d", resp.StatusCode)
	}

	out, err := d.osFacade.Create(binaryPath)
	if err != nil {
		return "", err
	}

	defer func() {
		closeErr := out.Close()

		if err == nil {
			err = closeErr
		}
	}()

	_, err = d.ioFacade.Copy(out, resp.Body)

	if err != nil {
		return "", err
	}

	return binaryPath, nil
}

func (d *GithubReleaseDownloader) binaryUrl(version string) string {
	return fmt.Sprintf("%s/v%s/%s", d.baseUrl, version, binaryName())
}
