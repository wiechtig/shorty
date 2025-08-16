package resolver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"go.wiechtig.com/shorty/internal/store"
)

func ResolveHandler(s *store.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		shortCode := strings.TrimPrefix(r.URL.Path, "/")
		shortenedLink, err := s.ResolveShortenedLink(ctx, shortCode)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				slog.WarnContext(ctx, "Shortened link not found", slog.String("short_code", shortCode))
				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintln(w, http.StatusText(http.StatusNotFound))
				return
			}
			slog.ErrorContext(ctx, "Error resolving shortened link", slog.Any("error", err), slog.String("short_code", shortCode))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		slog.InfoContext(ctx, "Redirecting", slog.String("short_code", shortCode), slog.String("original_url", shortenedLink.OriginalUrl))
		http.Redirect(w, r, shortenedLink.OriginalUrl, http.StatusFound)
	}
}
