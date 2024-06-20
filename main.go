package main

import (
	"log"
	"os"
	"ssh-tunnel-helper/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func main() {
	cobra.OnInitialize(initConfig)

	cmd.RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ssh_tunnel.yaml)")

	cmd.RootCmd.AddCommand(cmd.StartCmd)

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".ssh_tunnel")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatalf("Error reading config file: %v", err)
	}
}
