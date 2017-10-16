package auth

type Token struct {
	AccessToken string `json:"access_token"`
}

type Profile struct {
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
