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
		var config helpers.Configuration

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("unable to decode into struct: %v", err)
		}

		// Map servers by name for easy lookup
		servers := make(map[string]helpers.SshServerConfig)
		for _, server := range config.Servers {
			servers[server.Name] = server
		}

		// Display available connections
		fmt.Println("Available connections:")
		for i, conn := range config.SocksConnections {
			fmt.Printf("[%d] %s (SOCKS)\n", i, conn.Name)
		}
		offset := len(config.SocksConnections)
		for i, conn := range config.PortForwardConnections {
			fmt.Printf("[%d] %s (Port Forwarding)\n", offset+i, conn.Name)
		}

		// Prompt user for choice
		var choice int
		fmt.Print("Select a connection: ")
		fmt.Scan(&choice)

		if choice < 0 || choice >= (len(config.SocksConnections)+len(config.PortForwardConnections)) {
			log.Fatalf("invalid choice")
		}

		// Determine selected connection type and start tunnel
		if choice < len(config.SocksConnections) {
			selected := config.SocksConnections[choice]
			server, ok := servers[selected.SshServerConfig]
			if !ok {
				log.Fatalf("unknown server: %s", selected.SshServerConfig)
			}
			helpers.StartSocksTunnel(selected, server)
		} else {
			selected := config.PortForwardConnections[choice-len(config.SocksConnections)]
			server, ok := servers[selected.SshServerConfig]
			if !ok {
				log.Fatalf("unknown server: %s", selected.SshServerConfig)
			}
			helpers.StartPortForwarding(selected, server)
		}
	},
}
