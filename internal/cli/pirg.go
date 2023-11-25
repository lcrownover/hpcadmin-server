package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

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
		userName, _ := cmd.Flags().GetString("user")
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
		userName, _ := cmd.Flags().GetString("user")
		slog.Debug("removing user from PIRG", "method", "removeUserCmd.Run", "pirgName", pirgName, "userName", userName)
		fmt.Printf("[todo] User %s removed from PIRG %s\n", userName, pirgName)
	},
}

var setPICmd = &cobra.Command{
	Use:   "set-pi --pirg PIRGNAME --username USERNAME",
	Short: "Set the PIRG PI",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName, err := cmd.Flags().GetString("pirg")
		if err != nil {
			PrintAndExit("Required Argument: --pirg", 1)
		}
		userName, err := cmd.Flags().GetString("user")
		if err != nil {
			PrintAndExit("Required Argument: --user", 1)
		}
		slog.Debug("setting PIRG PI", "method", "setPICmd.Run", "pirgName", pirgName, "userName", userName)
		fmt.Printf("[todo] User %s set to PI for PIRG %s\n", userName, pirgName)
	},
}
