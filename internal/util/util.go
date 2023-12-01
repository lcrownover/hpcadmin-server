package util

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
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

func SuccessPrint(msg string) {
	color.Green(msg)
}

func InfoPrint(msg string) {
	color.Blue(msg)
}

func ErrorPrint(msg string) {
	color.Red(msg)
}
