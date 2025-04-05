package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var codeRunnerCmd = &cobra.Command{
	Use:   "code-runner",
	Short: "to start code-runner service",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("code runner is started")

	},
}
