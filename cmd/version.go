package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixy/gateway/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version info",
	Long:  `print version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		version.Print()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
