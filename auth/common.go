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

type endPointList struct {
	Google string
	Github string
}

var apiEndpoint = endPointList{
	Github: "https://api.github.com/user",
	Google: "https://www.googleapis.com/oauth2/v3/userinfo",
}

func HandleProviderLogin(w http.ResponseWriter, r *http.Request, apiURL string, conf oauth2.Config) {
	helper.SetHeader(w, r)
	switch r.Method {
	case "POST":
		handleProviderCallback(w, r, apiURL, conf)
	default:
		w.WriteHeader(204)
		return
	}
}

func handleProviderCallback(w http.ResponseWriter, r *http.Request, apiURL string, conf oauth2.Config) {
	type requestData struct {
		Code string `json:"code"`
	}
	var data requestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("error : %s", err.Error())
	}
	log.Printf("Request Data : %v", data)

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	accessToken, err := conf.Exchange(oauth2.NoContext, data.Code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		return
	}
	log.Printf("access_token: %s", accessToken)
	client := conf.Client(oauth2.NoContext, accessToken)
	response, err := client.Get(apiURL)
	if err != nil {
		fmt.Printf("People API request failed '%s'\n", err)
		return
	}
	defer response.Body.Close()
	var profile Profile
	err = json.NewDecoder(response.Body).Decode(&profile)
	if err != nil {
		log.Printf("Err : %s", err.Error())
	}

	token2 := jwt.New(jwt.SigningMethodHS256)
	claims := token2.Claims.(jwt.MapClaims)

	if profile.Email == "" {
		profile.Login = fmt.Sprintf("%s@github", profile.Login)
		claims["username"] = profile.Login
	} else {
		claims["email"] = profile.Email
		claims["username"] = profile.Email
	}
	claims["name"] = profile.Name
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(9000 * time.Hour).Unix()

	tokenString, err := token2.SignedString(helper.Secret)
	if err != nil {
		log.Printf("failed to generate a JWT token, error was: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jwt := Token{AccessToken: tokenString}
	json.NewEncoder(w).Encode(jwt)
}
