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
	"encoding/json"
	"fmt"
	"io"
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

type requestData struct {
	AuthorizationData struct {
		ClientID    string `json:"client_id"`
		RedirectURI string `json:"redirect_uri"`
	} `json:"authorizationData"`
	OAuthData struct {
		Code string `json:"code"`
	} `json:"oauthData"`
}

type Token struct {
	AccessToken string `json:"access_token"`
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
		return fmt.Errorf("HTTP methods wasn't POST")
	}
}

func handleProviderCallback(w http.ResponseWriter, r *http.Request, provider string, conf oauth2.Config) error {
	var data requestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("Couldn't decode request's body, maybe CORS issue ?\nError was: %s", err)
		}
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
	client.Timeout = time.Duration(5) * time.Second
	response, err := client.Get(apiEndpoint[provider])
	if err != nil {
		w.Write([]byte("Server internal error."))
		return fmt.Errorf("People API request failed.\nError was: %s", err)
	}
	defer response.Body.Close()

	var profile map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&profile)
	if err != nil {
		w.Write([]byte("Server internal error."))
		return fmt.Errorf("Couldn't decode %s's response.\nError was: %s", provider, err)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	setClaims(token, profile, provider)
	tokenString, err := token.SignedString(helper.Secret)
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

func setClaims(token *jwt.Token, profile map[string]interface{}, provider string) {
	switch provider {
	case "github":
		setGithubClaims(token, profile)
	case "google":
		setGoogleClaims(token, profile)
	}
}
