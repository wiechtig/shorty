package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	oapiTestutil "github.com/oapi-codegen/testutil"
	"wiechtig.com/shorty/internal/store"
	"wiechtig.com/shorty/internal/testutil"
)

// todo verify access control across paths; no need for db

func TestShortCodes(t *testing.T) {
	ctx := context.Background()

	testutil.WithDatabase(ctx, t, func(t *testing.T, db *pgxpool.Pool, s *store.Queries) {
		testutil.WithHTTPServer(ctx, t, db, s, func(t *testing.T, handler http.Handler) {

			// todo add happy path tests for creating, retrieving, deleting short codes

			t.Run("list short codes should return empty array", func(t *testing.T) {
				rr := oapiTestutil.NewRequest().
					Get("/api/purchases").
					GoWithHTTPHandler(t, handler).
					Recorder

				if rr.Code != http.StatusOK {
					t.Errorf("Expected status 200 OK, got %v", rr.Code)
				}

				var response []map[string]interface{}
				if err := json.NewDecoder(rr.Result().Body).Decode(&response); err != nil {
					t.Errorf("Failed to parse response: %v", err)
				}
				t.Logf("Response: %v", response)

				if response == nil {
					t.Errorf("Expected response to be present")
				}

				if got, want := len(response), 0; got != want {
					t.Errorf("count=%d, want=%d", got, want)
				}
			})

		})
	})
}
