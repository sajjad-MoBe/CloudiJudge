package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func ExecuteServer() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("couldn't execute app,", err)
		os.Exit(1)
	}
}

func init() {
	serveCmd.PersistentFlags().Int("listen", 80, "Application will be listening on this port")
	codeRunnerCmd.PersistentFlags().Int("listen", 2, "Application will be listening on this port")
	createAdminCmd.PersistentFlags().String("email", "-", "Admin will be created with this email")
	loadTestDataCmd.PersistentFlags().Bool("erase", false, "If true, all test data will be deleted")

	rootCmd.AddCommand(serveCmd, codeRunnerCmd, createAdminCmd, loadTestDataCmd)

}
