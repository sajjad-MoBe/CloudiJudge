package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		folder, err := cmd.Flags().GetString("folder")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filename, err := cmd.Flags().GetString("filename")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("dar hal test gereftan az code", folder+"/"+filename)
	},
}

func ExecuteServer() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("couldn't execute app,", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("folder", "/codes", "code, inputs and outputs should be placed here")
	rootCmd.PersistentFlags().String("filename", "test.go", "this file will be executed")

}
