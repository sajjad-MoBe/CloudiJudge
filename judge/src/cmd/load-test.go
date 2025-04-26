package cmd

import (
	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/load_test"

	"github.com/spf13/cobra"
)

var loadTestCmd = &cobra.Command{
	Use:   "load-test",
	Short: "generate test data",
	Run: func(cmd *cobra.Command, args []string) {
		var erase bool
		erase, err := cmd.Flags().GetBool("erase")
		if err != nil {
			erase = false
		}
		if !erase {
			load_test.GenerateAndFill()
		} else {
			load_test.Erase()
		}
	},
}
