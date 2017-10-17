package auth

import "fmt"

type Token struct {
	AccessToken string `json:"access_token"`
}

type Profile struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Picture   string `json:"picture"`
}

func (p Profile) Details() string {
	return fmt.Sprintf("Profile:\n\tlogin: %s\n\temail: %s\n\tfullname: %s\n\tavatar: %s", p.UserLogin(), p.Email, p.Name, p.Avatar())
}

func (p Profile) UserLogin() string {
	if p.Login != "" {
		return p.Login
	}
	return p.Email
}

func (p Profile) Avatar() string {
	if p.AvatarURL != "" {
		return p.AvatarURL
	}
	return p.Picture
}
