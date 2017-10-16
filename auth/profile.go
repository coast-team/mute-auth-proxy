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

func (p *Profile) FixUsername(provider string) {
	if p.Email == "" {
		p.Login = fmt.Sprintf("%s@%s", p.Login, provider)
	}
}

func (p Profile) Username() string {
	if p.Email != "" {
		return p.Email
	}
	return p.Login
}

func (p Profile) FullName() string {
	return p.Name
}

func (p Profile) Avatar() string {
	if p.AvatarURL != "" {
		return p.AvatarURL
	}
	return p.Picture
}
