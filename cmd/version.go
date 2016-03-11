package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of MMS",
	Long:  `All software has versions. This is MMS's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MMS File Manager -- HEAD")
	},
}
