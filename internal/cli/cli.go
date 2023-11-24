package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lcrownover/hpcadmin-server/internal/cli/cmd"
	"github.com/lcrownover/hpcadmin-server/internal/config"
	"github.com/spf13/cobra"
)

var debug bool

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "hpcadmin",
		Short: "HPCAdmin CLI",
		Long:  `HPCAdmin is a CMDB for hosting membership information for the Talapas HPC at University of Oregon`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			// Set up logging
			lvl := slog.LevelInfo
			if err != nil {
				fmt.Printf("Error reading debug flag: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("doing run stuff")
			if debug {
				lvl = slog.LevelDebug
			}
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: lvl,
			}))
			slog.SetDefault(logger)

			slog.Debug("Starting hpcadmin-cli")
			_, err = config.EnsureConfigDir()
			if err != nil {
				fmt.Printf("Error reading configuration directory: %v\n", err)
				os.Exit(1)
			}
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	cmd.PirgCmd.Flags().Bool("confirm", false, "Confirm pirg deletion")

	cmd.AddUserCmd.Flags().String("pirg", "", "Specify PIRG name")
	cmd.AddUserCmd.MarkFlagRequired("pirg")

	cmd.RemoveUserCmd.Flags().String("pirg", "", "Specify PIRG name")
	cmd.RemoveUserCmd.MarkFlagRequired("pirg")

	cmd.PirgCmd.AddCommand(cmd.CreatePirgCmd, cmd.DeletePirgCmd, cmd.AddUserCmd, cmd.RemoveUserCmd)
	rootCmd.AddCommand(cmd.LoginCmd, cmd.PirgCmd)
	return rootCmd.Execute()
}
