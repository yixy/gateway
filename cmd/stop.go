package cmd

import (
	"fmt"

	"github.com/yixy/gateway/server"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop http gateway service",
	Long:  `Stop the API gateway server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gateway service stop")
		err := server.Stop()
		if err != nil {
			fmt.Println("gateway service stop err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
