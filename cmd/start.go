package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"ssh-tunnel-helper/helpers"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an SSH tunnel",
	Long:  `Start an SSH tunnel for SOCKS proxy or port forwarding.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config struct {
			Connections []helpers.Connection `mapstructure:"connections"`
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("unable to decode into struct: %v", err)
		}

		fmt.Println("Available connections:")
		for i, conn := range config.Connections {
			fmt.Printf("[%d] %s\n", i, conn.Name)
		}

		var choice int
		fmt.Print("Select a connection: ")
		fmt.Scan(&choice)

		if choice < 0 || choice >= len(config.Connections) {
			log.Fatalf("invalid choice")
		}

		selected := config.Connections[choice]
		if selected.Type == "socks" {
			helpers.StartSOCKSTunnel(selected)
		} else if selected.Type == "portforward" {
			helpers.StartPortForwarding(selected)
		} else {
			log.Fatalf("unknown connection type: %s", selected.Type)
		}
	},
}
