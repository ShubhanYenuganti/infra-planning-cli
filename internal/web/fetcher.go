package web

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const allowedHost = "docs.aws.amazon.com"

type result struct {
	status    int
	fetchedAt time.Time
}

// Fetcher makes GET requests to docs.aws.amazon.com with a per-run
// in-memory cache, a 5-second timeout, and one retry on 5xx responses.
type Fetcher struct {
	client *http.Client
	cache  map[string]result
}

// New returns a ready-to-use Fetcher.
func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 5 * time.Second},
		cache:  make(map[string]result),
	}
}

// Fetch returns (statusCode, fetchedAt, nil) for allowlisted URLs.
// Non-allowlist URLs return an error without making any HTTP call.
// Results are cached for the lifetime of the Fetcher.
func (f *Fetcher) Fetch(rawURL string) (int, time.Time, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host != allowedHost {
		return 0, time.Time{}, fmt.Errorf("url not in allowlist: %s", rawURL)
	}
	if cached, ok := f.cache[rawURL]; ok {
		return cached.status, cached.fetchedAt, nil
	}
	status, err := f.fetchWithRetry(rawURL)
	if err != nil {
		return 0, time.Time{}, err
	}
	r := result{status: status, fetchedAt: time.Now().UTC()}
	f.cache[rawURL] = r
	return r.status, r.fetchedAt, nil
}

func (f *Fetcher) fetchWithRetry(rawURL string) (int, error) {
	do := func() (int, error) {
		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			return 0, err
		}
		req.Header.Set("User-Agent", "infra-planning-cli/0.1")
		resp, err := f.client.Do(req)
		if err != nil {
			return 0, err
		}
		resp.Body.Close()
		return resp.StatusCode, nil
	}

	status, err := do()
	if err != nil {
		return 0, err
	}
	if status >= 500 {
		status, err = do() // one retry on 5xx
		if err != nil {
			return 0, err
		}
	}
	return status, nil
}
