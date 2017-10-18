package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "mute-auth-proxy",
	Short: "Mute Authentication Proxy in Go.",
	Long:  `Mute Authentication Proxy in Go. It handles OAUTH login and proxies the ConiksClient request to the ConiksServer.`,
}

// Execute adds all child commands to the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
