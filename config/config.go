package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config represents the structure containing the information from the config file
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

// LoadConfig loads and parses the information from the config file and fill the Config struct
func LoadConfig(file string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	return &conf, nil
}
