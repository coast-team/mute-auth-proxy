// Copyright 2017-2018 Jean-Philippe Eisenbarth
//
// This file is part of Mute Authentication Proxy.
//
// Mute Authentication Proxy is free software: you can redistribute it and/or modify
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
// along with Mute Authentication Proxy. See the file COPYING.  If not, see <http://www.gnu.org/licenses/>.

package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coast-team/mute-auth-proxy/api"
	"github.com/coast-team/mute-auth-proxy/auth"
	"github.com/coast-team/mute-auth-proxy/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
		log.Fatalf("Couldn't load the config.\nError was: %s", err)
	}
	log.Println(conf)
	proxy := api.New("/botstorage", conf.BotStorageAddr)
	router := mux.NewRouter()
	router.HandleFunc("/auth/google", auth.MakeGoogleLoginHandler(conf))
	router.HandleFunc("/auth/github", auth.MakeGithubLoginHandler(conf))
	router.HandleFunc("/coniks", api.MakeConiksProxyHandler(conf))
	router.PathPrefix("/botstorage").HandlerFunc(api.MakeBotStorageProxyHandler(proxy))
	handlerFunc := handlers.CORS(handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}), handlers.AllowedOrigins(conf.AllowedOrigins))(router)
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), handlerFunc)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
