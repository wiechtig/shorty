package testutil

import (
	"context"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"wiechtig.com/shorty/internal/api"
	"wiechtig.com/shorty/internal/store"
)

func WithHTTPServer[TB testing.TB](ctx context.Context, tb TB, db *pgxpool.Pool, store *store.Queries, test func(t TB, handler http.Handler)) {
	mux := http.NewServeMux()
	server := api.New(db, store)
	handler := api.OpenAPIHandler(api.OpenAPIHandlerParams{
		Mux:      mux,
		Server:   server,
		BaseURL:  "/api",
		Verifier: nil,
	})

	test(tb, handler)
}
