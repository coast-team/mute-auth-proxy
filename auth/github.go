package auth

import (
	"net/http"

	"github.com/coast-team/mute-auth-proxy/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func MakeGithubLoginHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		githubOauthConfig := oauth2.Config{
			RedirectURL:  conf.OauthPrefs.RedirectURL,
			ClientID:     conf.OauthPrefs.GithubPrefs.ClientID,
			ClientSecret: conf.OauthPrefs.GithubPrefs.ClientSecret,
			Scopes:       []string{""},
			Endpoint:     github.Endpoint,
		}
		HandleProviderLogin(w, r, "github", githubOauthConfig)
	}
}
