package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

)
	var rootCmd = &cobra.Command{
	  Use: "Authorize is a Identity and Access application",
		Short: "Authorize contains an api for all Identity and Access Management",
		Long: "long version on the app",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing authorize '%v'\n", err)
		os.Exit(1)
	}
}


