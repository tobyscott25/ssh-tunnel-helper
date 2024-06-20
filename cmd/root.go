package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ssh-tunnel-helper",
	Short: "A CLI tool to manage SSH tunnels",
	Long:  `A CLI tool to manage SSH tunnels for SOCKS proxies and port forwarding using configuration from a YAML file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cfgFile string
var configDir string = "$HOME/.config/ssh-tunnel-helper"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+configDir+"/config.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "Toby Scott <hi@tobyscott.dev>", "Author name for copyright attribution")

	rootCmd.AddCommand(StartCmd)
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("\033[1;31mFailed to read configuration file!\033[0m\n")
		fmt.Printf("\033[1;31m%v\033[0m\n", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Please create a configuration file at %s/config.yaml\n", configDir)
		}
		os.Exit(1)
	}
}
