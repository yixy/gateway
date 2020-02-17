package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yixy/gateway/cfg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A Simple HTTP API gateway",
	Long:  `A Simple HTTP API gateway supporting RSA asymmetric key signature verification.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute", err)
	}
}

func init() {
	//OnInitialize sets the passed functions to be run when each command's Execute method is called.
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	path, err := os.Executable()
	if err != nil {
		fmt.Println("find execute binary dir error", err)
		os.Exit(1)
	}
	cfg.Dir = filepath.Dir(path)
	fmt.Println("================ print execute binary info ================")
	fmt.Println(path)    // for example /home/user/main
	fmt.Println(cfg.Dir) // for example /home/user
	fmt.Println("================            end            ================")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".gateway" (without extension).
		viper.AddConfigPath(cfg.Dir)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file not found.", err)
		} else {
			// Config file was found but another error was produced
			fmt.Println("Config file was found but another error was produced.", err)
		}
		os.Exit(1)
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err = cfg.ReadCfg()
	if err != nil {
		fmt.Println("config file is invalid.", err)
		os.Exit(1)
	}
}
