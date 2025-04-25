package cmd

import (
	"fmt"
	"os"

	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/code_runner"
	"github.com/spf13/cobra"
)

var codeRunnerCmd = &cobra.Command{
	Use:   "code-runner",
	Short: "to start code-runner service",
	Run: func(cmd *cobra.Command, args []string) {

		port, err := cmd.Flags().GetInt("listen")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if port < 1 || port > 79 {
			if port == 80 {
				fmt.Println("Error: please provide 'listen' flag with port between 1 and 79.")
			} else {
				fmt.Println("Error: The 'listen' port must be between 1 and 79.")
			}
			os.Exit(1)
		}

		fmt.Println("code runner is listening on port", port)
		code_runner.StartListening(port)
	},
}
