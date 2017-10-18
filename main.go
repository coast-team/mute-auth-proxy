package main

import (
	"fmt"
	"os"

	"github.com/coast-team/mute-auth-proxy/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
