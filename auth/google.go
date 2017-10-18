package auth

import (
	"net/http"

	"github.com/coast-team/mute-auth-proxy/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// MakeGoogleLoginHandler returns the handler for the Google login route
func MakeGoogleLoginHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		googleOauthConfig := oauth2.Config{
			RedirectURL:  conf.OauthPrefs.RedirectURL,
			ClientID:     conf.OauthPrefs.GooglePrefs.ClientID,
			ClientSecret: conf.OauthPrefs.GooglePrefs.ClientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
		HandleProviderLogin(w, r, "google", googleOauthConfig)
	}
}
