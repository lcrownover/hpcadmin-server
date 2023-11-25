package cli

import (
	"fmt"
	"log/slog"
	"os"
)

func PrintAndExit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

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


