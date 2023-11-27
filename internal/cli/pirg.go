package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

func init() {
	// Global flags for all pirg commands
	pirgCmd.PersistentFlags().StringP("pirg", "p", "", "Specify PIRG name")
	pirgCmd.MarkPersistentFlagRequired("pirg")

	// Create a pirg
	pirgCmd.AddCommand(createPirgCmd)

	// Delete a pirg
	deletePirgCmd.Flags().BoolP("confirm", "c", false, "Confirm pirg deletion")
	pirgCmd.AddCommand(deletePirgCmd)

	// Add a user to a pirg
	addUserCmd.PersistentFlags().StringP("username", "u", "", "Specify Username")
	addUserCmd.MarkFlagRequired("username")
	pirgCmd.AddCommand(addUserCmd)

	// Remove a user from a pirg
	removeUserCmd.PersistentFlags().StringP("username", "u", "", "Specify Username")
	removeUserCmd.MarkFlagRequired("username")
	pirgCmd.AddCommand(removeUserCmd)

	// Set the PI for a pirg
	setPICmd.Flags().StringP("username", "u", "", "Specify username")
	setPICmd.MarkFlagRequired("username")
	pirgCmd.AddCommand(setPICmd)

	// Add the pirg command to the root command
	rootCmd.AddCommand(pirgCmd)
}

var pirgCmd = &cobra.Command{
	Use:   "pirg",
	Short: "Manage PIRGs",
}

var createPirgCmd = &cobra.Command{
	Use:   "create --pirg PIRGNAME",
	Short: "Create a new PIRG",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName, _ := cmd.Flags().GetString("pirg")
		slog.Debug("creating pirg", "method", "createPirgCmd.Run", "pirgName", pirgName)
		fmt.Printf("[todo] PIRG created: %s\n", pirgName)
	},
}

var deletePirgCmd = &cobra.Command{
	Use:   "delete --pirg PIRGNAME [--confirm]",
	Short: "Delete a PIRG",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName, _ := cmd.Flags().GetString("pirg")
		confirm, _ := cmd.Flags().GetBool("confirm")
		slog.Debug("deleting PIRG", "method", "deletePirgCmd.Run", "pirgName", pirgName, "confirm", confirm)
		fmt.Printf("[todo] PIRG deleted: %s\n", pirgName)
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add-user --pirg PIRGNAME --username USERNAME",
	Short: "Add a user to a PIRG",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName, _ := cmd.Flags().GetString("pirg")
		userName, _ := cmd.Flags().GetString("username")
		slog.Debug("adding user to PIRG", "method", "addUserCmd.Run", "pirgName", pirgName, "userName", userName)
		fmt.Printf("[todo] User %s added to PIRG %s\n", userName, pirgName)
	},
}

var removeUserCmd = &cobra.Command{
	Use:   "remove-user --pirg PIRGNAME --username USERNAME",
	Short: "Remove a user from a PIRG",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName, _ := cmd.Flags().GetString("pirg")
		userName, _ := cmd.Flags().GetString("username")
		slog.Debug("removing user from PIRG", "method", "removeUserCmd.Run", "pirgName", pirgName, "userName", userName)
		fmt.Printf("[todo] User %s removed from PIRG %s\n", userName, pirgName)
	},
}

var setPICmd = &cobra.Command{
	Use:   "set-pi --pirg PIRGNAME --username USERNAME",
	Short: "Set the PIRG PI",
	// Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName, _ := cmd.Flags().GetString("pirg")
		userName, _ := cmd.Flags().GetString("username")
		slog.Debug("setting PIRG PI", "method", "setPICmd.Run", "pirgName", pirgName, "userName", userName)
		fmt.Printf("[todo] User %s set to PI for PIRG %s\n", userName, pirgName)
	},
}
