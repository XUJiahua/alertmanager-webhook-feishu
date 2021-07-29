package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "binary version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s, commit %s, built at %s by %s\n", version, commit, date, builtBy)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
