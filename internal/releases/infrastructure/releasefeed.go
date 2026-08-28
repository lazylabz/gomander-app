package infrastructure

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gomander/internal/facade"
)

const latestReleaseRequestTimeout = 30 * time.Second

type xmlFeed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []struct {
		Title string `xml:"title"`
	} `xml:"entry"`
}

// GithubReleaseFeed reads the Atom feed GitHub publishes for the repository's
// releases, whose first entry is the latest one.
type GithubReleaseFeed struct {
	ctx      context.Context
	ioFacade facade.IOFacade
	url      string
}

func NewGithubReleaseFeed(ctx context.Context, ioFacade facade.IOFacade, url string) *GithubReleaseFeed {
	return &GithubReleaseFeed{
		ctx:      ctx,
		ioFacade: ioFacade,
		url:      url,
	}
}

func (f *GithubReleaseFeed) GetLatestRelease() (version string, err error) {
	res, err := httpGet(f.ctx, f.url, latestReleaseRequestTimeout)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release: %w", err)
	}

	defer func(body io.ReadCloser) {
		bodyCloseError := body.Close()
		if err == nil {
			err = bodyCloseError
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release: received status code %d", res.StatusCode)
	}

	bodyBytes, err := f.ioFacade.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("failed to read response body: " + err.Error())
	}

	var releasesXml xmlFeed
	if err := xml.Unmarshal(bodyBytes, &releasesXml); err != nil {
		return "", errors.New("failed to unmarshal response body: " + err.Error())
	}

	if len(releasesXml.Entry) == 0 {
		return "", nil
	}

	// An entry with no title is a broken feed, not an empty one: reporting it as
	// the empty version would read as "you are up to date".
	if releasesXml.Entry[0].Title == "" {
		return "", errors.New("the latest entry in the feed has no version")
	}

	return releasesXml.Entry[0].Title, nil
}
