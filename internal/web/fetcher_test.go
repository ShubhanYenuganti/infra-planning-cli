package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetcherAllowlistRejectsNonAWS(t *testing.T) {
	f := New()
	_, _, err := f.Fetch("https://example.com/vpc")
	if err == nil {
		t.Fatalf("expected error for non-allowlist host")
	}
}

func TestFetcherCacheHitSkipsHTTP(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := New()
	const cachedURL = "https://docs.aws.amazon.com/vpc/"
	// Pre-seed the cache to avoid hitting the real network.
	f.cache[cachedURL] = result{status: 200, fetchedAt: time.Now()}

	status, _, err := f.Fetch(cachedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 200 {
		t.Fatalf("expected 200, got %d", status)
	}
	if callCount != 0 {
		t.Fatalf("expected no HTTP calls for cache hit, got %d", callCount)
	}
}

func TestFetcherRetriesOn5xx(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := &Fetcher{client: srv.Client(), cache: make(map[string]result)}
	status, err := f.fetchWithRetry(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected 200 after retry, got %d", status)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 calls (1 fail + 1 retry), got %d", callCount)
	}
}

func TestFetcherSetsUserAgent(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := &Fetcher{client: srv.Client(), cache: make(map[string]result)}
	_, err := f.fetchWithRetry(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotUA != "infra-planning-cli/0.1" {
		t.Fatalf("User-Agent = %q, want infra-planning-cli/0.1", gotUA)
	}
}
