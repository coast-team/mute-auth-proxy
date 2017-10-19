// Copyright 2017 Jean-Philippe Eisenbarth
//
// This file is part of Mute Authentication Proxy.
//
// Foobar is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Mute Authentication Proxy is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Foobar. See the file COPYING.  If not, see <http://www.gnu.org/licenses/>.

package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coast-team/mute-auth-proxy/api"
	"github.com/coast-team/mute-auth-proxy/auth"
	"github.com/coast-team/mute-auth-proxy/config"
	"github.com/spf13/cobra"
)

// RunCmd represents the run commands. It starts the web server.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Mute Authentication Proxy.",
	Long:  `Run the Mute Authentication Proxy. This is the main command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		run(cmd)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("config", "c", "config.toml", "The config file to load")
}

func run(cmd *cobra.Command) {
	confFilename := cmd.Flag("config").Value.String()
	conf, err := config.LoadConfig(confFilename)
	if err != nil {
		log.Fatalf("LoadConfig: %s", err)
	}
	log.Printf("Conf: %+v", conf)
	http.HandleFunc("/auth/google", auth.MakeGoogleLoginHandler(conf))
	http.HandleFunc("/auth/github", auth.MakeGithubLoginHandler(conf))
	http.HandleFunc("/coniks", api.MakeConiksProxyHandler(conf))
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
