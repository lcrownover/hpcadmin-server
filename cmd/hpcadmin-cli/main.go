package main

import (
	"fmt"
	"os"

	"github.com/lcrownover/hpcadmin-server/internal/cli"
)

func main() {
	var err error

	err = cli.Execute()
	if err != nil {
		fmt.Printf("Error executing cli: %v\n", err)
		os.Exit(1)
	}
}
