package cli

import (
	"fmt"
	"log/slog"

	"github.com/lcrownover/hpcadmin-server/internal/config"
	"github.com/lcrownover/hpcadmin-server/internal/util"
	"github.com/spf13/cobra"
)

var (
	err       error
	debug     bool
	configDir string
	rootCmd   = &cobra.Command{
		Use:   "hpcadmin",
		Short: "HPCAdmin CLI",
		Long:  `HPCAdmin is a CMDB for hosting membership information for the Talapas HPC at University of Oregon`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			util.ConfigureLogging(debug)
			slog.Debug("Starting hpcadmin-cli", "method", "Execute")

			configDir, err = config.EnsureCLIConfigDir()
			if err != nil {
				util.PrintAndExit(fmt.Sprintf("Error reading configuration directory: %v\n", err), 1)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
}

func Execute() error {
	return rootCmd.Execute()
}
