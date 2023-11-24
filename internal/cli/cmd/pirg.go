package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var PirgCmd = &cobra.Command{
	Use:   "pirg",
	Short: "Manage PIRGs",
}

var CreatePirgCmd = &cobra.Command{
	Use:   "create PIRGNAME",
	Short: "Create a new PIRG",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		// Add logic to create a group
		fmt.Printf("[todo] Group created: %s\n", groupName)
	},
}

var DeletePirgCmd = &cobra.Command{
	Use:   "delete PIRGNAME [--confirm]",
	Short: "Delete a PIRG",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		// Add logic to delete a group
		fmt.Printf("[todo] Group deleted: %s\n", groupName)
	},
}

var AddUserCmd = &cobra.Command{
	Use:   "add-user --pirg PIRGNAME username",
	Short: "Add a user to a PIRG",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName := args[0]
		userName := args[1]
		// Add logic to add a user to a group
		fmt.Printf("[todo] User %s added to PIRG %s\n", userName, pirgName)
	},
}

var RemoveUserCmd = &cobra.Command{
	Use:   "remove-user username --pirg PIRGNAME",
	Short: "Remove a user from a PIRG",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pirgName := args[1]
		userName := args[0]
		// Add logic to remove a user from a group
		fmt.Printf("[todo] User %s removed from PIRG %s\n", userName, pirgName)
	},
}

