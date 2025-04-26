package cmd

import (
	"fmt"
	"os"

	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/serve"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "to start server",
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("listen")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if port < 80 || port > 65535 {
			fmt.Println("Error: The 'listen' port must be between 80 and 65535.")
			os.Exit(1)
		}

		fmt.Println("server is listening on port", port)
		serve.StartListening(port)
	},
}
