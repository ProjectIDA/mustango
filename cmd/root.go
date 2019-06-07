package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	)

var RootCmd = &cobra.Command{
	Use:   "mustango",
	Short: "Query IRIS MUSTANG from command line",
	Long: `Command line program to query the IRIS MUSTANG Web Service`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

