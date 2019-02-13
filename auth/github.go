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

package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/coast-team/mute-auth-proxy/config"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// MakeGithubLoginHandler returns the handler for the Github login route
func MakeGithubLoginHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		githubOauthConfig := oauth2.Config{
			ClientSecret: conf.OauthPrefs.GithubPrefs.ClientSecret,
			Endpoint:     github.Endpoint,
		}
		err := handleProviderCallback(w, r, "github", githubOauthConfig)
		if err != nil {
			log.Println(err)
		}
	}
}

func setGithubClaims(claims jwt.MapClaims, profile map[string]interface{}) {
	claims["provider"] = "github"
	claims["login"] = profile["login"]
	claims["name"] = profile["name"]
	claims["email"] = profile["email"]
	claims["avatar"] = profile["avatar_url"]
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(9000 * time.Hour).Unix()
}
