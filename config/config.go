package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port             int
	ConiksServerAddr string      `toml:"coniksserver_addr"`
	OauthPrefs       OauthConfig `toml:"oauth"`
}

type OauthConfig struct {
	RedirectURL string        `toml:"redirect_url"`
	GooglePrefs ProviderPrefs `toml:"google"`
	GithubPrefs ProviderPrefs `toml:"github"`
}

type ProviderPrefs struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

func LoadConfig(file string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	return &conf, nil
}
