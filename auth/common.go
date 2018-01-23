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

package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coast-team/mute-auth-proxy/helper"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

var apiEndpoint = map[string]string{
	"github": "https://api.github.com/user",
	"google": "https://www.googleapis.com/oauth2/v3/userinfo",
}

// HandleProviderLogin is the generic handler for either Google and Github login route.
// It needs a oauth2.Config parameter
func HandleProviderLogin(w http.ResponseWriter, r *http.Request, provider string, conf oauth2.Config) error {
	helper.SetHeader(w, r)
	switch r.Method {
	case "POST":
		return handleProviderCallback(w, r, provider, conf)
	default:
		w.WriteHeader(204)
		return nil
	}
}

func handleProviderCallback(w http.ResponseWriter, r *http.Request, provider string, conf oauth2.Config) error {
	type requestData struct {
		AuthorizationData struct {
			ClientID    string `json:"client_id"`
			RedirectURI string `json:"redirect_uri"`
		} `json:"authorizationData"`
		OAuthData struct {
			Code string `json:"code"`
		} `json:"oauthData"`
	}

	var data requestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.Write([]byte("Couldn't decode request's body."))
		return fmt.Errorf("Couldn't decode request's body.\nError was: %s", err)
	}
	conf.ClientID = data.AuthorizationData.ClientID
	conf.RedirectURL = data.AuthorizationData.RedirectURI

	accessToken, err := conf.Exchange(oauth2.NoContext, data.OAuthData.Code)
	if err != nil {
		w.Write([]byte("Code exchange failed."))
		return fmt.Errorf("Code exchange failed.\nError was: %s", err)
	}
	client := conf.Client(oauth2.NoContext, accessToken)
	response, err := client.Get(apiEndpoint[provider])
	if err != nil {
		w.Write([]byte("Server internal error."))
		return fmt.Errorf("People API request failed.\nError was: %s", err)
	}
	defer response.Body.Close()
	var profile Profile
	err = json.NewDecoder(response.Body).Decode(&profile)
	if err != nil {
		w.Write([]byte("Server internal error."))
		return fmt.Errorf("Couldn't decode %s's response.\nError was: %s", provider, err)
	}
	log.Println(profile.Details())

	token2 := jwt.New(jwt.SigningMethodHS256)
	claims := token2.Claims.(jwt.MapClaims)
	claims["provider"] = provider
	claims["login"] = profile.UserLogin()
	claims["name"] = profile.Name
	claims["email"] = profile.Email
	claims["avatar"] = profile.Avatar()
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(9000 * time.Hour).Unix()

	tokenString, err := token2.SignedString(helper.Secret)
	if err != nil {
		w.Write([]byte("Server internal error."))
		return fmt.Errorf("Failed to generate a JWT token.\nError was: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jwt := Token{AccessToken: tokenString}
	json.NewEncoder(w).Encode(jwt)
	return nil
}
