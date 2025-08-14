package shared

import (
	"log/slog"
	"os"
)

func SetupLogger(level string) {
	var logLevel = new(slog.LevelVar)
	err := logLevel.UnmarshalText([]byte(level))
	if err != nil {
		slog.Error("Failed to set log level", slog.Any("error", err))
		panic(err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
