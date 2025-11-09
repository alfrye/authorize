package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


	 var rootCmd = &cobra.Command{
				Use: "Authorize is a Identity and Access application",
				Short: "Authorize contains an api for all Identity and Access Management",
				Long: "long version on the app",
		}
	 // startCmd *cobra.Command
	//1	err error
//    if startCmd,err = NewStartCommand(); err != nil {
//      return nil, err
// }
//rootCmd.AddCommand(startCmd)

//	return rootCmd, nil
//}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error occurred executing root command")
		os.Exit(1)
	}
}
