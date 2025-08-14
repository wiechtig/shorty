package shared

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type HTTPServerParams struct {
	Ctx    context.Context
	Wg     *sync.WaitGroup
	Server *http.Server
}

func HTTPServer(params HTTPServerParams) {
	ctx := params.Ctx
	defer params.Wg.Done()

	server := params.Server

	// Start the HTTP server in a separate goroutine
	go func() {
		slog.InfoContext(ctx, "Starting HTTP server", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
