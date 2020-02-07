package cmd

import (
	"fmt"

	"github.com/yixy/gateway/server"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start http gateway service",
	Long:  `Start a http server for receive API request.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gateway service start")
		err := server.Start()
		if err != nil {
			fmt.Println("gateway service start err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
