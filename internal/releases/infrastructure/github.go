package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const RepoOwnerAndName = "Lazylabz/gomander-app"

var DefaultLatestReleaseUrl = fmt.Sprintf("https://github.com/%s/releases.atom", RepoOwnerAndName)
var DefaultBinaryDownloadBaseUrl = fmt.Sprintf("https://github.com/%s/releases/download", RepoOwnerAndName)

// httpGet issues a GET bound to ctx, so app shutdown cancels a request still in
// flight, with a hard timeout of its own.
func httpGet(ctx context.Context, url string, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: timeout}

	return client.Do(req)
}
