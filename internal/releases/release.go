package releases

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Masterminds/semver"

	"gomander/internal/facade"
)

const (
	latestReleaseRequestTimeout = 30 * time.Second
	downloadRequestTimeout      = 10 * time.Minute
)

const CurrentRelease = "v1.4.0"

const RepoOwnerAndName = "Lazylabz/gomander-app"

var DefaultLatestReleaseUrl = fmt.Sprintf("https://github.com/%s/releases.atom", RepoOwnerAndName)
var DefaultBinaryDownloadBaseUrl = fmt.Sprintf("https://github.com/%s/releases/download", RepoOwnerAndName)

type XMLFeed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []struct {
		Title string `xml:"title"`
	} `xml:"entry"`
}

type ReleaseHelper struct {
	ctx                   context.Context
	runtimeFacade         facade.RuntimeFacade
	openFacade            facade.OpenFacade
	osFacade              facade.OSFacade
	ioFacade              facade.IOFacade
	latestReleaseUrl      string
	binaryDownloadBaseUrl string
}

func NewReleaseHelper(
	runtimeFacade facade.RuntimeFacade,
	openFacade facade.OpenFacade,
	osFacade facade.OSFacade,
	ioFacade facade.IOFacade,
	latestReleaseUrl string,
	binaryDownloadBaseUrl string,
) *ReleaseHelper {
	return &ReleaseHelper{
		runtimeFacade:         runtimeFacade,
		openFacade:            openFacade,
		osFacade:              osFacade,
		ioFacade:              ioFacade,
		latestReleaseUrl:      latestReleaseUrl,
		binaryDownloadBaseUrl: binaryDownloadBaseUrl,
	}
}

func SetReleaseHelperContext(rh *ReleaseHelper, ctx context.Context) {
	rh.ctx = ctx
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
	latestRelease, err := rh.getLatestRelease()
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

// httpGet issues a GET bound to rh.ctx (so app shutdown cancels in-flight
// requests) with a hard timeout. It falls back to context.Background when the
// context has not been set yet, since http.NewRequestWithContext panics on nil.
func (rh *ReleaseHelper) httpGet(url string, timeout time.Duration) (*http.Response, error) {
	ctx := rh.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: timeout}
	return client.Do(req)
}

func (rh *ReleaseHelper) getLatestRelease() (release *semver.Version, err error) {

	res, err := rh.httpGet(rh.latestReleaseUrl, latestReleaseRequestTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: " + err.Error())
	}

	defer func(Body io.ReadCloser) {
		bodyCloseError := Body.Close()
		if err == nil {
			err = bodyCloseError
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: received status code %d", res.StatusCode)
	}

	bodyBytes, err := rh.ioFacade.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	var releasesXml XMLFeed
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

func (rh *ReleaseHelper) DownloadLatestRelease(version string) (binaryPath string, err error) {
	binaryPath = filepath.Join(rh.osFacade.TempDir(), rh.getBinaryName(version))

	resp, err := rh.httpGet(rh.getBinaryUrl(version), downloadRequestTimeout)
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

	out, err := rh.osFacade.Create(binaryPath)
	if err != nil {
		return "", err
	}

	defer func() {
		closeErr := out.Close()

		if err == nil {
			err = closeErr
		}
	}()

	_, err = rh.ioFacade.Copy(out, resp.Body)

	if err != nil {
		return "", err
	}

	return binaryPath, nil
}

func (rh *ReleaseHelper) getBinaryUrl(version string) string {
	return fmt.Sprintf("%s/v%s/%s", rh.binaryDownloadBaseUrl, version, rh.getBinaryName(version))
}

func (rh *ReleaseHelper) InstallLatestReleaseAndQuit(binaryPath string) error {
	// Check binary exists
	if _, err := rh.osFacade.Stat(binaryPath); err != nil {
		return err
	}

	// Run binary path
	err := rh.runBinary(binaryPath)
	if err != nil {
		return err
	}

	// Close app
	rh.runtimeFacade.CloseApp(rh.ctx)
	return nil
}
