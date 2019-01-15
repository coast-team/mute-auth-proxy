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
	"os"
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
	Long: `Creates a file config.toml in the current working directory with the following content:

port = 4000
coniksserver_addr = "https://localhost:8400"
allowed_origins = ["http://localhost:4200"]

[oauth]
  [oauth.google]
    client_secret = "GOOGLE CLIENT SECRET"
  [oauth.github]
    client_secret = "GITHUB CLIENT SECRET"

Please fill this config file with the appropriate information.
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dest")
		if err != nil {
			log.Fatalf("Couldn't extract flag, error is : %s", err)
		}
		keyfilename, err := cmd.Flags().GetString("genkeyfile")
		if err != nil {
			log.Fatalf("Couldn't extract flag, error is : %s", err)
		}
		written := generateSymmetricKeyFile(dir, keyfilename)
		if written {
			fmt.Println("Symmetric key saved.")
		}
		written = generateConfigFile(dir)
		if written {
			fmt.Println("Please fill the generated config file.")
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("dest", "d", "./",
		"The path where to save the generated config file. (If the path denotes a directory then the config file path will be path/config.toml)")
	initCmd.Flags().StringP("genkeyfile", "k", "symmetric_key_file",
		"If this flag is specified, it will generate the symmetric key file (HMAC with SHA256 used for JWT signing) at the given location. The default location is ./symmetric_key_file")
}

// GenSymmetricKeyFile generates a key file with 256 bits symmetric key for HMAC.
func generateSymmetricKeyFile(dir, filepath string) bool {
	filepath = path.Join(dir, filepath)
	k := helper.GenerateRandomBytes()

	written, err := helper.WriteFile(filepath, k, 0600)
	if err != nil {
		log.Fatalf("Couldn't write keyfile.\nError was: %s\nMaybe all the directories in the path do not exist ?", err)
	}

	return written
}

func generateConfigFile(filepath string) bool {
	fileinfo, err := os.Stat(filepath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Couldn't get the file description of %s.\nError was: %s", filepath, err)
	}
	if !os.IsNotExist(err) && fileinfo.Mode().IsDir() {
		filepath = path.Join(filepath, "config.toml")
	}
	var conf = config.Config{
		Port:             4000,
		ConiksServerAddr: "http://localhost:8400",
		BotStorageAddr:   "http://localhost:4000",
		AllowedOrigins:   []string{"http://localhost:4200"},
		OauthPrefs: config.OauthConfig{
			GooglePrefs: config.ProviderPrefs{
				ClientSecret: "GOOGLE CLIENT SECRET",
			},
			GithubPrefs: config.ProviderPrefs{
				ClientSecret: "GITHUB CLIENT SECRET",
			},
		},
	}

	var confBuf bytes.Buffer
	enc := toml.NewEncoder(&confBuf)
	if err = enc.Encode(conf); err != nil {
		log.Fatalf("Couldn't encode config.\nError was: %s", err)
	}
	written, err := helper.WriteFile(filepath, confBuf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Couldn't write config file.\nError was: %s\nMaybe all the directories in the path do not exist ?", err)
	}
	return written
}
