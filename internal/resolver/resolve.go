package resolver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"wiechtig.com/shorty/internal/store"
)

var (
	resolvedLinksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shorty_resolved_links_total",
			Help: "Total number of successfully resolved links",
		},
		[]string{"short_code", "original_url", "client_ip", "user_agent", "referer", "browser", "device_type"},
	)
)

func ResolveHandler(s *store.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodGet {
			status := http.StatusMethodNotAllowed
			w.WriteHeader(status)
			_, _ = fmt.Fprintln(w, http.StatusText(status))
			return
		}

		shortCode := strings.TrimPrefix(r.URL.Path, "/")
		shortenedLink, err := s.ResolveShortenedLink(ctx, shortCode)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				status := http.StatusNotFound
				slog.WarnContext(ctx, "Shortened link not found", slog.String("short_code", shortCode))
				w.WriteHeader(status)
				_, _ = fmt.Fprintln(w, http.StatusText(status))
				return
			}

			status := http.StatusInternalServerError
			slog.ErrorContext(ctx, "Error resolving shortened link", slog.Any("error", err), slog.String("short_code", shortCode))
			w.WriteHeader(status)
			_, _ = fmt.Fprintln(w, http.StatusText(status))
			return
		}

		status := http.StatusFound

		// Extract analytical metadata
		clientIP := extractClientIP(r)
		userAgent := r.Header.Get("User-Agent")
		referer := r.Header.Get("Referer")
		browser := extractBrowser(userAgent)
		deviceType := extractDeviceType(userAgent)

		resolvedLinksTotal.WithLabelValues(
			shortCode,
			shortenedLink.OriginalUrl,
			clientIP,
			userAgent,
			referer,
			browser,
			deviceType,
		).Inc()

		slog.InfoContext(ctx, "Redirecting",
			slog.String("short_code", shortCode),
			slog.String("original_url", shortenedLink.OriginalUrl),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
			slog.String("referer", referer),
			slog.String("browser", browser),
			slog.String("device_type", deviceType),
		)

		http.Redirect(w, r, shortenedLink.OriginalUrl, status)
	}
}

func extractClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func extractBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "chrome") && !strings.Contains(ua, "edge"):
		return "chrome"
	case strings.Contains(ua, "firefox"):
		return "firefox"
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		return "safari"
	case strings.Contains(ua, "edge"):
		return "edge"
	case strings.Contains(ua, "opera"):
		return "opera"
	default:
		return "other"
	}
}

func extractDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "mobile"):
		return "mobile"
	case strings.Contains(ua, "tablet"):
		return "tablet"
	default:
		return "desktop"
	}
}
