package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.wiechtig.com/shorty/internal/store"
	"go.wiechtig.com/shorty/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenedLinkResolver(t *testing.T) {
	ctx := context.Background()

	testutil.WithDatabase(ctx, t, func(t *testing.T, db *pgxpool.Pool, s *store.Queries) {
		// create a new test server with a simple handler
		server := httptest.NewServer(resolveHandler(s))
		defer server.Close()

		t.Run("returns 302 when valid short code exists", func(t *testing.T) {
			// Create a client that doesn't follow redirects
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Get(server.URL + "/abc123")
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if got, want := resp.StatusCode, http.StatusFound; got != want {
				t.Fatalf("Expected status %v OK, got %v", want, got)
			}
		})

		t.Run("returns 404 when short code not found", func(t *testing.T) {
			resp, err := http.Get(server.URL + "/nonexistent")
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if got, want := resp.StatusCode, http.StatusNotFound; got != want {
				t.Fatalf("Expected status %v OK, got %v", want, got)
			}
		})

		t.Run("returns 405 for non-GET requests", func(t *testing.T) {
			resp, err := http.Post(server.URL+"/something", "", nil)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if got, want := resp.StatusCode, http.StatusMethodNotAllowed; got != want {
				t.Fatalf("Expected status %v OK, got %v", want, got)
			}
		})

		t.Run("handles empty short code", func(t *testing.T) {
			resp, err := http.Get(server.URL + "/")
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if got, want := resp.StatusCode, http.StatusNotFound; got != want {
				t.Fatalf("Expected status %v OK, got %v", want, got)
			}
		})

		//t.Run("redirects to correct URL", func(t *testing.T) {
		//	// test implementation
		//})
	})
}
