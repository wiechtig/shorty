package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"wiechtig.com/shorty/internal/resolver"
	"wiechtig.com/shorty/internal/shared"
	"wiechtig.com/shorty/internal/store"
)

func main() {
	// Create a context to handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a WaitGroup to keep track of running goroutines
	var wg sync.WaitGroup

	logLevel := os.Getenv("SHORTY_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	shared.SetupLogger(logLevel)

	databaseUrl := os.Getenv("SHORTY_DATABASE_URL")
	if databaseUrl == "" {
		panic("database url is required")
	}
	dbpool := shared.SetupDatabase(databaseUrl)
	shared.RunMigrations(dbpool, "db/migrations")
	s := store.New(dbpool)

	// Setup server and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", resolver.ResolveHandler(s))
	server := &http.Server{
		Addr:         ":4242",
		Handler:      shared.Telemetry(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	mm := http.NewServeMux()
	mm.Handle("/metrics", promhttp.Handler())
	metrics := &http.Server{
		Addr:         ":4343",
		Handler:      shared.Telemetry(mm),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the metrics server
	wg.Add(1)
	go shared.HTTPServer(shared.HTTPServerParams{
		Ctx:    ctx,
		Wg:     &wg,
		Server: metrics,
	})

	// Start the HTTP server
	wg.Add(1)
	go shared.HTTPServer(shared.HTTPServerParams{
		Ctx:    ctx,
		Wg:     &wg,
		Server: server,
	})

	// Listen for termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-signalCh

	// Start the graceful shutdown process
	slog.Info("Gracefully shutting down...")

	// Cancel the context to signal the HTTP server to stop
	cancel()

	// Wait for the HTTP server to finish
	wg.Wait()

	slog.Info("Shutdown complete.")
}
