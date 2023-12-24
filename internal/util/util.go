package util

import (
	"log/slog"
	"os"
)

func ConfigureLogging(debug bool) {
	// Set up logging
	lvl := slog.LevelInfo
	if debug {
		lvl = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
}
