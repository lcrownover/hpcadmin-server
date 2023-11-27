package cli

import (
	"fmt"
	"os"

	"github.com/lcrownover/hpcadmin-server/internal/auth"
	"github.com/lcrownover/hpcadmin-server/internal/util"
	"github.com/spf13/cobra"
)

const AZURE_TENANT_ID = "8f0b198f-f447-4cfe-ba03-526b46c661f8"
const AZURE_CLIENT_ID = "1951f213-c370-4a77-b7cd-7a4c303df45a"

func init() {
	rootCmd.AddCommand(LoginCmd)
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to HPCAdmin",
	Run: func(cmd *cobra.Command, args []string) {
		azureAuthOptions := auth.AzureAuthHandlerOptions{
			TenantID:  AZURE_TENANT_ID,
			ClientID:  AZURE_CLIENT_ID,
			ConfigDir: configDir,
		}
		ah := auth.NewAuthHandler(azureAuthOptions)

		err = ah.Authenticate()
		if err != nil {
			util.ErrorPrint(fmt.Sprintf("Error authenticating: %v\n", err))
			os.Exit(1)
		}
	},
}
