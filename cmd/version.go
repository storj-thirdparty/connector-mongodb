package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the cli",
	Long:  `Prints the version of the cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 1.0.5")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
