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
	"golang.org/x/oauth2/google"
)

// MakeGoogleLoginHandler returns the handler for the Google login route
func MakeGoogleLoginHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		googleOauthConfig := oauth2.Config{
			ClientSecret: conf.OauthPrefs.GooglePrefs.ClientSecret,
			Endpoint:     google.Endpoint,
		}
		err := handleProviderCallback(w, r, "google", googleOauthConfig)
		if err != nil {
			log.Println(err)
		}
	}
}

func setGoogleClaims(claims jwt.MapClaims, profile map[string]interface{}) {
	claims["provider"] = "google"
	claims["login"] = profile["email"]
	claims["name"] = profile["name"]
	claims["email"] = profile["email"]
	claims["avatar"] = profile["picture"]
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(9000 * time.Hour).Unix()
}
