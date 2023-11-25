package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lcrownover/hpcadmin-server/internal/config"
	"github.com/spf13/cobra"
)

var debug bool

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "hpcadmin",
		Short: "HPCAdmin CLI",
		Long:  `HPCAdmin is a CMDB for hosting membership information for the Talapas HPC at University of Oregon`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			ConfigureLogging(debug)
			slog.Debug("Starting hpcadmin-cli", "method", "Execute")

			_, err = config.EnsureConfigDir()
			if err != nil {
				PrintAndExit(fmt.Sprintf("Error reading configuration directory: %v\n", err), 1)
			}
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	pirgCmd.PersistentFlags().StringP("pirg", "p", "", "Specify PIRG name")
	pirgCmd.MarkPersistentFlagRequired("pirg")

	deletePirgCmd.Flags().Bool("confirm", false, "Confirm pirg deletion")

	addUserCmd.Flags().StringP("pirg", "p", "", "Specify PIRG name")
	addUserCmd.MarkFlagRequired("pirg")
	addUserCmd.PersistentFlags().StringP("user", "u", "", "Specify Username")
	addUserCmd.MarkPersistentFlagRequired("user")

	removeUserCmd.Flags().String("pirg", "", "Specify PIRG name")
	removeUserCmd.MarkFlagRequired("pirg")
	removeUserCmd.PersistentFlags().StringP("user", "u", "", "Specify Username")
	removeUserCmd.MarkPersistentFlagRequired("user")

	setPICmd.Flags().String("pirg", "", "Specify PIRG name")
	setPICmd.MarkFlagRequired("pirg")
	setPICmd.Flags().String("username", "", "Specify username")
	setPICmd.MarkFlagRequired("username")

	pirgCmd.AddCommand(createPirgCmd, deletePirgCmd, addUserCmd, removeUserCmd, setPICmd)
	rootCmd.AddCommand(LoginCmd, pirgCmd)
	return rootCmd.Execute()
}
