package source

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPFetcher retrieves live service configuration from an HTTP endpoint.
type HTTPFetcher struct {
	client *http.Client
}

// NewHTTPFetcher creates an HTTPFetcher with a sensible default timeout.
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &HTTPFetcher{
		client: &http.Client{Timeout: timeout},
	}
}

// Fetch calls the given URL and returns its JSON body as a flat key/value map.
func (f *HTTPFetcher) Fetch(url string) (map[string]string, error) {
	resp, err := f.client.Get(url) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("fetching %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %q", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing response JSON: %w", err)
	}

	return flattenMap("", raw), nil
}
