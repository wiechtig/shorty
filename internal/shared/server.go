package shared

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPServerParams struct {
	Ctx    context.Context
	Wg     *sync.WaitGroup
	Server *http.Server
}

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shorty_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"component"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shorty_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"component"},
	)
)

func HTTPServer(params HTTPServerParams) {
	ctx := params.Ctx
	defer params.Wg.Done()

	server := params.Server

	// Start the HTTP server in a separate goroutine
	go func() {
		slog.InfoContext(ctx, "Starting HTTP server", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "Could not listen and serve", slog.String("address", server.Addr), slog.Any("error", err))
		}
	}()

	// Wait for the context to be canceled
	select {
	case <-ctx.Done():
		slog.Debug("Shutting down HTTP server gracefully...")

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			slog.ErrorContext(ctx, "HTTP server shutdown error", slog.Any("error", err))
		}
	}

	slog.Info("HTTP server stopped.")
}

func Telemetry(next http.Handler, component string) http.Handler {
	return Metrics(component)(Logging()(next))
}

func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(next, w, r)

			statusCode := strconv.Itoa(m.Code)

			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"request",
				slog.String("method", r.Method),
				slog.String("status_code", statusCode),
				slog.String("url", r.URL.Path),
				slog.String("ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
		})
	}
}

func Metrics(component string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(next, w, r)

			requestsTotal.WithLabelValues(
				component,
			).Inc()

			requestDuration.WithLabelValues(
				component,
			).Observe(m.Duration.Seconds())
		})
	}
}
