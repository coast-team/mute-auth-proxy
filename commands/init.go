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
	"path"

	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/coast-team/mute-auth-proxy/config"
	"github.com/coast-team/mute-auth-proxy/helper"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a config file.",
	Long: `Creates a file config.toml in the current working directory with
the following content:

port = 4000
coniksserver_addr = "https://localhost:8400"

[oauth]
  redirect_url = "Web Client URL"
  [oauth.google]
    client_id = "GOOGLE CLIENT ID"
    client_secret = "GOOGLE CLIENT ID"
  [oauth.github]
    client_id = "GITHUB CLIENT ID"
    client_secret = "GITHUB CLIENT ID"

Please fill this config file with the appropriate information.
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := cmd.Flag("dest").Value.String()
		generateConfigFile(dir)
		fmt.Println("Please fill in the generated config file.")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("dest", "d", ".",
		"Location of the directory where to save the generated config file")
}

func generateConfigFile(dest string) {
	file := path.Join(dest, "config.toml")
	var conf = config.Config{
		Port:             4000,
		ConiksServerAddr: "http://localhost:8400",
		OauthPrefs: config.OauthConfig{
			RedirectURL: "Web Client URL",
			GooglePrefs: config.ProviderPrefs{
				ClientID:     "GOOGLE CLIENT ID",
				ClientSecret: "GOOGLE CLIENT SECRET",
			},
			GithubPrefs: config.ProviderPrefs{
				ClientID:     "GITHUB CLIENT ID",
				ClientSecret: "GITHUB CLIENT SECRET",
			},
		},
	}

	var confBuf bytes.Buffer
	enc := toml.NewEncoder(&confBuf)
	if err := enc.Encode(conf); err != nil {
		log.Fatalf("Coulnd't encode config. %s", err.Error())
	}
	if err := helper.WriteFile(file, confBuf.Bytes(), 0644); err != nil {
		log.Fatalf("Coulnd't write config: %s", err.Error())
	}
}