/*
Copyright © 2023 Blair Gillam <ns1h@airmada.net>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the flexclient version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
