package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "ssh-tunnel-helper",
	Short: "A CLI tool to manage SSH tunnels",
	Long:  `A CLI tool to manage SSH tunnels for SOCKS proxies and port forwarding using configuration from a YAML file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)
}
