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
	"log"

	"github.com/coast-team/mute-auth-proxy/auth"

	"github.com/coast-team/mute-auth-proxy/helper"
	"github.com/spf13/cobra"
)

// RunCmd represents the run commands. It starts the web server.
var genJWTCmd = &cobra.Command{
	Use:   "generate-jwt",
	Short: "Generate a JWT.",
	Long:  `Generate a valid signed JWT. To be used with the BotStorage for exemple.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		genjwt(cmd)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(genJWTCmd)
	genJWTCmd.Flags().StringP("botlogin", "l", "botlogin", "The login of the Bot (bot.storage for example)")
	genJWTCmd.Flags().StringP("keyfile", "k", "symmetric_key_file", "The key file (HMAC with SHA256 used for JWT signing) to load")
}

func genjwt(cmd *cobra.Command) {
	var err error
	keyfilepath, err := cmd.Flags().GetString("keyfile")
	if err != nil {
		log.Fatalf("Couldn't extract flag, error is : %s", err)
	}
	botlogin, err := cmd.Flags().GetString("botlogin")
	if err != nil {
		log.Fatalf("Couldn't extract flag, error is : %s", err)
	}
	keyData, err := helper.ReadFile(keyfilepath)
	if err != nil {
		log.Fatalf("Couldn't load the keyfile.\nError was: %s", err)
	}
	helper.SetSecret(keyData)
	token := helper.GenerateJWT()
	auth.SetClaims(token, map[string]interface{}{"login": botlogin}, "bot")
	tokenString, err := helper.GetSignedString(token)
	if err != nil {
		log.Fatalf("Couldn't sign the jwt, error is : %s", err)
	}
	log.Println(tokenString)
}
