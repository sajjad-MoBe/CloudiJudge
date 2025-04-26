package cmd

import (
	"fmt"
	"net/mail"
	"os"

	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/serve"
	"github.com/spf13/cobra"
)

var createAdminCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "create an admin",
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = mail.ParseAddress(email)
		if err != nil {
			fmt.Println("Please enter a valid email.")
			os.Exit(1)
		}

		serve.CreateAdmin(email)
	},
}
