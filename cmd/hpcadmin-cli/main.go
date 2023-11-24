package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lcrownover/hpcadmin-server/internal/cli"
)

func main() {
	var err error
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Debug("Starting hpcadmin-cli")

	err = cli.Execute()
	if err != nil {
		fmt.Printf("Error executing cli: %v\n", err)
		os.Exit(1)
	}
}
